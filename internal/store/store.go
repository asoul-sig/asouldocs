// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/asoul-sig/asouldocs/internal/conf"
	"github.com/asoul-sig/asouldocs/internal/osutil"
)

// Store is a store maintaining documentation hierarchies for multiple
// languages.
type Store struct {
	defaultLanguage string
	tocs            map[string]*TOC // Key is the language
}

// FirstDocPath returns the URL path of the first doc that has content in the
// default language.
func (s *Store) FirstDocPath() string {
	for _, dir := range s.tocs[s.defaultLanguage].Nodes {
		if len(dir.Content) > 0 {
			return dir.Path
		}

		for _, file := range dir.Nodes {
			return file.Path
		}
	}
	return "404"
}

// TOC returns the TOC of the given language. It returns the TOC of the default
// language if the given language is not found.
func (s *Store) TOC(language string) *TOC {
	toc, ok := s.tocs[language]
	if !ok {
		return s.tocs[s.defaultLanguage]
	}
	return toc
}

var ErrNoMatch = errors.New("no match for the path")

// Match matches a node with given path in given language. If the no such node
// exists, it fallbacks to use the node with same path in default language.
func (s *Store) Match(language, path string) (n *Node, fallback bool, err error) {
	toc := s.TOC(language)
	n, ok := toc.nodes[path]
	if ok && len(n.Content) > 0 {
		return n, false, nil
	}

	if toc.Language == s.defaultLanguage {
		return nil, false, ErrNoMatch
	}

	n, ok = s.tocs[s.defaultLanguage].nodes[path]
	if ok && len(n.Content) > 0 {
		return n, true, nil
	}
	return nil, false, ErrNoMatch
}

// Init initializes the documentation store from given type and target.
func Init(typ conf.DocType, target, dir string, languages []string, baseURLPath string) (*Store, error) {
	if len(languages) < 1 {
		return nil, errors.New("no languages")
	}

	root := filepath.Join(target, dir)
	if typ == conf.DocTypeRemote {
		// TODO: Fetch docs from remote
		return nil, errors.New("not implemented")
	}

	if !osutil.IsDir(root) {
		return nil, errors.Errorf("directory root %q does not exist", root)
	}

	tocs, err := initTocs(root, languages, baseURLPath)
	if err != nil {
		return nil, errors.Wrap(err, "init toc")
	}
	return &Store{
		defaultLanguage: languages[0],
		tocs:            tocs,
	}, nil
}
