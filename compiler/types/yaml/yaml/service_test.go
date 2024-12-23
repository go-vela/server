// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
)

func TestYaml_ServiceSlice_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		services *ServiceSlice
		want     *pipeline.ContainerSlice
	}{
		{
			services: &ServiceSlice{
				{
					Entrypoint:  []string{"/usr/local/bin/docker-entrypoint.sh"},
					Environment: map[string]string{"FOO": "bar"},
					Image:       "postgres:12-alpine",
					Name:        "postgres",
					Ports:       []string{"5432:5432"},
					Pull:        "not_present",
				},
			},
			want: &pipeline.ContainerSlice{
				{
					Detach:      true,
					Entrypoint:  []string{"/usr/local/bin/docker-entrypoint.sh"},
					Environment: map[string]string{"FOO": "bar"},
					Image:       "postgres:12-alpine",
					Name:        "postgres",
					Ports:       []string{"5432:5432"},
					Pull:        "not_present",
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.services.ToPipeline()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_ServiceSlice_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *ServiceSlice
	}{
		{
			failure: false,
			file:    "testdata/service.yml",
			want: &ServiceSlice{
				{
					Environment: raw.StringSliceMap{
						"POSTGRES_DB": "foo",
					},
					Image: "postgres:latest",
					Name:  "postgres",
					Ports: []string{"5432:5432"},
					Pull:  "not_present",
				},
				{
					Environment: raw.StringSliceMap{
						"MYSQL_DATABASE": "foo",
					},
					Image: "mysql:latest",
					Name:  "mysql",
					Ports: []string{"3061:3061"},
					Pull:  "not_present",
				},
			},
		},
		{
			failure: true,
			file:    "testdata/invalid.yml",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/service_nil.yml",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(ServiceSlice)

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

func TestYaml_Service_MergeEnv(t *testing.T) {
	// setup tests
	tests := []struct {
		service     *Service
		environment map[string]string
		failure     bool
	}{
		{
			service: &Service{
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:latest",
				Name:        "postgres",
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			service:     &Service{},
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			service:     nil,
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			service: &Service{
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:latest",
				Name:        "postgres",
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
			environment: nil,
			failure:     true,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.service.MergeEnv(test.environment)

		if test.failure {
			if err == nil {
				t.Errorf("MergeEnv should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("MergeEnv returned err: %v", err)
		}
	}
}
