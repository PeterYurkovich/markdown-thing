package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/PeterYurkovich/markdown-thing/templates"
)

//go:generate bun tw-build
//go:embed static/css
var staticCss embed.FS

func main() {
	mux := http.NewServeMux()

	// assuming you have a net/http#ServeMux called `mux`
	mux.Handle("/static/css/", http.FileServer(http.FS(staticCss)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			templa := templates.Hello("Peter")
			err := templa.Render(r.Context(), w)
			if err != nil {
				log.Println(err)
			}
			return
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, markdownLookup(r.URL.Path[1:]))
	})

	log.Fatal(http.ListenAndServe(":8080", mux))

}
