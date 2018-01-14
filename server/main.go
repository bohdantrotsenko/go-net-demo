package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type srv struct {
}

func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fl := w.(http.Flusher)
	for i := 0; i < 5; i++ {
		w.Write(([]byte)(fmt.Sprintf("%d\n", i)))
		fl.Flush()
		fmt.Printf("sending %d\n", i)
		time.Sleep(time.Second)
	}
}

func run(addr string) error {
	return http.ListenAndServe(addr, &srv{})
}

func main() {
	addr := flag.String("addr", ":48787", "listening address")
	flag.Parse()
	if err := run(*addr); err != nil {
		log.Fatal(err)
	}
}
