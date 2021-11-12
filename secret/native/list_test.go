// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database/sqlite"
	"github.com/go-vela/types/library"
)

func TestNative_List(t *testing.T) {
	// setup types
	sOne := new(library.Secret)
	sOne.SetID(1)
	sOne.SetOrg("foo")
	sOne.SetRepo("bar")
	sOne.SetTeam("")
	sOne.SetName("baz")
	sOne.SetValue("foob")
	sOne.SetType("repo")
	sOne.SetImages([]string{"foo", "bar"})
	sOne.SetEvents([]string{"foo", "bar"})
	sOne.SetAllowCommand(false)
	sOne.SetCreatedAt(1)
	sOne.SetCreatedBy("user")
	sOne.SetUpdatedAt(1)
	sOne.SetUpdatedBy("user2")

	sTwo := new(library.Secret)
	sTwo.SetID(2)
	sTwo.SetOrg("foo")
	sTwo.SetRepo("bar")
	sTwo.SetTeam("")
	sTwo.SetName("foob")
	sTwo.SetValue("baz")
	sTwo.SetType("repo")
	sTwo.SetImages([]string{"foo", "bar"})
	sTwo.SetEvents([]string{"foo", "bar"})
	sTwo.SetAllowCommand(false)
	sTwo.SetCreatedAt(1)
	sTwo.SetCreatedBy("user")
	sTwo.SetUpdatedAt(1)
	sTwo.SetUpdatedBy("user2")

	want := []*library.Secret{sTwo, sOne}

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

	_ = s.Create("repo", "foo", "bar", sOne)

	_ = s.Create("repo", "foo", "bar", sTwo)

	got, err := s.List("repo", "foo", "bar", 1, 10, []string{})
	if err != nil {
		t.Errorf("List returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("List is %v, want %v", got, want)
	}
}

func TestNative_List_Invalid(t *testing.T) {
	// setup database
	db, _ := sqlite.NewTest()
	_sql, _ := db.Sqlite.DB()
	_sql.Close()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List("repo", "foo", "bar", 1, 10, []string{})
	if err == nil {
		t.Errorf("List should have returned err")
	}

	if got != nil {
		t.Errorf("List is %v, want nil", got)
	}
}
