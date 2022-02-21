// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"reflect"
	"testing"

	"go.starlark.net/starlark"
)

func TestStarlark_toStarlark(t *testing.T) {
	dict := starlark.NewDict(16)

	err := dict.SetKey(starlark.String("foo"), starlark.String("bar"))
	if err != nil {
		t.Error(err)
	}

	a := make([]starlark.Value, 0)
	a = append(a, starlark.Value(starlark.String("foo")))
	a = append(a, starlark.Value(starlark.String("bar")))

	type args struct {
		value interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    starlark.Value
		wantErr bool
	}{
		{"string", args{value: "foo"}, starlark.String("foo"), false},
		{"byte array", args{value: []byte("array")}, starlark.String("array"), false},
		{"array", args{value: []string{"foo", "bar"}}, starlark.Tuple(a), false},
		{"bool", args{value: true}, starlark.Bool(true), false},
		{"float64", args{value: 0.1}, starlark.Float(0.1), false},
		{"float32", args{value: float32(0.1)}, starlark.Float(float32(0.1)), false},
		{"int", args{value: 1}, starlark.MakeInt(1), false},
		{"int32", args{value: int32(1)}, starlark.MakeInt(1), false},
		{"int64", args{value: int64(1)}, starlark.MakeInt(1), false},
		{"int16", args{value: int16(1)}, starlark.MakeInt(1), false},
		{"unit", args{value: uint(1)}, starlark.MakeInt(1), false},
		{"unit32", args{value: uint32(1)}, starlark.MakeInt(1), false},
		{"unit64", args{value: uint64(1)}, starlark.MakeInt(1), false},
		{"unit16", args{value: uint16(1)}, starlark.MakeInt(1), false},
		{"nil", args{value: nil}, starlark.None, false},
		{"map", args{value: map[string]string{"foo": "bar"}}, dict, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toStarlark(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("toStarlark() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toStarlark() got = %v, want %v", got, tt.want)
			}
		})
	}
}
