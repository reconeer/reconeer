package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const apiBase = "https://www.reconeer.com/api"

type SubdomainEntry struct {
	Subdomain string `json:"subdomain"`
	IP        string `json:"ip"`
}

func fetchSubdomains(domain string) ([]SubdomainEntry, error) {
	url := fmt.Sprintf("%s/domain/%s", apiBase, domain)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch: %s (%d)", domain, resp.StatusCode)
	}

	var result struct {
		Subdomains []SubdomainEntry `json:"subdomains"`
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Subdomains, nil
}

func handleDomain(domain string) {
	subs, err := fetchSubdomains(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s: %v\n", domain, err)
		return
	}
	for _, entry := range subs {
		fmt.Printf("%s,%s\n", entry.Subdomain, entry.IP)
	}
}

func handleDomainList(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[error] cannot open domain list file:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			handleDomain(domain)
		}
	}
}

func main() {
	domain := flag.String("d", "", "Fetch subdomains for a single domain")
	domainList := flag.String("dL", "", "Fetch subdomains for list of domains in file")
	flag.Parse()

	if *domain != "" {
		handleDomain(*domain)
	} else if *domainList != "" {
		handleDomainList(*domainList)
	} else {
		fmt.Println("Usage: reconeer -d <domain> OR reconeer -dL <domain_list>")
		os.Exit(1)
	}
}
