// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-vela/types/constants"
)

type (
	// SecretSlice is the pipeline representation
	// of the secrets block for a pipeline.
	//
	// swagger:model PipelineSecretSlice
	SecretSlice []*Secret

	// Secret is the pipeline representation of a
	// secret from the secrets block for a pipeline.
	//
	// swagger:model PipelineSecret
	Secret struct {
		Name   string     `json:"name,omitempty"   yaml:"name,omitempty"`
		Value  string     `json:"value,omitempty"  yaml:"value,omitempty"`
		Key    string     `json:"key,omitempty"    yaml:"key,omitempty"`
		Engine string     `json:"engine,omitempty" yaml:"engine,omitempty"`
		Type   string     `json:"type,omitempty"   yaml:"type,omitempty"`
		Origin *Container `json:"origin,omitempty" yaml:"origin,omitempty"`
		Pull   string     `json:"pull,omitempty"   yaml:"pull,omitempty"`
	}

	// StepSecretSlice is the pipeline representation
	// of the secrets block for a step in a pipeline.
	//
	// swagger:model PipelineStepSecretSlice
	StepSecretSlice []*StepSecret

	// StepSecret is the pipeline representation of a secret
	// from a secrets block for a step in a pipeline.
	//
	// swagger:model PipelineStepSecret
	StepSecret struct {
		Source string `json:"source,omitempty" yaml:"source,omitempty"`
		Target string `json:"target,omitempty" yaml:"target,omitempty"`
	}
)

var (
	// ErrInvalidEngine defines the error type when the
	// SecretEngine provided to the client is unsupported.
	ErrInvalidEngine = errors.New("invalid secret engine")
	// ErrInvalidOrg defines the error type when the
	// org in key does not equal the name of the organization.
	ErrInvalidOrg = errors.New("invalid organization in key")
	// ErrInvalidRepo defines the error type when the
	// repo in key does not equal the name of the repository.
	ErrInvalidRepo = errors.New("invalid repository in key")
	// ErrInvalidShared defines the error type when the
	// org in key does not equal the name of the team.
	ErrInvalidShared = errors.New("invalid team in key")
	// ErrInvalidPath defines the error type when the
	// path provided for a type (org, repo, shared) is invalid.
	ErrInvalidPath = errors.New("invalid secret path")
	// ErrInvalidName defines the error type when the name
	// contains restricted characters or is empty.
	ErrInvalidName = errors.New("invalid secret name")
)

// Purge removes the secrets that have a ruleset
// that do not match the provided ruledata.
func (s *SecretSlice) Purge(r *RuleData) (*SecretSlice, error) {
	counter := 1
	secrets := new(SecretSlice)

	// iterate through each Secret in the pipeline
	for _, secret := range *s {
		if secret.Origin.Empty() {
			// append the secret to the new slice of secrets
			*secrets = append(*secrets, secret)

			continue
		}

		match, err := secret.Origin.Ruleset.Match(r)
		if err != nil {
			return nil, fmt.Errorf("unable to process ruleset for secret %s: %w", secret.Name, err)
		}

		// verify ruleset matches
		if match {
			// overwrite the Container number with the Container counter
			secret.Origin.Number = counter

			// increment Container counter
			counter = counter + 1

			// append the secret to the new slice of secrets
			*secrets = append(*secrets, secret)
		}
	}

	return secrets, nil
}

// ParseOrg returns the parts (org, key) of the secret path
// when the secret is valid for a given organization.
func (s *Secret) ParseOrg(org string) (string, string, error) {
	path := s.Key

	// check if the secret is not a native or vault type
	if !strings.EqualFold(s.Engine, constants.DriverNative) &&
		!strings.EqualFold(s.Engine, constants.DriverVault) {
		return "", "", fmt.Errorf("%w: %s", ErrInvalidEngine, s.Engine)
	}

	// check if a path was provided
	if !strings.Contains(path, "/") {
		return "", "", fmt.Errorf("%w: %s ", ErrInvalidPath, path)
	}

	// split the full path into parts
	parts := strings.SplitN(path, "/", 2)

	// secret is invalid
	if len(parts) != 2 {
		return "", "", fmt.Errorf("%w: %s ", ErrInvalidPath, path)
	}

	// check if the org provided matches what we expect
	if !strings.EqualFold(parts[0], org) {
		return "", "", fmt.Errorf("%w: %s ", ErrInvalidOrg, parts[0])
	}

	// check if path segments empty
	if len(parts[1]) == 0 {
		return "", "", fmt.Errorf("%w: %s ", ErrInvalidPath, path)
	}

	// secret names can't be empty.
	if len(s.Name) == 0 {
		return "", "", fmt.Errorf("%w: %s missing name", ErrInvalidName, s.Key)
	}

	// environmental variables can't contain certain restricted characters.
	if strings.ContainsAny(s.Name, constants.SecretRestrictedCharacters) {
		return "", "", fmt.Errorf("%w (contains restricted characters): %s ", ErrInvalidName, s.Name)
	}

	return parts[0], parts[1], nil
}

// ParseRepo returns the parts (org, repo, key) of the secret path
// when the secret is valid for a given organization and repository.
func (s *Secret) ParseRepo(org, repo string) (string, string, string, error) {
	path := s.Key

	// check if the secret is not a native or vault type
	if !strings.EqualFold(s.Engine, constants.DriverNative) && !strings.EqualFold(s.Engine, constants.DriverVault) {
		return "", "", "", fmt.Errorf("%w: %s", ErrInvalidEngine, s.Engine)
	}

	// split the full path into parts
	parts := strings.SplitN(path, "/", 3)

	// secret is invalid
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("%w: %s ", ErrInvalidPath, path)
	}

	// check if the org provided matches what we expect
	if !strings.EqualFold(parts[0], org) {
		return "", "", "", fmt.Errorf("%w: %s ", ErrInvalidOrg, parts[0])
	}

	// check if the repo provided matches what we expect
	if !strings.EqualFold(parts[1], repo) {
		return "", "", "", fmt.Errorf("%w: %s ", ErrInvalidRepo, parts[1])
	}

	// check if path segments empty
	if len(parts[2]) == 0 {
		return "", "", "", fmt.Errorf("%w: %s ", ErrInvalidPath, path)
	}

	// secret names can't be empty.
	if len(s.Name) == 0 {
		return "", "", "", fmt.Errorf("%w: %s missing name", ErrInvalidName, s.Key)
	}

	// environmental variables can't contain certain restricted characters.
	if strings.ContainsAny(s.Name, constants.SecretRestrictedCharacters) {
		return "", "", "", fmt.Errorf("%w (contains restricted characters): %s ", ErrInvalidName, s.Name)
	}

	return parts[0], parts[1], parts[2], nil
}

// ParseShared returns the parts (org, team, key) of the secret path
// when the secret is valid for a given organization and team.
func (s *Secret) ParseShared() (string, string, string, error) {
	path := s.Key

	// check if the secret is not a native or vault type
	if !strings.EqualFold(s.Engine, constants.DriverNative) && !strings.EqualFold(s.Engine, constants.DriverVault) {
		return "", "", "", fmt.Errorf("%w: %s", ErrInvalidEngine, s.Engine)
	}

	// check if a path was provided
	if !strings.Contains(path, "/") {
		return "", "", "", fmt.Errorf("%w: %s ", ErrInvalidPath, path)
	}

	// split the full path into parts
	parts := strings.SplitN(path, "/", 3)

	// secret is invalid
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("%w: %s ", ErrInvalidPath, path)
	}

	// check if path segments empty
	if len(parts[1]) == 0 || len(parts[2]) == 0 {
		return "", "", "", fmt.Errorf("%w: %s ", ErrInvalidPath, path)
	}

	// secret names can't be empty.
	if len(s.Name) == 0 {
		return "", "", "", fmt.Errorf("%w: %s missing name", ErrInvalidName, s.Key)
	}

	// environmental variables can't contain certain restricted characters.
	if strings.ContainsAny(s.Name, constants.SecretRestrictedCharacters) {
		return "", "", "", fmt.Errorf("%w (contains restricted characters): %s ", ErrInvalidName, s.Name)
	}

	return parts[0], parts[1], parts[2], nil
}
