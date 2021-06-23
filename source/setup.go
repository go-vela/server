// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package source

import (
	"fmt"
	"strings"

	"github.com/go-vela/server/source/github"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// Setup represents the configuration necessary for
// creating a Vela service capable of integrating
// with a configured source system.
type Setup struct {
	// Source Configuration

	// specifies the driver to use for the source client
	Driver string
	// specifies the address to use for the source client
	Address string
	// specifies the OAuth client ID from the source system to use for the source client
	ClientID string
	// specifies the OAuth client secret from the source system to use for the source client
	ClientSecret string
	// specifies the Vela server address to use for the source client
	ServerAddress string
	// specifies the context for the commit status to use for the source client
	StatusContext string
	// specifies the Vela web UI address to use for the source client
	WebUIAddress string
	// specifies the OAuth scopes to use for the source client
	Scopes []string
}

// Github creates and returns a Vela service capable of
// integrating with a Github source system.
func (s *Setup) Github() (Service, error) {
	logrus.Trace("creating github source client from setup")

	// create new Github source service
	//
	// https://pkg.go.dev/github.com/go-vela/server/source/github?tab=doc#New
	return github.New(
		github.WithAddress(s.Address),
		github.WithClientID(s.ClientID),
		github.WithClientSecret(s.ClientSecret),
		github.WithServerAddress(s.ServerAddress),
		github.WithStatusContext(s.StatusContext),
		github.WithWebUIAddress(s.WebUIAddress),
		github.WithScopes(s.Scopes),
	)
}

// Gitlab creates and returns a Vela service capable of
// integrating with a Gitlab source system.
func (s *Setup) Gitlab() (Service, error) {
	logrus.Trace("creating gitlab source client from setup")

	return nil, fmt.Errorf("unsupported source driver: %s", constants.DriverGitlab)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating source setup for client")

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

	// verify a source OAuth client ID was provided
	if len(s.ClientID) == 0 {
		return fmt.Errorf("no source client id provided")
	}

	// verify a source OAuth client secret was provided
	if len(s.ClientSecret) == 0 {
		return fmt.Errorf("no source client secret provided")
	}

	// verify a source status context secret was provided
	if len(s.StatusContext) == 0 {
		return fmt.Errorf("no source status context provided")
	}

	if len(s.Scopes) == 0 {
		return fmt.Errorf("no source scopes provided")
	}

	// setup is valid
	return nil
}
