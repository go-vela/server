// SPDX-License-Identifier: Apache-2.0

package installation

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// UserInterface represents the Vela interface for user
// functions with the supported Database backends.
//

type InstallationInterface interface {
	// Installation Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateInstallationTable defines a function that creates the installations table.
	CreateInstallationTable(context.Context, string) error

	// Installation Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CreateInstallation defines a function that creates a new installation.
	CreateInstallation(context.Context, *api.Installation) (*api.Installation, error)
	// DeleteInstallation defines a function that deletes an existing installation.
	DeleteInstallation(context.Context, *api.Installation) error
	// GetInstallation defines a function that gets an installation by ID.
	GetInstallation(context.Context, string) (*api.Installation, error)
}
