// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/storage"
	"github.com/go-vela/server/tracing"
	"github.com/go-vela/server/version"
)

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

	cmd := cli.Command{
		Name:    "vela-server",
		Version: v.Semantic(),
		Action:  server,
	}

	// Add Core Flags
	cmd.Flags = Flags

	// Add Database Flags
	cmd.Flags = append(cmd.Flags, database.Flags...)

	// Add Queue Flags
	cmd.Flags = append(cmd.Flags, queue.Flags...)

	// Add Secret Flags
	cmd.Flags = append(cmd.Flags, secret.Flags...)

	// Add Source Flags
	cmd.Flags = append(cmd.Flags, scm.Flags...)

	// Add Tracing Flags
	cmd.Flags = append(cmd.Flags, tracing.Flags...)

	// Add S3 Flags
	app.Flags = append(app.Flags, storage.Flags...)

	if err = cmd.Run(context.Background(), os.Args); err != nil {
		logrus.Fatal(err)
	}
}
