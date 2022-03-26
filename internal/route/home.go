// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"github.com/flamego/flamego"
	"github.com/flamego/i18n"
	"github.com/flamego/template"

	"github.com/asoul-sig/asouldocs/internal/conf"
)

// GET /
func Home(c flamego.Context, t template.Template, data template.Data, l i18n.Locale) {
	if !conf.Page.HasLandingPage {
		c.Redirect(conf.Page.DocsBaseURL)
		return
	}

	data["Title"] = l.Translate("name") + " - " + l.Translate("tag_line")
	t.HTML(http.StatusOK, "home")
}
