// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/unknwon/com"
	"gopkg.in/ini.v1"
)

type protector struct {
	lock          sync.Mutex
	HasProtection bool
	Users         map[string]string
	Resources     map[string]map[string]bool
}

var (
	Protector = &protector{
		Users:     make(map[string]string),
		Resources: make(map[string]map[string]bool),
	}
)

func reloadProtects(localRoot string) error {
	Protector.lock.Lock()
	defer Protector.lock.Unlock()

	protectPath := path.Join(localRoot, "protect.ini")
	if !com.IsFile(protectPath) {
		return nil
	}

	Protector.HasProtection = true

	cfgs, err := ini.Load(protectPath)
	if err != nil {
		return fmt.Errorf("Fail to load protect.ini: %v", err)
	}

	for _, k := range cfgs.Section("user").Keys() {
		Protector.Users[k.Name()] = strings.ToLower(k.Value())
	}

	fmt.Println("\nProtected Resources:")
	for _, k := range cfgs.Section("auth").Keys() {
		fmt.Println("➜ ", k.Name())
		Protector.Resources[k.Name()] = make(map[string]bool)
		for _, name := range k.Strings(",") {
			fmt.Println("    ✓ ", name)
			Protector.Resources[k.Name()][name] = true
		}
	}

	return nil
}
