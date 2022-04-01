// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"fmt"
	"path"
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

	nodes map[string]*Node // Key is the Node.Path
}

// Node is a node in the documentation hierarchy.
type Node struct {
	Category  string // The category (name) of the node, for directories and single pages, categories are empty
	Path      string // The URL path
	LocalPath string // Full path with .md extension

	Content  []byte            // The content of the node
	Name     string            // The name in the given language
	Headings goldmarktoc.Items // Headings in the node

	Nodes []*Node // The list of sub-nodes
}

// Reload reloads and converts the content from local disk.
func (n *Node) Reload(baseURLPath string) error {
	pathPrefix := path.Join(baseURLPath, strings.SplitN(n.Path, "/", 2)[0])
	content, meta, headings, err := convertFile(pathPrefix, n.LocalPath)
	if err != nil {
		return err
	}
	n.Content = content
	n.Name = fmt.Sprintf("%s", meta["name"])
	n.Headings = headings
	return nil
}

const readme = "README"

// initTocs initializes documentation hierarchy for given languages in the given
// root directory.
func initTocs(root string, languages []string, baseURLPath string) (map[string]*TOC, error) {
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
			nodes:    make(map[string]*Node),
		}

		dirs := tocCfg.Section("").KeyStrings()
		toc.Nodes = make([]*Node, 0, len(dirs))
		for _, dir := range dirs {
			dirname := tocCfg.Section("").Key(dir).String()
			files := tocCfg.Section(dirname).KeyStrings()
			// Skip empty directory
			if len(files) == 0 {
				continue
			}
			fmt.Println(dirname + "/")

			dirNode := &Node{
				Path:  dirname,
				Nodes: make([]*Node, 0, len(files)-1),
			}
			toc.Nodes = append(toc.Nodes, dirNode)
			toc.nodes[dirNode.Path] = dirNode

			if tocCfg.Section(dirname).HasValue(readme) {
				localpath := filepath.Join(root, lang, dirNode.Path, readme+".md")
				if i > 0 && !osutil.IsFile(localpath) {
					continue // It is OK to have missing file for non-default language
				}

				dirNode.LocalPath = localpath
				err = dirNode.Reload(baseURLPath)
				if err != nil {
					return nil, errors.Wrapf(err, "reload node from %q", dirNode.LocalPath)
				}
			}

			for _, file := range files {
				filename := tocCfg.Section(dirname).Key(file).String()
				if filename == readme {
					continue
				}

				localpath := filepath.Join(root, lang, dirname, filename) + ".md"
				if i > 0 && !osutil.IsFile(localpath) {
					continue // It is OK to have missing file for non-default language
				}

				node := &Node{
					Category:  dirNode.Name,
					Path:      path.Join(dirname, filename),
					LocalPath: localpath,
				}
				dirNode.Nodes = append(dirNode.Nodes, node)
				toc.nodes[node.Path] = node

				err = node.Reload(baseURLPath)
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
			toc.nodes[node.Path] = node

			err = node.Reload("")
			if err != nil {
				return nil, errors.Wrapf(err, "reload node from %q", node.LocalPath)
			}
		}

		tocs[lang] = toc
	}
	return tocs, nil
}
