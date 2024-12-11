// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
)

// helper function to clean pending approval builds from the database.
func cleanupPendingApproval(c *cli.Context, settings *settings.Platform, db *database.Interface) (*internal.Metadata, error) {
	logrus.Debug("cleaning pending approval builds older than %s", c.Duration("pending-approval-timeout"))

	m := new(internal.Metadata)

	database, err := metadataDatabase(c)
	if err != nil {
		return nil, err
	}

	m.Database = database

	queue, err := metadataQueue(c)
	if err != nil {
		return nil, err
	}

	m.Queue = queue

	source, err := metadataSource(c)
	if err != nil {
		return nil, err
	}

	m.Source = source

	vela, err := metadataVela(c)
	if err != nil {
		return nil, err
	}

	m.Vela = vela

	return m, nil
}
