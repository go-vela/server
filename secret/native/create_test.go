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

func TestNative_Create_Org(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "*"
	team := ""
	name := "bar"
	value := "baz"
	typee := "org"
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

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create(typee, org, repo, want)
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}

	got, _ := s.Get(typee, org, repo, name)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Create is %v, want %v", got, want)
	}
}

func TestNative_Create_Repo(t *testing.T) {
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

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create(typee, org, repo, want)
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}

	got, _ := s.Get(typee, org, repo, name)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Create is %v, want %v", got, want)
	}
}

func TestNative_Create_Shared(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := ""
	team := "bar"
	name := "baz"
	value := "foob"
	typee := "shared"
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

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create(typee, org, team, want)
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}

	got, _ := s.Get(typee, org, team, name)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Create is %v, want %v", got, want)
	}
}

func TestNative_Create_Invalid(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "invalid"
	arr := []string{"foo", "bar"}
	sec := &library.Secret{
		ID:     &one,
		Org:    &org,
		Repo:   &repo,
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

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create(typee, org, repo, sec)
	if err == nil {
		t.Errorf("Create should have returned err")
	}
}
