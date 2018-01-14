package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func run(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("request: %s", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("calling endpoint: %s", err)
	}

	io.Copy(os.Stdout, resp.Body)
	return nil
}

func main() {
	url := flag.String("url", "http://localhost:48787", "server's endpoint")
	flag.Parse()
	if err := run(*url); err != nil {
		log.Fatal(err)
	}
}
