// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"net/http"

	"github.com/go-macaron/i18n"
	"github.com/go-macaron/pongo2"
	"github.com/unknwon/log"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"

	"github.com/asoul-go/asouldocs/models"
	"github.com/asoul-go/asouldocs/pkg/context"
	"github.com/asoul-go/asouldocs/pkg/setting"
	"github.com/asoul-go/asouldocs/routes"
)

var Web = cli.Command{
	Name:   "web",
	Usage:  "Start the web server",
	Action: runWeb,
	Flags: []cli.Flag{
		stringFlag("config, c", "custom/app.ini", "Custom configuration file path"),
	},
}

func runWeb(ctx *cli.Context) {
	if ctx.IsSet("config") {
		setting.CustomConf = ctx.String("config")
	}
	setting.NewContext()
	models.NewContext()

	log.Info("Peach %s", setting.AppVer)

	m := macaron.New()
	if !setting.ProdMode {
		m.Use(macaron.Logger())
	}
	m.Use(macaron.Recovery())
	m.Use(macaron.Statics(macaron.StaticOptions{
		SkipLogging: setting.ProdMode,
	}, "custom/public", "public", models.HTMLRoot))
	m.Use(i18n.I18n(i18n.Options{
		Files:       setting.Docs.Locales,
		DefaultLang: setting.Docs.Langs[0],
	}))
	tplDir := "templates"
	if setting.Page.UseCustomTpl {
		tplDir = "custom/templates"
	}
	m.Use(pongo2.Pongoer(pongo2.Options{
		Directory: tplDir,
	}))
	m.Use(context.Contexter())

	m.Get("/", routes.Home)
	m.Get(setting.Page.DocsBaseURL, routes.Docs)
	m.Get(setting.Page.DocsBaseURL+"/images/*", routes.DocsStatic)
	m.Get(setting.Page.DocsBaseURL+"/*", routes.Protect, routes.Docs)
	m.Post("/hook", routes.Hook)
	m.Get("/search", routes.Search)
	m.Get("/*", routes.Pages)

	listenAddr := fmt.Sprintf("%s:%d", setting.HTTPHost, setting.HTTPPort)
	log.Info("%s Listen on %s", setting.Site.Name, listenAddr)
	log.Fatal("Fail to start Peach: %v", http.ListenAndServe(listenAddr, m))
}
