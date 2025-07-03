package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type SubdomainRow struct {
	Subdomain       string  `json:"subdomain"`
	IP              string  `json:"ip"`
	Domain          string  `json:"domain"`
	Country         *string `json:"country"`
	ReverseResolves *bool   `json:"reverse_resolves"`
}

type DomainAPIResponse struct {
	Domain         string          `json:"domain"`
	SubdomainCount int             `json:"subdomainCount"`
	MostUsedIP     string          `json:"mostUsedIP"`
	MostUsedCount  interface{}     `json:"mostUsedCount"` // some are strings
	Subdomains     []SubdomainRow  `json:"subdomains"`
}

func fetchDomainData(domain string) {
	url := fmt.Sprintf("https://www.reconeer.com/api/domain/%s", domain)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[error] %s: %v\n", domain, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("[error] %s: HTTP %d\n", domain, resp.StatusCode)
		return
	}

	var apiResp DomainAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Printf("[error] decoding JSON for %s: %v\n", domain, err)
		return
	}

	fmt.Printf("== %s (%d subdomains) ==\n", apiResp.Domain, apiResp.SubdomainCount)
	for _, entry := range apiResp.Subdomains {
		fmt.Printf("%s,%s\n", entry.Subdomain, entry.IP)
	}
	fmt.Println()
}

func fetchDomainsFromFile(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			fetchDomainData(domain)
		}
	}
}

func main() {
	domain := flag.String("d", "", "Single domain to fetch")
	domainList := flag.String("dL", "", "File with list of domains")
	flag.Parse()

	if *domain != "" {
		fetchDomainData(*domain)
	} else if *domainList != "" {
		fetchDomainsFromFile(*domainList)
	} else {
		fmt.Println("Usage:")
		fmt.Println("  reconeer -d <domain>")
		fmt.Println("  reconeer -dL <domain_list_file>")
		os.Exit(1)
	}
}

