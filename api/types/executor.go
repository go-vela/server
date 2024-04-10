// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"strings"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// Executor is the library representation of an executor for a worker.
//
// swagger:model Executor
type Executor struct {
	ID           *int64          `json:"id,omitempty"`
	Host         *string         `json:"host,omitempty"`
	Runtime      *string         `json:"runtime,omitempty"`
	Distribution *string         `json:"distribution,omitempty"`
	Build        *library.Build  `json:"build,omitempty"`
	Repo         *Repo           `json:"repo,omitempty"`
	Pipeline     *pipeline.Build `json:"pipeline,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Executor type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (e *Executor) GetID() int64 {
	// return zero value if Executor type or ID field is nil
	if e == nil || e.ID == nil {
		return 0
	}

	return *e.ID
}

// GetHost returns the Host field.
//
// When the provided Executor type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (e *Executor) GetHost() string {
	// return zero value if Executor type or Host field is nil
	if e == nil || e.Host == nil {
		return ""
	}

	return *e.Host
}

// GetRuntime returns the Runtime field.
//
// When the provided Executor type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (e *Executor) GetRuntime() string {
	// return zero value if Executor type or Runtime field is nil
	if e == nil || e.Runtime == nil {
		return ""
	}

	return *e.Runtime
}

// GetDistribution returns the Distribution field.
//
// When the provided Executor type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (e *Executor) GetDistribution() string {
	// return zero value if Executor type or Distribution field is nil
	if e == nil || e.Distribution == nil {
		return ""
	}

	return *e.Distribution
}

// GetBuild returns the Build field.
//
// When the provided Executor type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (e *Executor) GetBuild() library.Build {
	// return zero value if Executor type or Build field is nil
	if e == nil || e.Build == nil {
		return library.Build{}
	}

	return *e.Build
}

// GetRepo returns the Repo field.
//
// When the provided Executor type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (e *Executor) GetRepo() Repo {
	// return zero value if Executor type or Repo field is nil
	if e == nil || e.Repo == nil {
		return Repo{}
	}

	return *e.Repo
}

// GetPipeline returns the Pipeline field.
//
// When the provided Executor type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (e *Executor) GetPipeline() pipeline.Build {
	// return zero value if Executor type or Pipeline field is nil
	if e == nil || e.Pipeline == nil {
		return pipeline.Build{}
	}

	return *e.Pipeline
}

// SetID sets the ID field.
//
// When the provided Executor type is nil, it
// will set nothing and immediately return.
func (e *Executor) SetID(v int64) {
	// return if Executor type is nil
	if e == nil {
		return
	}

	e.ID = &v
}

// SetHost sets the Host field.
//
// When the provided Executor type is nil, it
// will set nothing and immediately return.
func (e *Executor) SetHost(v string) {
	// return if Executor type is nil
	if e == nil {
		return
	}

	e.Host = &v
}

// SetRuntime sets the Runtime field.
//
// When the provided Executor type is nil, it
// will set nothing and immediately return.
func (e *Executor) SetRuntime(v string) {
	// return if Executor type is nil
	if e == nil {
		return
	}

	e.Runtime = &v
}

// SetDistribution sets the Distribution field.
//
// When the provided Executor type is nil, it
// will set nothing and immediately return.
func (e *Executor) SetDistribution(v string) {
	// return if Executor type is nil
	if e == nil {
		return
	}

	e.Distribution = &v
}

// SetBuild sets the Build field.
//
// When the provided Executor type is nil, it
// will set nothing and immediately return.
func (e *Executor) SetBuild(v library.Build) {
	// return if Executor type is nil
	if e == nil {
		return
	}

	e.Build = &v
}

// SetRepo sets the Repo field.
//
// When the provided Executor type is nil, it
// will set nothing and immediately return.
func (e *Executor) SetRepo(v Repo) {
	// return if Executor type is nil
	if e == nil {
		return
	}

	e.Repo = &v
}

// SetPipeline sets the pipeline Build field.
//
// When the provided Executor type is nil, it
// will set nothing and immediately return.
func (e *Executor) SetPipeline(v pipeline.Build) {
	// return if Executor type is nil
	if e == nil {
		return
	}

	e.Pipeline = &v
}

// String implements the Stringer interface for the Executor type.
func (e *Executor) String() string {
	return fmt.Sprintf(`{
  Build: %s,
  Distribution: %s,
  Host: %s,
  ID: %d,
  Repo: %v,
  Runtime: %s,
  Pipeline: %v,
}`,
		strings.ReplaceAll(e.Build.String(), " ", "  "),
		e.GetDistribution(),
		e.GetHost(),
		e.GetID(),
		strings.ReplaceAll(e.Repo.String(), " ", "  "),
		e.GetRuntime(),
		e.GetPipeline(),
	)
}
