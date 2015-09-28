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
	"io"
	"os"
	"path"
	"strings"

	"github.com/Unknwon/peach/models"
	"github.com/Unknwon/peach/modules/middleware"
	"github.com/Unknwon/peach/modules/setting"
)

func Docs(ctx *middleware.Context) {
	toc := models.Toc[ctx.Locale.Language()]
	if toc == nil {
		toc = models.Toc[setting.Docs.Langs[0]]
	}
	ctx.Data["Toc"] = toc

	nodeName := strings.TrimPrefix(ctx.Req.URL.Path, setting.Page.DocsBaseURL)

	title, content := toc.GetDoc(nodeName)
	if content == nil {
		NotFound(ctx)
		return
	}
	ctx.Data["Title"] = title
	ctx.Data["Content"] = string(content)
	ctx.HTML(200, "docs")
}

func DocsStatic(ctx *middleware.Context) {
	if len(ctx.Params("*")) > 0 {
		f, err := os.Open(path.Join(models.LocalRoot, "images", ctx.Params("*")))
		if err != nil {
			ctx.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		defer f.Close()

		_, err = io.Copy(ctx.RW(), f)
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
