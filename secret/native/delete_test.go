// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"testing"

	"github.com/go-vela/server/database"
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
	sec.SetCreatedAt(1)
	sec.SetUpdatedAt(1)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteSecret(context.TODO(), sec)
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

	err = s.Delete(context.TODO(), "repo", "foo", "bar", "baz")
	if err != nil {
		t.Errorf("Delete returned err: %v", err)
	}
}

func TestNative_Delete_Invalid(t *testing.T) {
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

	err = s.Delete(context.TODO(), "repo", "foo", "bar", "foob")
	if err == nil {
		t.Errorf("Delete should have returned err")
	}
}
