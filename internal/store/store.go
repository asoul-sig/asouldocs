// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/gogs/git-module"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/asoul-sig/asouldocs/internal/conf"
	"github.com/asoul-sig/asouldocs/internal/osutil"
)

// Store is a store maintaining documentation hierarchies for multiple
// languages.
type Store struct {
	// The list of config values
	typ         conf.DocType
	target      string
	targetDir   string
	languages   []string
	baseURLPath string

	// The list of inferred values
	rootDir         string
	defaultLanguage string
	tocs            atomic.Value
	reloadLock      sync.Mutex
}

// RootDir returns the root directory of documentation hierarchies.
func (s *Store) RootDir() string {
	return s.rootDir
}

func (s *Store) getTOCs() map[string]*TOC {
	return s.tocs.Load().(map[string]*TOC)
}

func (s *Store) setTOCs(tocs map[string]*TOC) {
	s.tocs.Store(tocs)
}

// FirstDocPath returns the URL path of the first doc that has content in the
// default language.
func (s *Store) FirstDocPath() string {
	for _, dir := range s.getTOCs()[s.defaultLanguage].Nodes {
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
	toc, ok := s.getTOCs()[language]
	if !ok {
		return s.getTOCs()[s.defaultLanguage]
	}
	return toc
}

var ErrNoMatch = errors.New("no match for the path")

// Match matches a node with given path in given language. If the no such node
// exists or the node content is empty, it fallbacks to use the node with same
// path in default language.
func (s *Store) Match(language, path string) (n *Node, fallback bool, err error) {
	toc := s.TOC(language)
	n, ok := toc.nodes[path]
	if ok && len(n.Content) > 0 {
		return n, false, nil
	}

	if toc.Language == s.defaultLanguage {
		return nil, false, ErrNoMatch
	}

	n, ok = s.getTOCs()[s.defaultLanguage].nodes[path]
	if ok && len(n.Content) > 0 {
		return n, true, nil
	}
	return nil, false, ErrNoMatch
}

// Reload re-initializes the documentation store.
func (s *Store) Reload() error {
	s.reloadLock.Lock()
	defer s.reloadLock.Unlock()

	log.Trace("Reloading %s...", s.target)

	root := filepath.Join(s.target, s.targetDir)
	if s.typ == conf.DocTypeRemote {
		localCache := filepath.Join("data", "docs")
		if !osutil.IsExist(localCache) {
			log.Trace("Cloning %s...", s.target)
			err := git.Clone(s.target, localCache, git.CloneOptions{Depth: 1})
			if err != nil {
				return errors.Wrapf(err, "clone %q", s.target)
			}
		} else {
			repo, err := git.Open(localCache)
			if err != nil {
				return errors.Wrapf(err, "open %q", localCache)
			}

			log.Trace("Pulling %s...", s.target)
			err = repo.Pull()
			if err != nil {
				return errors.Wrapf(err, "pull %q", s.target)
			}
		}

		root = filepath.Join(localCache, s.targetDir)
	}

	if !osutil.IsDir(root) {
		return errors.Errorf("directory root %q does not exist", root)
	}

	tocs, err := initTocs(root, s.languages, s.baseURLPath)
	if err != nil {
		return errors.Wrap(err, "init toc")
	}

	s.rootDir = root
	s.setTOCs(tocs)
	return nil
}

// Init initializes the documentation store from given type and target.
func Init(typ conf.DocType, target, targetDir string, languages []string, baseURLPath string) (*Store, error) {
	if len(languages) < 1 {
		return nil, errors.New("no languages")
	}

	s := &Store{
		typ:             typ,
		target:          target,
		targetDir:       targetDir,
		languages:       languages,
		baseURLPath:     baseURLPath,
		defaultLanguage: languages[0],
	}
	err := s.Reload()
	if err != nil {
		return nil, errors.Wrap(err, "reload")
	}
	return s, nil
}
