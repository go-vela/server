// SPDX-License-Identifier: Apache-2.0

package compiler

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/scm"
)

// Engine represents an interface for converting a yaml
// configuration to an executable pipeline for Vela.
type Engine interface {
	// Compiler Interface Functions

	// Compile defines a function that produces an executable
	// representation of a pipeline from an object. This calls
	// Parse internally to convert the object to a yaml configuration.
	Compile(context.Context, interface{}) (*pipeline.Build, *api.Pipeline, error)

	// CompileLite defines a function that produces an light executable
	// representation of a pipeline from an object. This calls
	// Parse internally to convert the object to a yaml configuration.
	CompileLite(context.Context, interface{}, *pipeline.RuleData, bool) (*yaml.Build, *api.Pipeline, error)

	// Duplicate defines a function that
	// creates a clone of the Engine.
	Duplicate() Engine

	// Parse defines a function that converts
	// an object to a yaml configuration.
	Parse(interface{}, string, *yaml.Template) (*yaml.Build, []byte, []string, error)

	// ParseRaw defines a function that converts
	// an object to a string.
	ParseRaw(interface{}) (string, error)

	// ValidateYAML defines a function that verifies
	// the yaml configuration is accurate.
	ValidateYAML(*yaml.Build) error

	// ValidatePipeline defines a function that verifies
	// the final pipeline build is accurate.
	ValidatePipeline(*pipeline.Build) error

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
	ExpandStages(context.Context, *yaml.Build, map[string]*yaml.Template, *pipeline.RuleData, []string) (*yaml.Build, []string, error)
	// ExpandSteps defines a function that injects the template
	// for each templated step in a yaml configuration with the provided template depth.
	ExpandSteps(context.Context, *yaml.Build, map[string]*yaml.Template, *pipeline.RuleData, []string, int) (*yaml.Build, []string, error)

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
	// the API build type in the Engine.
	WithBuild(*api.Build) Engine
	// WithComment defines a function that sets
	// the comment in the Engine.
	WithComment(string) Engine
	// WithCommit defines a function that sets
	// the commit in the Engine.
	WithCommit(string) Engine
	// WithFiles defines a function that sets
	// the changeset files in the Engine.
	WithFiles([]string) Engine
	// WithLocal defines a function that sets
	// the compiler local field in the Engine.
	WithLocal(bool) Engine
	// WithLocalTemplates defines a function that sets
	// the compiler local templates field in the Engine.
	WithLocalTemplates([]string) Engine
	// WithMetadata defines a function that sets
	// the compiler Metadata type in the Engine.
	WithMetadata(*internal.Metadata) Engine
	// WithRepo defines a function that sets
	// the API repo type in the Engine.
	WithRepo(*api.Repo) Engine
	// WithUser defines a function that sets
	// the API user type in the Engine.
	WithUser(*api.User) Engine
	// WithLabel defines a function that sets
	// the label(s) in the Engine.
	WithLabels([]string) Engine
	// WithSCM defines a function that sets
	// the scm in the Engine.
	WithSCM(scm.Service) Engine
	// WithDatabase defines a function that sets
	// the database in the Engine.
	WithDatabase(database.Interface) Engine
	// WithPrivateGitHub defines a function that sets
	// the private github client in the Engine.
	WithPrivateGitHub(context.Context, string, string) Engine
	// GetSettings defines a function that returns new api settings
	// with the compiler Engine fields filled.
	GetSettings() settings.Compiler
	// SetSettings defines a function that takes api settings
	// and updates the compiler Engine.
	SetSettings(*settings.Platform)
}
