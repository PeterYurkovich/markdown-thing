package md2html

import (
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type InternalLink struct {
	ast.Leaf
	Link string
	Name string
}

type Text struct {
	ast.Container
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
		io.WriteString(w, node.Link)
		io.WriteString(w, "\">")
		io.WriteString(w, node.Name)
		io.WriteString(w, "</a>")
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func NewCustomizedRender() *html.Renderer {
	htmlFlags := html.CommonFlags

	opts := html.RendererOptions{
		RenderNodeHook: RenderHook,
		Flags:          htmlFlags,
	}
	return html.NewRenderer(opts)
}

func newMarkdownParser() *parser.Parser {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Tables | parser.FencedCode
	parser := parser.NewWithExtensions(extensions)

	prev := parser.RegisterInline('[', nil)
	parser.RegisterInline('[', wikiLink(parser, prev))
	return parser
}

func wikiLink(_ *parser.Parser, fn parser.InlineParser) parser.InlineParser {
	return func(p *parser.Parser, original []byte, offset int) (int, ast.Node) {
		data := original[offset:]
		n := len(data)
		if n < 5 || data[1] != '[' {
			return fn(p, original, offset)
		}
		i := 2
		for i+1 < n && data[i] != ']' && data[i+1] != ']' {
			i++
		}
		text := data[2 : i+1]
		foundLink, found := GetLinkFromName(string(text))
		if !found {
			return fn(p, original, offset)
		}
		link := &ast.Link{
			Destination: []byte(foundLink.Link),
		}
		ast.AppendChild(link, &ast.Text{Leaf: ast.Leaf{Literal: []byte(foundLink.Name)}})
		return i + 4, link
	}
}
