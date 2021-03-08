// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// Setup represents the configuration necessary for
// creating a Vela service capable of integrating
// with a configured database.
type Setup struct {
	// Database Configuration

	// specifies the database driver to use
	Driver string
	// specifies the database address to use
	Address string
	// specifies the level of compression for logs stored in the database
	CompressionLevel int
	// specifies the number of idle connections to the database
	ConnectionIdle int
	// specifies the amount of time a connection may be reused for the database
	ConnectionLife time.Duration
	// specifies the number of open connections to the database
	ConnectionOpen int
	// specifies the key for encrypting and decrypting data using AES-256
	EncryptionKey string
}

// Linux creates and returns a Vela service capable of
// integrating with a Postgres database.
func (s *Setup) Postgres() (Service, error) {
	logrus.Trace("creating postgres database client from setup")

	// create new Postgres database service
	//
	// https://pkg.go.dev/github.com/go-vela/pkg-executor/executor/linux?tab=doc#New
	return New(s)
}

// Sqlite creates and returns a Vela service capable of
// integrating with a Sqlite database.
func (s *Setup) Sqlite() (Service, error) {
	logrus.Trace("creating sqlite database client from setup")

	// create new Sqlite database service
	//
	// https://pkg.go.dev/github.com/go-vela/pkg-executor/executor/local?tab=doc#New
	return New(s)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating executor setup for client")

	// verify a database driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no database driver provided")
	}

	// verify the database driver provided is valid
	switch s.Driver {
	case constants.DriverPostgres, "postgresql":
		fallthrough
	case constants.DriverSqlite, "sqlite":
		break
	default:
		return fmt.Errorf("invalid database driver provided: %s", s.Driver)
	}

	// verify a database address was provided
	if len(s.Address) == 0 {
		return fmt.Errorf("no database address provided")
	}

	// check if the database address has a scheme
	if !strings.Contains(s.Address, "://") {
		return fmt.Errorf("database address must be fully qualified (<scheme>://<host>)")
	}

	// check if the database address has a trailing slash
	if strings.HasSuffix(s.Address, "/") {
		return fmt.Errorf("database address must not have trailing slash")
	}

	// verify the compression level provided is valid
	switch s.CompressionLevel {
	case constants.CompressionNegOne:
		fallthrough
	case constants.CompressionZero:
		fallthrough
	case constants.CompressionOne:
		fallthrough
	case constants.CompressionTwo:
		fallthrough
	case constants.CompressionThree:
		fallthrough
	case constants.CompressionFour:
		fallthrough
	case constants.CompressionFive:
		fallthrough
	case constants.CompressionSix:
		fallthrough
	case constants.CompressionSeven:
		fallthrough
	case constants.CompressionEight:
		fallthrough
	case constants.CompressionNine:
		break
	default:
		return fmt.Errorf("invalid database compression level provided: %d", s.CompressionLevel)
	}

	// enforce AES-256, so check explicitly for 32 bytes on the key
	//
	// nolint: gomnd // ignore magic number
	if len(s.EncryptionKey) != 32 {
		return fmt.Errorf("invalid database encryption key provided: %d", len(s.EncryptionKey))
	}

	// setup is valid
	return nil
}
