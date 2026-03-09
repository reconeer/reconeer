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
	"sync"
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
	signupURL      = "https://www.reconeer.com/register"
	pricingURL     = "https://www.reconeer.com/pricing"
	updateAPI      = "https://api.github.com/repos/reconeer/reconeer/releases/latest"
)

func main() {

	var domains stringSlice
	var listFile string
	var outFile string
	var apiKey string
	var rateLimit int
	var threads int
	var silent bool
	var verbose bool
	var jsonl bool
	var showVersion bool
	var update bool

	flag.Var(&domains, "d", "Domain to enumerate (repeatable)")
	flag.StringVar(&listFile, "dL", "", "File with one domain per line")

	flag.StringVar(&apiKey, "k", "", "API key")
	flag.StringVar(&apiKey, "api-key", "", "API key")

	flag.IntVar(&rateLimit, "rl", 3, "Requests per second")
	flag.IntVar(&threads, "t", 5, "Number of concurrent workers")

	flag.StringVar(&outFile, "o", "", "Write output to file")
	flag.BoolVar(&jsonl, "jsonl", false, "JSONL output")

	flag.BoolVar(&silent, "silent", false, "Silent mode")
	flag.BoolVar(&verbose, "v", false, "Verbose logging")
	flag.BoolVar(&showVersion, "version", false, "Print version")
	flag.BoolVar(&update, "update", false, "Update Reconeer")

	flag.Parse()

	if update {
		runSelfUpdate()
		return
	}

	go checkLatestVersion(version, silent)

	if showVersion {
		fmt.Printf("reconeer %s (commit %s, built %s)\n", version, commit, date)
		return
	}

	inputDomains := make([]string, 0)
	inputDomains = append(inputDomains, domains...)

	if listFile != "" {
		ds, err := readDomainsFromFile(listFile)
		fatalIf(err, silent, "failed reading domain list: %v", err)
		inputDomains = append(inputDomains, ds...)
	}

	inputDomains = normalizeDomains(inputDomains)

	if len(inputDomains) == 0 {
		fmt.Println("Usage: reconeer -d example.com")
		os.Exit(1)
	}

	if apiKey == "" {
		apiKey = os.Getenv("RECONEER_API_KEY")
	}

	out, closeOut, err := outputWriter(outFile)
	fatalIf(err, silent, "failed opening output: %v", err)
	defer closeOut()

	limiter := newRateLimiter(rateLimit)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	ctx := context.Background()

	if !silent {
		fmt.Println("Reconeer CLI – passive subdomain enumeration via reconeer.com")
	}

	seen := make(map[string]struct{})
	var mu sync.Mutex

	jobs := make(chan string)
	wg := sync.WaitGroup{}

	for i := 0; i < threads; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()

			for domain := range jobs {

				limiter.Wait()

				if verbose {
					fmt.Println("querying", domain)
				}

				subs, status, err := fetchDomain(ctx, client, domain, apiKey)

				if err != nil {
					handleAPIError(err, status, silent)
					continue
				}

				for _, s := range subs {

					if s.Subdomain == "" {
						continue
					}

					mu.Lock()

					if _, ok := seen[s.Subdomain]; ok {
						mu.Unlock()
						continue
					}

					seen[s.Subdomain] = struct{}{}
					mu.Unlock()

					if jsonl {
						json.NewEncoder(out).Encode(s)
					} else {
						fmt.Fprintln(out, s.Subdomain)
					}
				}
			}
		}()
	}

	for _, d := range inputDomains {
		jobs <- d
	}

	close(jobs)
	wg.Wait()
}

func fetchDomain(ctx context.Context, client *http.Client, domain, apiKey string) ([]SubdomainData, int, error) {

	url := fmt.Sprintf("%s/api/domain/%s", defaultBaseURL, domain)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	var resp *http.Response

	for i := 0; i < 3; i++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf(string(body))
	}

	var parsed apiDomainResponse

	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, resp.StatusCode, err
	}

	return parsed.Subdomains, resp.StatusCode, nil
}

func runSelfUpdate() {

	fmt.Println("Checking latest version...")

	resp, err := http.Get(updateAPI)
	if err != nil {
		fmt.Println("Failed:", err)
		return
	}

	defer resp.Body.Close()

	var r struct {
		Tag string `json:"tag_name"`
	}

	json.NewDecoder(resp.Body).Decode(&r)

	latest := strings.TrimPrefix(r.Tag, "v")

	if latest == version {
		fmt.Println("Already up to date.")
		return
	}

	fmt.Println("New version:", latest)
	fmt.Println("Download:", "https://github.com/reconeer/reconeer/releases/latest")
}

func checkLatestVersion(current string, silent bool) {

	if silent || current == "dev" {
		return
	}

	resp, err := http.Get(updateAPI)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var r struct {
		Tag string `json:"tag_name"`
	}

	json.NewDecoder(resp.Body).Decode(&r)

	latest := strings.TrimPrefix(r.Tag, "v")

	if latest != "" && latest != current {

		fmt.Fprintf(os.Stderr,
			"\n[!] New version available: %s (current: %s)\n",
			latest,
			current,
		)

		fmt.Fprintf(os.Stderr,
			"Update: https://github.com/reconeer/reconeer/releases\n\n",
		)
	}
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

func readDomainsFromFile(path string) ([]string, error) {

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

func outputWriter(path string) (io.Writer, func(), error) {

	if path == "" {
		return os.Stdout, func() {}, nil
	}

	dir := filepath.Dir(path)

	if dir != "." {
		os.MkdirAll(dir, 0755)
	}

	f, err := os.Create(path)

	if err != nil {
		return nil, func() {}, err
	}

	return f, func() { f.Close() }, nil
}

func fatalIf(err error, silent bool, msg string, args ...any) {

	if err == nil {
		return
	}

	if !silent {
		fmt.Fprintf(os.Stderr, msg+"\n", args...)
	}

	os.Exit(1)
}

func handleAPIError(err error, status int, silent bool) {

	if silent {
		return
	}

	if status == 429 {
		fmt.Println("Rate limit reached. Upgrade:", pricingURL)
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
