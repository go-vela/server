// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
)

func TestNative_Get(t *testing.T) {
	// setup types
	want := new(api.Secret)
	want.SetID(1)
	want.SetOrg("foo")
	want.SetRepo("bar")
	want.SetTeam("")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("repo")
	want.SetImages([]string{"foo", "bar"})
	want.SetAllowEvents(api.NewEventsFromMask(1))
	want.SetAllowCommand(false)
	want.SetAllowSubstitution(false)
	want.SetRepoAllowlist([]string{})
	want.SetCreatedAt(1)
	want.SetCreatedBy("user")
	want.SetUpdatedAt(1)
	want.SetUpdatedBy("user2")

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}
	defer db.Close()

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

	_, _ = s.Create(context.TODO(), "repo", "foo", "bar", want)

	got, err := s.Get(context.TODO(), "repo", "foo", "bar", "baz")
	if err != nil {
		t.Errorf("Get returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get is %v, want %v", got, want)
	}
}

func TestNative_Get_Invalid(t *testing.T) {
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

	got, err := s.Get(context.TODO(), "repo", "foo", "bar", "foob")
	if err == nil {
		t.Errorf("Get should have returned err")
	}

	if got != nil {
		t.Errorf("Get is %v, want nil", got)
	}
}
