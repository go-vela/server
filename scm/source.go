// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package scm

import (
	"fmt"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// nolint: godot // top level comment ends in a list
//
// New creates and returns a Vela service capable of
// integrating with the configured source provider.
//
// Currently the following source providers are supported:
//
// * Github
func New(s *Setup) (Service, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/server/scm?tab=doc#Setup.Validate
	err := s.Validate()
	if err != nil {
		return nil, err
	}

	logrus.Debug("creating source service from setup")
	// process the source driver being provided
	switch s.Driver {
	case constants.DriverGithub:
		// handle the Github source driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/scm?tab=doc#Setup.Github
		return s.Github()
	case constants.DriverGitlab:
		// handle the Gitlab source driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/scm?tab=doc#Setup.Gitlab
		return s.Gitlab()
	default:
		// handle an invalid source driver being provided
		return nil, fmt.Errorf("invalid source driver provided: %s", s.Driver)
	}
}
