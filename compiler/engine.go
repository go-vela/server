// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiler

import (
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
)

// Engine represents an interface for converting a yaml
// configuration to an executable pipeline for Vela.
type Engine interface {
	// Compiler Interface Functions

	// Compile defines a function that produces an executable
	// representation of a pipeline from an object. This calls
	// Parse internally to convert the object to a yaml configuration.
	Compile(interface{}) (*pipeline.Build, *library.Pipeline, error)

	// CompileLite defines a function that produces an light executable
	// representation of a pipeline from an object. This calls
	// Parse internally to convert the object to a yaml configuration.
	CompileLite(interface{}, bool, bool, []string) (*yaml.Build, *library.Pipeline, error)

	// Duplicate defines a function that
	// creates a clone of the Engine.
	Duplicate() Engine

	// Parse defines a function that converts
	// an object to a yaml configuration.
	Parse(interface{}, string, *yaml.Template) (*yaml.Build, []byte, error)

	// ParseRaw defines a function that converts
	// an object to a string.
	ParseRaw(interface{}) (string, error)

	// Validate defines a function that verifies
	// the yaml configuration is accurate.
	Validate(*yaml.Build) error

	// Clone Compiler Interface Functions

	// CloneStage defines a function that injects the
	// clone stage process into a yaml configuration.
	CloneStage(*yaml.Build) (*yaml.Build, error)
	// CloneStep defines a function that injects the
	// clone step process into a yaml configuration.
	CloneStep(*yaml.Build) (*yaml.Build, error)

	// Environment Compiler Interface Functions

	// EnvironmentStages defines a function that injects the environment
	// variables for each step in every stage into a yaml configuration.
	EnvironmentStages(yaml.StageSlice, raw.StringSliceMap) (yaml.StageSlice, error)
	// EnvironmentSteps defines a function that injects the environment
	// variables for each step into a yaml configuration.
	EnvironmentSteps(yaml.StepSlice, raw.StringSliceMap) (yaml.StepSlice, error)
	// EnvironmentStep defines a function that injects the environment
	// variables for a single step into a yaml configuration.
	EnvironmentStep(*yaml.Step, raw.StringSliceMap) (*yaml.Step, error)
	// EnvironmentServices defines a function that injects the environment
	// variables for each service into a yaml configuration.
	EnvironmentServices(yaml.ServiceSlice, raw.StringSliceMap) (yaml.ServiceSlice, error)

	// Expand Compiler Interface Functions

	// ExpandStages defines a function that injects the template
	// for each templated step in every stage in a yaml configuration.
	ExpandStages(*yaml.Build, map[string]*yaml.Template) (*yaml.Build, error)
	// ExpandSteps defines a function that injects the template
	// for each templated step in a yaml configuration.
	ExpandSteps(*yaml.Build, map[string]*yaml.Template) (*yaml.Build, error)

	// Init Compiler Interface Functions

	// InitStage defines a function that injects the
	// init stage process into a yaml configuration.
	InitStage(*yaml.Build) (*yaml.Build, error)
	// InitStep step process into a yaml configuration.
	InitStep(*yaml.Build) (*yaml.Build, error)

	// Script Compiler Interface Functions

	// ScriptStages defines a function that injects the script
	// for each step in every stage in a yaml configuration.
	ScriptStages(yaml.StageSlice) (yaml.StageSlice, error)
	// ScriptSteps defines a function that injects the script
	// for each step in a yaml configuration.
	ScriptSteps(yaml.StepSlice) (yaml.StepSlice, error)

	// Substitute Compiler Interface Functions

	// SubstituteStages defines a function that replaces every
	// declared environment variable with it's corresponding
	// value for each step in every stage in a yaml configuration.
	SubstituteStages(yaml.StageSlice) (yaml.StageSlice, error)
	// SubstituteSteps defines a function that replaces every
	// declared environment variable with it's corresponding
	// value for each step in a yaml configuration.
	SubstituteSteps(yaml.StepSlice) (yaml.StepSlice, error)

	// Transform Compiler Interface Functions

	// TransformStages defines a function that converts a yaml
	// configuration with stages into an executable pipeline.
	TransformStages(*pipeline.RuleData, *yaml.Build) (*pipeline.Build, error)
	// TransformSteps defines a function that converts a yaml
	// configuration with steps into an executable pipeline.
	TransformSteps(*pipeline.RuleData, *yaml.Build) (*pipeline.Build, error)

	// With Compiler Interface Functions

	// WithBuild defines a function that sets
	// the library build type in the Engine.
	WithBuild(*library.Build) Engine
	// WithComment defines a function that sets
	// the comment in the Engine.
	WithComment(string) Engine
	// WithFiles defines a function that sets
	// the changeset files in the Engine.
	WithFiles([]string) Engine
	// WithLocal defines a function that sets
	// the compiler local field in the Engine.
	WithLocal(bool) Engine
	// WithMetadata defines a function that sets
	// the compiler Metadata type in the Engine.
	WithMetadata(*types.Metadata) Engine
	// WithRepo defines a function that sets
	// the library repo type in the Engine.
	WithRepo(*library.Repo) Engine
	// WithUser defines a function that sets
	// the library user type in the Engine.
	WithUser(*library.User) Engine
	// WithUser defines a function that sets
	// the private github client in the Engine.
	WithPrivateGitHub(string, string) Engine
	// WithLog defines a function that sets
	// the library log type in the Engine.
	WithLog(*library.Log) Engine
}
