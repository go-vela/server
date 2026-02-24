// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/storage"
)

func setupStorage(_ context.Context, c *cli.Command) (storage.Storage, error) {
	logrus.Debug("creating storage client from CLI configuration")

	if !c.Bool("storage.enable") {
		logrus.Debug("storage is not enabled from CLI configuration")

		return nil, nil
	}
	// storage configuration
	_setup := &storage.Setup{
		Enable:    c.Bool("storage.enable"),
		Driver:    c.String("storage.driver"),
		Endpoint:  c.String("storage.addr"),
		AccessKey: c.String("storage.access.key"),
		SecretKey: c.String("storage.secret.key"),
		Bucket:    c.String("storage.bucket.name"),
		Secure:    c.Bool("storage.use.ssl"),
	}
	// setup the storage
	//
	// https://pkg.go.dev/github.com/go-vela/server/storage#New
	return storage.New(_setup)
}
