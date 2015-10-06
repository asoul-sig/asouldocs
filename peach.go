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

// Peach is a web server for multi-language, real-time synchronization and searchable documentation.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"

	"github.com/Unknwon/log"
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/i18n"
	"github.com/macaron-contrib/pongo2"

	"github.com/peachdocs/peach/models"
	"github.com/peachdocs/peach/modules/middleware"
	"github.com/peachdocs/peach/modules/setting"
	"github.com/peachdocs/peach/routers"
)

const APP_VER = "0.5.2.1005"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setting.AppVer = APP_VER

	config := flag.String("config", "custom/app.ini", "custom config path")
	flag.Parse()

	setting.CustomConf = *config

	setting.NewContext()
	models.NewContext()
}

func main() {
	log.Info("Peach %s", APP_VER)

	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Statics(macaron.StaticOptions{
		SkipLogging: setting.ProdMode,
	}, "custom/public", "public"))
	m.Use(i18n.I18n(i18n.Options{
		Files: setting.Docs.Locales,
	}))
	tplDir := "templates"
	if setting.Page.UseCustomTpl {
		tplDir = "custom/templates"
	}
	m.Use(pongo2.Pongoer(pongo2.Options{
		Directory: tplDir,
	}))
	m.Use(middleware.Contexter())

	m.Get("/", routers.Home)
	m.Get("/docs", routers.Docs)
	m.Get("/docs/images/*", routers.DocsStatic)
	m.Get("/docs/*", routers.Docs)
	m.Post("/hook", routers.Hook)
	m.Get("/search", routers.Search)
	m.Get("/*", routers.Pages)

	m.NotFound(routers.NotFound)

	listenAddr := fmt.Sprintf("0.0.0.0:%d", setting.HTTPPort)
	log.Info("%s Listen on %s", setting.Site.Name, listenAddr)
	log.Fatal("Fail to start Peach: %v", http.ListenAndServe(listenAddr, m))
}
