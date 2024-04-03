package main

import (
	"fmt"
	"log"
	"net/http"
)

// go mux hello world
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	mux.HandleFunc(`/{tree}`, func(w http.ResponseWriter, r *http.Request) {
		markdownLookup(r.URL.Path[1:])
	})

	log.Fatal(http.ListenAndServe(":8080", mux))

}
