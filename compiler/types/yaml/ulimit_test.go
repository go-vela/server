// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"os"
	"reflect"
	"testing"

	"github.com/buildkite/yaml"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestYaml_UlimitSlice_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		ulimits *UlimitSlice
		want    *pipeline.UlimitSlice
	}{
		{
			ulimits: &UlimitSlice{
				{
					Name: "foo",
					Soft: 1024,
					Hard: 2048,
				},
			},
			want: &pipeline.UlimitSlice{
				{
					Name: "foo",
					Soft: 1024,
					Hard: 2048,
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.ulimits.ToPipeline()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_UlimitSlice_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *UlimitSlice
	}{
		{
			failure: false,
			file:    "testdata/ulimit_slice.yml",
			want: &UlimitSlice{
				{
					Name: "foo",
					Soft: 1024,
					Hard: 1024,
				},
				{
					Name: "bar",
					Soft: 1024,
					Hard: 2048,
				},
			},
		},
		{
			failure: false,
			file:    "testdata/ulimit_string.yml",
			want: &UlimitSlice{
				{
					Name: "foo",
					Soft: 1024,
					Hard: 1024,
				},
				{
					Name: "bar",
					Soft: 1024,
					Hard: 2048,
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
			file:    "testdata/ulimit_equal_error.yml",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/ulimit_colon_error.yml",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/ulimit_softlimit_error.yml",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/ulimit_hardlimit1_error.yml",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/ulimit_hardlimit2_error.yml",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(UlimitSlice)

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
