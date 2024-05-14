package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/PeterYurkovich/markdown-thing/md2html"
)

//go:generate bun tw-build
//go:embed static/css
var staticCss embed.FS

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/css/", http.FileServer(http.FS(staticCss)))
	mux.HandleFunc("/", md2html.Server)

	log.Fatal(http.ListenAndServe(":8080", mux))

}
