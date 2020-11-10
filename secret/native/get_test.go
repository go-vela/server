// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database"

	"github.com/go-vela/types/library"
)

func TestNative_Get(t *testing.T) {
	// setup types
	want := new(library.Secret)
	want.SetID(1)
	want.SetOrg("foo")
	want.SetRepo("bar")
	want.SetTeam("")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("repo")
	want.SetImages([]string{"foo", "bar"})
	want.SetEvents([]string{"foo", "bar"})
	want.SetAllowCommand(false)

	passphrase := "go-vela"

	// setup database
	d, _ := database.NewTest()

	defer func() {
		d.Database.Exec("delete from secrets;")
		d.Database.Close()
	}()

	// run test
	s, err := New(d, passphrase)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	_ = s.Create("repo", "foo", "bar", want)

	value, _ := decrypt([]byte(want.GetValue()), passphrase)
	want.Value = &value

	got, err := s.Get("repo", "foo", "bar", "baz")
	if err != nil {
		t.Errorf("Get returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get is %v, want %v", got, want)
	}
}

func TestNative_Get_Invalid(t *testing.T) {
	// setup database
	d, _ := database.NewTest()
	d.Database.Close()

	passphrase := "go-vela"

	// run test
	s, err := New(d, passphrase)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get("repo", "foo", "bar", "foob")
	if err == nil {
		t.Errorf("Get should have returned err")
	}

	if got != nil {
		t.Errorf("Get is %v, want nil", got)
	}
}
