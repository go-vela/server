// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	api "github.com/go-vela/server/api/types"
)

type (
	// TemplateSlice is the yaml representation
	// of the templates block for a pipeline.
	TemplateSlice []*Template

	// Template is the yaml representation of a template
	// from the templates block for a pipeline.
	Template struct {
		Name      string                 `yaml:"name,omitempty"   json:"name,omitempty"  jsonschema:"required,minLength=1,description=Unique identifier for the template.\nReference: https://go-vela.github.io/docs/reference/yaml/templates/#the-name-key"`
		Source    string                 `yaml:"source,omitempty" json:"source,omitempty" jsonschema:"required,minLength=1,description=Path to template in remote system.\nReference: https://go-vela.github.io/docs/reference/yaml/templates/#the-source-key"`
		Format    string                 `yaml:"format,omitempty" json:"format,omitempty" jsonschema:"enum=starlark,enum=golang,enum=go,default=go,minLength=1,description=language used within the template file \nReference: https://go-vela.github.io/docs/reference/yaml/templates/#the-format-key"`
		Type      string                 `yaml:"type,omitempty"   json:"type,omitempty" jsonschema:"minLength=1,example=github,description=Type of template provided from the remote system.\nReference: https://go-vela.github.io/docs/reference/yaml/templates/#the-type-key"`
		Variables map[string]interface{} `yaml:"vars,omitempty"   json:"vars,omitempty" jsonschema:"description=Variables injected into the template.\nReference: https://go-vela.github.io/docs/reference/yaml/templates/#the-variables-key"`
	}

	// StepTemplate is the yaml representation of the
	// template block for a step in a pipeline.
	StepTemplate struct {
		Name      string                 `yaml:"name,omitempty" json:"name,omitempty" jsonschema:"required,minLength=1,description=Unique identifier for the template.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-template-key"`
		Variables map[string]interface{} `yaml:"vars,omitempty" json:"vars,omitempty" jsonschema:"description=Variables injected into the template.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-template-key"`
	}
)

// UnmarshalYAML implements the Unmarshaler interface for the TemplateSlice type.
func (t *TemplateSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// template slice we try unmarshalling to
	templateSlice := new([]*Template)

	// attempt to unmarshal as a template slice type
	err := unmarshal(templateSlice)
	if err != nil {
		return err
	}

	// overwrite existing TemplateSlice
	*t = TemplateSlice(*templateSlice)

	return nil
}

// ToAPI converts the Template type
// to an API Template type.
func (t *Template) ToAPI() *api.Template {
	template := new(api.Template)

	template.SetName(t.Name)
	template.SetSource(t.Source)
	template.SetType(t.Type)

	return template
}

// TemplateFromAPI converts the API Template type
// to a yaml Template type.
func TemplateFromAPI(t *api.Template) *Template {
	template := &Template{
		Name:   t.GetName(),
		Source: t.GetSource(),
		Type:   t.GetType(),
	}

	return template
}

// Map helper function that creates a map of templates from a slice of templates.
func (t *TemplateSlice) Map() map[string]*Template {
	m := make(map[string]*Template)

	if t == nil {
		return m
	}

	for _, tmpl := range *t {
		m[tmpl.Name] = tmpl
	}

	return m
}
