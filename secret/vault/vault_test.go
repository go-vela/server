// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

	type args struct {
		version string
		prefix  string
	}
	tests := []struct {
		name string
		args args
	}{
		{"v1", args{version: "1", prefix: ""}},
		{"v2", args{version: "2", prefix: ""}},
		{"v2 with prefix", args{version: "2", prefix: "prefix"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(
				WithAddress(fake.URL),
				WithAuthMethod(""),
				WithAWSRole(""),
				WithPrefix(tt.args.prefix),
				WithToken("foo"),
				WithTokenDuration(0),
				WithVersion(tt.args.version),
			)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			if s == nil {
				t.Error("New returned nil client")
			}
		})
	}
}

func TestVault_New_Error(t *testing.T) {
	type args struct {
		version string
		prefix  string
	}
	tests := []struct {
		name string
		args args
	}{
		{"v1", args{version: "1", prefix: ""}},
		{"v2", args{version: "2", prefix: ""}},
		{"v2 with prefix", args{version: "2", prefix: "prefix"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(
				WithAddress("!@#$%^&*()"),
				WithAuthMethod(""),
				WithAWSRole(""),
				WithPrefix(tt.args.prefix),
				WithToken("foo"),
				WithTokenDuration(0),
				WithVersion(tt.args.version),
			)
			if err == nil {
				t.Errorf("New should have returned err")
			}

			if s != nil {
				t.Error("New should have returned nil client")
			}
		})
	}
}

func TestVault_secretFromVault(t *testing.T) {
	// setup types
	inputV1 := &api.Secret{
		Data: map[string]interface{}{
			"events":        []interface{}{"foo", "bar"},
			"images":        []interface{}{"foo", "bar"},
			"name":          "bar",
			"org":           "foo",
			"repo":          "*",
			"team":          "foob",
			"type":          "org",
			"value":         "baz",
			"allow_command": true,
			"created_at":    int64(1563474077),
			"created_by":    "octocat",
			"updated_at":    int64(1563474079),
			"updated_by":    "octocat2",
		},
	}

	inputV2 := &api.Secret{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				"events":        []interface{}{"foo", "bar"},
				"images":        []interface{}{"foo", "bar"},
				"name":          "bar",
				"org":           "foo",
				"repo":          "*",
				"team":          "foob",
				"type":          "org",
				"value":         "baz",
				"allow_command": true,
				"created_at":    int64(1563474077),
				"created_by":    "octocat",
				"updated_at":    int64(1563474079),
				"updated_by":    "octocat2",
			},
		},
	}

	want := new(library.Secret)
	want.SetOrg("foo")
	want.SetRepo("*")
	want.SetTeam("foob")
	want.SetName("bar")
	want.SetValue("baz")
	want.SetType("org")
	want.SetEvents([]string{"foo", "bar"})
	want.SetImages([]string{"foo", "bar"})
	want.SetAllowCommand(true)
	want.SetCreatedAt(1563474077)
	want.SetCreatedBy("octocat")
	want.SetUpdatedAt(1563474079)
	want.SetUpdatedBy("octocat2")

	type args struct {
		secret *api.Secret
	}
	tests := []struct {
		name string
		args args
	}{
		{"v1", args{secret: inputV1}},
		{"v2", args{secret: inputV2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := secretFromVault(tt.args.secret)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("secretFromVault is %v, want %v", got, want)
			}
		})
	}
}

func TestVault_vaultFromSecret(t *testing.T) {
	// setup types
	s := new(library.Secret)
	s.SetOrg("foo")
	s.SetRepo("*")
	s.SetTeam("foob")
	s.SetName("bar")
	s.SetValue("baz")
	s.SetType("org")
	s.SetEvents([]string{"foo", "bar"})
	s.SetImages([]string{"foo", "bar"})
	s.SetAllowCommand(true)
	s.SetCreatedAt(1563474077)
	s.SetCreatedBy("octocat")
	s.SetUpdatedAt(1563474079)
	s.SetUpdatedBy("octocat2")

	want := &api.Secret{
		Data: map[string]interface{}{
			"events":        []string{"foo", "bar"},
			"images":        []string{"foo", "bar"},
			"name":          "bar",
			"org":           "foo",
			"repo":          "*",
			"team":          "foob",
			"type":          "org",
			"value":         "baz",
			"allow_command": true,
			"created_at":    int64(1563474077),
			"created_by":    "octocat",
			"updated_at":    int64(1563474079),
			"updated_by":    "octocat2",
		},
	}

	// run test
	got := vaultFromSecret(s)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("vaultFromSecret is %v, want %v", got, want)
	}
}
