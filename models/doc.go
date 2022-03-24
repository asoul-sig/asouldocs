// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"os"

	"github.com/unknwon/com"
	"github.com/unknwon/log"

	"github.com/asoul-go/asouldocs/pkg/setting"
)

func initLangDocs(tocs map[string]*Toc, localRoot, lang string) {
	toc := tocs[lang]

	for _, dir := range toc.Nodes {
		if !com.IsFile(dir.FileName) {
			continue
		}

		if err := dir.ReloadContent(); err != nil {
			log.Error("Fail to load doc file: %v", err)
			continue
		}

		for _, file := range dir.Nodes {
			if !com.IsFile(file.FileName) {
				continue
			}

			if err := file.ReloadContent(); err != nil {
				log.Error("Fail to load doc file: %v", err)
				continue
			}
		}
	}

	for _, page := range toc.Pages {
		if !com.IsFile(page.FileName) {
			continue
		}

		if err := page.ReloadContent(); err != nil {
			log.Error("Fail to load doc file: %v", err)
			continue
		}
	}
}

func initDocs(tocs map[string]*Toc, localRoot string) {
	for _, lang := range setting.Docs.Langs {
		initLangDocs(tocs, localRoot, lang)
	}
}

func NewContext() {
	if com.IsExist(HTMLRoot) {
		if err := os.RemoveAll(HTMLRoot); err != nil {
			log.Fatal("Fail to clean up HTMLRoot: %v", err)
		}
	}

	if err := ReloadDocs(); err != nil {
		log.Fatal("Fail to init docs: %v", err)
	}
}
