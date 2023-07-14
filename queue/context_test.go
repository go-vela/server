// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

import (
	"context"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestExecutor_FromContext(t *testing.T) {
	// setup types
	_service, _ := New(&Setup{})

	// setup tests
	tests := []struct {
		context context.Context
		want    Service
	}{
		{
			//nolint:staticcheck,revive // ignore using string with context value
			context: context.WithValue(context.Background(), key, _service),
			want:    _service,
		},
		{
			context: context.Background(),
			want:    nil,
		},
		{
			//nolint:staticcheck,revive // ignore using string with context value
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

func TestExecutor_FromGinContext(t *testing.T) {
	// setup types
	_service, _ := New(&Setup{})

	// setup tests
	tests := []struct {
		context *gin.Context
		value   interface{}
		want    Service
	}{
		{
			context: new(gin.Context),
			value:   _service,
			want:    _service,
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

func TestExecutor_WithContext(t *testing.T) {
	// setup types
	_service, _ := New(&Setup{})

	//nolint:staticcheck,revive // ignore using string with context value
	want := context.WithValue(context.Background(), key, _service)

	// run test
	got := WithContext(context.Background(), _service)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("WithContext is %v, want %v", got, want)
	}
}

func TestExecutor_WithGinContext(t *testing.T) {
	// setup types
	_service, _ := New(&Setup{})

	want := new(gin.Context)
	want.Set(key, _service)

	// run test
	got := new(gin.Context)
	WithGinContext(got, _service)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("WithGinContext is %v, want %v", got, want)
	}
}
