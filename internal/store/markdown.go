// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"bytes"
	"os"

	goldmarktoc "github.com/abhinav/goldmark-toc"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	goldmarkmeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

func convertFile(file string) (content []byte, meta map[string]interface{}, headings goldmarktoc.Items, err error) {
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
		),
	)

	ctx := parser.NewContext()
	doc := md.Parser().Parse(text.NewReader(body), parser.WithContext(ctx))
	tree, err := goldmarktoc.Inspect(doc, body)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inspect TOC")
	}

	var buf bytes.Buffer
	err = md.Renderer().Render(&buf, body, doc)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "render")
	}

	headings = tree.Items
	if len(headings) > 0 {
		headings = headings[0].Items
	}
	return buf.Bytes(), goldmarkmeta.Get(ctx), headings, nil
}
