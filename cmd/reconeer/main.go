package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SubdomainData struct {
	Subdomain string `json:"subdomain"`
	IP        string `json:"ip,omitempty"`
}

type apiDomainResponse struct {
	Subdomains []SubdomainData `json:"subdomains"`
}

type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ",") }
func (s *stringSlice) Set(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	*s = append(*s, v)
	return nil
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const (
	defaultBaseURL = "https://www.reconeer.com"
	signupURL      = "https://www.reconeer.com/signup?utm_source=cli&utm_medium=prompt&utm_campaign=oss_funnel"
	pricingURL     = "https://www.reconeer.com/pricing?utm_source=cli&utm_medium=prompt&utm_campaign=oss_funnel"
)

func main() {
	var domains stringSlice
	var listFile string
	var outFile string
	var apiKey string
	var rateLimit int
	var silent bool
	var verbose bool
	var jsonl bool
	var showVersion bool

	// Support both short and long flags
	flag.Var(&domains, "d", "Domain to enumerate (repeatable)")
	flag.Var(&domains, "domain", "Domain to enumerate (repeatable)")
	flag.StringVar(&listFile, "dL", "", "File with one domain per line; supports '-' for STDIN")
	flag.StringVar(&listFile, "list", "", "File with one domain per line; supports '-' for STDIN")

	flag.StringVar(&apiKey, "k", "", "API key (or set RECONEER_API_KEY)")
	flag.StringVar(&apiKey, "api-key", "", "API key (or set RECONEER_API_KEY)")

	flag.IntVar(&rateLimit, "rl", 3, "Max requests per second (client-side). Default: 3")
	flag.IntVar(&rateLimit, "rate-limit", 3, "Max requests per second (client-side). Default: 3")

	flag.StringVar(&outFile, "o", "", "Write output to file (default: STDOUT)")
	flag.StringVar(&outFile, "output", "", "Write output to file (default: STDOUT)")
	flag.BoolVar(&jsonl, "jsonl", false, "Output JSON Lines (one object per subdomain)")

	flag.BoolVar(&silent, "silent", false, "Print only subdomains (no banners/logs)")
	flag.BoolVar(&verbose, "v", false, "Verbose logging")
	flag.BoolVar(&verbose, "verbose", false, "Verbose logging")
	flag.BoolVar(&showVersion, "version", false, "Print version")
	flag.Parse()

	if showVersion {
		fmt.Printf("reconeer %s (commit %s, built %s)\n", version, commit, date)
		return
	}

	// Domain sources: flags + list + stdin (if piped and no listFile)
	inputDomains := make([]string, 0, 32)
	inputDomains = append(inputDomains, domains...)

	if listFile != "" {
		ds, err := readDomainsFromFileOrStdin(listFile)
		fatalIf(err, silent, "failed reading domain list: %v", err)
		inputDomains = append(inputDomains, ds...)
	} else {
		// If stdin is piped, read domains from stdin too
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			ds, err := readDomainsFromReader(os.Stdin)
			fatalIf(err, silent, "failed reading domains from stdin: %v", err)
			inputDomains = append(inputDomains, ds...)
		}
	}

	inputDomains = normalizeDomains(inputDomains)
	if len(inputDomains) == 0 {
		printHelpHint()
		os.Exit(2)
	}

	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("RECONEER_API_KEY"))
	}

	out, closeOut, err := outputWriter(outFile)
	fatalIf(err, silent, "failed opening output: %v", err)
	defer closeOut()

	limiter := newRateLimiter(rateLimit)
	client := &http.Client{Timeout: 30 * time.Second}
	ctx := context.Background()

	if !silent {
		fmt.Fprintf(os.Stderr, "Reconeer CLI – passive subdomain enumeration via %s\n", defaultBaseURL)
		if apiKey == "" {
			fmt.Fprintf(os.Stderr, "Tip: set RECONEER_API_KEY for rate-limit isolation: %s\n", signupURL)
		}
	}

	seen := make(map[string]struct{}, 4096)

	for _, d := range inputDomains {
		limiter.Wait()
		if verbose && !silent {
			fmt.Fprintf(os.Stderr, "[*] querying %s\n", d)
		}

		subs, status, err := fetchDomain(ctx, client, d, apiKey)
		if err != nil {
			handleAPIError(err, status, silent)
			continue
		}

		for _, s := range subs {
			if s.Subdomain == "" {
				continue
			}
			if _, ok := seen[s.Subdomain]; ok {
				continue
			}
			seen[s.Subdomain] = struct{}{}

			if jsonl {
				enc := json.NewEncoder(out)
				_ = enc.Encode(s)
			} else {
				fmt.Fprintln(out, s.Subdomain)
			}
		}
	}

	if !silent && len(seen) == 0 {
		fmt.Fprintf(os.Stderr, "No results. If this is unexpected, check docs: https://www.reconeer.com/docs.html\n")
	}
}

func printHelpHint() {
	fmt.Fprintln(os.Stderr, "No domains provided.\nUsage: reconeer -d example.com | reconeer -dL domains.txt | cat domains.txt | reconeer\nRun: reconeer -h")
}

func fatalIf(err error, silent bool, format string, args ...any) {
	if err == nil {
		return
	}
	if !silent {
		fmt.Fprintf(os.Stderr, "[!] "+format+"\n", args...)
	}
	os.Exit(1)
}

func outputWriter(path string) (io.Writer, func(), error) {
	if strings.TrimSpace(path) == "" {
		return os.Stdout, func() {}, nil
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, func() {}, err
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, func() {}, err
	}
	return f, func() { _ = f.Close() }, nil
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
	sc.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

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
	out := make([]string, 0, len(in))
	seen := make(map[string]struct{}, len(in))
	for _, d := range in {
		d = strings.TrimSpace(d)
		d = strings.TrimPrefix(d, "http://")
		d = strings.TrimPrefix(d, "https://")
		d = strings.TrimSuffix(d, "/")
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

func fetchDomain(ctx context.Context, client *http.Client, domain, apiKey string) ([]SubdomainData, int, error) {
	url := fmt.Sprintf("%s/api/domain/%s", defaultBaseURL, domain)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return nil, resp.StatusCode, fmt.Errorf("api error (%d): %s", resp.StatusCode, msg)
	}

	var parsed apiDomainResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("parse failed: %w", err)
	}
	return parsed.Subdomains, resp.StatusCode, nil
}

func handleAPIError(err error, status int, silent bool) {
	if silent {
		return
	}
	switch status {
	case http.StatusTooManyRequests:
		fmt.Fprintf(os.Stderr, "[!] Rate limit reached (HTTP 429).\n")
		fmt.Fprintf(os.Stderr, "    Free tier is limited (10 queries/day). Upgrade for unlimited queries:\n")
		fmt.Fprintf(os.Stderr, "    %s\n", pricingURL)
		return
	case http.StatusUnauthorized, http.StatusForbidden:
		fmt.Fprintf(os.Stderr, "[!] Access denied (%d). If you're using an API key, verify it.\n", status)
		fmt.Fprintf(os.Stderr, "    Get a free key: %s\n", signupURL)
		return
	}
	fmt.Fprintf(os.Stderr, "[!] %v\n", err)
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

func (r *rateLimiter) Wait() { <-r.ch }
