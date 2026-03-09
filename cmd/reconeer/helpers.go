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
	"strings"
	"time"
)

func fetchDomain(ctx context.Context, client *http.Client, domain, apiKey string) ([]SubdomainData, int, error) {

	url := fmt.Sprintf("%s/api/domain/%s", defaultBaseURL, domain)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

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
	fmt.Fprintln(os.Stderr, "Usage: reconeer -d example.com | reconeer -dL domains.txt")
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

	if status == http.StatusTooManyRequests {
		fmt.Println("Rate limit reached.")
		fmt.Println("Upgrade:", pricingURL)
		return
	}

	fmt.Println(err)
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
