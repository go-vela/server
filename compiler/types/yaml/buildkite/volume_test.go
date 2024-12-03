// SPDX-License-Identifier: Apache-2.0

package buildkite

import (
	"os"
	"reflect"
	"testing"

	"github.com/buildkite/yaml"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestYaml_VolumeSlice_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		volumes *VolumeSlice
		want    *pipeline.VolumeSlice
	}{
		{
			volumes: &VolumeSlice{
				{
					Source:      "/foo",
					Destination: "/bar",
					AccessMode:  "ro",
				},
			},
			want: &pipeline.VolumeSlice{
				{
					Source:      "/foo",
					Destination: "/bar",
					AccessMode:  "ro",
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.volumes.ToPipeline()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_VolumeSlice_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *VolumeSlice
	}{
		{
			failure: false,
			file:    "testdata/volume_slice.yml",
			want: &VolumeSlice{
				{
					Source:      "/foo",
					Destination: "/foo",
					AccessMode:  "ro",
				},
				{
					Source:      "/foo",
					Destination: "/bar",
					AccessMode:  "ro",
				},
				{
					Source:      "/foo",
					Destination: "/foobar",
					AccessMode:  "ro",
				},
			},
		},
		{
			failure: false,
			file:    "testdata/volume_string.yml",
			want: &VolumeSlice{
				{
					Source:      "/foo",
					Destination: "/foo",
					AccessMode:  "ro",
				},
				{
					Source:      "/foo",
					Destination: "/bar",
					AccessMode:  "ro",
				},
				{
					Source:      "/foo",
					Destination: "/foobar",
					AccessMode:  "ro",
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
			file:    "testdata/volume_error.yml",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(VolumeSlice)

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
