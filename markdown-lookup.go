package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func markdownLookup(path string) string {
	path = fmt.Sprintf("markdown/%s.md", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "no"
	}
	file, err := os.Open(path)
	if err != nil {
		return "no"
	}
	defer file.Close()
	buffer := make([]byte, 1024)
	output := []byte{}
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "no"
		}
		output = append(output, buffer[:n]...)
	}
	return string(mdToHTML(output))
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
