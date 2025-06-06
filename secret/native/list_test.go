// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
)

func TestNative_List(t *testing.T) {
	// setup types
	sOne := new(api.Secret)
	sOne.SetID(1)
	sOne.SetOrg("foo")
	sOne.SetRepo("bar")
	sOne.SetTeam("")
	sOne.SetName("baz")
	sOne.SetValue("foob")
	sOne.SetType("repo")
	sOne.SetImages([]string{"foo", "bar"})
	sOne.SetAllowEvents(api.NewEventsFromMask(1))
	sOne.SetAllowCommand(false)
	sOne.SetAllowSubstitution(false)
	sOne.SetRepoAllowlist([]string{})
	sOne.SetCreatedAt(1)
	sOne.SetCreatedBy("user")
	sOne.SetUpdatedAt(1)
	sOne.SetUpdatedBy("user2")

	sTwo := new(api.Secret)
	sTwo.SetID(2)
	sTwo.SetOrg("foo")
	sTwo.SetRepo("bar")
	sTwo.SetTeam("")
	sTwo.SetName("foob")
	sTwo.SetValue("baz")
	sTwo.SetType("repo")
	sTwo.SetImages([]string{"foo", "bar"})
	sTwo.SetAllowEvents(api.NewEventsFromMask(1))
	sTwo.SetAllowCommand(false)
	sTwo.SetAllowSubstitution(false)
	sTwo.SetRepoAllowlist([]string{})
	sTwo.SetCreatedAt(1)
	sTwo.SetCreatedBy("user")
	sTwo.SetUpdatedAt(1)
	sTwo.SetUpdatedBy("user2")

	want := []*api.Secret{sTwo, sOne}

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteSecret(context.TODO(), sOne)
		_ = db.DeleteSecret(context.TODO(), sTwo)
		db.Close()
	}()

	// run test
	s, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	_, _ = s.Create(context.TODO(), "repo", "foo", "bar", sOne)

	_, _ = s.Create(context.TODO(), "repo", "foo", "bar", sTwo)

	got, err := s.List(context.TODO(), "repo", "foo", "bar", 1, 10, []string{})
	if err != nil {
		t.Errorf("List returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("List is %v, want %v", got, want)
	}
}

func TestNative_List_Empty(t *testing.T) {
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

	got, err := s.List(context.TODO(), "repo", "foo", "bar", 1, 10, []string{})
	if err != nil {
		t.Errorf("List returned err: %v", err)
	}

	if len(got) > 0 {
		t.Errorf("List is %v, want []", got)
	}
}
