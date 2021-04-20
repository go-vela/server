// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/secret/native"
	"github.com/go-vela/server/secret/vault"
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
	// specifies the token to use for the secret client
	Token string
	// specifies the token duration to use for the secret client
	TokenDuration time.Duration
	// specifies the version to use for the secret client
	Version string
}

// Native creates and returns a Vela service capable of
// integrating with a Native (Database) secret system.
func (s *Setup) Native() (Service, error) {
	logrus.Trace("creating native secret client from setup")

	// create new native secret service
	//
	// https://pkg.go.dev/github.com/go-vela/server/secret/native?tab=doc#New
	return native.New(
		native.WithDatabase(s.Database),
	)
}

// Vault creates and returns a Vela service capable of
// integrating with a Hashicorp Vault secret system.
func (s *Setup) Vault() (Service, error) {
	logrus.Trace("creating vault secret client from setup")

	// create new Vault secret service
	//
	// https://pkg.go.dev/github.com/go-vela/server/secret/vault?tab=doc#New
	return vault.New(
		vault.WithAddress(s.Address),
		vault.WithAuthMethod(s.AuthMethod),
		vault.WithAWSRole(s.AwsRole),
		vault.WithPrefix(s.Prefix),
		vault.WithToken(s.Token),
		vault.WithTokenDuration(s.TokenDuration),
		vault.WithVersion(s.Version),
	)
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

		// verify a secret token or authentication method was provided
		if len(s.Token) == 0 && len(s.AuthMethod) == 0 {
			return fmt.Errorf("no secret token or authentication method provided")
		}

		// check if the secret token is empty
		if len(s.Token) == 0 {
			// process the secret authentication method being provided
			switch s.AuthMethod {
			case "aws":
				// verify a secret AWS role was provided
				if len(s.AwsRole) == 0 {
					return fmt.Errorf("no secret AWS role provided")
				}
			default:
				return fmt.Errorf("invalid secret authentication method provided: %s", s.AuthMethod)
			}
		}
	}

	// setup is valid
	return nil
}
