package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("req: %s %s %s, reqid: %s", req.Host, req.URL.Path, req.UserAgent(), req.Header.Get("X-Reqid"))
		switch req.URL.Path {
		case "/flush":
			flusher, ok := w.(http.Flusher)
			if !ok {
				panic("expected http.ResponseWriter to be an http.Flusher")
			}

			for i := 0; i < 16; i++ {
				fmt.Fprintf(w, "chunk [%02d]: %v\n", i, time.Now())
				flusher.Flush()
				time.Sleep(time.Millisecond * 10)
			}
			log.Printf("send ok")
		case "/bao.jpg":
			w.WriteHeader(200)
			for i := 0; i < 1000*10; i++ {
				time.Sleep(100 * time.Millisecond)
				w.Write(make([]byte, 1000))
			}
		}

	})

	err := http.ListenAndServe(":8090", nil)
	log.Println(err)
}
