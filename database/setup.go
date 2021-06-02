// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-vela/server/database/postgres"
	"github.com/go-vela/server/database/sqlite"
	"github.com/sirupsen/logrus"
)

// Setup represents the configuration necessary for
// creating a Vela service capable of integrating
// with a configured database system.
type Setup struct {
	// Database Configuration

	// specifies the driver to use for the database client
	Driver string
	// specifies the address to use for the database client
	Address string
	// specifies the level of compression to use for the database client
	CompressionLevel int
	// specifies the connection duration to use for the database client
	ConnectionLife time.Duration
	// specifies the maximum idle connections for the database client
	ConnectionIdle int
	// specifies the maximum open connections for the database client
	ConnectionOpen int
	// specifies the encryption key to use for the database client
	EncryptionKey string
}

// Postgres creates and returns a Vela service capable of
// integrating with a Postgres database system.
func (s *Setup) Postgres() (Service, error) {
	logrus.Trace("creating postgres database client from setup")

	// create new Postgres database service
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/postgres?tab=doc#New
	return postgres.New(
		postgres.WithAddress(s.Address),
		postgres.WithCompressionLevel(s.CompressionLevel),
		postgres.WithConnectionLife(s.ConnectionLife),
		postgres.WithConnectionIdle(s.ConnectionIdle),
		postgres.WithConnectionOpen(s.ConnectionOpen),
		postgres.WithEncryptionKey(s.EncryptionKey),
	)
}

// Sqlite creates and returns a Vela service capable of
// integrating with a Sqlite database system.
func (s *Setup) Sqlite() (Service, error) {
	logrus.Trace("creating sqlite database client from setup")

	// create new Sqlite database service
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/sqlite?tab=doc#New
	return sqlite.New(
		sqlite.WithAddress(s.Address),
		sqlite.WithCompressionLevel(s.CompressionLevel),
		sqlite.WithConnectionLife(s.ConnectionLife),
		sqlite.WithConnectionIdle(s.ConnectionIdle),
		sqlite.WithConnectionOpen(s.ConnectionOpen),
		sqlite.WithEncryptionKey(s.EncryptionKey),
	)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating database setup for client")

	// verify a database driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no database driver provided")
	}

	// verify a database address was provided
	if len(s.Address) == 0 {
		return fmt.Errorf("no database address provided")
	}

	// check if the database address has a trailing slash
	if strings.HasSuffix(s.Address, "/") {
		return fmt.Errorf("database address must not have trailing slash")
	}

	// verify a database encryption key was provided
	if len(s.EncryptionKey) == 0 {
		return fmt.Errorf("no database encryption key provided")
	}

	// enforce AES-256 for the encryption key - explicitly check for 32 characters in the key
	//
	// nolint: gomnd // ignore magic number
	if len(s.EncryptionKey) != 32 {
		// nolint: lll // ignore long line length due to long error message
		return fmt.Errorf("database encryption key must have 32 characters - provided length: %d", len(s.EncryptionKey))
	}

	// setup is valid
	return nil
}
