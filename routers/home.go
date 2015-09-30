// Copyright 2015 Unknwon
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

package routers

import (
	"strings"

	"github.com/Unknwon/peach/models"
	"github.com/Unknwon/peach/modules/middleware"
	"github.com/Unknwon/peach/modules/setting"
)

func Home(ctx *middleware.Context) {
	if !setting.Page.HasLandingPage {
		ctx.Redirect(setting.Page.DocsBaseURL)
		return
	}

	ctx.HTML(200, "home")
}

func Pages(ctx *middleware.Context) {
	toc := models.Tocs[ctx.Locale.Language()]
	if toc == nil {
		toc = models.Tocs[setting.Docs.Langs[0]]
	}

	pageName := strings.ToLower(ctx.Req.URL.Path[1:])
	for i := range toc.Pages {
		if toc.Pages[i].Name == pageName {
			ctx.Data["Title"] = toc.Pages[i].Title
			ctx.Data["Content"] = string(toc.Pages[i].Content())
			ctx.Data["Pages"] = toc.Pages
			ctx.HTML(200, "docs")
			return
		}
	}

	NotFound(ctx)
}

func NotFound(ctx *middleware.Context) {
	ctx.Data["Title"] = "404"
	ctx.HTML(404, "404")
}
