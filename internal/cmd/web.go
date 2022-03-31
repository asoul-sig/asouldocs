// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	gotemplate "html/template"
	"net/http"
	"time"

	"github.com/flamego/flamego"
	"github.com/flamego/i18n"
	"github.com/flamego/template"
	"github.com/urfave/cli"
	log "unknwon.dev/clog/v2"

	"github.com/asoul-sig/asouldocs/conf/locale"
	"github.com/asoul-sig/asouldocs/internal/conf"
	"github.com/asoul-sig/asouldocs/internal/store"
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
	log.Info("ASoulDocs %s", conf.App.Version)

	tocs, err := store.Init(conf.Docs.Type, conf.Docs.Target, conf.Docs.TargetDir, conf.I18n.Languages)
	if err != nil {
		log.Fatal("Failed to init store: %v", err)
	}

	f := flamego.New()
	f.Use(flamego.Recovery())

	// Custom assets should be served first to support overwrite
	f.Use(flamego.Static(
		flamego.StaticOptions{
			Directory: conf.Asset.CustomDirectory,
			SetETag:   true,
		},
	))

	// Load assets and templates directly from disk in development
	funcMaps := []gotemplate.FuncMap{{
		"Year": func() int { return time.Now().Year() },
		"Safe": func(p []byte) gotemplate.HTML { return gotemplate.HTML(p) },
	}}
	if flamego.Env() == flamego.EnvTypeDev {
		f.Use(flamego.Static())
		f.Use(template.Templater(
			template.Options{
				AppendDirectories: []string{conf.Page.CustomDirectory},
				FuncMaps:          funcMaps,
			},
		))
	} else {
		f.Use(flamego.Static(
			flamego.StaticOptions{
				FileSystem: http.FS(public.Files),
				SetETag:    true,
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
				FuncMaps:          funcMaps,
			},
		))
	}

	languages := make([]i18n.Language, len(conf.I18n.Languages))
	for i := range conf.I18n.Languages {
		languages[i] = i18n.Language{
			Name:        conf.I18n.Languages[i],
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
		data["URL"] = r.URL.Path
		data["Site"] = conf.Site
		data["Page"] = conf.Page
		data["Languages"] = languages
	})

	notFound := func(t template.Template, data template.Data, l i18n.Locale) {
		data["Title"] = l.Translate("status::404")
		t.HTML(http.StatusNotFound, "404")
	}

	f.Get("/",
		func(c flamego.Context, t template.Template, data template.Data, l i18n.Locale) {
			if !conf.Page.HasLandingPage {
				c.Redirect(conf.Page.DocsBasePath)
				return
			}

			data["Title"] = l.Translate("name") + " - " + l.Translate("tag_line")
			t.HTML(http.StatusOK, "home")
		},
	)
	f.Get(conf.Page.DocsBasePath+"/?{**}",
		func(c flamego.Context, t template.Template, data template.Data, l i18n.Locale) {
			toc, ok := tocs[l.Lang()]
			if !ok {
				toc = tocs[conf.I18n.Languages[0]]
			}

			current := c.Param("**")
			if current == "" || current == "/" {
				c.Redirect(conf.Page.DocsBasePath + "/" + toc.Nodes[0].Path)
				return
			}

			data["Current"] = current
			data["TOC"] = toc

			// TODO: fallback to default language if given page does not exist in the current language, and display notice
			var node *store.Node
		loop:
			for _, dir := range toc.Nodes {
				data["Category"] = dir.Name
				if dir.Path == current {
					node = dir
					break loop
				}

				for _, file := range dir.Nodes {
					if file.Path == current {
						node = file
						break loop
					}
				}
			}

			if node == nil {
				notFound(t, data, l)
				return
			}

			if flamego.Env() == flamego.EnvTypeDev {
				err = node.Reload()
				if err != nil {
					panic("reload node: " + err.Error())
				}
			}

			data["Title"] = node.Name + " - " + l.Translate("name")
			data["Node"] = node
			t.HTML(http.StatusOK, "docs")
		},
	)

	f.NotFound(notFound)

	listenAddr := fmt.Sprintf("%s:%d", conf.App.HTTPHost, conf.App.HTTPPort)
	log.Info("Listen on http://%s", listenAddr)
	if err := http.ListenAndServe(listenAddr, f); err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
