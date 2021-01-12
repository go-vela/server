// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executors

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

func TestExecutors_Retrieve(t *testing.T) {
	// setup types
	eID := int64(1)
	e := library.Executor{ID: &eID}
	want := []library.Executor{e}

	// setup context
	gin.SetMode(gin.TestMode)

	context, _ := gin.CreateTestContext(nil)
	ToContext(context, want)

	// run test
	got := Retrieve(context)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}
