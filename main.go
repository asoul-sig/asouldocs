// Copyright 2015 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// ASoulDocs is a web server for multilingual, real-time synchronized and searchable documentation.
package main

import (
	"os"

	"github.com/urfave/cli"
	log "unknwon.dev/clog/v2"

	"github.com/asoul-sig/asouldocs/internal/cmd"
	"github.com/asoul-sig/asouldocs/internal/conf"
)

func init() {
	conf.App.Version = "1.0.0+dev"
}

func main() {
	app := cli.NewApp()
	app.Name = "ASoulDocs"
	app.Usage = "Ellien's documentation server"
	app.Version = conf.App.Version
	app.Commands = []cli.Command{
		cmd.Web,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("Failed to start application: %v", err)
	}
}
