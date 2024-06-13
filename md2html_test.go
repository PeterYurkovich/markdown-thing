package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/PeterYurkovich/markdown-thing/md2html"
)

func TestInnerLink(t *testing.T) {
	input := []byte(`[[000 Index]]`)
	expected := []byte("<p><a class=\"text-vesper-highlight\" href=\"/000 Index.md\">000 Index</a></p>\n")
	output := md2html.MdToHTML(input)
	if !bytes.Equal(output, expected) {
		t.Errorf("expected %s, got %s", expected, output)
	}
}
func TestInnerLinkLookup(t *testing.T) {
	input := []byte(`[[Kubernetes]]`)
	expected := []byte("<p><a class=\"text-vesper-highlight\" href=\"/Wiki/Kubernetes/Kubernetes.md\">Kubernetes</a></p>\n")
	output := md2html.MdToHTML(input)
	if !bytes.Equal(output, expected) {
		t.Errorf("expected %s, got %s", expected, output)
	}
}

func TestInnerLinkLookupWithinText(t *testing.T) {
	input := []byte(`Primary Reference: [[Kubernetes]]`)
	expected := []byte("<p>Primary Reference: <a class=\"text-vesper-highlight\" href=\"/Wiki/Kubernetes/Kubernetes.md\">Kubernetes</a></p>\n")
	output := md2html.MdToHTML(input)
	if string(output) != string(expected) {
		t.Errorf("expected %s, got %s", expected, output)
	}
}

func TestCluster(t *testing.T) {
	input, err := md2html.MarkdownLookup("Wiki/Kubernetes/Cluster.md")
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected := "<h1 id=\"cluster\">Cluster</h1>\n\n<p>Primary Reference: <a class=\"text-vesper-highlight\" href=\"/Wiki/Kubernetes/Kubernetes.md\">Kubernetes</a></p>\n\n<hr>\n\n<p>When deploying <a class=\"text-vesper-highlight\" href=\"/Wiki/Kubernetes/Kubernetes.md\">Kubernetes</a> you receive a cluster. This cluster is made up of <a class=\"text-vesper-highlight\" href=\"/Wiki/Kubernetes/Objects/Node.md\">Node</a> which host <a class=\"text-vesper-highlight\" href=\"/Wiki/Kubernetes/Objects/Pod.md\">Pod</a> as well as a <a class=\"text-vesper-highlight\" href=\"/Wiki/Kubernetes/Control Pane.md\">Control Pane</a></p>\n"
	if string(input) != string(expected) {
		t.Errorf("expected %s, got %s", expected, input)
	}
}

func Test404(t *testing.T) {
	_, err := md2html.MarkdownLookup("DoesntExist.md")
	if err == nil || !errors.Is(err, md2html.Error404) {
		t.Errorf("expected no error, got %s", err)
	}
}
func TestBlocked(t *testing.T) {
	if !md2html.BlockedLink("Business/Accrisoft/Accrisoft MOC.md") {
		t.Errorf("expected link to be blocked")
	}
}

func Test404WithIgnoredLink(t *testing.T) {
	_, err := md2html.MarkdownLookup(".obsidian/types.json")
	if err == nil || !errors.Is(err, md2html.Error404) {
		t.Errorf("expected no error, got %s", err)
	}
}
