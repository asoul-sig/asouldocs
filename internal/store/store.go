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

// Init initializes the documentation store from given type and target.
func Init(typ conf.DocType, target, dir string, languages []string) (map[string]*TOC, error) {
	root := filepath.Join(target, dir)
	if typ == conf.DocTypeRemote {
		// TODO: Fetch docs from remote
		return nil, errors.New("not implemented")
	}

	if !osutil.IsDir(root) {
		return nil, errors.Errorf("directory root %q does not exist", root)
	}

	tocs, err := initTocs(root, languages)
	if err != nil {
		return nil, errors.Wrap(err, "init toc")
	}
	return tocs, nil
}
