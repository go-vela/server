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

func TestNative_Update(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	team := ""
	name := "baz"
	value := "foob"
	typee := "repo"
	arr := []string{"foo", "bar"}
	want := &library.Secret{
		ID:     &one,
		Org:    &org,
		Repo:   &repo,
		Team:   &team,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
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

	err = s.Update(typee, org, repo, want)
	if err != nil {
		t.Errorf("Update returned err: %v", err)
	}

	got, _ := s.Get(typee, org, repo, name)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Update is %v, want %v", got, want)
	}
}

func TestNative_Update_Invalid(t *testing.T) {
	// setup types
	name := "baz"
	value := "foob"
	sec := &library.Secret{
		Name:  &name,
		Value: &value,
	}

	// setup database
	d, _ := database.NewTest()
	d.Database.Close()

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Update("repo", "foo", "bar", sec)
	if err == nil {
		t.Errorf("Update should have returned err")
	}
}
