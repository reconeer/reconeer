
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
)

type SubdomainData struct {
    Subdomain       string `json:"subdomain"`
    IP              string `json:"ip"`
    Domain          string `json:"domain"`
    Country         string `json:"country"`
    ReverseResolves *bool  `json:"reverse_resolves"`
}

type APIResponse struct {
    Domain         string          `json:"domain"`
    SubdomainCount int             `json:"subdomainCount"`
    MostUsedIP     string          `json:"mostUsedIP"`
    MostUsedCount  string          `json:"mostUsedCount"`
    Subdomains     []SubdomainData `json:"subdomains"`
}

func fetchSubdomains(domain string) {
    fmt.Println("Fetching for:", domain)
    url := fmt.Sprintf("https://reconeer.com/api/domain/%s", domain)

    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error fetching data:", err)
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response:", err)
        return
    }

    var result APIResponse
    err = json.Unmarshal(body, &result)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    for _, sub := range result.Subdomains {
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
        file, err := os.Open(*domainList)
        if err != nil {
            fmt.Println("Error reading domain list:", err)
            return
        }
        defer file.Close()

        var domains []string
        buf := make([]byte, 1024)
        for {
            n, err := file.Read(buf)
            if n > 0 {
                content := string(buf[:n])
                for _, line := range splitLines(content) {
                    domains = append(domains, line)
                }
            }
            if err != nil {
                break
            }
        }

        for _, d := range domains {
            fetchSubdomains(d)
        }
    } else {
        flag.Usage()
    }
}

func splitLines(s string) []string {
    var lines []string
    current := ""
    for _, r := range s {
        if r == '\n' || r == '\r' {
            if current != "" {
                lines = append(lines, current)
                current = ""
            }
        } else {
            current += string(r)
        }
    }
    if current != "" {
        lines = append(lines, current)
    }
    return lines
}
