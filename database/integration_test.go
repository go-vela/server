// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"fmt"
	"testing"
)

func TestDatabase_Integration(t *testing.T) {
	fmt.Println("Hello")

	t.Error("demonstrating failure")
}
