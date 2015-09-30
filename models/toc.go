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
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"sync"

	"github.com/Unknwon/com"
	"github.com/Unknwon/log"
	"gopkg.in/ini.v1"

	"github.com/Unknwon/peach/modules/setting"
)

type Node struct {
	Name    string // Name in TOC
	Title   string // Name in given language
	content []byte

	Plain    bool   // Root node without content
	FileName string // Full path with .md extension
	Nodes    []*Node
}

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

func (n *Node) ReloadContent() error {
	data, err := ioutil.ReadFile(n.FileName)
	if err != nil {
		return err
	}

	n.Title, data = parseNodeName(n.Name, data)
	n.Plain = len(bytes.TrimSpace(data)) == 0

	if !n.Plain {
		n.content = markdown(data)
	}
	return nil
}

func (n *Node) Content() []byte {
	if !setting.ProdMode {
		if err := n.ReloadContent(); err != nil {
			log.Error("Fail to reload content: %v", err)
		}
	}

	return n.content
}

// Toc represents table of content in a specific language.
type Toc struct {
	RootPath string
	Nodes    []*Node
}

// GetDoc should only be called by top level toc.
func (t *Toc) GetDoc(name string) (string, []byte) {
	name = strings.TrimPrefix(name, "/")

	// Returns first node whatever avaiable as default.
	if len(name) == 0 {
		if len(t.Nodes) == 0 ||
			t.Nodes[0].Plain {
			return "", nil
		}
		return t.Nodes[0].Title, t.Nodes[0].Content()
	}

	infos := strings.Split(name, "/")

	// Dir node.
	if len(infos) == 1 {
		for i := range t.Nodes {
			if t.Nodes[i].Name == infos[0] {
				return t.Nodes[i].Title, t.Nodes[i].Content()
			}
		}
		return "", nil
	}

	// File node.
	for i := range t.Nodes {
		if t.Nodes[i].Name == infos[0] {
			for j := range t.Nodes[i].Nodes {
				if t.Nodes[i].Nodes[j].Name == infos[1] {
					return t.Nodes[i].Nodes[j].Title, t.Nodes[i].Nodes[j].Content()
				}
			}
		}
	}

	return "", nil
}

var (
	tocLocker = sync.RWMutex{}
	Tocs      map[string]*Toc
)

func initToc(localRoot string) error {
	tocPath := path.Join(localRoot, "TOC.ini")
	if !com.IsFile(tocPath) {
		return fmt.Errorf("TOC not found: %s", tocPath)
	}

	// Generate Toc.
	tocCfg, err := ini.Load(tocPath)
	if err != nil {
		return fmt.Errorf("Fail to load TOC.ini: %v", err)
	}

	Tocs = make(map[string]*Toc)
	for _, lang := range setting.Docs.Langs {
		toc := &Toc{
			RootPath: localRoot,
		}
		dirs := tocCfg.Section("").KeyStrings()
		toc.Nodes = make([]*Node, 0, len(dirs))
		for _, dir := range dirs {
			dirName := tocCfg.Section("").Key(dir).String()
			fmt.Println(dirName + "/")
			files := tocCfg.Section(dirName).KeyStrings()

			// Skip empty directory.
			if len(files) == 0 {
				continue
			}

			dirNode := &Node{
				Name:     dirName,
				FileName: path.Join(localRoot, lang, dirName, tocCfg.Section(dirName).Key(files[0]).String()) + ".md",
				Nodes:    make([]*Node, 0, len(files)-1),
			}
			toc.Nodes = append(toc.Nodes, dirNode)

			for _, file := range files[1:] {
				fileName := tocCfg.Section(dirName).Key(file).String()
				fmt.Println(strings.Repeat(" ", len(dirName))+"|__", fileName)

				node := &Node{
					Name:     fileName,
					FileName: path.Join(localRoot, lang, dirName, fileName) + ".md",
				}
				dirNode.Nodes = append(dirNode.Nodes, node)
			}
		}

		Tocs[lang] = toc
	}
	return nil
}

func ReloadDocs(localRoot string) {
	tocLocker.Lock()
	defer tocLocker.Unlock()

	if err := initToc(localRoot); err != nil {
		log.Error("init.Toc: %v", err)
		return
	}
	initDocs(localRoot)
}

func init() {
	if setting.Docs.Type == "local" {
		if !com.IsDir(setting.Docs.Target) {
			log.Fatal("Local documentation not found: %s", setting.Docs.Target)
			return
		}
		ReloadDocs(setting.Docs.Target)
	}
}
