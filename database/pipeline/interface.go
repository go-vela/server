// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// PipelineInterface represents the Vela interface for pipeline
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type PipelineInterface interface {
	// Pipeline Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreatePipelineIndexes defines a function that creates the indexes for the pipelines table.
	CreatePipelineIndexes(context.Context) error
	// CreatePipelineTable defines a function that creates the pipelines table.
	CreatePipelineTable(context.Context, string) error

	// Pipeline Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountPipelines defines a function that gets the count of all pipelines.
	CountPipelines(context.Context) (int64, error)
	// CountPipelinesForRepo defines a function that gets the count of pipelines by repo ID.
	CountPipelinesForRepo(context.Context, *api.Repo) (int64, error)
	// CreatePipeline defines a function that creates a new pipeline.
	CreatePipeline(context.Context, *api.Pipeline) (*api.Pipeline, error)
	// DeletePipeline defines a function that deletes an existing pipeline.
	DeletePipeline(context.Context, *api.Pipeline) error
	// GetPipeline defines a function that gets a pipeline by ID.
	GetPipeline(context.Context, int64) (*api.Pipeline, error)
	// GetPipelineForRepo defines a function that gets a pipeline by commit SHA and repo ID.
	GetPipelineForRepo(context.Context, string, *api.Repo) (*api.Pipeline, error)
	// ListPipelines defines a function that gets a list of all pipelines.
	ListPipelines(context.Context) ([]*api.Pipeline, error)
	// ListPipelinesForRepo defines a function that gets a list of pipelines by repo ID.
	ListPipelinesForRepo(context.Context, *api.Repo, int, int) ([]*api.Pipeline, int64, error)
	// UpdatePipeline defines a function that updates an existing pipeline.
	UpdatePipeline(context.Context, *api.Pipeline) (*api.Pipeline, error)
}
