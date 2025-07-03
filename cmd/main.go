package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type SubdomainData struct {
	Subdomain       string  `json:"subdomain"`
	IP              string  `json:"ip"`
	Domain          string  `json:"domain"`
	Country         *string `json:"country"`
	ReverseResolves *bool   `json:"reverse_resolves"`
}

type DomainResponse struct {
	Domain         string          `json:"domain"`
	SubdomainCount int             `json:"subdomainCount"`
	MostUsedIP     string          `json:"mostUsedIP"`
	MostUsedCount  string          `json:"mostUsedCount"`
	Subdomains     []SubdomainData `json:"subdomains"`
}

func fetchSubdomains(domain string) {
	fmt.Printf("Fetching for: %s\n", domain)
	url := fmt.Sprintf("https://www.reconeer.com/api/domain/%s", domain)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch: %v", err)
	}
	defer resp.Body.Close()

	var result DomainResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for _, sub := range result.Subdomains {
		fmt.Println(sub.Subdomain)
	}
}

func fetchFromList(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}
	domains := strings.Split(string(content), "\n")
	for _, domain := range domains {
		domain = strings.TrimSpace(domain)
		if domain != "" {
			fetchSubdomains(domain)
		}
	}
}

func main() {
	domain := flag.String("d", "", "Domain to fetch subdomains for")
	domainList := flag.String("dL", "", "File with list of domains")
	flag.Parse()

	if *domain != "" {
		fetchSubdomains(*domain)
	} else if *domainList != "" {
		fetchFromList(*domainList)
	} else {
		fmt.Println("Usage of reconeer:")
		fmt.Println("  -d string")
		fmt.Println("    \tDomain to fetch subdomains for")
		fmt.Println("  -dL string")
		fmt.Println("    \tFile with list of domains")
	}
}

