package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// fetchDomain — GET /api/domain/{domain} returning parsed subdomain list
func fetchDomain(ctx context.Context, client *http.Client, domain, apiKey string) ([]SubdomainData, int, error) {
	url := fmt.Sprintf("%s/api/domain/%s", defaultBaseURL, domain)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Accept", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("%s", string(body))
	}
	var parsed apiDomainResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, resp.StatusCode, err
	}
	return parsed.Subdomains, resp.StatusCode, nil
}

// fetchDomainRaw — GET /api/domain/{domain}?format=csv|tsv returning raw bytes (Premium)
func fetchDomainRaw(ctx context.Context, client *http.Client, domain, apiKey, format string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/api/domain/%s?format=%s", defaultBaseURL, domain, format)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", userAgent())
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("%s", string(body))
	}
	return body, resp.StatusCode, nil
}

// postBulk — POST /api/domain/bulk with up to 50 domains (Premium)
func postBulk(ctx context.Context, client *http.Client, apiKey string, domains []string) (*apiBulkResponse, int, error) {
	body, _ := json.Marshal(map[string][]string{"domains": domains})
	url := fmt.Sprintf("%s/api/domain/bulk", defaultBaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("%s", string(respBody))
	}
	var parsed apiBulkResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, resp.StatusCode, err
	}
	return &parsed, resp.StatusCode, nil
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
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
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
		// Strip path and port if present
		if i := strings.IndexAny(d, "/:?#"); i > 0 {
			d = d[:i]
		}
		d = strings.TrimSuffix(d, ".")
		d = strings.ToLower(d)
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

func outputWriter(path string) (io.Writer, func(), error) {
	if path == "" {
		return os.Stdout, func() {}, nil
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, func() {}, err
	}
	return f, func() { _ = f.Close() }, nil
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

// handleAPIError — pretty-prints common API errors with action-oriented hints.
func handleAPIError(err error, status int, silent bool) {
	if silent {
		return
	}
	switch status {
	case http.StatusUnauthorized:
		fmt.Fprintln(os.Stderr, "[!] 401 — Invalid or missing API key.")
		fmt.Fprintln(os.Stderr, "    Get yours: "+signupURL)
	case http.StatusPaymentRequired:
		fmt.Fprintln(os.Stderr, "[!] 402 — Premium feature (CSV/TSV export or bulk endpoint).")
		fmt.Fprintln(os.Stderr, "    Upgrade: "+pricingURL)
	case http.StatusTooManyRequests:
		fmt.Fprintln(os.Stderr, "[!] 429 — Daily quota reached on free tier.")
		fmt.Fprintln(os.Stderr, "    Upgrade for unlimited: "+pricingURL)
	case http.StatusNotFound:
		fmt.Fprintln(os.Stderr, "[!] 404 — No data observed yet for this target.")
	default:
		// Truncate noisy HTML/JSON error bodies
		msg := err.Error()
		if len(msg) > 240 {
			msg = msg[:240] + "..."
		}
		fmt.Fprintf(os.Stderr, "[!] %d — %s\n", status, msg)
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
