package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// fetchDomain queries /api/domain/{domain}. Retries on transient (5xx / network)
// failures up to 3 times with exponential backoff. Surfaces 401/402/429 cleanly.
func fetchDomain(ctx context.Context, client *http.Client, domain, apiKey string) ([]SubdomainData, int, error) {
	url := fmt.Sprintf("%s/api/domain/%s", defaultBaseURL, domain)
	ua := fmt.Sprintf("reconeer-cli/%s (%s/%s)", version, runtime.GOOS, runtime.GOARCH)

	var lastErr error
	var lastStatus int
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s (clipped to 5s)
			wait := time.Duration(1<<attempt) * time.Second
			if wait > 5*time.Second {
				wait = 5 * time.Second
			}
			select {
			case <-ctx.Done():
				return nil, lastStatus, ctx.Err()
			case <-time.After(wait):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("User-Agent", ua)
		req.Header.Set("Accept", "application/json")
		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue // retry on network error
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		lastStatus = resp.StatusCode

		// Retry on 5xx
		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			lastErr = fmt.Errorf("server error %d: %s", resp.StatusCode, truncate(string(body), 200))
			continue
		}
		// 4xx is fatal — don't retry
		if resp.StatusCode >= 400 {
			return nil, resp.StatusCode, fmt.Errorf("%s", truncate(string(body), 400))
		}

		var parsed apiDomainResponse
		if err := json.Unmarshal(body, &parsed); err != nil {
			return nil, resp.StatusCode, fmt.Errorf("invalid JSON response: %v", err)
		}
		return parsed.Subdomains, resp.StatusCode, nil
	}
	return nil, lastStatus, lastErr
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func readDomainsFromFileOrStdin(path string) ([]string, error) {

	if path == "-" {
		return readDomainsFromReader(os.Stdin)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return readDomainsFromReader(f)
}

func readDomainsFromReader(r io.Reader) ([]string, error) {

	sc := bufio.NewScanner(r)

	var out []string

	for sc.Scan() {

		line := strings.TrimSpace(sc.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		out = append(out, line)
	}

	return out, sc.Err()
}

func normalizeDomains(in []string) []string {

	out := []string{}
	seen := map[string]struct{}{}

	for _, d := range in {

		d = strings.TrimSpace(d)
		d = strings.TrimPrefix(d, "http://")
		d = strings.TrimPrefix(d, "https://")
		d = strings.TrimSuffix(d, "/")

		if d == "" {
			continue
		}

		if _, ok := seen[d]; ok {
			continue
		}

		seen[d] = struct{}{}
		out = append(out, d)
	}

	return out
}

func printHelpHint() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  reconeer -d example.com")
	fmt.Fprintln(os.Stderr, "  reconeer -d example.com -k YOUR_API_KEY")
	fmt.Fprintln(os.Stderr, "  RECONEER_API_KEY=YOUR_KEY reconeer -d example.com")
	fmt.Fprintln(os.Stderr, "  reconeer -dL scope.txt -silent -o out.txt")
	fmt.Fprintln(os.Stderr, "  reconeer -d example.com -jsonl | jq -r '.subdomain'")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Get a free API key: https://www.reconeer.com/register")
	fmt.Fprintln(os.Stderr, "Full docs:          https://www.reconeer.com/docs")
}

func outputWriter(path string) (io.Writer, func(), error) {

	if path == "" {
		return os.Stdout, func() {}, nil
	}

	dir := filepath.Dir(path)

	if dir != "." && dir != "" {
		os.MkdirAll(dir, 0755)
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, func() {}, err
	}

	return f, func() { f.Close() }, nil
}

func fatalIf(err error, silent bool, format string, args ...any) {

	if err == nil {
		return
	}

	if !silent {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}

	os.Exit(1)
}

func handleAPIError(err error, status int, silent bool) {
	if silent {
		return
	}
	switch status {
	case http.StatusUnauthorized:
		fmt.Fprintln(os.Stderr, "[!] Invalid or missing API key.")
		fmt.Fprintln(os.Stderr, "    Get one at:", signupURL)
	case http.StatusPaymentRequired:
		// Server uses 402 for Premium-only endpoints (CSV export, bulk).
		fmt.Fprintln(os.Stderr, "[!] This endpoint requires Premium.")
		fmt.Fprintln(os.Stderr, "    Upgrade at:", pricingURL)
	case http.StatusTooManyRequests:
		fmt.Fprintln(os.Stderr, "[!] Daily quota reached.")
		fmt.Fprintln(os.Stderr, "    Upgrade for unlimited:", pricingURL)
	case http.StatusNotFound:
		// Per-domain 404 just means no data yet — not an error worth dwelling on.
		// Quiet by default; the main loop prints a summary at the end.
	default:
		fmt.Fprintln(os.Stderr, "[!] API error:", err)
	}
}

type rateLimiter struct {
	ch <-chan time.Time
}

func newRateLimiter(rps int) *rateLimiter {

	if rps <= 0 {
		rps = 3
	}

	interval := time.Second / time.Duration(rps)

	t := time.NewTicker(interval)

	return &rateLimiter{ch: t.C}
}

func (r *rateLimiter) Wait() {
	<-r.ch
}
