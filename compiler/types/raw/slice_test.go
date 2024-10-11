// SPDX-License-Identifier: Apache-2.0

package raw

import (
	"os"
	"reflect"
	"testing"

	"github.com/buildkite/yaml"
)

func TestRaw_StringSlice_UnmarshalJSON(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *StringSlice
	}{
		{
			failure: false,
			file:    "testdata/string.json",
			want:    &StringSlice{"foo"},
		},
		{
			failure: false,
			file:    "testdata/slice.json",
			want:    &StringSlice{"foo", "bar"},
		},
		{
			failure: false,
			file:    "",
			want:    new(StringSlice),
		},
		{
			failure: true,
			file:    "testdata/invalid.json",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		var (
			err error

			b   = []byte{}
			got = new(StringSlice)
		)

		if len(test.file) > 0 {
			b, err = os.ReadFile(test.file)
			if err != nil {
				t.Errorf("unable to read %s file: %v", test.file, err)
			}
		}

		err = got.UnmarshalJSON(b)

		if test.failure {
			if err == nil {
				t.Errorf("UnmarshalJSON should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UnmarshalJSON returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("UnmarshalJSON is %v, want %v", got, test.want)
		}
	}
}

func TestRaw_StringSlice_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *StringSlice
	}{
		{
			failure: false,
			file:    "testdata/string.yml",
			want:    &StringSlice{"foo"},
		},
		{
			failure: false,
			file:    "testdata/slice.yml",
			want:    &StringSlice{"foo", "bar"},
		},
		{
			failure: true,
			file:    "testdata/invalid.yml",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(StringSlice)

		b, err := os.ReadFile(test.file)
		if err != nil {
			t.Errorf("unable to read %s file: %v", test.file, err)
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
