package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/PeterYurkovich/markdown-thing/md2html"
	"github.com/PeterYurkovich/markdown-thing/templates"
)

func Server(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if path == "" {
		http.Redirect(w, r, "README.md", http.StatusMovedPermanently)
		return
	}
	if md2html.BlockedLink(path) {
		w.Header().Set("Content-Type", "text/html")
		template := templates.Blocked(path)
		err := template.Render(r.Context(), w)
		if err != nil {
			log.Println(err)
		}
		return
	}
	w.Header().Set("Content-Type", "text/html")
	markdown, err := md2html.MarkdownLookup(path)
	if err != nil {
		if errors.Is(err, md2html.Error404) {
			w.WriteHeader(http.StatusNotFound)
			template := templates.ErrorPage(http.StatusNotFound, "Not Found")
			err := template.Render(r.Context(), w)
			if err != nil {
				log.Println(err)
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		template := templates.ErrorPage(http.StatusInternalServerError, "Internal Server Error")
		err := template.Render(r.Context(), w)
		if err != nil {
			log.Println(err)
		}
		return
	}
	template := templates.RawBody(markdown)
	err = template.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
	return
}
