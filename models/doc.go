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

func initLangDocs(localRoot, lang string) {
	toc := Tocs[lang]
	for _, dir := range toc.Nodes {
		data, err := ioutil.ReadFile(dir.FileName)
		if err != nil {
			log.Error("Fail to read doc file: %v", err)
			continue
		}

		dir.Title, data = parseNodeName(dir.Name, data)
		dir.Plain = len(bytes.TrimSpace(data)) == 0

		if !dir.Plain {
			dir.content = markdown(data)
		}

		for _, file := range dir.Nodes {
			if com.IsFile(file.FileName) {
				data, err = ioutil.ReadFile(file.FileName)
				if err != nil {
					log.Error("Fail to read doc file: %v", err)
					continue
				}
			} else {
				data = []byte("")
			}

			file.Title, data = parseNodeName(file.Name, data)
			file.content = markdown(data)
		}
	}
}

func initDocs(localRoot string) {
	for _, lang := range setting.Docs.Langs {
		initLangDocs(localRoot, lang)
	}
}
