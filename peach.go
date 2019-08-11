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

// Peach is a web server for multi-language, real-time synchronization and searchable documentation.
package main

import (
	"os"
	"runtime"

	"github.com/urfave/cli"

	"github.com/peachdocs/peach/cmd"
	"github.com/peachdocs/peach/pkg/setting"
)

const APP_VER = "0.9.8.0810"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "Peach"
	app.Usage = "Modern Documentation Server"
	app.Version = APP_VER
	app.Author = "unknwon"
	app.Email = "u@gogs.io"
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.New,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
