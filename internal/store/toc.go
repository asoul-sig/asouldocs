// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"bytes"
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
	Title    string            // The title of the document in the given language
	Headings goldmarktoc.Items // Headings in the node

	Nodes    []*Node   // The list of sub-nodes
	Previous *PageLink // The previous page
	Next     *PageLink // The next page
}

// PageLink is a link to another page.
type PageLink struct {
	Title string // The title of the page
	Path  string // the path to the page
}

// Reload reloads and converts the content from local disk.
func (n *Node) Reload(baseURLPath string) error {
	pathPrefix := path.Join(baseURLPath, strings.SplitN(n.Path, "/", 2)[0])
	content, meta, headings, err := convertFile(pathPrefix, n.LocalPath)
	if err != nil {
		return err
	}
	n.Content = content
	n.Title = fmt.Sprintf("%v", meta["title"])
	n.Headings = headings

	previous, ok := meta["previous"].(map[any]any)
	if ok {
		n.Previous = &PageLink{
			Title: fmt.Sprintf("%v", previous["title"]),
			Path:  string(convertRelativeLink(pathPrefix, []byte(fmt.Sprintf("%v", previous["path"])))),
		}
	}
	next, ok := meta["next"].(map[any]any)
	if ok {
		n.Next = &PageLink{
			Title: fmt.Sprintf("%v", next["title"]),
			Path:  string(convertRelativeLink(pathPrefix, []byte(fmt.Sprintf("%v", next["path"])))),
		}
	}
	return nil
}

const readme = "README"

// initTocs initializes documentation hierarchy for given languages in the given
// root directory. The language is the key in the returned map.
func initTocs(root string, languages []string, baseURLPath string) (map[string]*TOC, error) {
	tocPath := filepath.Join(root, "toc.ini")
	tocCfg, err := ini.Load(tocPath)
	if err != nil {
		return nil, errors.Wrapf(err, "load %q", tocPath)
	}

	var tocprint bytes.Buffer
	tocs := make(map[string]*TOC)
	for i, lang := range languages {
		tocprint.WriteString(lang)
		tocprint.WriteString(":\n")

		toc := &TOC{
			Language: lang,
			nodes:    make(map[string]*Node),
		}

		var previous *Node
		setPrevious := func(n *Node) {
			defer func() {
				previous = n
			}()
			if previous == nil {
				return
			}

			if n.Previous == nil {
				n.Previous = &PageLink{
					Title: previous.Title,
					Path:  string(convertRelativeLink(baseURLPath, []byte(previous.Path))),
				}
			}
			if previous.Next == nil {
				previous.Next = &PageLink{
					Title: n.Title,
					Path:  string(convertRelativeLink(baseURLPath, []byte(n.Path))),
				}
			}
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
			tocprint.WriteString(dirname)
			tocprint.WriteString("/\n")

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

				if len(dirNode.Content) > 0 {
					setPrevious(dirNode)
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
					Category:  dirNode.Title,
					Path:      path.Join(dirname, filename),
					LocalPath: localpath,
				}
				dirNode.Nodes = append(dirNode.Nodes, node)
				toc.nodes[node.Path] = node

				err = node.Reload(baseURLPath)
				if err != nil {
					return nil, errors.Wrapf(err, "reload node from %q", node.LocalPath)
				}

				setPrevious(node)
				tocprint.WriteString(strings.Repeat(" ", len(dirname)))
				tocprint.WriteString("|__")
				tocprint.WriteString(filename)
				tocprint.WriteString("\n")
			}
		}

		// Single pages
		pages := tocCfg.Section("pages").KeysHash()
		toc.Pages = make([]*Node, 0, len(pages))
		for _, page := range pages {
			tocprint.WriteString(page)
			tocprint.WriteString("\n")

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

	fmt.Print(tocprint.String())
	return tocs, nil
}
