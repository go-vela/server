// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/types"

	"github.com/sirupsen/logrus"
)

// Metadata creates the metadata for the Database.
func (d *Database) Metadata() *types.Database {
	// log a message indicating the metadata creation
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Trace
	logrus.Trace("creating database metadata")

	return &types.Database{
		Driver: d.Config.Driver,
		Host:   d.Url.Host,
	}
}

// Metadata creates the metadata for the Queue.
func (q *Queue) Metadata() *types.Queue {
	// log a message indicating the metadata creation
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Trace
	logrus.Trace("creating queue metadata")

	return &types.Queue{
		Driver: q.Config.Driver,
		Host:   q.Url.Host,
	}
}

// Metadata creates the metadata for the Source.
func (s *Source) Metadata() *types.Source {
	// log a message indicating the metadata creation
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Trace
	logrus.Trace("creating source metadata")

	return &types.Source{
		Driver: s.Config.Driver,
		Host:   s.Url.Host,
	}
}
