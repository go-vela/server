// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"

	"github.com/hashicorp/vault/api"
)

func TestVault_New(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if s == nil {
		t.Error("New returned nil client")
	}
}

func TestVault_New_Error(t *testing.T) {
	// run test
	s, err := New("!@#$%^&*()", "")
	if err == nil {
		t.Errorf("New should have returned err")
	}

	if s != nil {
		t.Error("New should have returned nil client")
	}
}

func TestVault_secretFromVault(t *testing.T) {
	// setup types
	v := &api.Secret{
		Data: map[string]interface{}{
			"events": []interface{}{"foo", "bar"},
			"images": []interface{}{"foo", "bar"},
			"name":   "bar",
			"org":    "foo",
			"repo":   "*",
			"team":   "foob",
			"type":   "org",
			"value":  "baz",
		},
	}

	org := "foo"
	repo := "*"
	team := "foob"
	name := "bar"
	value := "baz"
	typee := "org"
	arr := []string{"foo", "bar"}
	want := &library.Secret{
		Org:    &org,
		Repo:   &repo,
		Team:   &team,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	// run test
	got := secretFromVault(v)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("secretFromVault is %v, want %v", got, want)
	}
}

func TestVault_vaultFromSecret(t *testing.T) {
	// setup types
	org := "foo"
	repo := "*"
	team := "foob"
	name := "bar"
	value := "baz"
	typee := "org"
	arr := []string{"foo", "bar"}
	s := &library.Secret{
		Org:    &org,
		Repo:   &repo,
		Team:   &team,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	want := &api.Secret{
		Data: map[string]interface{}{
			"events": []string{"foo", "bar"},
			"images": []string{"foo", "bar"},
			"name":   "bar",
			"org":    "foo",
			"repo":   "*",
			"team":   "foob",
			"type":   "org",
			"value":  "baz",
		},
	}

	// run test
	got := vaultFromSecret(s)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("vaultFromSecret is %v, want %v", got, want)
	}
}
