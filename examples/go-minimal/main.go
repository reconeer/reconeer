package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	domain := "example.com"
	req, _ := http.NewRequest("GET", "https://www.reconeer.com/api/domain/"+domain, nil)

	if k := os.Getenv("RECONEER_API_KEY"); k != "" {
		req.Header.Set("X-API-Key", k)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("status:", resp.Status)
	// parse JSON per docs: https://www.reconeer.com/docs.html
}
