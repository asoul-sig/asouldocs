// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"strings"
	"time"

	"gopkg.in/macaron.v1"

	"github.com/asoul-go/asouldocs/pkg/setting"
)

type Context struct {
	*macaron.Context
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context) {
		ctx := &Context{
			Context: c,
		}
		c.Map(ctx)

		ctx.Data["Link"] = strings.TrimSuffix(ctx.Req.URL.Path, ".html")
		ctx.Data["AppVer"] = setting.AppVer
		ctx.Data["Site"] = setting.Site
		ctx.Data["Page"] = setting.Page
		ctx.Data["Navbar"] = setting.Navbar
		ctx.Data["Asset"] = setting.Asset
		ctx.Data["Extension"] = setting.Extension
		ctx.Data["Year"] = time.Now().Year()
	}
}
