// SPDX-License-Identifier: Apache-2.0

package scm

import (
	"fmt"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// New creates and returns a Vela service capable of
// integrating with the configured scm provider.
//
// Currently the following scm providers are supported:
//
// * Github
// .
func New(s *Setup) (Service, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/server/scm?tab=doc#Setup.Validate
	err := s.Validate()
	if err != nil {
		return nil, err
	}

	logrus.Debug("creating scm service from setup")
	// process the scm driver being provided
	switch s.Driver {
	case constants.DriverGithub:
		// handle the Github scm driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/scm?tab=doc#Setup.Github
		return s.Github()
	case constants.DriverGitlab:
		// handle the Gitlab scm driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/scm?tab=doc#Setup.Gitlab
		return s.Gitlab()
	default:
		// handle an invalid scm driver being provided
		return nil, fmt.Errorf("invalid scm driver provided: %s", s.Driver)
	}
}
