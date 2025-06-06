// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"context"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestExecutor_FromContext(t *testing.T) {
	// setup types
	_service, _ := New(context.Background(), &Setup{})

	// setup tests
	tests := []struct {
		context context.Context
		want    Service
	}{
		{
			//nolint:revive // ignore using string with context value
			context: context.WithValue(context.Background(), key, _service),
			want:    _service,
		},
		{
			context: context.Background(),
			want:    nil,
		},
		{
			//nolint:revive // ignore using string with context value
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
	_service, _ := New(context.Background(), &Setup{})

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

func TestExecutor_WithGinContext(t *testing.T) {
	// setup types
	_service, _ := New(context.Background(), &Setup{})

	want := new(gin.Context)
	want.Set(key, _service)

	// run test
	got := new(gin.Context)
	WithGinContext(got, _service)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("WithGinContext is %v, want %v", got, want)
	}
}
