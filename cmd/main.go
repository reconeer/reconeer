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
    "time"
)

type ApiResponse struct {
    Domain         string          `json:"domain"`
    SubdomainCount int             `json:"subdomainCount"`
    MostUsedIP     string          `json:"mostUsedIP"`
    MostUsedCount  string          `json:"mostUsedCount"`
    Subdomains     []SubdomainData `json:"subdomains"`
}

type SubdomainData struct {
    Subdomain       string  `json:"subdomain"`
    IP              string  `json:"ip"`
    Domain          string  `json:"domain"`
    Country         *string `json:"country"`
    ReverseResolves *bool   `json:"reverse_resolves"`
}

func fetchSubdomains(domain string) {
    fmt.Println("Fetching for:", domain)
    url := fmt.Sprintf("https://www.reconeer.com/api/domain/%s", domain)

    client := &http.Client{Timeout: 60 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        log.Fatalf("Failed to fetch data: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read response: %v", err)
    }

    var apiResp ApiResponse
    if err := json.Unmarshal(body, &apiResp); err != nil {
        log.Fatalf("Error: json: %v", err)
    }

    for _, sub := range apiResp.Subdomains {
        fmt.Printf("%s -> %s\n", sub.Subdomain, sub.IP)
    }
}

func main() {
    domain := flag.String("d", "", "Domain to fetch subdomains for")
    domainList := flag.String("dL", "", "File with list of domains")
    flag.Parse()

    if *domain != "" {
        fetchSubdomains(*domain)
    } else if *domainList != "" {
        file, err := os.ReadFile(*domainList)
        if err != nil {
            log.Fatalf("Failed to read file: %v", err)
        }

        lines := strings.Split(string(file), "\n")
        for _, line := range lines {
            line = strings.TrimSpace(line)
            if line != "" {
                fetchSubdomains(line)
            }
        }
    } else {
        fmt.Println("Usage of reconeer:")
        fmt.Println("  -d string")
        fmt.Println("        Domain to fetch subdomains for")
        fmt.Println("  -dL string")
        fmt.Println("        File with list of domains")
    }
}
