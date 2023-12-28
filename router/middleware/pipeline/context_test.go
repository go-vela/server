// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types/library"
)

func TestPipeline_FromContext(t *testing.T) {
	// setup types
	_pipeline := new(library.Pipeline)

	gin.SetMode(gin.TestMode)
	_context, _ := gin.CreateTestContext(nil)
	_context.Set(key, _pipeline)

	_emptyContext, _ := gin.CreateTestContext(nil)

	_nilContext, _ := gin.CreateTestContext(nil)
	_nilContext.Set(key, nil)

	_typeContext, _ := gin.CreateTestContext(nil)
	_typeContext.Set(key, 1)

	// setup tests
	tests := []struct {
		name    string
		context *gin.Context
		want    *library.Pipeline
	}{
		{
			name:    "context",
			context: _context,
			want:    _pipeline,
		},
		{
			name:    "context with no value",
			context: _emptyContext,
			want:    nil,
		},
		{
			name:    "context with nil value",
			context: _nilContext,
			want:    nil,
		},
		{
			name:    "context with wrong value type",
			context: _typeContext,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := FromContext(test.context)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("FromContext for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}

func TestPipeline_ToContext(t *testing.T) {
	// setup types
	_pipeline := new(library.Pipeline)

	gin.SetMode(gin.TestMode)
	_context, _ := gin.CreateTestContext(nil)

	// setup tests
	tests := []struct {
		name    string
		context *gin.Context
		want    *library.Pipeline
	}{
		{
			name:    "context",
			context: _context,
			want:    _pipeline,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ToContext(test.context, test.want)

			got := test.context.Value(key)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ToContext for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
