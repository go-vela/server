// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database/sqlite"
	"github.com/go-vela/types/library"
)

func TestNative_Create_Org(t *testing.T) {
	// setup types
	want := new(library.Secret)
	want.SetID(1)
	want.SetOrg("foo")
	want.SetRepo("*")
	want.SetTeam("")
	want.SetName("bar")
	want.SetValue("baz")
	want.SetType("org")
	want.SetImages([]string{"foo", "bar"})
	want.SetEvents([]string{"foo", "bar"})
	want.SetAllowCommand(false)
	want.SetCreatedAt(1)
	want.SetCreatedBy("user")
	want.SetUpdatedAt(1)
	want.SetUpdatedBy("user2")

	// setup database
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from secrets;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create("org", "foo", "*", want)
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}

	got, _ := s.Get("org", "foo", "*", "bar")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Create is %v, want %v", got, want)
	}
}

func TestNative_Create_Repo(t *testing.T) {
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
	want.SetCreatedAt(1)
	want.SetCreatedBy("user")
	want.SetUpdatedAt(1)
	want.SetUpdatedBy("user2")

	// setup database
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from secrets;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create("repo", "foo", "bar", want)
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}

	got, _ := s.Get("repo", "foo", "bar", "baz")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Create is %v, want %v", got, want)
	}
}

func TestNative_Create_Shared(t *testing.T) {
	// setup types
	want := new(library.Secret)
	want.SetID(1)
	want.SetOrg("foo")
	want.SetRepo("")
	want.SetTeam("bar")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("shared")
	want.SetImages([]string{"foo", "bar"})
	want.SetEvents([]string{"foo", "bar"})
	want.SetAllowCommand(false)
	want.SetCreatedAt(1)
	want.SetCreatedBy("user")
	want.SetUpdatedAt(1)
	want.SetUpdatedBy("user2")

	// setup database
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from secrets;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create("shared", "foo", "bar", want)
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}

	got, _ := s.Get("shared", "foo", "bar", "baz")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Create is %v, want %v", got, want)
	}
}

func TestNative_Create_Invalid(t *testing.T) {
	// setup types
	sec := new(library.Secret)
	sec.SetID(1)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetTeam("")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("invalid")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})
	sec.SetAllowCommand(false)
	sec.SetCreatedAt(1)
	sec.SetCreatedBy("user")
	sec.SetUpdatedAt(1)
	sec.SetUpdatedBy("user2")

	// setup database
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from secrets;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create("invalid", "foo", "bar", sec)
	if err == nil {
		t.Errorf("Create should have returned err")
	}
}
