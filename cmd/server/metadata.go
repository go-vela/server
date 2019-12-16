// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"net/url"

	"github.com/go-vela/types"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

// helper function to setup the metadata from the CLI arguments.
func setupMetadata(c *cli.Context) (*types.Metadata, error) {
	logrus.Debug("Creating metadata from CLI configuration")

	m := new(types.Metadata)

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

// helper function to capture the database metadata from the CLI arguments.
func metadataDatabase(c *cli.Context) (*types.Database, error) {
	logrus.Trace("Creating database metadata from CLI configuration")

	u, err := url.Parse(c.String("database.config"))
	if err != nil {
		return nil, err
	}

	return &types.Database{
		Driver: c.String("database.driver"),
		Host:   u.Host,
	}, nil
}

// helper function to capture the queue metadata from the CLI arguments.
func metadataQueue(c *cli.Context) (*types.Queue, error) {
	logrus.Trace("Creating queue metadata from CLI configuration")

	u, err := url.Parse(c.String("queue-config"))
	if err != nil {
		return nil, err
	}

	return &types.Queue{
		Driver: c.String("queue-driver"),
		Host:   u.Host,
	}, nil
}

// helper function to capture the source metadata from the CLI arguments.
func metadataSource(c *cli.Context) (*types.Source, error) {
	logrus.Trace("Creating source metadata from CLI configuration")

	u, err := url.Parse(c.String("source-url"))
	if err != nil {
		return nil, err
	}

	return &types.Source{
		Driver: c.String("source-driver"),
		Host:   u.Host,
	}, nil
}

// helper function to capture the Vela metadata from the CLI arguments.
func metadataVela(c *cli.Context) (*types.Vela, error) {
	logrus.Trace("Creating Vela metadata from CLI configuration")

	vela := new(types.Vela)

	if len(c.String("server-addr")) > 0 {
		vela.Address = c.String("server-addr")
	}

	if len(c.String("webui-addr")) > 0 {
		vela.WebAddress = c.String("webui-addr")
	}

	return vela, nil
}
