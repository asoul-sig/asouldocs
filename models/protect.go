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
