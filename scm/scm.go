// SPDX-License-Identifier: Apache-2.0

package scm

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
)

// New creates and returns a Vela service capable of
// integrating with the configured scm provider.
//
// Currently the following scm providers are supported:
//
// * Github
// .
func New(ctx context.Context, s *Setup) (Service, error) {
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
		return s.Github(ctx)
	case constants.DriverGitlab:
		// handle the Gitlab scm driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/scm?tab=doc#Setup.Gitlab
		return s.Gitlab(ctx)
	default:
		// handle an invalid scm driver being provided
		return nil, fmt.Errorf("invalid scm driver provided: %s", s.Driver)
	}
}
