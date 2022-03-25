// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package conf

var (
	// Application settings
	App struct {
		// ⚠️ WARNING: Should only be set by the main package (i.e. "main.go").
		Version string `ini:"-"`

		RunMode  string
		HTTPHost string `ini:"HTTP_HOST"`
		HTTPPort int    `ini:"HTTP_PORT"`
	}

	// Site settings
	Site struct {
		Description string
		ExternalURL string `ini:"EXTERNAL_URL"`
	}

	// Asset settings
	Asset struct {
		CustomDirectory string
	}

	// Page settings
	Page struct {
		HasLandingPage  bool
		DocsBaseURL     string `ini:"DOCS_BASE_URL"`
		CustomDirectory string

		// todo
		// NavbarTplPath  string
		// HomeTplPath    string
		// DocsTplPath    string
		// FooterTplPath  string
		// DisqusTplPath  string
		// DuoShuoTplPath string
	}

	// I18n settings
	I18n struct {
		Langs           []string
		Names           []string
		CustomDirectory string
	}
)

// todo

var (
	Navbar struct {
		Items []*NavbarItem
	}

	Docs struct {
		Type      DocType
		Target    string
		TargetDir string
		Secret    string
		Langs     []string

		// Only used for languages are not en-US or zh-CN to bypass error panic.
		Locales map[string][]byte
	}

	Extension struct {
		EnableEditPage       bool
		EditPageLinkFormat   string
		EnableDisqus         bool
		DisqusShortName      string
		EnableDuoShuo        bool
		DuoShuoShortName     string
		HighlightJSCustomCSS string
		EnableSearch         bool
		GABlock              string
	}
)

type NavbarItem struct {
	Icon         string
	Locale, Link string
	Blank        bool
	Enabled      bool
}

const (
	LOCAL  = "local"
	REMOTE = "remote"
)

type DocType string

func (t DocType) IsLocal() bool {
	return t == LOCAL
}

func (t DocType) IsRemote() bool {
	return t == REMOTE
}
