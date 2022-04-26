// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"encoding/json"
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
	"github.com/asoul-sig/asouldocs/internal/osutil"
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

	docstore, err := store.Init(conf.Docs.Type, conf.Docs.Target, conf.Docs.TargetDir, conf.I18n.Languages, conf.Page.DocsBasePath)
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

	// Serve assets of the documentation
	f.Use(flamego.Static(
		flamego.StaticOptions{
			Directory: docstore.RootDir(),
		},
	))

	// Load assets from disk if in development and the local directory exists
	if flamego.Env() == flamego.EnvTypeDev &&
		osutil.IsDir("public") {
		f.Use(flamego.Static())
	} else {
		f.Use(flamego.Static(
			flamego.StaticOptions{
				FileSystem: http.FS(public.Files),
				SetETag:    true,
			},
		))
	}

	// Load templates from disk if in development and the local directory exists
	funcMaps := []gotemplate.FuncMap{{
		"Year": func() int { return time.Now().Year() },
		"Safe": func(p []byte) gotemplate.HTML { return gotemplate.HTML(p) },
	}}
	if flamego.Env() == flamego.EnvTypeDev &&
		osutil.IsDir("templates") {
		f.Use(template.Templater(
			template.Options{
				AppendDirectories: []string{conf.Page.CustomDirectory},
				FuncMaps:          funcMaps,
			},
		))
	} else {
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
		data["BuildCommit"] = conf.BuildCommit
		data["Site"] = conf.Site
		data["Page"] = conf.Page
		data["Extension"] = conf.Extension

		data["Tr"] = l.Translate
		data["Lang"] = l.Lang()
		data["Languages"] = languages

		data["URL"] = r.URL.Path
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
			current := c.Param("**")
			if current == "" || current == "/" {
				c.Redirect(conf.Page.DocsBasePath + "/" + docstore.FirstDocPath())
				return
			}

			if flamego.Env() == flamego.EnvTypeDev {
				err = docstore.Reload()
				if err != nil {
					panic("reload store: " + err.Error())
				}
			}

			data["Current"] = current
			data["TOC"] = docstore.TOC(l.Lang())

			node, fallback, err := docstore.Match(l.Lang(), current)
			if err != nil {
				notFound(t, data, l)
				return
			}

			data["Fallback"] = fallback
			data["Category"] = node.Category
			data["Title"] = node.Title + " - " + l.Translate("name")
			data["Node"] = node

			if conf.Docs.EditPageLinkFormat != "" {
				blob := strings.TrimPrefix(node.LocalPath, docstore.RootDir()+"/")
				data["EditLink"] = strings.Replace(conf.Docs.EditPageLinkFormat, "{blob}", blob, 1)
			}
			t.HTML(http.StatusOK, "docs/page")
		},
	)
	f.Any("/webhook", func(w http.ResponseWriter) {
		err := docstore.Reload()
		if err != nil {
			log.Error("Failed to reload store triggered by webhook: %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	f.NotFound(notFound)

	listenAddr := fmt.Sprintf("%s:%d", conf.App.HTTPHost, conf.App.HTTPPort)
	log.Info("Listen on http://%s", listenAddr)
	if err := http.ListenAndServe(listenAddr, f); err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
