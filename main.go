package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/PeterYurkovich/markdown-thing/md2html"
	"github.com/PeterYurkovich/markdown-thing/templates"
)

//go:generate bun tw-build
//go:embed static/css
var staticCss embed.FS

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/css/", http.FileServer(http.FS(staticCss)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "README.md", http.StatusMovedPermanently)
			return
		}
		if md2html.BlockedLink(r.URL.Path) {
			w.Header().Set("Content-Type", "text/html")
			template := templates.Blocked(r.URL.Path)
			err := template.Render(r.Context(), w)
			if err != nil {
				log.Println(err)
			}
			return
		}
		w.Header().Set("Content-Type", "text/html")
		template := templates.Body(md2html.MarkdownLookup(r.URL.Path[1:]))
		err := template.Render(r.Context(), w)
		if err != nil {
			log.Println(err)
		}
		return
	})

	log.Fatal(http.ListenAndServe(":8080", mux))

}
