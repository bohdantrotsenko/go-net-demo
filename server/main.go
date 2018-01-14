package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

// feel free to uncomment some logs to play

type srv struct {
	counter int32
	flag    chan struct{}
	closer  sync.Once
}

func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fl := w.(http.Flusher)
		res := atomic.AddInt32(&s.counter, 1) - 1

		//log.Printf("request %d accepted.\n", res)

		fl.Flush() // this will force issuing a reponse with a chunked transfer

		// wait for DELETE
		// this is an important part, the following will block until
		// the channel is closed (in our case)
		// that implies a hanging connection
		_, _ = <-s.flag

		//log.Printf("request %d finishes.\n", res)
		w.Write(([]byte)(fmt.Sprintf("%d", res)))
		atomic.AddInt32(&s.counter, -1)
		return
	} else if r.Method == "DELETE" {
		log.Println("received DELETE request")
		close(s.flag) // this closes the 'flag' channel
		s.flag = make(chan struct{})
		log.Println("all hanging POST calls will now reply")
		return
	}
	log.Printf("unexpected request %q %q", r.Method, r.URL)
	http.Error(w, "bad request", http.StatusBadRequest)
}

func run(addr string) error {
	s := &srv{
		flag: make(chan struct{}),
	}
	return http.ListenAndServe(addr, s)
}

func main() {
	addr := flag.String("addr", ":48787", "listening address")
	flag.Parse()
	if err := run(*addr); err != nil {
		log.Fatal(err)
	}
}
