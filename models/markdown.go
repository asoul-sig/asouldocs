// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"bytes"

	"github.com/russross/blackfriday"
)

var (
	tab    = []byte("\t")
	spaces = []byte("    ")
)

type MarkdownRender struct {
	blackfriday.Renderer
}

func (mr *MarkdownRender) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	var tmp bytes.Buffer
	mr.Renderer.BlockCode(&tmp, text, lang)
	out.Write(bytes.Replace(tmp.Bytes(), tab, spaces, -1))
}

func markdown(raw []byte) []byte {
	htmlFlags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	renderer := &MarkdownRender{
		Renderer: blackfriday.HtmlRenderer(htmlFlags, "", ""),
	}

	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS

	return blackfriday.Markdown(raw, renderer, extensions)
}
