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
	"fmt"
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
	Content []byte

	Plain    bool   // Root node without content
	FileName string // For root node only
	Nodes    []*Node
}

// GetDoc should only be called by top level toc.
func (n *Node) GetDoc(name string) (string, []byte) {
	name = strings.TrimPrefix(name, "/")

	// Returns first node whatever avaiable as default.
	if len(name) == 0 {
		if len(n.Nodes) == 0 ||
			n.Nodes[0].Plain {
			return "", nil
		}
		return n.Nodes[0].Title, n.Nodes[0].Content
	}

	infos := strings.Split(name, "/")

	// Dir node.
	if len(infos) == 1 {
		for i := range n.Nodes {
			if n.Nodes[i].Name == infos[0] {
				return n.Nodes[i].Title, n.Nodes[i].Content
			}
		}
		return "", nil
	}

	// File node.
	for i := range n.Nodes {
		if n.Nodes[i].Name == infos[0] {
			for j := range n.Nodes[i].Nodes {
				if n.Nodes[i].Nodes[j].Name == infos[1] {
					return n.Nodes[i].Nodes[j].Title, n.Nodes[i].Nodes[j].Content
				}
			}
		}
	}

	return "", nil
}

var (
	tocLocker = sync.RWMutex{}
	Toc       map[string]*Node
	LocalRoot string
)

func initToc() error {
	tocPath := path.Join(LocalRoot, "TOC.ini")
	if !com.IsFile(tocPath) {
		return fmt.Errorf("TOC not found: %s", tocPath)
	}

	// Generate Toc.
	toc, err := ini.Load(tocPath)
	if err != nil {
		return fmt.Errorf("Fail to load TOC.ini: %v", err)
	}

	Toc = make(map[string]*Node)
	for _, lang := range setting.Docs.Langs {
		rootNode := &Node{}
		dirs := toc.Section("").KeyStrings()
		rootNode.Nodes = make([]*Node, 0, len(dirs))
		for _, dir := range dirs {
			dirName := toc.Section("").Key(dir).String()
			fmt.Println(dirName + "/")
			files := toc.Section(dirName).KeyStrings()

			// Skip empty directory.
			if len(files) == 0 {
				continue
			}

			dirNode := &Node{
				Name:     dirName,
				FileName: toc.Section(dirName).Key(files[0]).String(),
				Nodes:    make([]*Node, 0, len(files)-1),
			}
			rootNode.Nodes = append(rootNode.Nodes, dirNode)

			for _, file := range files[1:] {
				fileName := toc.Section(dirName).Key(file).String()
				fmt.Println(strings.Repeat(" ", len(dirName))+"|__", fileName)

				node := &Node{
					Name: fileName,
				}
				dirNode.Nodes = append(dirNode.Nodes, node)
			}
		}

		Toc[lang] = rootNode
	}
	return nil
}

func ReloadDocs() {
	tocLocker.Lock()
	defer tocLocker.Unlock()

	if err := initToc(); err != nil {
		log.Error("init.Toc: %v", err)
		return
	}
	initDocs()
}

func init() {
	if setting.Docs.Type == "local" {
		if !com.IsDir(setting.Docs.Target) {
			log.Error("Local documentation not found: %s", setting.Docs.Target)
			return
		}
		LocalRoot = setting.Docs.Target
		ReloadDocs()
	}
}
