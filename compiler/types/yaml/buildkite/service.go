// SPDX-License-Identifier: Apache-2.0

package buildkite

import (
	"fmt"
	"maps"
	"strings"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
	"github.com/go-vela/server/constants"
)

type (
	// ServiceSlice is the yaml representation
	// of the Services block for a pipeline.
	ServiceSlice []*Service

	// Service is the yaml representation
	// of a Service in a pipeline.
	Service struct {
		Image       string             `yaml:"image,omitempty"       json:"image,omitempty"       jsonschema:"required,minLength=1,description=Docker image used to create ephemeral container.\nReference: https://go-vela.github.io/docs/reference/yaml/services/#the-image-key"`
		Name        string             `yaml:"name,omitempty"        json:"name,omitempty"        jsonschema:"required,minLength=1,description=Unique identifier for the container in the pipeline.\nReference: https://go-vela.github.io/docs/reference/yaml/services/#the-name-key"`
		Entrypoint  raw.StringSlice    `yaml:"entrypoint,omitempty"  json:"entrypoint,omitempty"  jsonschema:"description=Commands to execute inside the container.\nReference: https://go-vela.github.io/docs/reference/yaml/services/#the-entrypoint-key"`
		Environment raw.StringSliceMap `yaml:"environment,omitempty" json:"environment,omitempty" jsonschema:"description=Variables to inject into the container environment.\nReference: https://go-vela.github.io/docs/reference/yaml/services/#the-environment-key"`
		Ports       raw.StringSlice    `yaml:"ports,omitempty"       json:"ports,omitempty"       jsonschema:"description=List of ports to map for the container in the pipeline.\nReference: https://go-vela.github.io/docs/reference/yaml/services/#the-ports-key"`
		Pull        string             `yaml:"pull,omitempty"        json:"pull,omitempty"        jsonschema:"enum=always,enum=not_present,enum=on_start,enum=never,default=not_present,description=Declaration to configure if and when the Docker image is pulled.\nReference: https://go-vela.github.io/docs/reference/yaml/services/#the-pul-key"`
		Ulimits     UlimitSlice        `yaml:"ulimits,omitempty"     json:"ulimits,omitempty"     jsonschema:"description=Set the user limits for the container.\nReference: https://go-vela.github.io/docs/reference/yaml/services/#the-ulimits-key"`
		User        string             `yaml:"user,omitempty"        json:"user,omitempty"        jsonschema:"description=Set the user for the container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-user-key"`
	}
)

// ToPipeline converts the ServiceSlice type
// to a pipeline ContainerSlice type.
func (s *ServiceSlice) ToPipeline() *pipeline.ContainerSlice {
	// service slice we want to return
	serviceSlice := new(pipeline.ContainerSlice)

	// iterate through each element in the service slice
	for _, service := range *s {
		// append the element to the pipeline container slice
		*serviceSlice = append(*serviceSlice, &pipeline.Container{
			Detach:      true,
			Image:       service.Image,
			Name:        service.Name,
			Entrypoint:  service.Entrypoint,
			Environment: service.Environment,
			Ports:       service.Ports,
			Pull:        service.Pull,
			Ulimits:     *service.Ulimits.ToPipeline(),
			User:        service.User,
		})
	}

	return serviceSlice
}

// UnmarshalYAML implements the Unmarshaler interface for the ServiceSlice type.
func (s *ServiceSlice) UnmarshalYAML(unmarshal func(any) error) error {
	// service slice we try unmarshalling to
	serviceSlice := new([]*Service)

	// attempt to unmarshal as a service slice type
	err := unmarshal(serviceSlice)
	if err != nil {
		return err
	}

	// iterate through each element in the service slice
	for _, service := range *serviceSlice {
		// handle nil service to avoid panic
		if service == nil {
			return fmt.Errorf("invalid service with nil content found")
		}

		// implicitly set `pull` field if empty
		if len(service.Pull) == 0 {
			service.Pull = constants.PullNotPresent
		}

		// TODO: remove this in a future release
		//
		// handle true deprecated pull policy
		//
		// a `true` pull policy equates to `always`
		if strings.EqualFold(service.Pull, "true") {
			service.Pull = constants.PullAlways
		}

		// TODO: remove this in a future release
		//
		// handle false deprecated pull policy
		//
		// a `false` pull policy equates to `not_present`
		if strings.EqualFold(service.Pull, "false") {
			service.Pull = constants.PullNotPresent
		}
	}

	// overwrite existing ServiceSlice
	*s = ServiceSlice(*serviceSlice)

	return nil
}

// MergeEnv takes a list of environment variables and attempts
// to set them in the service environment. If the environment
// variable already exists in the service, than this will
// overwrite the existing environment variable.
func (s *Service) MergeEnv(environment map[string]string) error {
	// check if the service container is empty
	if s == nil || s.Environment == nil {
		// TODO: evaluate if we should error here
		//
		// immediately return and do nothing
		//
		// treated as a no-op
		return nil
	}

	// check if the environment provided is empty
	if environment == nil {
		return fmt.Errorf("empty environment provided for service %s", s.Name)
	}

	// apply environment to service environment
	maps.Copy(s.Environment, environment)

	return nil
}

func (s *Service) ToYAML() *yaml.Service {
	if s == nil {
		return nil
	}

	return &yaml.Service{
		Image:       s.Image,
		Name:        s.Name,
		Entrypoint:  s.Entrypoint,
		Environment: s.Environment,
		Ports:       s.Ports,
		Pull:        s.Pull,
		Ulimits:     *s.Ulimits.ToYAML(),
		User:        s.User,
	}
}

func (s *ServiceSlice) ToYAML() *yaml.ServiceSlice {
	// service slice we want to return
	serviceSlice := new(yaml.ServiceSlice)

	// iterate through each element in the service slice
	for _, service := range *s {
		// append the element to the yaml service slice
		*serviceSlice = append(*serviceSlice, service.ToYAML())
	}

	return serviceSlice
}
