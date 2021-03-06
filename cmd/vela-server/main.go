// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-vela/server/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

// hostname stores the server host name reported by the kernel.
var hostname string

// create an init function to set the hostname for the server.
//
// https://golang.org/doc/effective_go.html#init
func init() {
	// attempt to capture the hostname for the server
	hostname, _ = os.Hostname()
	// check if a hostname is set
	if len(hostname) == 0 {
		// default the hostname to localhost
		hostname = "localhost"
	}
}

func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	app := cli.NewApp()

	// Server Information

	app.Name = "vela-server"
	app.HelpName = "vela-server"
	app.Usage = "Vela server designed for creating builds from pipelines"
	app.Copyright = "Copyright (c) 2021 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Server Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = v.Semantic()

	// Server Flags

	app.Flags = flags()

	// Server Start

	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
