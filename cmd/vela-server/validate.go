// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Validate verifies the API is properly configured.
func (a *API) Validate() error {
	// verify a server address was provided
	if len(a.Address) == 0 {
		return fmt.Errorf("no server address provided")
	}

	// check if the server address has a scheme
	if !strings.Contains(a.Address, "://") {
		return fmt.Errorf("server address must be fully qualified (<scheme>://<host>)")
	}

	// check if the server address has a trailing slash
	if strings.HasSuffix(a.Address, "/") {
		return fmt.Errorf("server address must not have trailing slash")
	}

	// parse the server address into a url structure
	parsed, err := url.Parse(a.Address)
	if err != nil {
		return fmt.Errorf("unable to parse server address: %v", err)
	}

	// save the parsed server address
	a.Url = parsed

	// verify a server secret was provided
	if len(a.Secret) == 0 {
		return fmt.Errorf("no server secret provided")
	}

	return nil
}

// Validate verifies the Build is properly configured.
func (b *Build) Validate() error {
	// verify a build timeout was provided
	if b.Timeout <= 0 {
		return fmt.Errorf("no build timeout provided")
	}

	return nil
}

// Validate verifies the Compiler is properly configured.
func (c *Compiler) Validate() error {
	logrus.Trace("validating compiler configuration")

	// check if the github driver is enabled
	if c.Github.Driver {
		// verify a github address was provided
		if len(c.Github.Address) == 0 {
			return fmt.Errorf("no compiler GitHub address provided")
		}

		// check if the compiler GitHub address has a scheme
		if !strings.Contains(c.Github.Address, "://") {
			return fmt.Errorf("compiler GitHub address must be fully qualified (<scheme>://<host>)")
		}

		// check if the compiler GitHub address has a trailing slash
		if strings.HasSuffix(c.Github.Address, "/") {
			return fmt.Errorf("compiler GitHub address must not have trailing slash")
		}

		// parse the compiler GitHub address into a url structure
		parsed, err := url.Parse(c.Github.Address)
		if err != nil {
			return fmt.Errorf("unable to parse compiler GitHub address: %v", err)
		}

		// save the parsed compiler GitHub address
		c.Github.Url = parsed

		// verify a github token was provided
		if len(c.Github.Token) == 0 {
			return fmt.Errorf("no compiler GitHub token provided")
		}
	}

	// check if an address for the modification service is provided
	if len(c.Modification.Address) > 0 {
		// check if the compiler modification address has a scheme
		if !strings.Contains(c.Modification.Address, "://") {
			return fmt.Errorf("compiler modification address must be fully qualified (<scheme>://<host>)")
		}

		// check if the compiler modification address has a trailing slash
		if strings.HasSuffix(c.Modification.Address, "/") {
			return fmt.Errorf("compiler modification address must not have trailing slash")
		}

		// parse the compiler modification address into a url structure
		parsed, err := url.Parse(c.Modification.Address)
		if err != nil {
			return fmt.Errorf("unable to parse compiler modification address: %v", err)
		}

		// save the parsed compiler modification address
		c.Modification.Url = parsed

		// verify a secret for the modification service was provided
		if len(c.Modification.Secret) == 0 {
			return fmt.Errorf("no compiler modification secret provided")
		}
	}

	return nil
}

// Validate verifies the Database is properly configured.
func (d *Database) Validate() error {
	logrus.Trace("validating database configuration")

	// verify a database driver was provided
	if len(d.Driver) == 0 {
		return fmt.Errorf("no database driver provided")
	}

	// verify a database address was provided
	if len(d.Address) == 0 {
		return fmt.Errorf("no database address provided")
	}

	// check if the database address has a scheme
	if !strings.Contains(d.Address, "://") {
		return fmt.Errorf("database address must be fully qualified (<scheme>://<host>)")
	}

	// check if the database address has a trailing slash
	if strings.HasSuffix(d.Address, "/") {
		return fmt.Errorf("database address must not have trailing slash")
	}

	// parse the database address into a url structure
	parsed, err := url.Parse(d.Address)
	if err != nil {
		return fmt.Errorf("unable to parse database address: %v", err)
	}

	// save the parsed database address
	d.Url = parsed

	// verify the compression level provided is valid
	switch d.CompressionLevel {
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
		return fmt.Errorf("invalid database compression level provided: %d", d.CompressionLevel)
	}

	// enforce AES-256, so check explicitly for 32 bytes on the key
	//
	// nolint: gomnd // ignore magic number
	if len(d.EncryptionKey) != 32 {
		return fmt.Errorf("invalid database encryption key provided: %d", len(d.EncryptionKey))
	}

	return nil
}

// helper function to validate the secret CLI configuration.
//
// nolint:lll // ignoring line length check to avoid breaking up error messages
func validateSecret(c *cli.Context) error {
	logrus.Trace("Validating secret CLI configuration")

	if c.Bool("vault-driver") {
		if len(c.String("vault-addr")) == 0 {
			return fmt.Errorf("vault-addr (VELA_SECRET_VAULT_ADDR or SECRET_VAULT_ADDR) flag not specified")
		}

		if len(c.String("vault-token")) == 0 && len(c.String("vault-auth-method")) == 0 {
			return fmt.Errorf("vault-token (VELA_SECRET_VAULT_TOKEN or SECRET_VAULT_TOKEN) or vault-auth-method (VELA_SECRET_VAULT_AUTH_METHOD or SECRET_VAULT_AUTH_METHOD) flag not specified")
		}

		if len(c.String("vault-token")) == 0 {
			switch c.String("vault-auth-method") {
			case "aws":
			default:
				return fmt.Errorf("vault auth method of '%s' is unsupported", c.String("vault-auth-method"))
			}

			if c.String("vault-auth-method") == "aws" {
				if len(c.String("vault-aws-role")) == 0 {
					return fmt.Errorf("vault-aws-role (VELA_SECRET_VAULT_AWS_ROLE or SECRET_VAULT_AWS_ROLE) flag not specified")
				}
			}
		}
	}

	return nil
}

// Validate verifies the Security is properly configured.
func (s *Security) Validate() error {
	logrus.Trace("validating security configuration")

	// verify the refresh tokens last longer than access tokens
	if s.RefreshToken <= s.AccessToken {
		return fmt.Errorf("refresh token duration must be larger than access token duration")
	}

	// check if secure cookies are disabled
	if !s.SecureCookie {
		logrus.Warning("secure cookies are disabled - running with insecure mode")
	}

	// check if webhook validation is disabled
	if !s.WebhookValidation {
		logrus.Warning("webhook validation is disabled - running with insecure mode")
	}

	return nil
}

// Validate verifies the Source is properly configured.
func (s *Source) Validate() error {
	logrus.Trace("validating source configuration")

	// verify a source driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no source driver provided")
	}

	// verify a source address was provided
	if len(s.Address) == 0 {
		return fmt.Errorf("no source address provided")
	}

	// check if the source address has a scheme
	if !strings.Contains(s.Address, "://") {
		return fmt.Errorf("source address must be fully qualified (<scheme>://<host>)")
	}

	// check if the source address has a trailing slash
	if strings.HasSuffix(s.Address, "/") {
		return fmt.Errorf("source address must not have trailing slash")
	}

	// parse the source address into a url structure
	parsed, err := url.Parse(s.Address)
	if err != nil {
		return fmt.Errorf("unable to parse source address: %v", err)
	}

	// save the parsed source address
	s.Url = parsed

	// verify a source OAuth client ID was provided
	if len(s.ClientID) == 0 {
		return fmt.Errorf("no source client id provided")
	}

	// verify a source OAuth client secret was provided
	if len(s.ClientSecret) == 0 {
		return fmt.Errorf("no source client secret provided")
	}

	return nil
}

// Validate verifies the WebUI is properly configured.
func (w *WebUI) Validate() error {
	// verify a web UI address was provided
	if len(w.Address) == 0 {
		logrus.Warning("no web UI address provided - running in headless mode")

		return nil
	}

	// check if the web UI address has a scheme
	if !strings.Contains(w.Address, "://") {
		return fmt.Errorf("web UI address must be fully qualified (<scheme>://<host>)")
	}

	// check if the web UI address has a trailing slash
	if strings.HasSuffix(w.Address, "/") {
		return fmt.Errorf("web UI address must not have trailing slash")
	}

	// parse the web UI address into a url structure
	parsed, err := url.Parse(w.Address)
	if err != nil {
		return fmt.Errorf("unable to parse web UI address: %v", err)
	}

	// save the parsed web UI address
	w.Url = parsed

	// verify a web UI OAuth endpoint was provided
	if len(w.OAuthEndpoint) == 0 {
		return fmt.Errorf("no web UI OAuth endpoint provided")
	}

	return nil
}

// Validate verifies the Server is properly configured.
func (s *Server) Validate() error {
	// log a message indicating the configuration verification
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
	logrus.Info("validating server configuration")

	// check that hostname was properly populated
	if len(s.Config.API.Address.Hostname()) == 0 {
		switch strings.ToLower(s.Config.API.Address.Scheme) {
		case "http", "https":
			retErr := "server address invalid: %s"
			return fmt.Errorf(retErr, w.Config.API.Address.String())
		default:
			// hostname will be empty if a scheme is not provided
			retErr := "server address invalid, no scheme: %s"
			return fmt.Errorf(retErr, s.Config.API.Address.String())
		}
	}

	// verify the API configuration
	err := s.Config.API.Validate()
	if err != nil {
		return err
	}

	// verify the build configuration
	err = s.Config.Build.Validate()
	if err != nil {
		return err
	}

	// verify the compiler configuration
	err = s.Config.Compiler.Validate()
	if err != nil {
		return err
	}

	// verify the database configuration
	err = s.Config.Database.Validate()
	if err != nil {
		return err
	}

	// verify the queue configuration
	//
	// https://godoc.org/github.com/go-vela/pkg-queue/queue#Setup.Validate
	err = s.Config.Queue.Validate()
	if err != nil {
		return err
	}

	// verify the security configuration
	err = s.Config.Security.Validate()
	if err != nil {
		return err
	}

	// verify the source configuration
	err = s.Config.Source.Validate()
	if err != nil {
		return err
	}

	// verify the web UI configuration
	err = s.Config.WebUI.Validate()
	if err != nil {
		return err
	}

	return nil
}
