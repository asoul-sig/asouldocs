// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	gotemplate "html/template"
	"net/http"
	"strings"
	"time"

	"github.com/flamego/flamego"
	"github.com/flamego/i18n"
	"github.com/flamego/template"
	"github.com/urfave/cli"
	log "unknwon.dev/clog/v2"

	"github.com/asoul-sig/asouldocs/conf/locale"
	"github.com/asoul-sig/asouldocs/internal/conf"
	"github.com/asoul-sig/asouldocs/internal/route"
	"github.com/asoul-sig/asouldocs/public"
	"github.com/asoul-sig/asouldocs/templates"
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
	err := conf.Init(ctx.String("config"))
	if err != nil {
		log.Fatal("Failed to init configuration: %v", err)
	}

	// models.NewContext() // todo

	log.Info("ASoulDocs %s", conf.App.Version)

	f := flamego.New()
	f.Use(flamego.Recovery())
	f.Use(flamego.Static(
		flamego.StaticOptions{
			FileSystem: http.FS(public.Files),
			SetETag:    true,
		},
	))
	f.Use(flamego.Static(
		flamego.StaticOptions{
			Directory: conf.Asset.CustomDirectory,
			SetETag:   true,
		},
	))

	fs, err := template.EmbedFS(templates.Files, ".", []string{".html"})
	if err != nil {
		log.Fatal("Failed to convert template files to embed.FS: %v", err)
	}
	f.Use(template.Templater(
		template.Options{
			FileSystem:        fs,
			AppendDirectories: []string{conf.Page.CustomDirectory},
			FuncMaps: []gotemplate.FuncMap{{
				"Year": func() int { return time.Now().Year() },
			}},
		},
	))

	languages := make([]i18n.Language, len(conf.I18n.Langs))
	for i := range conf.I18n.Langs {
		languages[i] = i18n.Language{
			Name:        conf.I18n.Langs[i],
			Description: conf.I18n.Names[i],
		}
	}
	f.Use(i18n.I18n(
		i18n.Options{
			FileSystem:        http.FS(locale.Files),
			AppendDirectories: []string{conf.I18n.CustomDirectory},
			Languages:         languages,
		},
	))

	f.Use(func(r *http.Request, data template.Data, l i18n.Locale) {
		data["Tr"] = l.Translate
		data["URL"] = strings.TrimSuffix(r.URL.Path, ".html")
		data["Site"] = conf.Site
		data["Page"] = conf.Page
		data["Languages"] = languages
		// todo
		// data["Extension"] = setting.Extension
	})

	f.Get("/", route.Home)

	f.NotFound(func(t template.Template, data template.Data, l i18n.Locale) {
		data["Title"] = l.Translate("status::404")
		t.HTML(http.StatusNotFound, "404")
	})

	// m.Get(setting.Page.DocsBaseURL, routes.Docs)
	// m.Get(setting.Page.DocsBaseURL+"/images/*", routes.DocsStatic)
	// m.Get(setting.Page.DocsBaseURL+"/*", routes.Protect, routes.Docs)
	// m.Post("/hook", routes.Hook)
	// m.Get("/*", routes.Pages)

	listenAddr := fmt.Sprintf("%s:%d", conf.App.HTTPHost, conf.App.HTTPPort)
	log.Info("Listen on http://%s", listenAddr)
	if err := http.ListenAndServe(listenAddr, f); err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
