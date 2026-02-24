// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/constants"
)

// Build is the pipeline representation of a build for a pipeline.
//
// swagger:model PipelineBuild
type Build struct {
	ID          string             `json:"id,omitempty"          yaml:"id,omitempty"`
	Token       string             `json:"token,omitempty"       yaml:"token,omitempty"`
	TokenExp    int64              `json:"token_exp,omitempty"   yaml:"token_exp,omitempty"`
	Version     string             `json:"version,omitempty"     yaml:"version,omitempty"`
	Metadata    Metadata           `json:"metadata,omitempty"    yaml:"metadata,omitempty"`
	Environment raw.StringSliceMap `json:"environment,omitempty" yaml:"environment,omitempty"`
	Worker      Worker             `json:"worker,omitempty"      yaml:"worker,omitempty"`
	Deployment  Deployment         `json:"deployment,omitempty"  yaml:"deployment,omitempty"`
	Secrets     SecretSlice        `json:"secrets,omitempty"     yaml:"secrets,omitempty"`
	Services    ContainerSlice     `json:"services,omitempty"    yaml:"services,omitempty"`
	Stages      StageSlice         `json:"stages,omitempty"      yaml:"stages,omitempty"`
	Steps       ContainerSlice     `json:"steps,omitempty"       yaml:"steps,omitempty"`
}

// Purge removes the steps, in every stage, that contain a ruleset
// that do not match the provided ruledata. If all steps from a
// stage are removed, then the entire stage is removed from the
// pipeline. If no stages are provided in the pipeline, then the
// function will remove the steps that have a ruleset that do not
// match the provided ruledata. If both stages and steps are
// provided, then an empty pipeline is returned.
func (b *Build) Purge(r *RuleData) (*Build, error) {
	// return an empty pipeline if both stages and steps are provided
	if len(b.Stages) > 0 && len(b.Steps) > 0 {
		return nil, fmt.Errorf("cannot have both stages and steps at the top level of pipeline")
	}

	// purge stages pipeline if stages are provided
	if len(b.Stages) > 0 {
		pStages, err := b.Stages.Purge(r)
		if err != nil {
			return nil, err
		}

		b.Stages = *pStages
	}

	// purge steps pipeline if steps are provided
	if len(b.Steps) > 0 {
		pSteps, err := b.Steps.Purge(r)
		if err != nil {
			return nil, err
		}

		b.Steps = *pSteps
	}

	// purge services in pipeline if services are provided
	if len(b.Services) > 0 {
		pServices, err := b.Services.Purge(r)
		if err != nil {
			return nil, err
		}

		b.Services = *pServices
	}

	// purge secrets in pipeline if secrets are provided
	if len(b.Secrets) > 0 {
		pSecrets, err := b.Secrets.Purge(r)
		if err != nil {
			return nil, err
		}

		b.Secrets = *pSecrets
	}

	// return the purged pipeline
	return b, nil
}

// Sanitize cleans the fields for every step in each stage so they
// can be safely executed on the worker. If no stages are provided
// in the pipeline, then the function will sanitize the fields for
// every step in the pipeline. The fields are sanitized based off
// of the provided runtime driver which is setup on every worker.
// Currently, this function supports the following runtimes:
//
//   - Docker
//   - Kubernetes
func (b *Build) Sanitize(driver string) *Build {
	// return an empty pipeline if both stages and steps are provided
	if len(b.Stages) > 0 && len(b.Steps) > 0 {
		return nil
	}

	// sanitize stages pipeline if they are provided
	if len(b.Stages) > 0 {
		b.Stages = *b.Stages.Sanitize(driver)
	}

	// sanitize steps pipeline if they are provided
	if len(b.Steps) > 0 {
		b.Steps = *b.Steps.Sanitize(driver)
	}

	// sanitize services pipeline if they are provided
	if len(b.Services) > 0 {
		b.Services = *b.Services.Sanitize(driver)
	}

	// sanitize secret plugins pipeline if they are provided
	for i, secret := range b.Secrets {
		if secret.Origin.Empty() {
			continue
		}

		b.Secrets[i].Origin = secret.Origin.Sanitize(driver)
	}

	switch driver {
	// sanitize pipeline for Docker
	case constants.DriverDocker:
		if strings.Contains(b.ID, " ") {
			b.ID = strings.ReplaceAll(b.ID, " ", "-")
		}
	// sanitize pipeline for Kubernetes
	case constants.DriverKubernetes:
		if strings.Contains(b.ID, " ") {
			b.ID = strings.ReplaceAll(b.ID, " ", "-")
		}

		if strings.Contains(b.ID, "_") {
			b.ID = strings.ReplaceAll(b.ID, "_", "-")
		}

		if strings.Contains(b.ID, ".") {
			b.ID = strings.ReplaceAll(b.ID, ".", "-")
		}

		// Kubernetes requires DNS compatible names (lowercase, <= 63 chars)
		b.ID = strings.ToLower(b.ID)

		const dnsMaxLength = 63
		if utf8.RuneCountInString(b.ID) > dnsMaxLength {
			const randomSuffixLength = 6

			rs := []rune(b.ID)
			b.ID = fmt.Sprintf(
				"%s-%s",
				string(rs[:dnsMaxLength-1-randomSuffixLength]),
				dnsSafeRandomString(randomSuffixLength),
			)
		}
	}

	// return the purged pipeline
	return b
}

// Prepare sets up the pipeline for execution by populating
// the container ID a directory values based on repo and build data.
// It is used by the worker after deserializing the executable.
func (p *Build) Prepare(org, name string, number int64, local bool) {
	// check if the compiler is setup for a local pipeline
	if local {
		// check if the org provided is empty
		if len(org) == 0 {
			// set a default for the org
			org = "localOrg"
		}

		// check if the repo provided is empty
		if len(name) == 0 {
			// set a default for the repo
			name = "localRepo"
		}

		// check if the number provided is empty
		if number == 0 {
			// set a default for the number
			number = 1
		}
	}

	p.ID = fmt.Sprintf(constants.PipelineIDPattern, org, name, number)

	// set the unique ID for each step in each stage of the executable pipeline
	for _, stage := range p.Stages {
		for _, step := range stage.Steps {
			// create pattern for steps
			pattern := fmt.Sprintf(
				constants.StageIDPattern,
				org,
				name,
				number,
				stage.Name,
				step.Name,
			)

			// set id to the pattern
			step.ID = pattern

			// set the workspace directory
			step.Directory = step.Environment["VELA_WORKSPACE"]
		}
	}

	for _, step := range p.Steps {
		// create pattern for steps
		pattern := fmt.Sprintf(
			constants.StepIDPattern,
			org,
			name,
			number,
			step.Name,
		)

		// set id to the pattern
		step.ID = pattern

		// set the workspace directory
		step.Directory = step.Environment["VELA_WORKSPACE"]
	}

	// set the unique ID for each service in the executable pipeline
	for _, service := range p.Services {
		// create pattern for services
		pattern := fmt.Sprintf(
			constants.ServiceIDPattern,
			org,
			name,
			number,
			service.Name,
		)

		// set id to the pattern
		service.ID = pattern
	}

	// set the unique ID for each secret in the executable pipeline
	for _, secret := range p.Secrets {
		// skip non plugin secrets
		if secret.Origin.Empty() {
			continue
		}

		// create pattern for secrets
		pattern := fmt.Sprintf(
			constants.SecretIDPattern,
			org,
			name,
			number,
			secret.Origin.Name,
		)

		// set id to the pattern
		secret.Origin.ID = pattern
	}
}
