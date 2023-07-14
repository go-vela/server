// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/secret/native"
)

func TestSecret_FromContext(t *testing.T) {
	// setup types
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}
	defer db.Close()

	want, err := native.New(
		native.WithDatabase(db),
	)
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
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}
	defer db.Close()

	want, err := native.New(
		native.WithDatabase(db),
	)
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
