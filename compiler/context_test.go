// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiler

import (
	"context"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCompiler_FromContext(t *testing.T) {
	// setup types
	var _engine Engine

	// setup tests
	tests := []struct {
		context context.Context
		want    Engine
	}{
		{
			//nolint: staticcheck, revive // ignore using string with context value
			context: context.WithValue(context.Background(), key, _engine),
			want:    _engine,
		},
		{
			context: context.Background(),
			want:    nil,
		},
		{
			//nolint: staticcheck, revive // ignore using string with context value
			context: context.WithValue(context.Background(), key, "foo"),
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := FromContext(test.context)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("FromContext is %v, want %v", got, test.want)
		}
	}
}

func TestCompiler_FromGinContext(t *testing.T) {
	// setup types
	var _engine Engine

	// setup tests
	tests := []struct {
		context *gin.Context
		value   interface{}
		want    Engine
	}{
		{
			context: new(gin.Context),
			value:   _engine,
			want:    _engine,
		},
		{
			context: new(gin.Context),
			value:   nil,
			want:    nil,
		},
		{
			context: new(gin.Context),
			value:   "foo",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.value != nil {
			test.context.Set(key, test.value)
		}

		got := FromGinContext(test.context)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("FromGinContext is %v, want %v", got, test.want)
		}
	}
}

func TestCompiler_WithContext(t *testing.T) {
	// setup types
	var _engine Engine

	//nolint: staticcheck, revive // ignore using string with context value
	want := context.WithValue(context.Background(), key, _engine)

	// run test
	got := WithContext(context.Background(), _engine)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("WithContext is %v, want %v", got, want)
	}
}

func TestCompiler_WithGinContext(t *testing.T) {
	// setup types
	var _engine Engine

	want := new(gin.Context)
	want.Set(key, _engine)

	// run test
	got := new(gin.Context)
	WithGinContext(got, _engine)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("WithGinContext is %v, want %v", got, want)
	}
}
