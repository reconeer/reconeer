// Reconeer CLI — passive subdomain enumeration via the Reconeer API.
// Single-domain or bulk; multiple output formats; auth via env var or -k flag.
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type SubdomainData struct {
	Subdomain       string `json:"subdomain"`
	IP              string `json:"ip,omitempty"`
	Country         string `json:"country,omitempty"`
	ReverseResolves bool   `json:"reverse_resolves,omitempty"`
}

type apiDomainResponse struct {
	Subdomains []SubdomainData `json:"subdomains"`
}

type apiBulkResponse struct {
	Count   int                        `json:"count"`
	Domains map[string][]SubdomainData `json:"domains"`
}

type apiError struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
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
	signupURL      = "https://www.reconeer.com/register?utm_source=cli&utm_medium=prompt&utm_campaign=oss_funnel"
	pricingURL     = "https://www.reconeer.com/pricing?utm_source=cli&utm_medium=prompt&utm_campaign=oss_funnel"
	updateAPI      = "https://api.github.com/repos/reconeer/reconeer/releases/latest"
)

func userAgent() string {
	return fmt.Sprintf("reconeer-cli/%s (+https://www.reconeer.com)", version)
}

func main() {
	var (
		domains     stringSlice
		listFile    string
		outFile     string
		apiKey      string
		format      string
		rateLimit   int
		silent      bool
		verbose     bool
		jsonl       bool
		bulk        bool
		fields      bool
		showVersion bool
	)

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
	flag.StringVar(&format, "format", "", "Output format: txt (default) | jsonl | json | csv | tsv. csv/tsv require Premium.")
	flag.BoolVar(&jsonl, "jsonl", false, "Shortcut for -format jsonl")
	flag.BoolVar(&bulk, "bulk", false, "Use the Premium /api/domain/bulk endpoint (up to 50 domains per request)")
	flag.BoolVar(&fields, "fields", false, "Include ip,country,reverse_resolves columns in txt output")

	flag.BoolVar(&silent, "silent", false, "Print only subdomains (no banners/logs)")
	flag.BoolVar(&verbose, "v", false, "Verbose logging")
	flag.BoolVar(&verbose, "verbose", false, "Verbose logging")
	flag.BoolVar(&showVersion, "version", false, "Print version")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Reconeer CLI — passive subdomain enumeration")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  reconeer -d example.com -k YOUR_API_KEY")
		fmt.Fprintln(os.Stderr, "  export RECONEER_API_KEY=...; reconeer -d example.com")
		fmt.Fprintln(os.Stderr, "  reconeer -dL scope.txt -silent -o all_subs.txt")
		fmt.Fprintln(os.Stderr, "  reconeer -d example.com -format jsonl | jq -r '.subdomain'")
		fmt.Fprintln(os.Stderr, "  reconeer -d example.com -format csv -o subs.csv          # Premium")
		fmt.Fprintln(os.Stderr, "  reconeer -dL scope.txt -bulk -o all_subs.json             # Premium")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Subfinder users: set RECONEER_API_KEY — Reconeer is a default source since subfinder/PR#1694.")
		fmt.Fprintln(os.Stderr, "Get a free key: https://www.reconeer.com/register")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	checkLatestVersion(version, silent)

	if showVersion {
		fmt.Printf("reconeer %s (commit %s, built %s)\n", version, commit, date)
		return
	}

	// Resolve format from --jsonl shortcut
	if jsonl && format == "" {
		format = "jsonl"
	}
	if format == "" {
		format = "txt"
	}
	format = strings.ToLower(format)
	switch format {
	case "txt", "jsonl", "json", "csv", "tsv":
	default:
		fmt.Fprintf(os.Stderr, "Unknown -format %q (must be txt|jsonl|json|csv|tsv)\n", format)
		os.Exit(2)
	}

	inputDomains := make([]string, 0, 32)
	inputDomains = append(inputDomains, domains...)

	if listFile != "" {
		ds, err := readDomainsFromFileOrStdin(listFile)
		fatalIf(err, silent, "failed reading domain list: %v", err)
		inputDomains = append(inputDomains, ds...)
	} else {
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			ds, err := readDomainsFromReader(os.Stdin)
			fatalIf(err, silent, "failed reading domains from stdin: %v", err)
			inputDomains = append(inputDomains, ds...)
		}
	}

	inputDomains = normalizeDomains(inputDomains)
	if len(inputDomains) == 0 {
		flag.Usage()
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
		fmt.Fprintf(os.Stderr, "Reconeer CLI %s — passive subdomain enumeration via %s\n", version, defaultBaseURL)
		if apiKey == "" {
			fmt.Fprintf(os.Stderr, "Tip: set RECONEER_API_KEY for higher rate limits: %s\n", signupURL)
		}
	}

	// === Bulk path ===
	// One POST request handles up to 50 domains; chunk if more.
	if bulk {
		if format == "txt" {
			format = "json" // bulk default to json since txt-per-domain is ambiguous
		}
		bulkRun(ctx, client, apiKey, inputDomains, out, format, silent, verbose)
		return
	}

	// === Single-domain mode ===
	seen := make(map[string]struct{}, 4096)
	var csvHeaderWritten bool

	for _, d := range inputDomains {
		limiter.Wait()

		if verbose && !silent {
			fmt.Fprintf(os.Stderr, "[*] querying %s (format=%s)\n", d, format)
		}

		if format == "csv" || format == "tsv" {
			// Stream Premium CSV/TSV straight from API; concat output across domains
			body, status, err := fetchDomainRaw(ctx, client, d, apiKey, format)
			if err != nil {
				handleAPIError(err, status, silent)
				continue
			}
			// Strip header on subsequent domains
			lines := bytes.SplitN(body, []byte("\n"), 2)
			if len(lines) == 2 {
				if !csvHeaderWritten {
					out.Write(lines[0])
					out.Write([]byte{'\n'})
					csvHeaderWritten = true
				}
				out.Write(lines[1])
			} else {
				out.Write(body)
			}
			continue
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

			switch format {
			case "jsonl":
				_ = json.NewEncoder(out).Encode(s)
			case "json":
				// Accumulate then flush at end — for now treat like jsonl
				_ = json.NewEncoder(out).Encode(s)
			case "txt":
				if fields {
					fmt.Fprintf(out, "%s\t%s\t%s\t%v\n", s.Subdomain, s.IP, s.Country, s.ReverseResolves)
				} else {
					fmt.Fprintln(out, s.Subdomain)
				}
			}
		}
	}

	if !silent && len(seen) == 0 && (format == "txt" || format == "jsonl" || format == "json") {
		fmt.Fprintf(os.Stderr, "No results. Docs: https://www.reconeer.com/docs\n")
	}
}

func bulkRun(ctx context.Context, client *http.Client, apiKey string, domains []string, out io.Writer, format string, silent, verbose bool) {
	const batchSize = 50
	all := apiBulkResponse{Domains: map[string][]SubdomainData{}}
	for i := 0; i < len(domains); i += batchSize {
		end := i + batchSize
		if end > len(domains) {
			end = len(domains)
		}
		batch := domains[i:end]
		if verbose && !silent {
			fmt.Fprintf(os.Stderr, "[*] bulk request %d–%d (%d domains)\n", i+1, end, len(batch))
		}
		resp, status, err := postBulk(ctx, client, apiKey, batch)
		if err != nil {
			handleAPIError(err, status, silent)
			continue
		}
		all.Count += resp.Count
		for k, v := range resp.Domains {
			all.Domains[k] = v
		}
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	_ = enc.Encode(all)
}

func checkLatestVersion(current string, silent bool) {
	if silent || current == "dev" {
		return
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", updateAPI, nil)
	req.Header.Set("User-Agent", userAgent())
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var r struct {
		Tag string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return
	}
	latest := strings.TrimPrefix(r.Tag, "v")
	if latest != "" && latest != current {
		fmt.Fprintf(os.Stderr, "\n[!] New version available: %s (current: %s)\n", latest, current)
		fmt.Fprintf(os.Stderr, "    Update: https://github.com/reconeer/reconeer/releases\n\n")
	}
}

// keep linkage so bufio is used if main.go is the only file the linker sees
var _ = bufio.ScanLines
