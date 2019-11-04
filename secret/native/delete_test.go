// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"testing"

	"github.com/go-vela/server/database"

	"github.com/go-vela/types/library"
)

func TestNative_Delete(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "org"
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
	_ = d.CreateSecret(sec)

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Delete(typee, org, repo, name)
	if err != nil {
		t.Errorf("Delete returned err: %v", err)
	}
}

func TestNative_Delete_Invalid(t *testing.T) {
	// setup database
	d, _ := database.NewTest()
	d.Database.Close()

	// run test
	s, err := New(d)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Delete("repo", "foo", "bar", "foob")
	if err == nil {
		t.Errorf("Delete should have returned err")
	}
}
