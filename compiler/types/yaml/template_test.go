// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"os"
	"reflect"
	"testing"

	"github.com/buildkite/yaml"

	"github.com/go-vela/types/library"
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

func TestYAML_Template_ToLibrary(t *testing.T) {
	// setup types
	want := new(library.Template)
	want.SetName("docker_build")
	want.SetSource("github.com/go-vela/atlas/stable/docker_build")
	want.SetType("github")

	tmpl := &Template{
		Name:   "docker_build",
		Source: "github.com/go-vela/atlas/stable/docker_build",
		Type:   "github",
	}

	// run test
	got := tmpl.ToLibrary()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToLibrary is %v, want %v", got, want)
	}
}

func TestYAML_TemplateFromLibrary(t *testing.T) {
	// setup types
	want := &Template{
		Name:   "docker_build",
		Source: "github.com/go-vela/atlas/stable/docker_build",
		Type:   "github",
	}

	tmpl := new(library.Template)
	tmpl.SetName("docker_build")
	tmpl.SetSource("github.com/go-vela/atlas/stable/docker_build")
	tmpl.SetType("github")

	// run test
	got := TemplateFromLibrary(tmpl)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TemplateFromLibrary is %v, want %v", got, want)
	}
}
