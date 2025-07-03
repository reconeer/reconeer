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

type SubdomainData struct {
	Subdomain       string  `json:"subdomain"`
	IP              string  `json:"ip"`
	Domain          string  `json:"domain"`
	Country         *string `json:"country"`
	ReverseResolves *bool   `json:"reverse_resolves"`
}

type APIResponse struct {
	Domain         string          `json:"domain"`
	SubdomainCount int             `json:"subdomainCount"`
	MostUsedIP     string          `json:"mostUsedIP"`
	MostUsedCount  string          `json:"mostUsedCount"`
	Subdomains     []SubdomainData `json:"subdomains"`
}

func fetchDomain(domain string) {
	fmt.Printf("Fetching for: %s\n", domain)
	url := fmt.Sprintf("https://reconeer.com/api/domain/%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Error: HTTP %d returned for %s\n", resp.StatusCode, domain)
		return
	}

	var result APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("JSON decode error: %v\n", err)
		return
	}

	fmt.Printf("Found %d subdomains for %s\n", result.SubdomainCount, result.Domain)
	for _, sub := range result.Subdomains {
		fmt.Printf("%s -> %s\n", sub.Subdomain, sub.IP)
	}
	fmt.Println()
}

func main() {
	domain := flag.String("d", "", "Domain to fetch subdomains for")
	domainList := flag.String("dL", "", "File containing list of domains")

	flag.Parse()

	if *domain != "" {
		fetchDomain(*domain)
	} else if *domainList != "" {
		file, err := os.Open(*domainList)
		if err != nil {
			log.Fatalf("Could not open file: %v\n", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				fetchDomain(line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading file: %v\n", err)
		}
	} else {
		fmt.Println("Usage: reconeer -d <domain> OR reconeer -dL <domain_file>")
	}
}

