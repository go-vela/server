// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"fmt"
	"strings"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// Setup represents the configuration necessary for
// creating a Vela service capable of integrating
// with a configured secret system.
type Setup struct {
	// Secret Configuration

	// specifies the driver to use for the secret client
	Driver string

	// specifies the database service to use for the secret client
	Database database.Service

	// specifies the address to use for the secret client
	Address string
	// specifies the authentication method to use for the secret client
	AuthMethod string
	// specifies the AWS role to use for the secret client
	AwsRole string
	// specifies the prefix to use for the secret client
	Prefix string
	// specifies the prefix to use for the secret client
	Prefix string
	// specifies the token to use for the secret client
	Token string
	// specifies the version to use for the secret client
	Version string
}

// Native creates and returns a Vela service capable of
// integrating with a Native (Database) secret system.
func (s *Setup) Native() (Service, error) {
	logrus.Trace("creating native secret client from setup")

	return nil, fmt.Errorf("unsupported secret driver: %s", constants.DriverNative)
}

// Vault creates and returns a Vela service capable of
// integrating with a Hashicorp Vault secret system.
func (s *Setup) Vault() (Service, error) {
	logrus.Trace("creating vault secret client from setup")

	return nil, fmt.Errorf("unsupported secret driver: %s", constants.DriverVault)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating secret setup for client")

	// verify a secret driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no secret driver provided")
	}

	// process the secret driver being provided
	switch s.Driver {
	case constants.DriverNative:
		// verify a secret database was provided
		if s.Database == nil {
			return fmt.Errorf("no secret database service provided")
		}
	case constants.DriverVault:
		fallthrough
	default:
		// verify a secret address was provided
		if len(s.Address) == 0 {
			return fmt.Errorf("no secret address provided")
		}

		// check if the secret address has a scheme
		if !strings.Contains(s.Address, "://") {
			return fmt.Errorf("secret address must be fully qualified (<scheme>://<host>)")
		}

		// check if the secret address has a trailing slash
		if strings.HasSuffix(s.Address, "/") {
			return fmt.Errorf("secret address must not have trailing slash")
		}
	}

	// setup is valid
	return nil
}
