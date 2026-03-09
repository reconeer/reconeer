package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
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
	signupURL      = "https://www.reconeer.com/register?utm_source=cli&utm_medium=prompt&utm_campaign=oss_funnel"
	pricingURL     = "https://www.reconeer.com/pricing?utm_source=cli&utm_medium=prompt&utm_campaign=oss_funnel"

	updateAPI = "https://api.github.com/repos/reconeer/reconeer/releases/latest"
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

	checkLatestVersion(version, silent)

	if showVersion {
		fmt.Printf("reconeer %s (commit %s, built %s)\n", version, commit, date)
		return
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
		fmt.Fprintf(os.Stderr, "No results. Docs: https://www.reconeer.com/docs.html\n")
	}
}

func checkLatestVersion(current string, silent bool) {
	if silent || current == "dev" {
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(updateAPI)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	type release struct {
		Tag string `json:"tag_name"`
	}

	var r release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return
	}

	latest := strings.TrimPrefix(r.Tag, "v")

	if latest != "" && latest != current {
		fmt.Fprintf(os.Stderr, "\n[!] New version available: %s (current: %s)\n", latest, current)
		fmt.Fprintf(os.Stderr, "    Update: https://github.com/reconeer/reconeer/releases\n\n")
	}
}
