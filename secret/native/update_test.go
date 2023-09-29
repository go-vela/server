// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/library"
)

func TestNative_Update(t *testing.T) {
	// setup types
	original := new(library.Secret)
	original.SetID(1)
	original.SetOrg("foo")
	original.SetRepo("bar")
	original.SetTeam("")
	original.SetName("baz")
	original.SetValue("secretValue")
	original.SetType("repo")
	original.SetImages([]string{"foo", "baz"})
	original.SetEvents([]string{"foob", "bar"})
	original.SetAllowCommand(true)
	original.SetCreatedAt(1)
	original.SetCreatedBy("user")
	original.SetUpdatedAt(time.Now().UTC().Unix())
	original.SetUpdatedBy("user")

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
	want.SetUpdatedAt(time.Now().UTC().Unix())
	want.SetUpdatedBy("user2")

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		db.DeleteSecret(context.TODO(), original)
		db.Close()
	}()

	_, _ = db.CreateSecret(context.TODO(), original)

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Update(context.TODO(), "repo", "foo", "bar", want)
	if err != nil {
		t.Errorf("Update returned err: %v", err)
	}

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
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}
	defer db.Close()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	_, err = s.Update(context.TODO(), "repo", "foo", "bar", sec)
	if err == nil {
		t.Errorf("Update should have returned err")
	}
}
