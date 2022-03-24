// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// ASoulDocs is a web server for multilingual, real-time synchronized and searchable documentation.
package main

import (
	"os"
	"runtime"

	"github.com/urfave/cli"

	"github.com/asoul-sig/asouldocs/cmd"
	"github.com/asoul-sig/asouldocs/pkg/setting"
)

const APP_VER = "0.9.8.0810"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "ASoulDocs"
	app.Usage = "Ellien's documentation server"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.New,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	_ = app.Run(os.Args)
}
