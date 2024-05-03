package main

import (
	"bytes"
	"testing"

	"github.com/PeterYurkovich/markdown-thing/md2html"
)

func TestInnerLink(t *testing.T) {
	input := []byte(`[[000 Index]]`)
	expected := []byte(`<a href="/000 Index">000 Index</a>`)
	output := md2html.MdToHTML(input)
	if !bytes.Equal(output, expected) {
		t.Errorf("expected %s, got %s", expected, output)
	}
}
