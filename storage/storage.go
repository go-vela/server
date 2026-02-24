// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
)

// New creates and returns a Vela service capable of
// integrating with the configured storage environment.
// Currently, the following storages are supported:
//
// * minio
// .
func New(s *Setup) (Storage, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/server/storage#Setup.Validate
	if s.Enable {
		err := s.Validate()
		if err != nil {
			return nil, fmt.Errorf("unable to validate storage setup: %w", err)
		}

		logrus.Debug("creating storage client from setup")
		// process the storage driver being provided
		switch s.Driver {
		case constants.DriverMinio:
			// handle the storage driver being provided
			//
			// https://pkg.go.dev/github.com/go-vela/server/storage?tab=doc#Setup.Minio
			return s.Minio()
		default:
			// handle an invalid storage driver being provided
			return nil, fmt.Errorf("invalid storage driver provided: %s", s.Driver)
		}
	}

	return nil, nil
}
