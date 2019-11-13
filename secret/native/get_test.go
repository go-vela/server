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

func TestNative_Get(t *testing.T) {
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
	want := &library.Secret{
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

	// setup database
	d, _ := database.NewTest()
	defer func() {
		d.Database.Exec("delete from secrets;")
		d.Database.Close()
	}()
	_ = d.CreateSecret(want)

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get(typee, org, repo, name)
	if err != nil {
		t.Errorf("Get returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get is %v, want %v", got, want)
	}
}

func TestNative_Get_Invalid(t *testing.T) {
	// setup database
	d, _ := database.NewTest()
	d.Database.Close()

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get("repo", "foo", "bar", "foob")
	if err == nil {
		t.Errorf("Get should have returned err")
	}

	if got != nil {
		t.Errorf("Get is %v, want nil", got)
	}
}
