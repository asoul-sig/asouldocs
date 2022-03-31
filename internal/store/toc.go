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

	"github.com/asoul-sig/asouldocs/internal/osutil"
)

// TOC represents documentation hierarchy for a specific language.
type TOC struct {
	Language string  // The language of the documentation
	Nodes    []*Node // Directories of the documentation
	Pages    []*Node // Individuals pages of the documentation
}

// Node is a node in the documentation hierarchy.
type Node struct {
	Path      string // The URL path
	LocalPath string // Full path with .md extension

	Content  []byte            // The content of the node
	Name     string            // The name in the given language
	Headings goldmarktoc.Items // Headings in the node

	Nodes []*Node // The list of sub-nodes
}

// Reload reloads and converts the content from local disk.
func (n *Node) Reload() error {
	content, meta, headings, err := convertFile(n.LocalPath)
	if err != nil {
		return err
	}
	n.Content = content
	n.Name = fmt.Sprintf("%s", meta["name"])
	n.Headings = headings
	return nil
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
	for i, lang := range languages {
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
				localpath := filepath.Join(root, lang, dirNode.Path, readme+".md")
				if i > 0 && !osutil.IsFile(localpath) {
					continue // It is OK to have missing file for non-default language
				}

				dirNode.LocalPath = localpath
				err = dirNode.Reload()
				if err != nil {
					return nil, errors.Wrapf(err, "reload node from %q", dirNode.LocalPath)
				}
			}

			for _, file := range files {
				filename := tocCfg.Section(dirname).Key(file).String()
				if filename == readme {
					continue
				}

				docpath := filepath.Join(dirname, filename)
				localpath := filepath.Join(root, lang, docpath) + ".md"
				if i > 0 && !osutil.IsFile(localpath) {
					continue // It is OK to have missing file for non-default language
				}

				node := &Node{
					Path:      docpath,
					LocalPath: localpath,
				}
				dirNode.Nodes = append(dirNode.Nodes, node)

				err = node.Reload()
				if err != nil {
					return nil, errors.Wrapf(err, "reload node from %q", node.LocalPath)
				}
				fmt.Println(strings.Repeat(" ", len(dirname))+"|__", filename)
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

			err = node.Reload()
			if err != nil {
				return nil, errors.Wrapf(err, "reload node from %q", node.LocalPath)
			}
		}

		tocs[lang] = toc
	}
	return tocs, nil
}
