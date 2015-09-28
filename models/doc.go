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

package models

import (
	"bytes"
	"io/ioutil"
	"path"
	"strings"

	"github.com/Unknwon/com"
	"github.com/Unknwon/log"

	"github.com/Unknwon/peach/modules/setting"
)

func parseNodeName(name string, data []byte) (string, []byte) {
	data = bytes.TrimSpace(data)
	startIdx := bytes.Index(data, []byte("---"))
	if startIdx == -1 {
		return name, []byte("")
	}
	endIdx := bytes.Index(data[startIdx+1:], []byte("---")) + startIdx
	if endIdx == -1 {
		return name, []byte("")
	}

	opts := strings.Split(strings.TrimSpace(string(string(data[startIdx+3:endIdx]))), "\n")

	title := name
	for _, opt := range opts {
		infos := strings.SplitN(opt, ":", 2)
		if len(infos) != 2 {
			continue
		}

		switch strings.TrimSpace(infos[0]) {
		case "name":
			title = strings.TrimSpace(infos[1])
		}
	}

	return title, data[endIdx+4:]
}

func initLangDocs(lang string) {
	toc := Toc[lang]
	for _, dir := range toc.Nodes {
		docPath := path.Join(LocalRoot, lang, dir.Name, dir.FileName+".md")
		data, err := ioutil.ReadFile(docPath)
		if err != nil {
			log.Error("Fail to read doc file: %v", err)
			continue
		}

		dir.Title, data = parseNodeName(dir.Name, data)
		dir.Plain = len(bytes.TrimSpace(data)) == 0

		if !dir.Plain {
			dir.Content = markdown(data)
		}

		for _, file := range dir.Nodes {
			docPath := path.Join(LocalRoot, lang, dir.Name, file.Name+".md")

			if com.IsFile(docPath) {
				data, err = ioutil.ReadFile(docPath)
				if err != nil {
					log.Error("Fail to read doc file: %v", err)
					continue
				}
			} else {
				data = []byte("")
			}

			file.Title, data = parseNodeName(file.Name, data)
			file.Content = markdown(data)
		}
	}
}

func initDocs() {
	for _, lang := range setting.Docs.Langs {
		initLangDocs(lang)
	}
}
