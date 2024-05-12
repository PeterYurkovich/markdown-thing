package main

import (
	"bytes"
	"testing"

	"github.com/PeterYurkovich/markdown-thing/md2html"
)

func TestInnerLink(t *testing.T) {
	input := []byte(`[[000 Index]]`)
	expected := []byte("<p><a href=\"/000 Index.md\">000 Index</a></p>\n")
	output := md2html.MdToHTML(input)
	if !bytes.Equal(output, expected) {
		t.Errorf("expected %s, got %s", expected, output)
	}
}
func TestInnerLinkLookup(t *testing.T) {
	input := []byte(`[[Kubernetes]]`)
	expected := []byte("<p><a href=\"/Business/Red Hat/Notes/General/OpenShift/Kubernetes/Kubernetes.md\">Kubernetes</a></p>\n")
	output := md2html.MdToHTML(input)
	if !bytes.Equal(output, expected) {
		t.Errorf("expected %s, got %s", expected, output)
	}
}

func TestInnerLinkLookupWithinText(t *testing.T) {
	input := []byte(`Primary Reference: [[Kubernetes]]`)
	expected := []byte("<p>Primary Reference: <a href=\"/Business/Red Hat/Notes/General/OpenShift/Kubernetes/Kubernetes.md\">Kubernetes</a></p>\n")
	output := md2html.MdToHTML(input)
	if string(output) != string(expected) {
		t.Errorf("expected %s, got %s", expected, output)
	}
}

func TestCluster(t *testing.T) {
	input := md2html.MarkdownLookup("Business/Red Hat/Notes/General/OpenShift/Kubernetes/Cluster.md")
	expected := "<h1 id=\"cluster\">Cluster</h1>\n\n<p>Primary Reference: <a href=\"/Business/Red Hat/Notes/General/OpenShift/Kubernetes/Kubernetes.md\">Kubernetes</a></p>\n\n<hr>\n\n<p>When deploying <a href=\"/Business/Red Hat/Notes/General/OpenShift/Kubernetes/Kubernetes.md\">Kubernetes</a> you receive a cluster. This cluster is made up of <a href=\"/Business/Red Hat/Notes/General/OpenShift/Kubernetes/Objects/Node.md\">Node</a> which host <a href=\"/Business/Red Hat/Notes/General/OpenShift/Kubernetes/Objects/Pod.md\">Pod</a> as well as a <a href=\"/Business/Red Hat/Notes/General/OpenShift/Kubernetes/Control Pane.md\">Control Pane</a></p>\n"
	if string(input) != string(expected) {
		t.Errorf("expected %s, got %s", expected, input)
	}
}
