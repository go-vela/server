// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/raw"
)

// Build is the yaml representation of a build for a pipeline.
type Build struct {
	Version     string             `yaml:"version,omitempty"     json:"version,omitempty"     jsonschema:"required,minLength=1,description=Provide syntax version used to evaluate the pipeline.\nReference: https://go-vela.github.io/docs/reference/yaml/version/"`
	Metadata    Metadata           `yaml:"metadata,omitempty"    json:"metadata,omitempty"    jsonschema:"description=Pass extra information.\nReference: https://go-vela.github.io/docs/reference/yaml/metadata/"`
	Environment raw.StringSliceMap `yaml:"environment,omitempty" json:"environment,omitempty" jsonschema:"description=Provide global environment variables injected into the container environment.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-environment-key"`
	Worker      Worker             `yaml:"worker,omitempty"      json:"worker,omitempty"      jsonschema:"description=Limit the pipeline to certain types of workers.\nReference: https://go-vela.github.io/docs/reference/yaml/worker/"`
	Secrets     SecretSlice        `yaml:"secrets,omitempty"     json:"secrets,omitempty"     jsonschema:"description=Provide sensitive information.\nReference: https://go-vela.github.io/docs/reference/yaml/secrets/"`
	Services    ServiceSlice       `yaml:"services,omitempty"    json:"services,omitempty"    jsonschema:"description=Provide detached (headless) execution instructions.\nReference: https://go-vela.github.io/docs/reference/yaml/services/"`
	Stages      StageSlice         `yaml:"stages,omitempty"      json:"stages,omitempty"      jsonschema:"oneof_required=stages,description=Provide parallel execution instructions.\nReference: https://go-vela.github.io/docs/reference/yaml/stages/"`
	Steps       StepSlice          `yaml:"steps,omitempty"       json:"steps,omitempty"       jsonschema:"oneof_required=steps,description=Provide sequential execution instructions.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/"`
	Templates   TemplateSlice      `yaml:"templates,omitempty"   json:"templates,omitempty"   jsonschema:"description=Provide the name of templates to expand.\nReference: https://go-vela.github.io/docs/reference/yaml/templates/"`
	Git         Git                `yaml:"git,omitempty"       json:"git,omitempty"      jsonschema:"description=Provide the git access specifications.\nReference: https://go-vela.github.io/docs/reference/yaml/git/"`
}

// ToPipelineAPI converts the Build type to an API Pipeline type.
func (b *Build) ToPipelineAPI() *api.Pipeline {
	pipeline := new(api.Pipeline)

	pipeline.SetFlavor(b.Worker.Flavor)
	pipeline.SetPlatform(b.Worker.Platform)
	pipeline.SetVersion(b.Version)
	pipeline.SetServices(len(b.Services) > 0)
	pipeline.SetStages(len(b.Stages) > 0)
	pipeline.SetSteps(len(b.Steps) > 0)
	pipeline.SetTemplates(len(b.Templates) > 0)

	// set default for external and internal secrets
	external := false
	internal := false

	// iterate through all secrets in the build
	for _, secret := range b.Secrets {
		// check if external and internal secrets have been found
		if external && internal {
			// exit the loop since both secrets have been found
			break
		}

		// check if the secret origin is empty
		if secret.Origin.Empty() {
			// origin was empty so an internal secret was found
			internal = true
		} else {
			// origin was not empty so an external secret was found
			external = true
		}
	}

	pipeline.SetExternalSecrets(external)
	pipeline.SetInternalSecrets(internal)

	return pipeline
}

// UnmarshalYAML implements the Unmarshaler interface for the Build type.
func (b *Build) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// build we try unmarshalling to
	build := new(struct {
		Git         Git
		Version     string
		Metadata    Metadata
		Environment raw.StringSliceMap
		Worker      Worker
		Secrets     SecretSlice
		Services    ServiceSlice
		Stages      StageSlice
		Steps       StepSlice
		Templates   TemplateSlice
	})

	// attempt to unmarshal as a build type
	err := unmarshal(build)
	if err != nil {
		return err
	}

	// give the documented default value to metadata environment
	if build.Metadata.Environment == nil {
		build.Metadata.Environment = []string{"steps", "services", "secrets"}
	}

	// override the values
	b.Git = build.Git
	b.Version = build.Version
	b.Metadata = build.Metadata
	b.Environment = build.Environment
	b.Worker = build.Worker
	b.Secrets = build.Secrets
	b.Services = build.Services
	b.Stages = build.Stages
	b.Steps = build.Steps
	b.Templates = build.Templates

	return nil
}
