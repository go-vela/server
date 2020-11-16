// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"testing"

	"github.com/go-vela/server/database"
)

func TestNative_New(t *testing.T) {
	// setup types
	d, _ := database.NewTest()
	defer d.Database.Close()
	passphrase := "C639A572E14D5075C526FDDD43E4ECF6"

	// run test
	s, err := New(d, passphrase)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if s == nil {
		t.Error("New returned nil client")
	}
}

func TestNative_New_Error(t *testing.T) {
	// run test
	s, err := New(nil, "")
	if err == nil {
		t.Errorf("New should have returned err")
	}

	if s != nil {
		t.Error("New should have returned nil client")
	}
}
