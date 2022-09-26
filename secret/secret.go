// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"fmt"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// New creates and returns a Vela service capable of
// integrating with the configured secret provider.
//
// Currently the following secret providers are supported:
//
// * Native
// * Vault
// .
func New(s *Setup) (Service, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/server/secret?tab=doc#Setup.Validate
	err := s.Validate()
	if err != nil {
		return nil, err
	}

	logrus.Debug("creating secret service from setup")
	// process the secret driver being provided
	switch s.Driver {
	case constants.DriverNative:
		// handle the Native secret driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/secret?tab=doc#Setup.Native
		return s.Native()
	case constants.DriverVault:
		// handle the Vault secret driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/secret?tab=doc#Setup.Vault
		return s.Vault()
	default:
		// handle an invalid secret driver being provided
		return nil, fmt.Errorf("invalid secret driver provided: %s", s.Driver)
	}
}
