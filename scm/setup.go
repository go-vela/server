// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package scm

import (
	"fmt"
	"strings"

	"github.com/go-vela/server/scm/native"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// Setup represents the configuration necessary for
// creating a Vela service capable of integrating
// with a configured scm system.
type Setup struct {
	// scm Configuration

	// specifies the driver to use for the scm client
	Driver string
	// specifies the address to use for the scm client
	Address string
	// specifies the OAuth client ID from the scm system to use for the scm client
	ClientID string
	// specifies the OAuth client secret from the scm system to use for the scm client
	ClientSecret string
	// specifies the Vela server address to use for the scm client
	ServerAddress string
	// specifies the Vela server address that the scm provider should use to send Vela webhooks
	ServerWebhookAddress string
	// specifies the context for the commit status to use for the scm client
	StatusContext string
	// specifies the Vela web UI address to use for the scm client
	WebUIAddress string
	// specifies the OAuth scopes to use for the scm client
	Scopes []string
}

// BitBucket creates and returns a Vela service capable of
// integrating with a Gitlab scm system.
func (s *Setup) BitBucket() (Service, error) {
	logrus.Trace("creating bitbucket scm client from setup")

	// create new BitBucket scm service based on Native SCM implementation
	//
	// https://pkg.go.dev/github.com/go-vela/server/scm/github?tab=doc#New
	return native.New(
		native.WithAddress(s.Address),
		native.WithClientID(s.ClientID),
		native.WithClientSecret(s.ClientSecret),
		native.WithServerAddress(s.ServerAddress),
		native.WithServerWebhookAddress(s.ServerWebhookAddress),
		native.WithStatusContext(s.StatusContext),
		native.WithWebUIAddress(s.WebUIAddress),
		native.WithScopes(s.Scopes),
	)
}

// Github creates and returns a Vela service capable of
// integrating with a Github scm system.
func (s *Setup) Github() (Service, error) {
	logrus.Trace("creating github scm client from setup")

	// create new Github scm service
	//
	// https://pkg.go.dev/github.com/go-vela/server/scm/github?tab=doc#New
	// return github.New(
	// 	github.WithAddress(s.Address),
	// 	github.WithClientID(s.ClientID),
	// 	github.WithClientSecret(s.ClientSecret),
	// 	github.WithServerAddress(s.ServerAddress),
	// 	github.WithServerWebhookAddress(s.ServerWebhookAddress),
	// 	github.WithStatusContext(s.StatusContext),
	// 	github.WithWebUIAddress(s.WebUIAddress),
	// 	github.WithScopes(s.Scopes),
	// )
	// TODO: uncomment this to swap to the native implementation for GitHub
	// Consider using this in future versions consolidate on Native package.
	return native.New(
		native.WithAddress(s.Address),
		native.WithClientID(s.ClientID),
		native.WithClientSecret(s.ClientSecret),
		native.WithServerAddress(s.ServerAddress),
		native.WithServerWebhookAddress(s.ServerWebhookAddress),
		native.WithStatusContext(s.StatusContext),
		native.WithWebUIAddress(s.WebUIAddress),
		native.WithScopes(s.Scopes),
	)
}

// Gitlab creates and returns a Vela service capable of
// integrating with a Gitlab scm system.
func (s *Setup) Gitlab() (Service, error) {
	logrus.Trace("creating gitlab scm client from setup")

	return nil, fmt.Errorf("unsupported scm driver: %s", constants.DriverGitlab)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating scm setup for client")

	// verify a scm driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no scm driver provided")
	}

	// verify a scm address was provided
	if len(s.Address) == 0 {
		return fmt.Errorf("no scm address provided")
	}

	// check if the scm address has a scheme
	if !strings.Contains(s.Address, "://") {
		return fmt.Errorf("scm address must be fully qualified (<scheme>://<host>)")
	}

	// check if the scm address has a trailing slash
	if strings.HasSuffix(s.Address, "/") {
		return fmt.Errorf("scm address must not have trailing slash")
	}

	// verify a scm OAuth client ID was provided
	if len(s.ClientID) == 0 {
		return fmt.Errorf("no scm client id provided")
	}

	// verify a scm OAuth client secret was provided
	if len(s.ClientSecret) == 0 {
		return fmt.Errorf("no scm client secret provided")
	}

	// verify a scm status context secret was provided
	if len(s.StatusContext) == 0 {
		return fmt.Errorf("no scm status context provided")
	}

	if len(s.Scopes) == 0 {
		return fmt.Errorf("no scm scopes provided")
	}

	// setup is valid
	return nil
}
