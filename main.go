package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Structs matching the API response
type SubdomainEntry struct {
	Subdomain       string `json:"subdomain"`
	IP              string `json:"ip"`
	Domain          string `json:"domain"`
	ReverseResolves *bool  `json:"reverse_resolves"`
}

type APIResponse struct {
	Domain         string           `json:"domain"`
	SubdomainCount int              `json:"subdomainCount"`
	MostUsedIP     string           `json:"mostUsedIP"`
	MostUsedCount  string           `json:"mostUsedCount"`
	Subdomains     []SubdomainEntry `json:"subdomains"`
}

func fetchSubdomains(domain string) {
	url := fmt.Sprintf("https://www.reconeer.com/api/domain/%s", domain)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching domain data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Error: Received status code %d from server", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	var result APIResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	for _, sub := range result.Subdomains {
		fmt.Println(sub.Subdomain)
	}
}

func main() {
	domain := flag.String("d", "", "Fetch subdomains for a single domain")
	domainList := flag.String("dL", "", "File with list of domains to fetch subdomains for")
	flag.Parse()

	if *domain == "" && *domainList == "" {
		fmt.Println("Usage:")
		fmt.Println("  reconeer -d example.com")
		fmt.Println("  reconeer -dL domains.txt")
		os.Exit(1)
	}

	if *domain != "" {
		fmt.Println("Fetching for:", *domain)
		fetchSubdomains(*domain)
	}

	if *domainList != "" {
		data, err := os.ReadFile(*domainList)
		if err != nil {
			log.Fatalf("Error reading domain list: %v", err)
		}
		domains := strings.Split(string(data), "\n")
		for _, d := range domains {
			d = strings.TrimSpace(d)
			if d == "" {
				continue
			}
			fmt.Printf("\n--- %s ---\n", d)
			fetchSubdomains(d)
		}
	}
}

