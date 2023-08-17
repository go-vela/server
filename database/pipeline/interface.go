// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"context"

	"github.com/go-vela/types/library"
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
	CountPipelinesForRepo(context.Context, *library.Repo) (int64, error)
	// CreatePipeline defines a function that creates a new pipeline.
	CreatePipeline(context.Context, *library.Pipeline) (*library.Pipeline, error)
	// DeletePipeline defines a function that deletes an existing pipeline.
	DeletePipeline(context.Context, *library.Pipeline) error
	// GetPipeline defines a function that gets a pipeline by ID.
	GetPipeline(context.Context, int64) (*library.Pipeline, error)
	// GetPipelineForRepo defines a function that gets a pipeline by commit SHA and repo ID.
	GetPipelineForRepo(context.Context, string, *library.Repo) (*library.Pipeline, error)
	// ListPipelines defines a function that gets a list of all pipelines.
	ListPipelines(context.Context) ([]*library.Pipeline, error)
	// ListPipelinesForRepo defines a function that gets a list of pipelines by repo ID.
	ListPipelinesForRepo(context.Context, *library.Repo, int, int) ([]*library.Pipeline, int64, error)
	// UpdatePipeline defines a function that updates an existing pipeline.
	UpdatePipeline(context.Context, *library.Pipeline) (*library.Pipeline, error)
}
