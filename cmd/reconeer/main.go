package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type SubdomainData struct {
	Subdomain string `json:"subdomain"`
	IP        string `json:"ip"`
}

func fetchSubdomains(domain string) ([]SubdomainData, error) {
	url := fmt.Sprintf("https://www.reconeer.com/api/domain/%s", domain)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var result struct {
		Subdomains []SubdomainData `json:"subdomains"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return result.Subdomains, nil
}

func main() {
	domain := flag.String("d", "", "Domain to fetch subdomains for")
	domainList := flag.String("dL", "", "File with list of domains")
	jsonOutput := flag.Bool("json", false, "Output as JSON")
	flag.Parse()

	if *domain == "" && *domainList == "" {
		fmt.Println("Usage of reconeer:")
		flag.PrintDefaults()
		return
	}

	domains := []string{}
	if *domain != "" {
		domains = append(domains, *domain)
	}
	if *domainList != "" {
		file, err := os.Open(*domainList)
		if err != nil {
			fmt.Printf("Failed to open domain list: %v\n", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			d := strings.TrimSpace(scanner.Text())
			if d != "" {
				domains = append(domains, d)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading domain list: %v\n", err)
			return
		}
	}

	for _, d := range domains {
		//fmt.Printf("Fetching for: %s\n", d)
		subdomains, err := fetchSubdomains(d)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if *jsonOutput {
			output, err := json.MarshalIndent(subdomains, "", "  ")
			if err != nil {
				fmt.Printf("Failed to marshal JSON: %v\n", err)
				continue
			}
			fmt.Println(string(output))
		} else {
			for _, sd := range subdomains {
				fmt.Println(sd.Subdomain)
			}
		}
	}
}
