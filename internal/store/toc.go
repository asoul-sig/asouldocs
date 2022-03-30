// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"fmt"
	"path/filepath"
	"strings"

	goldmarktoc "github.com/abhinav/goldmark-toc"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

// TOC represents documentation hierarchy for a specific language.
type TOC struct {
	Language string  // The language of the documentation
	Nodes    []*Node // Directories of the documentation
	Pages    []*Node // Individuals pages of the documentation
}

// Node is a node in the documentation hierarchy.
type Node struct {
	Name      string            // The name in the given language
	Path      string            // The URL path
	LocalPath string            // Full path with .md extension
	Content   []byte            // The content of the node
	Headings  goldmarktoc.Items // Headings in the node

	Nodes []*Node // The list of sub-nodes
}

// initTocs initializes documentation hierarchy for given languages in the given
// root directory.
func initTocs(root string, languages []string) (map[string]*TOC, error) {
	tocPath := filepath.Join(root, "toc.ini")
	tocCfg, err := ini.Load(tocPath)
	if err != nil {
		return nil, errors.Wrapf(err, "load %q", tocPath)
	}

	tocs := make(map[string]*TOC)
	for _, lang := range languages {
		fmt.Println("***", lang, "***")

		toc := &TOC{
			Language: lang,
		}

		dirs := tocCfg.Section("").KeyStrings()
		toc.Nodes = make([]*Node, 0, len(dirs))
		for _, dir := range dirs {
			dirname := tocCfg.Section("").Key(dir).String()
			fmt.Println(dirname + "/")

			files := tocCfg.Section(dirname).KeyStrings()
			// Skip empty directory
			if len(files) == 0 {
				continue
			}

			dirNode := &Node{
				Path:  dirname,
				Nodes: make([]*Node, 0, len(files)-1),
			}
			toc.Nodes = append(toc.Nodes, dirNode)

			const readme = "README"
			if tocCfg.Section(dirname).HasValue(readme) {
				dirNode.LocalPath = filepath.Join(root, lang, dirNode.Path, readme+".md")

				content, meta, headings, err := convertFile(dirNode.LocalPath)
				if err != nil {
					return nil, errors.Wrapf(err, "convert file %q", dirNode.LocalPath)
				}
				dirNode.Name = fmt.Sprintf("%s", meta["name"])
				dirNode.Content = content
				dirNode.Headings = headings
			}

			for _, file := range files {
				filename := tocCfg.Section(dirname).Key(file).String()
				if filename == readme {
					continue
				}
				fmt.Println(strings.Repeat(" ", len(dirname))+"|__", filename)

				docpath := filepath.Join(dirname, filename)
				node := &Node{
					Path:      docpath,
					LocalPath: filepath.Join(root, lang, docpath) + ".md",
				}
				dirNode.Nodes = append(dirNode.Nodes, node)

				content, meta, headings, err := convertFile(node.LocalPath)
				if err != nil {
					return nil, errors.Wrapf(err, "convert file %q", node.LocalPath)
				}
				node.Name = fmt.Sprintf("%s", meta["name"])
				node.Content = content
				node.Headings = headings
			}
		}

		// Single pages
		pages := tocCfg.Section("pages").KeysHash()
		toc.Pages = make([]*Node, 0, len(pages))
		for _, page := range pages {
			fmt.Println(page)

			node := &Node{
				Path:      page,
				LocalPath: filepath.Join(root, lang, page) + ".md",
			}
			toc.Pages = append(toc.Pages, node)

			content, meta, headings, err := convertFile(node.LocalPath)
			if err != nil {
				return nil, errors.Wrapf(err, "read and render file %q", node.LocalPath)
			}
			node.Name = fmt.Sprintf("%s", meta["name"])
			node.Content = content
			node.Headings = headings
		}

		tocs[lang] = toc
	}
	return tocs, nil
}
