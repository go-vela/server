// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"github.com/go-vela/types/library"
)

// PipelineService represents the Vela interface for pipeline
// functions with the supported Database backends.
//
// nolint: revive // ignore name stutter
type PipelineService interface {
	// Pipeline Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateIndexes creates the indexes for the pipelines table.
	CreateIndexes() error
	// CreateTable defines a function that creates the pipelines table.
	CreateTable(string) error

	// Pipeline Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountPipelines defines a function that gets the count of all pipelines.
	CountPipelines() (int64, error)
	// CountPipelinesForRepo defines a function that gets the count of pipelines by repo ID.
	CountPipelinesForRepo(*library.Repo) (int64, error)
	// CreatePipeline defines a function that creates a new pipeline.
	CreatePipeline(*library.Pipeline) error
	// DeletePipeline defines a function that deletes an existing pipeline.
	DeletePipeline(*library.Pipeline) error
	// GetPipeline defines a function that gets a pipeline by ID.
	GetPipeline(int64) (*library.Pipeline, error)
	// GetPipelineForRepo defines a function that gets a pipeline by commit SHA and repo ID.
	GetPipelineForRepo(string, *library.Repo) (*library.Pipeline, error)
	// LastPipelineForRepo defines a function that gets the last pipeline by repo ID.
	LastPipelineForRepo(*library.Repo) (*library.Pipeline, error)
	// ListPipelines defines a function that gets a list of all pipelines.
	ListPipelines() ([]*library.Pipeline, error)
	// ListPipelinesForRepo defines a function that gets a list of pipelines by repo ID.
	ListPipelinesForRepo(*library.Repo, int, int) ([]*library.Pipeline, int64, error)
	// UpdatePipeline defines a function that updates an existing pipeline.
	UpdatePipeline(*library.Pipeline) error
}
