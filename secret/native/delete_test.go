// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"testing"

	"github.com/go-vela/server/database/sqlite"
	"github.com/go-vela/types/library"
)

func TestNative_Delete(t *testing.T) {
	// setup types
	sec := new(library.Secret)
	sec.SetID(1)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetTeam("")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})
	sec.SetAllowCommand(false)

	// setup database
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from secrets;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	_ = db.CreateSecret(sec)

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Delete("repo", "foo", "bar", "baz")
	if err != nil {
		t.Errorf("Delete returned err: %v", err)
	}
}

func TestNative_Delete_Invalid(t *testing.T) {
	// setup database
	db, _ := sqlite.NewTest()

	defer func() { _sql, _ := db.Sqlite.DB(); _sql.Close() }()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Delete("repo", "foo", "bar", "foob")
	if err == nil {
		t.Errorf("Delete should have returned err")
	}
}
