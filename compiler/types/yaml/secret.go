// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/types/constants"
)

type (
	// SecretSlice is the yaml representation
	// of the secrets block for a pipeline.
	SecretSlice []*Secret

	// Secret is the yaml representation of a secret
	// from the secrets block for a pipeline.
	Secret struct {
		Name   string `yaml:"name,omitempty"   json:"name,omitempty" jsonschema:"required,minLength=1,description=Name of secret to reference in the pipeline.\nReference: https://go-vela.github.io/docs/reference/yaml/secrets/#the-name-key"`
		Key    string `yaml:"key,omitempty"    json:"key,omitempty" jsonschema:"minLength=1,description=Path to secret to fetch from storage backend.\nReference: https://go-vela.github.io/docs/reference/yaml/secrets/#the-key-key"`
		Engine string `yaml:"engine,omitempty" json:"engine,omitempty" jsonschema:"enum=native,enum=vault,default=native,description=Name of storage backend to fetch secret from.\nReference: https://go-vela.github.io/docs/reference/yaml/secrets/#the-engine-key"`
		Type   string `yaml:"type,omitempty"   json:"type,omitempty" jsonschema:"enum=repo,enum=org,enum=shared,default=repo,description=Type of secret to fetch from storage backend.\nReference: https://go-vela.github.io/docs/reference/yaml/secrets/#the-type-key"`
		Origin Origin `yaml:"origin,omitempty" json:"origin,omitempty" jsonschema:"description=Declaration to pull secrets from non-internal secret providers.\nReference: https://go-vela.github.io/docs/reference/yaml/secrets/#the-origin-key"`
		Pull   string `yaml:"pull,omitempty"   json:"pull,omitempty" jsonschema:"enum=step_start,enum=build_start,default=build_start,description=When to pull in secrets from storage backend.\nReference: https://go-vela.github.io/docs/reference/yaml/secrets/#the-pull-key"`
	}

	// Origin is the yaml representation of a method
	// for looking up secrets with a secret plugin.
	Origin struct {
		Environment raw.StringSliceMap     `yaml:"environment,omitempty" json:"environment,omitempty" jsonschema:"description=Variables to inject into the container environment.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-environment-key"`
		Image       string                 `yaml:"image,omitempty"       json:"image,omitempty" jsonschema:"required,minLength=1,description=Docker image to use to create the ephemeral container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-image-key"`
		Name        string                 `yaml:"name,omitempty"        json:"name,omitempty" jsonschema:"required,minLength=1,description=Unique name for the secret origin.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-name-key"`
		Parameters  map[string]interface{} `yaml:"parameters,omitempty"  json:"parameters,omitempty" jsonschema:"description=Extra configuration variables for the secret plugin.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-parameters-key"`
		Secrets     StepSecretSlice        `yaml:"secrets,omitempty"     json:"secrets,omitempty" jsonschema:"description=Secrets to inject that are necessary to retrieve the secrets.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-secrets-key"`
		Pull        string                 `yaml:"pull,omitempty"        json:"pull,omitempty" jsonschema:"enum=always,enum=not_present,enum=on_start,enum=never,default=not_present,description=Declaration to configure if and when the Docker image is pulled.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-pull-key"`
		Ruleset     Ruleset                `yaml:"ruleset,omitempty"     json:"ruleset,omitempty" jsonschema:"description=Conditions to limit the execution of the container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
	}
)

// ToPipeline converts the SecretSlice type
// to a pipeline SecretSlice type.
func (s *SecretSlice) ToPipeline() *pipeline.SecretSlice {
	// secret slice we want to return
	secretSlice := new(pipeline.SecretSlice)

	// iterate through each element in the secret slice
	for _, secret := range *s {
		// append the element to the pipeline secret slice
		*secretSlice = append(*secretSlice, &pipeline.Secret{
			Name:   secret.Name,
			Key:    secret.Key,
			Engine: secret.Engine,
			Type:   secret.Type,
			Origin: secret.Origin.ToPipeline(),
			Pull:   secret.Pull,
		})
	}

	return secretSlice
}

// UnmarshalYAML implements the Unmarshaler interface for the SecretSlice type.
func (s *SecretSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// secret slice we try unmarshalling to
	secretSlice := new([]*Secret)

	// attempt to unmarshal as a secret slice type
	err := unmarshal(secretSlice)
	if err != nil {
		return err
	}

	tmp := SecretSlice{}

	// iterate through each element in the secret slice
	for _, secret := range *secretSlice {
		if secret.Origin.Empty() && len(secret.Name) == 0 {
			continue
		}

		if secret.Origin.Empty() && len(secret.Key) == 0 {
			secret.Key = secret.Name
		}

		// implicitly set `engine` field if empty
		if secret.Origin.Empty() && len(secret.Engine) == 0 {
			secret.Engine = constants.DriverNative
		}

		// implicitly set `type` field if empty
		if secret.Origin.Empty() && len(secret.Type) == 0 {
			secret.Type = constants.SecretRepo
		}

		// implicitly set `type` field if empty
		if secret.Origin.Empty() && len(secret.Pull) == 0 {
			secret.Pull = constants.SecretPullBuild
		}

		// implicitly set `pull` field if empty
		if !secret.Origin.Empty() && len(secret.Origin.Pull) == 0 {
			secret.Origin.Pull = constants.PullNotPresent
		}

		// TODO: remove this in a future release
		//
		// handle true deprecated pull policy
		//
		// a `true` pull policy equates to `always`
		if !secret.Origin.Empty() && strings.EqualFold(secret.Origin.Pull, "true") {
			secret.Origin.Pull = constants.PullAlways
		}

		// TODO: remove this in a future release
		//
		// handle false deprecated pull policy
		//
		// a `false` pull policy equates to `not_present`
		if !secret.Origin.Empty() && strings.EqualFold(secret.Origin.Pull, "false") {
			secret.Origin.Pull = constants.PullNotPresent
		}

		tmp = append(tmp, secret)
	}

	// overwrite existing SecretSlice
	*s = tmp

	return nil
}

// Empty returns true if the provided origin is empty.
func (o *Origin) Empty() bool {
	// return true if the origin is nil
	if o == nil {
		return true
	}

	// return true if every origin field is empty
	if o.Environment == nil &&
		len(o.Image) == 0 &&
		len(o.Name) == 0 &&
		o.Parameters == nil &&
		len(o.Secrets) == 0 &&
		len(o.Pull) == 0 {
		return true
	}

	return false
}

// MergeEnv takes a list of environment variables and attempts
// to set them in the secret environment. If the environment
// variable already exists in the secret, than this will
// overwrite the existing environment variable.
func (o *Origin) MergeEnv(environment map[string]string) error {
	// check if the secret container is empty
	if o.Empty() {
		// TODO: evaluate if we should error here
		//
		// immediately return and do nothing
		//
		// treated as a no-op
		return nil
	}

	// check if the environment provided is empty
	if environment == nil {
		return fmt.Errorf("empty environment provided for secret %s", o.Name)
	}

	// iterate through all environment variables provided
	for key, value := range environment {
		// set or update the secret environment variable
		o.Environment[key] = value
	}

	return nil
}

// ToPipeline converts the Origin type
// to a pipeline Container type.
func (o *Origin) ToPipeline() *pipeline.Container {
	return &pipeline.Container{
		Environment: o.Environment,
		Image:       o.Image,
		Name:        o.Name,
		Pull:        o.Pull,
		Ruleset:     *o.Ruleset.ToPipeline(),
		Secrets:     *o.Secrets.ToPipeline(),
	}
}

type (
	// StepSecretSlice is the yaml representation of
	// the secrets block for a step in a pipeline.
	StepSecretSlice []*StepSecret

	// StepSecret is the yaml representation of a secret
	// from a secrets block for a step in a pipeline.
	StepSecret struct {
		Source string `yaml:"source,omitempty"`
		Target string `yaml:"target,omitempty"`
	}
)

// ToPipeline converts the StepSecretSlice type
// to a pipeline StepSecretSlice type.
func (s *StepSecretSlice) ToPipeline() *pipeline.StepSecretSlice {
	// step secret slice we want to return
	secretSlice := new(pipeline.StepSecretSlice)

	// iterate through each element in the step secret slice
	for _, secret := range *s {
		// append the element to the pipeline step secret slice
		*secretSlice = append(*secretSlice, &pipeline.StepSecret{
			Source: secret.Source,
			Target: secret.Target,
		})
	}

	return secretSlice
}

// UnmarshalYAML implements the Unmarshaler interface for the StepSecretSlice type.
func (s *StepSecretSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// string slice we try unmarshalling to
	stringSlice := new(raw.StringSlice)

	// attempt to unmarshal as a string slice type
	err := unmarshal(stringSlice)
	if err == nil {
		// iterate through each element in the string slice
		for _, secret := range *stringSlice {
			// append the element to the step secret slice
			*s = append(*s, &StepSecret{
				Source: secret,
				Target: strings.ToUpper(secret),
			})
		}

		return nil
	}

	// step secret slice we try unmarshalling to
	secrets := new([]*StepSecret)

	// attempt to unmarshal as a step secret slice type
	err = unmarshal(secrets)
	if err == nil {
		// check for secret source and target
		for _, secret := range *secrets {
			if len(secret.Source) == 0 || len(secret.Target) == 0 {
				return fmt.Errorf("no secret source or target found")
			}

			secret.Target = strings.ToUpper(secret.Target)
		}

		// overwrite existing StepSecretSlice
		*s = StepSecretSlice(*secrets)

		return nil
	}

	return errors.New("failed to unmarshal StepSecretSlice")
}
