// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"bytes"
	"fmt"
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
	"github.com/yuin/goldmark/util"
)

func convertFile(pathPrefix, file string) (content []byte, meta map[string]any, headings goldmarktoc.Items, err error) {
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
				highlighting.WithGuessLanguage(true),
			),
			extension.NewFootnote(),
		),
	)

	ctx := parser.NewContext(
		func(cfg *parser.ContextConfig) {
			cfg.IDs = newIDs()
		},
	)
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

func convertRelativeLink(pathPrefix string, link []byte) []byte {
	var anchor []byte
	if i := bytes.IndexByte(link, '#'); i > -1 {
		if i == 0 {
			return link
		}

		anchor = link[i:]
		link = link[:i]
	}

	// Example: README.md => /docs/introduction
	if bytes.EqualFold(link, []byte(readme+".md")) {
		link = append([]byte(pathPrefix), anchor...)
		return link
	}

	// Example: "installation.md" => "installation"
	link = bytes.TrimSuffix(link, []byte(".md"))

	// Example: "../howto/README" => "../howto/"
	link = bytes.TrimSuffix(link, []byte(readme))

	// Example: ("/docs", "../howto/") => "/docs/howto"
	link = []byte(path.Join(pathPrefix, string(link)))

	link = append(link, anchor...)
	return link
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

		link.Destination = convertRelativeLink(pathPrefix, link.Destination)
		return ast.WalkSkipChildren, nil
	})
}

// ids is a modified version to allow any non-whitespace characters instead of
// just alphabets or numerics from
// https://github.com/yuin/goldmark/blob/113ae87dd9e662b54012a596671cb38f311a8e9c/parser/parser.go#L65.
type ids struct {
	values map[string]bool
}

func newIDs() parser.IDs {
	return &ids{
		values: map[string]bool{},
	}
}

func (s *ids) Generate(value []byte, kind ast.NodeKind) []byte {
	value = util.TrimLeftSpace(value)
	value = util.TrimRightSpace(value)
	if len(value) == 0 {
		if kind == ast.KindHeading {
			value = []byte("heading")
		} else {
			value = []byte("id")
		}
	}
	if _, ok := s.values[util.BytesToReadOnlyString(value)]; !ok {
		s.values[util.BytesToReadOnlyString(value)] = true
		return value
	}
	for i := 1; ; i++ {
		newResult := fmt.Sprintf("%s-%d", value, i)
		if _, ok := s.values[newResult]; !ok {
			s.values[newResult] = true
			return []byte(newResult)
		}
	}
}

func (s *ids) Put(value []byte) {
	s.values[util.BytesToReadOnlyString(value)] = true
}
