// SPDX-License-Identifier: Apache-2.0

package types

import (
	"reflect"
	"testing"
)

func TestTypes_ToString(t *testing.T) {
	// setup tests
	tests := []struct {
		parameter interface{}
		want      interface{}
	}{
		{parameter: "string", want: "string"},  // string
		{parameter: true, want: "true"},        // bool
		{parameter: []byte{1}, want: "AQ=="},   // []byte
		{parameter: float32(1.1), want: "1.1"}, // float32
		{parameter: float64(1.1), want: "1.1"}, // float64
		{parameter: 1, want: "1"},              // int
		{parameter: int8(1), want: "1"},        // int8
		{parameter: int16(1), want: "1"},       // int16
		{parameter: int32(1), want: "1"},       // int32
		{parameter: int64(1), want: "1"},       // int64
		{parameter: uint(1), want: "1"},        // uint
		{parameter: uint8(1), want: "1"},       // uint8
		{parameter: uint16(1), want: "1"},      // uint16
		{parameter: uint32(1), want: "1"},      // uint32
		{parameter: uint64(1), want: "1"},      // uint64
		{ // map
			parameter: map[string]string{"hello": "world"},
			want:      "{\"hello\":\"world\"}",
		},
		{ // slice
			parameter: []interface{}{1, 2, 3},
			want:      "1,2,3",
		},
		{ // slice complex
			parameter: []interface{}{struct{ Foo string }{Foo: "bar"}},
			want:      "[{\"foo\":\"bar\"}]",
		},
		{ // complex
			parameter: []struct{ Foo string }{{"bar"}, {"baz"}},
			want:      "[{\"foo\":\"bar\"},{\"foo\":\"baz\"}]",
		},
	}

	// run tests
	for _, test := range tests {
		got := ToString(test.parameter)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToString is %v, want %v", got, test.want)
		}
	}
}
