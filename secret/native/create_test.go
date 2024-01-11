// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-vela/server/database"
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
	want.SetAllowEvents(library.NewEventsFromMask(1))
	want.SetAllowCommand(false)
	want.SetCreatedAt(1)
	want.SetCreatedBy("user")
	want.SetUpdatedAt(1)
	want.SetUpdatedBy("user2")

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteSecret(context.TODO(), want)
		db.Close()
	}()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Create(context.TODO(), "org", "foo", "*", want)
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}

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
	want.SetAllowEvents(library.NewEventsFromMask(1))
	want.SetAllowCommand(false)
	want.SetCreatedAt(1)
	want.SetCreatedBy("user")
	want.SetUpdatedAt(1)
	want.SetUpdatedBy("user2")

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteSecret(context.TODO(), want)
		db.Close()
	}()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Create(context.TODO(), "repo", "foo", "bar", want)
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}

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
	want.SetAllowEvents(library.NewEventsFromMask(1))
	want.SetAllowCommand(false)
	want.SetCreatedAt(1)
	want.SetCreatedBy("user")
	want.SetUpdatedAt(1)
	want.SetUpdatedBy("user2")

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteSecret(context.TODO(), want)
		db.Close()
	}()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Create(context.TODO(), "shared", "foo", "bar", want)
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}

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
	sec.SetAllowEvents(library.NewEventsFromMask(1))
	sec.SetAllowCommand(false)
	sec.SetCreatedAt(1)
	sec.SetCreatedBy("user")
	sec.SetUpdatedAt(1)
	sec.SetUpdatedBy("user2")

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteSecret(context.TODO(), sec)
		db.Close()
	}()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	_, err = s.Create(context.TODO(), "invalid", "foo", "bar", sec)
	if err == nil {
		t.Errorf("Create should have returned err")
	}
}
