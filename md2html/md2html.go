package md2html

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type InternalLink struct {
	ast.Leaf
	Literal string
}

func MarkdownLookup(path string) string {
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
	return string(MdToHTML(output))
}

func MdToHTML(md []byte) []byte {
	p := newMarkdownParser()
	doc := p.Parse(md)
	renderer := NewCustomizedRender()
	return markdown.Render(doc, renderer)
}

func RenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if leafNode, ok := node.(*InternalLink); ok {
		RenderInternalLink(w, leafNode, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func RenderInternalLink(w io.Writer, node *InternalLink, entering bool) (ast.WalkStatus, bool) {
	if entering {
		io.WriteString(w, "<a href=\"/")
		io.WriteString(w, node.Literal)
		io.WriteString(w, "\">")
		io.WriteString(w, node.Literal)
		io.WriteString(w, "</a>")
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func NewCustomizedRender() *html.Renderer {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank

	opts := html.RendererOptions{
		RenderNodeHook: RenderHook,
		Flags:          htmlFlags,
	}
	return html.NewRenderer(opts)
}

func ParserHook(data []byte) (ast.Node, []byte, int) {
	if node, d, n := ParseInternalLink(data); node != nil {
		return node, d, n
	}
	return nil, nil, 0
}

func newMarkdownParser() *parser.Parser {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Tables | parser.FencedCode
	p := parser.NewWithExtensions(extensions)
	p.Opts.ParserHook = ParserHook
	return p
}

func ParseInternalLink(data []byte) (ast.Node, []byte, int) {
	var node *InternalLink
	var n int
	// find an internal link which starts with [[
	if i := bytes.Index(data, []byte("[[")); i != -1 {
		if i+2 < len(data) && data[i+1] == '[' {
			if j := bytes.Index(data[i+2:], []byte("]]")); j != -1 {
				node = &InternalLink{Literal: string(data[i+2 : i+2+j])}
				n = i + 4 + j
			}
		}
	}

	return node, []byte{}, n
}
