
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "net/http"
    "os"
)

type SubdomainData struct {
    Subdomain       string  `json:"subdomain"`
    IP              string  `json:"ip"`
    Domain          string  `json:"domain"`
    Country         *string `json:"country"`
    ReverseResolves *bool   `json:"reverse_resolves"`
}

type APIResponse struct {
    Domain          string          `json:"domain"`
    SubdomainCount  int             `json:"subdomainCount"`
    MostUsedIP      string          `json:"mostUsedIP"`
    MostUsedCount   string          `json:"mostUsedCount"`
    Subdomains      []SubdomainData `json:"subdomains"`
}

func fetchSubdomains(domain string) ([]SubdomainData, error) {
    url := fmt.Sprintf("https://www.reconeer.com/api/domain/%s", domain)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result APIResponse
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.Subdomains, nil
}

func main() {
    domain := flag.String("d", "", "Domain to fetch subdomains for")
    domainList := flag.String("dL", "", "File with list of domains")
    flag.Parse()

    if *domain == "" && *domainList == "" {
        fmt.Println("Usage of reconeer:")
        fmt.Println("  -d string
	Domain to fetch subdomains for")
        fmt.Println("  -dL string
	File with list of domains")
        return
    }

    var domains []string

    if *domain != "" {
        domains = append(domains, *domain)
    } else if *domainList != "" {
        data, err := os.ReadFile(*domainList)
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }
        for _, line := range string(data) {
            if line != '\n' && line != '\r' {
                domains = append(domains, string(line))
            }
        }
    }

    for _, d := range domains {
        fmt.Printf("Fetching for: %s\n", d)
        subdomains, err := fetchSubdomains(d)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        for _, s := range subdomains {
            fmt.Println(s.Subdomain)
        }
    }
}
