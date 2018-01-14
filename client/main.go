package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func run(url string, N int) error {
	errChan := make(chan error, 4)
	resChan := make(chan int, N)

	coordChan := make(chan struct{}, N)

	for i := 0; i < N; i++ {
		// this is a goroutine, it simulates a long-running process
		go func(reqNo int) {
			resp, err := http.Post(url, "", nil)
			if resp != nil {
				defer resp.Body.Close()
			}
			if err != nil {
				errChan <- fmt.Errorf("[%d] doing POST: %s", reqNo, err)
				return
			}

			// this is just a sign to help coordinate DELETE request,
			// which I want to happen after N POST requests are sent
			coordChan <- struct{}{}

			reply, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errChan <- fmt.Errorf("[%d] read: %s", reqNo, err)
			}
			connNo, err := strconv.ParseInt(string(reply), 10, 32)
			if err != nil {
				errChan <- fmt.Errorf("[%d] parsing reply: %s", reqNo, err)
			}
			resChan <- int(connNo)
		}(i)
	}

	// wait for POST requests to be sent
	for i := 0; i < N; i++ {
		select {
		case _ = <-coordChan:
			if i > 0 && i%100 == 99 {
				log.Printf("sent %d POST requests\n", i+1)
			}
		case err := <-errChan:
			return err // some goroutine has an error
		}
	}
	log.Printf("sent %d POST requests\n", N)

	delReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("create DELETE request: %s", err)
	}

	delResp, err := http.DefaultClient.Do(delReq)
	if delResp != nil {
		defer delResp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("send DELETE request: %s", err)
	}

	// DELETE is sent, wait for results
	flags := make([]int, N)

	for i := 0; i < N; i++ {
		select {
		case data := <-resChan:
			flags[data] = flags[data] + 1
			if i > 0 && i%100 == 99 {
				log.Printf("received replies from %d POST requests\n", i+1)
			}
		case err := <-errChan:
			return err
		}
	}
	log.Printf("received replies from all %d POST requests\n", N)

	// let's verify the replies: each flag must be 1
	for i := 0; i < N; i++ {
		if flags[i] != 1 {
			log.Printf("expected flag[%d] to be 1, got %d\n", i, flags[i])
		}
	}

	log.Println("done")
	return nil
}

func main() {
	url := flag.String("url", "http://localhost:48787", "server's endpoint")
	N := flag.Int("N", 5, "number of requests")
	flag.Parse()
	if err := run(*url, int(*N)); err != nil {
		log.Fatal(err)
	}
}
