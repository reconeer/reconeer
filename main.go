package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "io"
    "encoding/json"
)

type SubdomainData struct {
    Subdomain string `json:"subdomain"`
    IP        string `json:"ip"`
}

func fetchSubdomains(domain string) ([]SubdomainData, error) {
    url := fmt.Sprintf("https://www.reconeer.com/api/domain/%s", domain)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var results []SubdomainData
    err = json.Unmarshal(body, &results)
    return results, err
}

func processDomains(domains []string) {
    for _, domain := range domains {
        fmt.Printf("Fetching for: %s\n", domain)
        results, err := fetchSubdomains(domain)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            continue
        }

        for _, entry := range results {
            fmt.Printf("%s,%s\n", entry.Subdomain, entry.IP)
        }
    }
}

func main() {
    singleDomain := flag.String("d", "", "Domain to fetch subdomains for")
    domainList := flag.String("dL", "", "File with list of domains")
    flag.Parse()

    if *singleDomain == "" && *domainList == "" {
        flag.Usage()
        os.Exit(1)
    }

    var domains []string
    if *singleDomain != "" {
        domains = append(domains, *singleDomain)
    }

    if *domainList != "" {
        file, err := os.Open(*domainList)
        if err != nil {
            log.Fatalf("Failed to open domain list: %v", err)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            domains = append(domains, strings.TrimSpace(scanner.Text()))
        }

        if err := scanner.Err(); err != nil {
            log.Fatalf("Scanner error: %v", err)
        }
    }

    processDomains(domains)
}