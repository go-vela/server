// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	config struct {
		// specifies the address to use for the Sqlite client
		Address string
		// specifies the level of compression to use for the Sqlite client
		CompressionLevel int
		// specifies the connection duration to use for the Sqlite client
		ConnectionLife time.Duration
		// specifies the maximum idle connections for the Sqlite client
		ConnectionIdle int
		// specifies the maximum open connections for the Sqlite client
		ConnectionOpen int
		// specifies the encryption key to use for the Sqlite client
		EncryptionKey string
	}

	client struct {
		config *config
		Sqlite *gorm.DB
	}
)

// New returns a Database implementation that integrates with a Sqlite instance.
//
// nolint: golint // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new Sqlite client
	c := new(client)

	// create new fields
	c.config = new(config)
	c.Sqlite = new(gorm.DB)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// create the new Sqlite database client
	//
	// https://pkg.go.dev/gorm.io/gorm#Open
	_sqlite, err := gorm.Open(sqlite.Open(c.config.Address), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// set the Sqlite database client in the Sqlite client
	c.Sqlite = _sqlite

	return c, nil
}

// NewTest returns a Database implementation that integrates with a fake Sqlite instance.
//
// This function is intended for running tests only.
//
// nolint: golint // ignore returning unexported client
func NewTest() (*client, error) {
	// create new Sqlite client
	c := new(client)

	// create new fields
	c.config = &config{
		CompressionLevel: 3,
		EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
	}
	c.Sqlite = new(gorm.DB)

	// create the new Sqlite database client
	//
	// https://pkg.go.dev/gorm.io/gorm#Open
	_sqlite, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		return nil, err
	}

	c.Sqlite = _sqlite

	return c, nil
}
