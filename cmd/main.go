package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
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

func fetchSubdomains(domain string, outputJSON bool, silent bool, withIP bool) {
	client := &http.Client{Timeout: 60 * time.Second}
	url := fmt.Sprintf("https://reconeer.com/api/domain/%s", domain)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s: %v\n", domain, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Reading response: %v\n", err)
		return
	}

	var result APIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Decoding response: %v\n", err)
		return
	}

	if outputJSON {
		data, _ := json.MarshalIndent(result.Subdomains, "", "  ")
		fmt.Println(string(data))
		return
	}

	for _, sub := range result.Subdomains {
		if silent {
			fmt.Println(sub.Subdomain)
		} else if withIP {
			fmt.Printf("%s (%s)\n", sub.Subdomain, sub.IP)
		} else {
			fmt.Printf("%s\n", sub.Subdomain)
		}
	}
}

func main() {
	domain := flag.String("d", "", "Single domain to fetch")
	domainList := flag.String("dL", "", "File with list of domains")
	output := flag.String("o", "", "Output file")
	outputJSON := flag.Bool("json", false, "Output in JSON format")
	withIP := flag.Bool("ip", false, "Include IP in output")
	silent := flag.Bool("silent", false, "Only print subdomains")
	timeout := flag.Int("timeout", 30, "Timeout per request (seconds)")

	flag.Parse()

	results := []string{}
	client := &http.Client{Timeout: time.Duration(*timeout) * time.Second}

	if *domain != "" {
		url := fmt.Sprintf("https://reconeer.com/api/domain/%s", *domain)
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var result APIResponse
		json.Unmarshal(body, &result)

		for _, sub := range result.Subdomains {
			if *outputJSON {
				s, _ := json.Marshal(sub)
				results = append(results, string(s))
			} else {
				line := sub.Subdomain
				if *withIP {
					line += " (" + sub.IP + ")"
				}
				results = append(results, line)
			}
		}
	} else if *domainList != "" {
		data, err := ioutil.ReadFile(*domainList)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		for _, line := range strings.Split(string(data), "\n") {
			if line != "" {
				fetchSubdomains(line, *outputJSON, *silent, *withIP)
			}
		}
		os.Exit(0)
	}

	if *output != "" {
		ioutil.WriteFile(*output, []byte(strings.Join(results, "\n")), 0644)
	} else {
		for _, line := range results {
			fmt.Println(line)
		}
	}
}

