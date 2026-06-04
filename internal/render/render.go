package render

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type Options struct {
	Highlight bool
	Chroma    string
}

var (
	cached   goldmark.Markdown
	cachedOp Options
)

func Markdown(src []byte, opt Options) ([]byte, error) {
	md := getRenderer(opt)
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}
	return buf.Bytes(), nil
}

func getRenderer(opt Options) goldmark.Markdown {
	if cached != nil && opt == cachedOp {
		return cached
	}
	exts := []goldmark.Extender{
		extension.GFM,
		extension.Linkify,
		extension.Footnote,
		extension.DefinitionList,
		extension.Typographer,
	}
	if opt.Highlight {
		cfg := []highlighting.Option{highlighting.WithGuessLanguage(true)}
		if opt.Chroma != "" {
			cfg = append(cfg, highlighting.WithStyle(opt.Chroma))
		}
		exts = append(exts, highlighting.NewHighlighting(cfg...))
	}
	md := goldmark.New(
		goldmark.WithExtensions(exts...),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	cached = md
	cachedOp = opt
	return md
}
