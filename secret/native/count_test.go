// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"testing"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/library"
)

func TestNative_Count(t *testing.T) {
	// setup types
	sec := new(library.Secret)
	sec.SetID(1)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})
	sec.SetCreatedAt(1)
	sec.SetUpdatedAt(1)

	want := 1

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		db.DeleteSecret(context.TODO(), sec)
		db.Close()
	}()

	_, _ = db.CreateSecret(context.TODO(), sec)

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Count(context.TODO(), "repo", "foo", "bar", []string{})
	if err != nil {
		t.Errorf("Count returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("Count is %v, want %v", got, want)
	}
}

func TestNative_Count_Empty(t *testing.T) {
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

	got, err := s.Count(context.TODO(), "repo", "foo", "bar", []string{})
	if err != nil {
		t.Errorf("Count returned err: %v", err)
	}

	if got != 0 {
		t.Errorf("Count is %v, want 0", got)
	}
}
