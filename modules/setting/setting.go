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

package setting

import (
	"github.com/Unknwon/com"
	"github.com/Unknwon/log"
	"github.com/Unknwon/macaron"
	"gopkg.in/ini.v1"
)

type NavbarItem struct {
	Icon         string
	Locale, Link string
	Blank        bool
}

var (
	AppVer   string
	ProdMode bool
	HTTPPort int

	Site struct {
		Name   string
		Desc   string
		UseCDN bool
		URL    string
	}

	Page struct {
		HasLandingPage bool
		DocsBaseURL    string

		UseCustomTpl  bool
		NavbarTplPath string
		HomeTplPath   string
		DocsTplPath   string
		FooterTplPath string
		DisqusTplPath string
	}

	Navbar struct {
		Items []*NavbarItem
	}

	Asset struct {
		CustomCSS string
	}

	Docs struct {
		Type   string
		Target string
		Secret string
		Langs  []string

		// Only used for languages are not en-US or zh-CN to bypass error panic.
		Locales map[string][]byte
	}

	Extension struct {
		EnableDisqus         bool
		DisqusShortName      string
		HighlightJSCustomCSS string
		EnableSearch         bool
	}

	Cfg *ini.File
)

func init() {
	log.Prefix = "[Peach]"

	sources := []interface{}{"conf/app.ini"}
	if com.IsFile("custom/app.ini") {
		sources = append(sources, "custom/app.ini")
	}

	var err error
	Cfg, err = macaron.SetConfig(sources[0], sources[1:]...)
	if err != nil {
		log.Fatal("Fail to load config: %v", err)
	}

	sec := Cfg.Section("")
	if sec.Key("RUN_MODE").String() == "prod" {
		ProdMode = true
		macaron.Env = macaron.PROD
		macaron.ColorLog = false
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(5555)

	sec = Cfg.Section("site")
	Site.Name = sec.Key("NAME").MustString("Peach Server")
	Site.Desc = sec.Key("DESC").String()
	Site.UseCDN = sec.Key("USE_CDN").MustBool()
	Site.URL = sec.Key("URL").String()

	sec = Cfg.Section("page")
	Page.HasLandingPage = sec.Key("HAS_LANDING_PAGE").MustBool()
	Page.DocsBaseURL = sec.Key("DOCS_BASE_URL").Validate(func(in string) string {
		if len(in) == 0 {
			return "/docs"
		} else if in[0] != '/' {
			return "/" + in
		}
		return in
	})

	Page.UseCustomTpl = sec.Key("USE_CUSTOM_TPL").MustBool()
	Page.NavbarTplPath = "navbar.html"
	Page.HomeTplPath = "home.html"
	Page.DocsTplPath = "docs.html"
	Page.FooterTplPath = "footer.html"
	Page.DisqusTplPath = "disqus.html"

	sec = Cfg.Section("navbar")
	list := sec.KeyStrings()
	Navbar.Items = make([]*NavbarItem, len(list))
	for i, name := range list {
		secName := "navbar." + sec.Key(name).String()
		Navbar.Items[i] = &NavbarItem{
			Icon:   Cfg.Section(secName).Key("ICON").String(),
			Locale: Cfg.Section(secName).Key("LOCALE").MustString(secName),
			Link:   Cfg.Section(secName).Key("LINK").MustString("/"),
			Blank:  Cfg.Section(secName).Key("BLANK").MustBool(),
		}
	}

	sec = Cfg.Section("asset")
	Asset.CustomCSS = sec.Key("CUSTOM_CSS").String()

	sec = Cfg.Section("docs")
	Docs.Type = sec.Key("TYPE").In("local", []string{"local", "remote"})
	Docs.Target = sec.Key("TARGET").String()
	Docs.Secret = sec.Key("SECRET").String()
	Docs.Langs = Cfg.Section("i18n").Key("LANGS").Strings(",")
	Docs.Locales = make(map[string][]byte)
	for _, lang := range Docs.Langs {
		if lang != "en-US" || lang != "zh-CN" {
			Docs.Locales["locale_"+lang+".ini"] = []byte("")
		}
	}

	sec = Cfg.Section("extension")
	Extension.EnableDisqus = sec.Key("ENABLE_DISQUS").MustBool()
	Extension.DisqusShortName = sec.Key("DISQUS_SHORT_NAME").String()
	Extension.HighlightJSCustomCSS = sec.Key("HIGHLIGHTJS_CUSTOM_CSS").String()
	Extension.EnableSearch = sec.Key("ENABLE_SEARCH").MustBool()
}
