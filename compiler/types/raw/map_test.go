// SPDX-License-Identifier: Apache-2.0

package raw

import (
	"database/sql/driver"
	"os"
	"reflect"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestRaw_StringSliceMap_UnmarshalJSON(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *StringSliceMap
	}{
		{
			failure: false,
			file:    "testdata/string_map.json",
			want:    &StringSliceMap{"foo": "bar"},
		},
		{
			failure: false,
			file:    "testdata/slice_map.json",
			want:    &StringSliceMap{"foo": "bar"},
		},
		{
			failure: false,
			file:    "testdata/map.json",
			want:    &StringSliceMap{"foo": "bar"},
		},
		{
			failure: false,
			file:    "",
			want:    new(StringSliceMap),
		},
		{
			failure: true,
			file:    "testdata/invalid.json",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/invalid_2.json",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		var (
			err error

			b   = []byte{}
			got = new(StringSliceMap)
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

func TestRaw_StringSliceMap_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *StringSliceMap
	}{
		{
			failure: false,
			file:    "testdata/string_map.yml",
			want:    &StringSliceMap{"foo": "bar"},
		},
		{
			failure: false,
			file:    "testdata/slice_map.yml",
			want:    &StringSliceMap{"foo": "bar"},
		},
		{
			failure: false,
			file:    "testdata/map.yml",
			want:    &StringSliceMap{"foo": "bar"},
		},
		{
			failure: true,
			file:    "testdata/invalid.yml",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/invalid_2.yml",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(StringSliceMap)

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

func TestStringSliceMap_Value(t *testing.T) {
	tests := []struct {
		name    string
		s       StringSliceMap
		want    driver.Value
		wantErr bool
	}{
		{"valid", StringSliceMap{"foo": "test1"}, "{\"foo\":\"test1\"}", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("StringSliceMap.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSliceMap.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceMap_Scan(t *testing.T) {
	type args struct {
		value interface{}
	}

	tests := []struct {
		name    string
		s       *StringSliceMap
		args    args
		wantErr bool
	}{
		{"valid", &StringSliceMap{"foo": "test1"}, args{value: "{\"foo\":\"test1\"}"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("StringSliceMap.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
