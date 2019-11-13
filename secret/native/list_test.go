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

func TestNative_List(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	team := ""
	name := "baz"
	value := "foob"
	typee := "repo"
	arr := []string{"foo", "bar"}
	booL := false
	sOne := &library.Secret{
		ID:           &one,
		Org:          &org,
		Repo:         &repo,
		Team:         &team,
		Name:         &name,
		Value:        &value,
		Type:         &typee,
		Images:       &arr,
		Events:       &arr,
		AllowCommand: &booL,
	}
	two := int64(2)
	sTwo := &library.Secret{
		ID:           &two,
		Org:          &org,
		Repo:         &repo,
		Team:         &team,
		Name:         &value,
		Value:        &name,
		Type:         &typee,
		Images:       &arr,
		Events:       &arr,
		AllowCommand: &booL,
	}
	want := []*library.Secret{sTwo, sOne}

	// setup database
	d, _ := database.NewTest()
	defer func() {
		d.Database.Exec("delete from secrets;")
		d.Database.Close()
	}()
	_ = d.CreateSecret(sOne)
	_ = d.CreateSecret(sTwo)

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List(typee, org, repo, 1, 10)
	if err != nil {
		t.Errorf("List returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("List is %v, want %v", got, want)
	}
}

func TestNative_List_Invalid(t *testing.T) {
	// setup database
	d, _ := database.NewTest()
	d.Database.Close()

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List("repo", "foo", "bar", 1, 10)
	if err == nil {
		t.Errorf("List should have returned err")
	}

	if got != nil {
		t.Errorf("List is %v, want nil", got)
	}
}
