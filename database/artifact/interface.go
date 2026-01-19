// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// ArtifactInterface represents the Vela interface for artifacts
// functions with the supported Database backends.
//

type ArtifactInterface interface {
	// Artifact Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateArtifactIndexes defines a function that creates the indexes for the artifacts table.
	CreateArtifactIndexes(context.Context) error
	// CreateArtifactTable defines a function that creates the artifacts table.
	CreateArtifactTable(context.Context, string) error

	// Artifact Management Functions

	// CountArtifacts returns the count of all artifacts.
	CountArtifacts(context.Context) (int64, error)

	// CreateArtifact creates a new artifact.
	CreateArtifact(context.Context, *api.Artifact) (*api.Artifact, error)

	// DeleteArtifact removes an artifact by ID.
	DeleteArtifact(context.Context, *api.Artifact) error

	// GetArtifact returns an artifact by ID.
	GetArtifact(context.Context, int64) (*api.Artifact, error)

	// GetArtifactForBuild defines a function that gets an artifact by number and build ID.
	GetArtifactForBuild(context.Context, *api.Build) (*api.Artifact, error)

	// ListArtifacts returns a list of all artifacts.
	ListArtifacts(context.Context) ([]*api.Artifact, error)

	// ListArtifactsByBuildID returns a list of artifacts by build ID.
	ListArtifactsByBuildID(context.Context, int64) ([]*api.Artifact, error)

	// UpdateArtifact updates an artifact by ID.
	UpdateArtifact(context.Context, *api.Artifact) (*api.Artifact, error)
}
