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
	"fmt"
	"strings"

	"github.com/unknwon/com"

	"github.com/peachdocs/peach/models"
	"github.com/peachdocs/peach/pkg/context"
	"github.com/peachdocs/peach/pkg/setting"
)

func Home(ctx *context.Context) {
	if !setting.Page.HasLandingPage {
		ctx.Redirect(setting.Page.DocsBaseURL)
		return
	}

	ctx.HTML(200, "home")
}

func Pages(ctx *context.Context) {
	toc := models.Tocs[ctx.Locale.Language()]
	if toc == nil {
		toc = models.Tocs[setting.Docs.Langs[0]]
	}

	pageName := strings.ToLower(strings.TrimSuffix(ctx.Req.URL.Path[1:], ".html"))
	for i := range toc.Pages {
		if toc.Pages[i].Name == pageName {
			page := toc.Pages[i]
			langVer := toc.Lang
			if !com.IsFile(page.FileName) {
				ctx.Data["IsShowingDefault"] = true
				langVer = setting.Docs.Langs[0]
				page = models.Tocs[langVer].Pages[i]
			}
			if !setting.ProdMode {
				page.ReloadContent()
			}

			ctx.Data["Title"] = page.Title
			ctx.Data["Content"] = fmt.Sprintf(`<script type="text/javascript" src="/%s/%s?=%d"></script>`, langVer, page.DocumentPath+".js", page.LastBuildTime)
			ctx.Data["Pages"] = toc.Pages

			renderEditPage(ctx, page.DocumentPath)
			ctx.HTML(200, "docs")
			return
		}
	}

	NotFound(ctx)
}

func NotFound(ctx *context.Context) {
	ctx.Data["Title"] = "404"
	ctx.HTML(404, "404")
}
