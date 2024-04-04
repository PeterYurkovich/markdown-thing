package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PeterYurkovich/markdown-thing/templates"
)

// go mux hello world
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			templa := templates.Hello("Peter")
			err := templa.Render(r.Context(), w)
			if err != nil {
				log.Println(err)
			}
			return
		}
		fmt.Fprintf(w, markdownLookup(r.URL.Path[1:]))
	})

	log.Fatal(http.ListenAndServe(":8080", mux))

}
