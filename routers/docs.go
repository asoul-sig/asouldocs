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
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/Unknwon/com"
	"github.com/Unknwon/log"

	"github.com/peachdocs/peach/models"
	"github.com/peachdocs/peach/modules/middleware"
	"github.com/peachdocs/peach/modules/setting"
)

func renderEditPage(ctx *middleware.Context, documentPath string) {
	if setting.Extension.EnableEditPage {
		ctx.Data["EditPageLink"] = com.Expand(setting.Extension.EditPageLinkFormat, map[string]string{
			"lang": ctx.Locale.Language(),
			"blob": documentPath + ".md",
		})
	}
}

func Docs(ctx *middleware.Context) {
	toc := models.Tocs[ctx.Locale.Language()]
	if toc == nil {
		toc = models.Tocs[setting.Docs.Langs[0]]
	}
	ctx.Data["Toc"] = toc

	nodeName := strings.TrimPrefix(strings.ToLower(strings.TrimSuffix(ctx.Req.URL.Path, ".html")), setting.Page.DocsBaseURL)
	node, isDefault := toc.GetDoc(nodeName)
	if node == nil {
		NotFound(ctx)
		return
	}
	if !setting.ProdMode {
		node.ReloadContent()
	}

	langVer := toc.Lang
	if isDefault {
		ctx.Data["IsShowingDefault"] = isDefault
		langVer = setting.Docs.Langs[0]
	}
	ctx.Data["Title"] = node.Title
	ctx.Data["Content"] = fmt.Sprintf(`<script type="text/javascript" src="/%s/%s?=%d"></script>`, langVer, node.DocumentPath+".js", node.LastBuildTime)

	renderEditPage(ctx, node.DocumentPath)
	ctx.HTML(200, "docs")
}

func DocsStatic(ctx *middleware.Context) {
	if len(ctx.Params("*")) > 0 {
		f, err := os.Open(path.Join(models.Tocs[setting.Docs.Langs[0]].RootPath, "images", ctx.Params("*")))
		if err != nil {
			ctx.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		defer f.Close()

		_, err = io.Copy(ctx.Resp, f)
		if err != nil {
			ctx.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		return
	}
	ctx.Error(404)
}

func Hook(ctx *middleware.Context) {
	if ctx.Query("secret") != setting.Docs.Secret {
		ctx.Error(403)
		return
	}

	log.Info("Incoming hook update request")
	if err := models.ReloadDocs(); err != nil {
		ctx.Error(500)
		return
	}
	ctx.Status(200)
}
