// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database"

	"github.com/go-vela/types/library"
)

func TestNative_Update(t *testing.T) {
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

	// setup database
	d, _ := database.NewTest()

	defer func() {
		d.Database.Exec("delete from secrets;")
		d.Database.Close()
	}()

	_ = d.CreateSecret(want)

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Update("repo", "foo", "bar", want)
	if err != nil {
		t.Errorf("Update returned err: %v", err)
	}

	got, _ := s.Get("repo", "foo", "bar", "baz")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Update is %v, want %v", got, want)
	}
}

func TestNative_Update_Invalid(t *testing.T) {
	// setup types
	sec := new(library.Secret)
	sec.SetName("baz")
	sec.SetValue("foob")

	// setup database
	d, _ := database.NewTest()
	d.Database.Close()

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Update("repo", "foo", "bar", sec)
	if err == nil {
		t.Errorf("Update should have returned err")
	}
}
