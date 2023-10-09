// SPDX-License-Identifier: Apache-2.0

package database

import (
	"fmt"
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// Validate verifies the required fields from the provided configuration are populated correctly.
func (c *config) Validate() error {
	logrus.Trace("validating database configuration for engine")

	// verify a database driver was provided
	if len(c.Driver) == 0 {
		return fmt.Errorf("no database driver provided")
	}

	// verify a database address was provided
	if len(c.Address) == 0 {
		return fmt.Errorf("no database address provided")
	}

	// check if the database address has a trailing slash
	if strings.HasSuffix(c.Address, "/") {
		return fmt.Errorf("invalid database address provided: address must not have trailing slash")
	}

	// verify a database encryption key was provided
	if len(c.EncryptionKey) == 0 {
		return fmt.Errorf("no database encryption key provided")
	}

	// check the database encryption key length - enforce AES-256 by forcing 32 characters in the key
	if len(c.EncryptionKey) != 32 {
		return fmt.Errorf("invalid database encryption key provided: key length (%d) must be 32 characters", len(c.EncryptionKey))
	}

	// verify the database compression level is valid
	switch c.CompressionLevel {
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
		return fmt.Errorf("invalid database compression level provided: level (%d) must be between %d and %d",
			c.CompressionLevel, constants.CompressionNegOne, constants.CompressionNine,
		)
	}

	return nil
}
