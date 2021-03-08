// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
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

	// verify the database setup
	//
	// https://godoc.org/github.com/go-vela/server/database#Setup.Validate
	return d.Config.Validate()
}

// Validate verifies the Queue is properly configured.
func (q *Queue) Validate() error {
	logrus.Trace("validating queue configuration")

	// verify the queue setup
	//
	// https://godoc.org/github.com/go-vela/pkg-queue/queue#Setup.Validate
	return q.Config.Validate()
}

// Validate verifies the Secrets is properly configured.
func (s *Secrets) Validate() error {
	logrus.Trace("validating secrets configuration")

	// check if the vault driver is enabled
	if s.Vault.Driver {
		// verify a vault address was provided
		if len(s.Vault.Address) == 0 {
			return fmt.Errorf("no secrets Vault address provided")
		}

		// check if the vault address has a scheme
		if !strings.Contains(s.Vault.Address, "://") {
			return fmt.Errorf("secrets Vault address must be fully qualified (<scheme>://<host>)")
		}

		// check if the vault address has a trailing slash
		if strings.HasSuffix(s.Vault.Address, "/") {
			return fmt.Errorf("secrets Vault address must not have trailing slash")
		}

		// verify a vault token or authentication method was provided
		if len(s.Vault.Token) == 0 && len(s.Vault.AuthMethod) == 0 {
			return fmt.Errorf("no secrets Vault token or authentication method provided")
		}

		// check if a vault token was provided
		if len(s.Vault.Token) == 0 {
			// verify the vault authentication method provided is valid
			switch s.Vault.AuthMethod {
			case "aws":
				// verify a vault AWS role is provided
				if len(s.Vault.AwsRole) == 0 {
					return fmt.Errorf("no secrets Vault AWS role provided")
				}
			default:
				return fmt.Errorf("invalid secrets vault authentication method provided: %s", s.Vault.AuthMethod)
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
		logrus.Warning("secure cookies are disabled - running in insecure mode")
	}

	// check if webhook validation is disabled
	if !s.WebhookValidation {
		logrus.Warning("webhook validation is disabled - running in insecure mode")
	}

	return nil
}

// Validate verifies the Source is properly configured.
func (s *Source) Validate() error {
	logrus.Trace("validating source configuration")

	// verify the source setup
	//
	// https://godoc.org/github.com/go-vela/server/source#Setup.Validate
	return s.Config.Validate()
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
	err = s.Config.Queue.Validate()
	if err != nil {
		return err
	}

	// verify the secrets configuration
	err = s.Config.Secrets.Validate()
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
