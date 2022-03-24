// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routes

import (
	"github.com/asoul-sig/asouldocs/models"
	"github.com/asoul-sig/asouldocs/pkg/context"
	"github.com/asoul-sig/asouldocs/pkg/setting"
)

func Search(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("search")

	toc := models.Tocs[ctx.Locale.Language()]
	if toc == nil {
		toc = models.Tocs[setting.Docs.Langs[0]]
	}

	q := ctx.Query("q")
	if len(q) == 0 {
		ctx.Redirect(setting.Page.DocsBaseURL)
		return
	}

	ctx.Data["Keyword"] = q
	ctx.Data["Results"] = toc.Search(q)

	ctx.HTML(200, "search")
}
