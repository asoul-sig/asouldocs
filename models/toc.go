// Copyright 2015 unknwon
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
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/unknwon/com"
	"github.com/mschoch/blackfriday-text"
	"github.com/russross/blackfriday"
	"gopkg.in/ini.v1"

	"github.com/peachdocs/peach/pkg/setting"
)

type Node struct {
	Name  string // Name in TOC
	Title string // Name in given language
	text  []byte // Clean text without formatting
	runes []rune

	Plain         bool // Root node without content
	DocumentPath  string
	FileName      string // Full path with .md extension
	Nodes         []*Node
	LastBuildTime int64
}

func (n *Node) SetText(text []byte) {
	n.text = text
	n.runes = []rune(string(n.text))
}

func (n *Node) Text() []byte {
	return n.text
}

var textRender = blackfridaytext.TextRenderer()
var (
	docsRoot = "data/docs"
	HTMLRoot = "data/html"
)

func parseNodeName(name string, data []byte) (string, []byte) {
	data = bytes.TrimSpace(data)
	if len(data) < 3 || string(data[:3]) != "---" {
		return name, []byte("")
	}
	endIdx := bytes.Index(data[3:], []byte("---")) + 3
	if endIdx == -1 {
		return name, []byte("")
	}

	opts := strings.Split(strings.TrimSpace(string(string(data[3:endIdx]))), "\n")

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

	return title, data[endIdx+3:]
}

func (n *Node) ReloadContent() error {
	data, err := ioutil.ReadFile(n.FileName)
	if err != nil {
		return err
	}

	n.Title, data = parseNodeName(n.Name, data)
	n.Plain = len(bytes.TrimSpace(data)) == 0

	if !n.Plain {
		n.SetText(bytes.ToLower(blackfriday.Markdown(data, textRender, 0)))
		data = markdown(data)
	}

	return n.GenHTML(data)
}

// HTML2JS converts []byte type of HTML content into JS format.
func HTML2JS(data []byte) []byte {
	s := string(data)
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\"", `\"`, -1)
	return []byte(s)
}

func (n *Node) GenHTML(data []byte) error {
	var htmlPath string
	if setting.Docs.Type.IsLocal() {
		htmlPath = path.Join(HTMLRoot, strings.TrimPrefix(n.FileName, setting.Docs.Target))
	} else {
		htmlPath = path.Join(HTMLRoot, strings.TrimPrefix(n.FileName, filepath.Join(docsRoot, setting.Docs.TargetDir)))
	}
	htmlPath = strings.Replace(htmlPath, ".md", ".js", 1)

	buf := new(bytes.Buffer)
	buf.WriteString("document.write(\"")
	buf.Write(HTML2JS(data))
	buf.WriteString("\")")

	n.LastBuildTime = time.Now().Unix()
	if err := com.WriteFile(htmlPath+".tmp", buf.Bytes()); err != nil {
		return err
	}
	os.Remove(htmlPath)
	return os.Rename(htmlPath+".tmp", htmlPath)
}

// Toc represents table of content in a specific language.
type Toc struct {
	RootPath string
	Lang     string
	Nodes    []*Node
	Pages    []*Node
}

// GetDoc should only be called by top level toc.
func (t *Toc) GetDoc(name string) (*Node, bool) {
	name = strings.TrimPrefix(name, "/")

	// Returns first node whatever avaiable as default.
	if len(name) == 0 {
		if len(t.Nodes) == 0 ||
			t.Nodes[0].Plain {
			return nil, false
		}
		return t.Nodes[0], false
	}

	infos := strings.Split(name, "/")

	// Dir node.
	if len(infos) == 1 {
		for i := range t.Nodes {
			if t.Nodes[i].Name == infos[0] {
				return t.Nodes[i], false
			}
		}
		return nil, false
	}

	// File node.
	for i := range t.Nodes {
		if t.Nodes[i].Name == infos[0] {
			for j := range t.Nodes[i].Nodes {
				if t.Nodes[i].Nodes[j].Name == infos[1] {
					if com.IsFile(t.Nodes[i].Nodes[j].FileName) {
						return t.Nodes[i].Nodes[j], false
					}

					// If not default language, try again.
					n, _ := Tocs[setting.Docs.Langs[0]].GetDoc(name)
					return n, true
				}
			}
		}
	}

	return nil, false
}

type SearchResult struct {
	Title string
	Path  string
	Match string
}

func (n *Node) adjustRange(start int) (int, int) {
	start -= 20
	if start < 0 {
		start = 0
	}

	length := len(n.runes)
	end := start + 230
	if end > length {
		end = length
	}
	return start, end
}

func (t *Toc) Search(q string) []*SearchResult {
	if len(q) == 0 {
		return nil
	}
	q = strings.ToLower(q)

	results := make([]*SearchResult, 0, 5)

	// Dir node.
	for i := range t.Nodes {
		if idx := bytes.Index(t.Nodes[i].Text(), []byte(q)); idx > -1 {
			start, end := t.Nodes[i].adjustRange(utf8.RuneCount(t.Nodes[i].Text()[:idx]))
			results = append(results, &SearchResult{
				Title: t.Nodes[i].Title,
				Path:  t.Nodes[i].Name,
				Match: string(t.Nodes[i].runes[start:end]),
			})
		}
	}

	// File node.
	for i := range t.Nodes {
		for j := range t.Nodes[i].Nodes {
			if idx := bytes.Index(t.Nodes[i].Nodes[j].Text(), []byte(q)); idx > -1 {
				start, end := t.Nodes[i].Nodes[j].adjustRange(utf8.RuneCount(t.Nodes[i].Nodes[j].Text()[:idx]))
				results = append(results, &SearchResult{
					Title: t.Nodes[i].Nodes[j].Title,
					Path:  path.Join(t.Nodes[i].Name, t.Nodes[i].Nodes[j].Name),
					Match: string(t.Nodes[i].Nodes[j].runes[start:end]),
				})
			}
		}
	}

	return results
}

var (
	tocLocker = sync.Mutex{}
	Tocs      map[string]*Toc
)

func initToc(localRoot string) (map[string]*Toc, error) {
	tocPath := path.Join(localRoot, "TOC.ini")
	if !com.IsFile(tocPath) {
		return nil, fmt.Errorf("TOC not found: %s", tocPath)
	}

	// Generate Toc.
	tocCfg, err := ini.Load(tocPath)
	if err != nil {
		return nil, fmt.Errorf("Fail to load TOC.ini: %v", err)
	}

	tocs := make(map[string]*Toc)
	for _, lang := range setting.Docs.Langs {
		toc := &Toc{
			RootPath: localRoot,
			Lang:     lang,
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

			documentPath := path.Join(dirName, tocCfg.Section(dirName).Key(files[0]).String())
			dirNode := &Node{
				Name:         dirName,
				DocumentPath: documentPath,
				FileName:     path.Join(localRoot, lang, documentPath) + ".md",
				Nodes:        make([]*Node, 0, len(files)-1),
			}
			toc.Nodes = append(toc.Nodes, dirNode)

			for _, file := range files[1:] {
				fileName := tocCfg.Section(dirName).Key(file).String()
				fmt.Println(strings.Repeat(" ", len(dirName))+"|__", fileName)

				documentPath = path.Join(dirName, fileName)
				node := &Node{
					Name:         fileName,
					DocumentPath: documentPath,
					FileName:     path.Join(localRoot, lang, documentPath) + ".md",
				}
				dirNode.Nodes = append(dirNode.Nodes, node)
			}
		}

		// Single pages.
		pages := tocCfg.Section("pages").KeyStrings()
		toc.Pages = make([]*Node, 0, len(pages))
		for _, page := range pages {
			pageName := tocCfg.Section("pages").Key(page).String()
			fmt.Println(pageName)

			toc.Pages = append(toc.Pages, &Node{
				Name:         pageName,
				DocumentPath: pageName,
				FileName:     path.Join(localRoot, lang, pageName) + ".md",
			})
		}

		tocs[lang] = toc
	}
	return tocs, nil
}

func ReloadDocs() error {
	tocLocker.Lock()
	defer tocLocker.Unlock()

	localRoot := setting.Docs.Target

	// Fetch docs from remote.
	if setting.Docs.Type.IsRemote() {
		localRoot = docsRoot

		absRoot, err := filepath.Abs(localRoot)
		if err != nil {
			return fmt.Errorf("filepath.Abs: %v", err)
		}

		// Clone new or pull to update.
		if com.IsDir(absRoot) {
			stdout, stderr, err := com.ExecCmdDir(absRoot, "git", "pull")
			if err != nil {
				return fmt.Errorf("Fail to update docs from remote source(%s): %v - %s", setting.Docs.Target, err, stderr)
			}
			fmt.Println(stdout)
		} else {
			os.MkdirAll(filepath.Dir(absRoot), os.ModePerm)
			stdout, stderr, err := com.ExecCmd("git", "clone", setting.Docs.Target, absRoot)
			if err != nil {
				return fmt.Errorf("Fail to clone docs from remote source(%s): %v - %s", setting.Docs.Target, err, stderr)
			}
			fmt.Println(stdout)
		}

		// Append subdir to root as needed
		localRoot = path.Join(localRoot, setting.Docs.TargetDir)
	}

	if !com.IsDir(localRoot) {
		return fmt.Errorf("Documentation not found: %s - %s", setting.Docs.Type, localRoot)
	}

	tocs, err := initToc(localRoot)
	if err != nil {
		return fmt.Errorf("initToc: %v", err)
	}
	initDocs(tocs, localRoot)
	Tocs = tocs
	return reloadProtects(localRoot)
}
