// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/internal"
)

// helper function to setup the metadata from the CLI arguments.
func setupMetadata(c *cli.Command) (*internal.Metadata, error) {
	logrus.Debug("creating metadata from CLI configuration")

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

	storage, err := metadataStorage(c)
	if err != nil {
		return nil, err
	}

	m.Storage = storage

	return m, nil
}

// helper function to capture the database metadata from the CLI arguments.
func metadataDatabase(c *cli.Command) (*internal.Database, error) {
	logrus.Trace("creating database metadata from CLI configuration")

	u, err := url.Parse(c.String("database.addr"))
	if err != nil {
		return nil, err
	}

	return &internal.Database{
		Driver: c.String("database.driver"),
		Host:   u.Host,
	}, nil
}

// helper function to capture the queue metadata from the CLI arguments.
func metadataQueue(c *cli.Command) (*internal.Queue, error) {
	logrus.Trace("creating queue metadata from CLI configuration")

	u, err := url.Parse(c.String("queue.addr"))
	if err != nil {
		return nil, err
	}

	return &internal.Queue{
		Driver: c.String("queue.driver"),
		Host:   u.Host,
	}, nil
}

// helper function to capture the source metadata from the CLI arguments.
func metadataSource(c *cli.Command) (*internal.Source, error) {
	logrus.Trace("creating source metadata from CLI configuration")

	u, err := url.Parse(c.String("scm.addr"))
	if err != nil {
		return nil, err
	}

	return &internal.Source{
		Driver: c.String("scm.driver"),
		Host:   u.Host,
	}, nil
}

// helper function to capture the queue metadata from the CLI arguments.
func metadataStorage(c *cli.Context) (*internal.Storage, error) {
	logrus.Trace("creating storage metadata from CLI configuration")

	u, err := url.Parse(c.String("storage.addr"))
	if err != nil {
		return nil, err
	}

	return &internal.Storage{
		Driver: c.String("storage.driver"),
		Host:   u.Host,
	}, nil
}

// helper function to capture the Vela metadata from the CLI arguments.
//
//nolint:unparam // ignore unparam for now
func metadataVela(c *cli.Command) (*internal.Vela, error) {
	logrus.Trace("creating Vela metadata from CLI configuration")

	vela := new(internal.Vela)

	if len(c.String("server-addr")) > 0 {
		vela.Address = c.String("server-addr")
	}

	if len(c.String("webui-addr")) > 0 {
		vela.WebAddress = c.String("webui-addr")
	}

	if len(c.StringSlice("cors-allow-origins")) > 0 {
		vela.CorsAllowOrigins = c.StringSlice("cors-allow-origins")
	}

	if len(c.String("webui-oauth-callback")) > 0 {
		vela.WebOauthCallbackPath = c.String("webui-oauth-callback")
	}

	if c.Duration("access-token-duration").Seconds() > 0 {
		vela.AccessTokenDuration = c.Duration("access-token-duration")
	}

	if c.Duration("refresh-token-duration").Seconds() > 0 {
		vela.RefreshTokenDuration = c.Duration("refresh-token-duration")
	}

	return vela, nil
}
