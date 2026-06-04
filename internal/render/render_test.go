package render

import (
	"strings"
	"testing"
)

func TestMarkdown(t *testing.T) {
	in := []byte(`# Hello

A paragraph with **bold** and *italic* and a [link](https://example.com).

- one
- two
- three

` + "```go" + `
package main

func main() {}
` + "```" + `
`)
	out, err := Markdown(in, Options{Highlight: true, Chroma: "github"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"<h1",
		"Hello",
		"<strong>bold</strong>",
		"<em>italic</em>",
		`href="https://example.com"`,
		"<ul>",
		"<li>one</li>",
		"<pre",
		"<code",
		"package",
		"font-weight:bold", // chroma inline style for "package" keyword
	} {
		if !strings.Contains(string(out), want) {
			t.Errorf("output missing %q\n---\n%s\n---", want, out)
		}
	}
}

func TestGFMTable(t *testing.T) {
	in := []byte("| a | b |\n|---|---|\n| 1 | 2 |\n")
	out, err := Markdown(in, Options{Highlight: false})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "<table>") {
		t.Errorf("expected <table>, got:\n%s", out)
	}
}

func TestStrikethrough(t *testing.T) {
	in := []byte("~~gone~~")
	out, err := Markdown(in, Options{Highlight: false})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "<del>gone</del>") {
		t.Errorf("expected <del>, got: %s", out)
	}
}
