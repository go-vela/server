// SPDX-License-Identifier: Apache-2.0

package buildkite

import (
	"os"
	"reflect"
	"testing"

	"github.com/buildkite/yaml"

	api "github.com/go-vela/server/api/types"
)

func TestBuild_TemplateSlice_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *TemplateSlice
	}{
		{
			failure: false,
			file:    "testdata/template.yml",
			want: &TemplateSlice{
				{
					Name:   "docker_build",
					Source: "github.com/go-vela/atlas/stable/docker_create",
					Type:   "github",
				},
				{
					Name:   "docker_build",
					Source: "github.com/go-vela/atlas/stable/docker_build",
					Format: "go",
					Type:   "github",
				},
				{
					Name:   "docker_publish",
					Source: "github.com/go-vela/atlas/stable/docker_publish",
					Format: "starlark",
					Type:   "github",
				},
			},
		},
		{
			failure: true,
			file:    "testdata/invalid.yml",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(TemplateSlice)

		b, err := os.ReadFile(test.file)
		if err != nil {
			t.Errorf("unable to read file: %v", err)
		}

		err = yaml.Unmarshal(b, got)

		if test.failure {
			if err == nil {
				t.Errorf("UnmarshalYAML should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UnmarshalYAML returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("UnmarshalYAML is %v, want %v", got, test.want)
		}
	}
}

func TestYAML_Template_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Template)
	want.SetName("docker_build")
	want.SetSource("github.com/go-vela/atlas/stable/docker_build")
	want.SetType("github")

	tmpl := &Template{
		Name:   "docker_build",
		Source: "github.com/go-vela/atlas/stable/docker_build",
		Type:   "github",
	}

	// run test
	got := tmpl.ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestYAML_TemplateFromAPI(t *testing.T) {
	// setup types
	want := &Template{
		Name:   "docker_build",
		Source: "github.com/go-vela/atlas/stable/docker_build",
		Type:   "github",
	}

	tmpl := new(api.Template)
	tmpl.SetName("docker_build")
	tmpl.SetSource("github.com/go-vela/atlas/stable/docker_build")
	tmpl.SetType("github")

	// run test
	got := TemplateFromAPI(tmpl)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TemplateFromAPI is %v, want %v", got, want)
	}
}
