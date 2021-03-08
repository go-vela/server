// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"testing"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/secret/native"

	"github.com/gin-gonic/gin"
)

func TestSecret_FromContext(t *testing.T) {
	// setup types
	d, _ := database.NewTest()
	defer d.Database.Close()

	want, err := native.New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	context.Set("native", want)

	// run test
	got := FromContext(context, "native")

	if got != want {
		t.Errorf("FromContext is %v, want %v", got, want)
	}
}

func TestSecret_FromContext_Bad(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	context.Set("native", nil)

	// run test
	got := FromContext(context, "native")

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestSecret_FromContext_WrongType(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	context.Set("native", 1)

	// run test
	got := FromContext(context, "native")

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestSecret_FromContext_Empty(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)

	// run test
	got := FromContext(context, "native")

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestSecret_ToContext(t *testing.T) {
	// setup types
	d, _ := database.NewTest()
	defer d.Database.Close()

	want, err := native.New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	ToContext(context, "native", want)

	// run test
	got := context.Value("native")

	if got != want {
		t.Errorf("ToContext is %v, want %v", got, want)
	}
}
