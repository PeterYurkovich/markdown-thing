package md2html

import (
	"errors"
	"log"
	"net/http"

	"github.com/PeterYurkovich/markdown-thing/templates"
)

func Server(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if path == "" {
		http.Redirect(w, r, "README.md", http.StatusMovedPermanently)
		return
	}
	if BlockedLink(path) {
		w.Header().Set("Content-Type", "text/html")
		template := templates.Blocked(path)
		err := template.Render(r.Context(), w)
		if err != nil {
			log.Println(err)
		}
		return
	}
	w.Header().Set("Content-Type", "text/html")
	markdown, err := MarkdownLookup(path)
	if err != nil {
		if errors.Is(err, Error404) {
			w.WriteHeader(http.StatusNotFound)
			template := templates.Error(http.StatusNotFound, "Not Found")
			err := template.Render(r.Context(), w)
			if err != nil {
				log.Println(err)
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		template := templates.Error(http.StatusInternalServerError, "Internal Server Error")
		err := template.Render(r.Context(), w)
		if err != nil {
			log.Println(err)
		}
		return
	}
	template := templates.Body(markdown)
	err = template.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
	return
}
