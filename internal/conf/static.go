// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package conf

var (
	// Application settings
	App struct {
		// ⚠️ WARNING: Should only be set by the main package (i.e. "main.go").
		Version string `ini:"-"`

		Env      string
		HTTPHost string `ini:"HTTP_ADDR"`
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
		DocsBasePath    string
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
		Languages       []string
		Names           []string
		CustomDirectory string
	}

	// Documentation settings
	Docs struct {
		Type      DocType
		Target    string
		TargetDir string
	}
)

type DocType string

const (
	DocTypeLocal  DocType = "local"
	DocTypeRemote DocType = "remote"
)
