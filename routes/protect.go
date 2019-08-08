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

package routes

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/peachdocs/peach/models"
	"github.com/peachdocs/peach/pkg/context"
)

func authRequired(ctx *context.Context) {
	ctx.Resp.Header().Set("WWW-Authenticate", "Basic realm=\".\"")
	ctx.Error(401)
}

func basicAuthDecode(encoded string) (string, string, error) {
	s, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}

	auth := strings.SplitN(string(s), ":", 2)
	return auth[0], auth[1], nil
}

func encodeMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func Protect(ctx *context.Context) {
	if !models.Protector.HasProtection {
		return
	}

	// Check if resource is protected.
	allows, yes := models.Protector.Resources[strings.TrimPrefix(ctx.Req.URL.Path, "/docs/")]
	if !yes {
		return
	}

	// Check if auth is presented.
	authHead := ctx.Req.Header.Get("Authorization")
	if len(authHead) == 0 {
		authRequired(ctx)
		return
	}

	auths := strings.Fields(authHead)
	if len(auths) != 2 || auths[0] != "Basic" {
		ctx.Error(401)
		return
	}

	uname, passwd, err := basicAuthDecode(auths[1])
	if err != nil {
		ctx.Error(401)
		return
	}

	// Check if auth is valid.
	if models.Protector.Users[uname] != encodeMd5(passwd) {
		ctx.Error(401)
		return
	}

	// Check if user has access to the resource.
	if !allows[uname] {
		ctx.Error(403)
		return
	}
}
