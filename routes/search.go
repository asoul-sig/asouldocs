// Copyright 2015 unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package routes

import (
	"github.com/peachdocs/peach/models"
	"github.com/peachdocs/peach/pkg/context"
	"github.com/peachdocs/peach/pkg/setting"
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
