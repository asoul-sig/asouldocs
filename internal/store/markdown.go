// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"bytes"
	"net/url"
	"os"
	"path"

	goldmarktoc "github.com/abhinav/goldmark-toc"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	goldmarkmeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

func convertFile(pathPrefix, file string) (content []byte, meta map[string]interface{}, headings goldmarktoc.Items, err error) {
	body, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "read")
	}

	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
			goldmarkhtml.WithXHTML(),
			goldmarkhtml.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			goldmarkmeta.Meta,
			emoji.Emoji,
			highlighting.NewHighlighting(
				highlighting.WithStyle("base16-snazzy"),
			),
			extension.NewFootnote(),
		),
	)

	ctx := parser.NewContext()
	doc := md.Parser().Parse(text.NewReader(body), parser.WithContext(ctx))

	// Headings
	tree, err := goldmarktoc.Inspect(doc, body)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inspect headings")
	}
	headings = tree.Items
	if len(headings) > 0 {
		headings = headings[0].Items
	}

	// Links
	err = inspectLinks(pathPrefix, doc)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inspect links")
	}

	var buf bytes.Buffer
	err = md.Renderer().Render(&buf, body, doc)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "render")
	}

	return buf.Bytes(), goldmarkmeta.Get(ctx), headings, nil
}

func inspectLinks(pathPrefix string, doc ast.Node) error {
	return ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}

		dest, err := url.Parse(string(link.Destination))
		if err != nil {
			return ast.WalkContinue, nil
		}

		if dest.Scheme == "http" || dest.Scheme == "https" {
			// TODO: external links adds an SVG
			return ast.WalkSkipChildren, nil
		} else if dest.Scheme != "" {
			return ast.WalkContinue, nil
		}

		// Example: README.md => /docs/introduction
		if bytes.EqualFold(link.Destination, []byte(readme+".md")) {
			link.Destination = []byte(pathPrefix)
			return ast.WalkSkipChildren, nil
		}

		// Example: "installation.md" => "installation"
		link.Destination = bytes.TrimSuffix(link.Destination, []byte(".md"))

		// Example: "../howto/README" => "../howto/"
		link.Destination = bytes.TrimSuffix(link.Destination, []byte(readme))

		// Example: ("/docs", "../howto/") => "/docs/howto"
		link.Destination = []byte(path.Join(pathPrefix, string(link.Destination)))
		return ast.WalkSkipChildren, nil
	})
}
