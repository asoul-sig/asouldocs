// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package conf

import (
	"path/filepath"
	"strings"

	"github.com/flamego/flamego"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	log "unknwon.dev/clog/v2"

	"github.com/asoul-sig/asouldocs/conf"
	"github.com/asoul-sig/asouldocs/internal/osutil"
)

func init() {
	// Initialize the primary logger until logging service is up.
	err := log.NewConsole()
	if err != nil {
		panic("init console logger: " + err.Error())
	}
}

// File is the configuration object.
var File *ini.File

// Init initializes configuration from conf assets and given custom
// configuration file. If `customConf` is empty, it falls back to default
// location, i.e. "<WORK DIR>/custom".
//
// It is safe to call this function multiple times with desired `customConf`,
// but it is not concurrent safe.
//
// NOTE: The order of loading configuration sections matters as one may depend
// on another.
func Init(customConf string) (err error) {
	if customConf == "" {
		customConf = filepath.Join("custom", "conf", "app.ini")
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return errors.Wrap(err, "get absolute path")
		}
	}

	if !osutil.IsFile(customConf) {
		return errors.Errorf("no custom configuration found at %q", customConf)
	}

	data, err := conf.Files.ReadFile("app.ini")
	if err != nil {
		return errors.Wrap(err, `read default "app.ini"`)
	}

	File, err = ini.LoadSources(
		ini.LoadOptions{
			IgnoreInlineComment: true,
		},
		data, customConf,
	)
	if err != nil {
		return errors.Wrap(err, "load configuration sources")
	}
	File.NameMapper = ini.SnackCase

	if err = File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return errors.Wrap(err, "mapping default section")
	} else if err = File.Section("site").MapTo(&Site); err != nil {
		return errors.Wrap(err, "mapping [site] section")
	} else if err = File.Section("asset").MapTo(&Asset); err != nil {
		return errors.Wrap(err, "mapping [asset] section")
	} else if err = File.Section("page").MapTo(&Page); err != nil {
		return errors.Wrap(err, "mapping [page] section")
	} else if err = File.Section("i18n").MapTo(&I18n); err != nil {
		return errors.Wrap(err, "mapping [i18n] section")
	} else if err = File.Section("docs").MapTo(&Docs); err != nil {
		return errors.Wrap(err, "mapping [docs] section")
	}

	if err = File.Section("extension.plausible").MapTo(&Extension.Plausible); err != nil {
		return errors.Wrap(err, "mapping [extension.plausible] section")
	} else if err = File.Section("extension.google_analytics").MapTo(&Extension.GoogleAnalytics); err != nil {
		return errors.Wrap(err, "mapping [extension.google_analytics] section")
	} else if err = File.Section("extension.disqus").MapTo(&Extension.Disqus); err != nil {
		return errors.Wrap(err, "mapping [extension.disqus] section")
	} else if err = File.Section("extension.utterances").MapTo(&Extension.Utterances); err != nil {
		return errors.Wrap(err, "mapping [extension.utterances] section")
	}

	Page.DocsBasePath = strings.TrimRight(Page.DocsBasePath, "/")

	if App.Env == "prod" {
		flamego.SetEnv(flamego.EnvTypeProd)
	}
	return nil
}
