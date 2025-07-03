
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
)

type SubdomainData struct {
    Subdomain       string `json:"subdomain"`
    IP              string `json:"ip"`
    Domain          string `json:"domain"`
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
    fmt.Printf("Fetching for: %s\n", domain)
    url := fmt.Sprintf("https://reconeer.com/api/domain/%s", domain)
    resp, err := http.Get(url)
    if err != nil {
        log.Fatalf("Error fetching data: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Error reading response body: %v", err)
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
    domain := flag.String("d", "", "Target domain")
    domainList := flag.String("dL", "", "File with list of domains")
    flag.Parse()

    if *domain != "" {
        fetchSubdomains(*domain)
    } else if *domainList != "" {
        content, err := ioutil.ReadFile(*domainList)
        if err != nil {
            log.Fatalf("Could not read domain list: %v", err)
        }
        lines := strings.Split(string(content), "\n")
        for _, d := range lines {
            if d != "" {
                fetchSubdomains(d)
            }
        }
    } else {
        flag.Usage()
        os.Exit(1)
    }
}
