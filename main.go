package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/PeterYurkovich/markdown-thing/md2html"
	"github.com/PeterYurkovich/markdown-thing/server"
	"github.com/PeterYurkovich/markdown-thing/templates"
)

//go:generate bun tw-build
//go:embed static/css
var staticCss embed.FS

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/css/", http.FileServer(http.FS(staticCss)))
	mux.HandleFunc("/tree", func(w http.ResponseWriter, r *http.Request) {
		tree, err := md2html.GetMarkdownTree()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			template := templates.Error(http.StatusInternalServerError, "Internal Server Error")
			err := template.Render(r.Context(), w)
			if err != nil {
				log.Println(err)
			}
			return
		}
		w.Header().Set("Content-Type", "text/html")
		template := templates.TreeTemplate(tree)
		err = template.Render(r.Context(), w)
		if err != nil {
			log.Println(err)
		}
		return
	})
	mux.HandleFunc("/", server.Server)

	log.Fatal(http.ListenAndServe(":8080", mux))

}
