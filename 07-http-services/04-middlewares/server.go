package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func logger(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("request received")
		handler(w, r)
		log.Println("response served")
	}
}

func profile(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			end := time.Now()
			elapsed := end.Sub(start) / time.Millisecond
			fmt.Printf("%q took %dms\n", r.URL.Path, elapsed)
		}()
		handler(w, r)
	}
}

func foo(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	w.Write([]byte("Foo"))
}

func bar(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	w.Write([]byte("Bar"))
}

func main() {
	fooWithLogger := logger(foo)
	barWithLogger := logger(bar)
	fooWithLoggerAndProfile := profile(fooWithLogger)
	barWithLoggerAndProfile := profile(barWithLogger)
	http.HandleFunc("/foo", fooWithLoggerAndProfile)
	http.HandleFunc("/bar", barWithLoggerAndProfile)
	http.ListenAndServe(":8080", nil)
}
