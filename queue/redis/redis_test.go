// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	api "github.com/go-vela/server/api/types"
)

// Ptr is a helper routine that allocates a new T value
// to store v and returns a pointer to it.
//
//go:fix inline
func Ptr[T any](v T) *T {
	return new(v)
}

// setup global variables used for testing.
var (
	_signingPrivateKey = "tCIevHOBq6DdN5SSBtteXUusjjd0fOqzk2eyi0DMq04NewmShNKQeUbbp3vkvIckb4pCxc+vxUo+mYf/vzOaSg=="
	_signingPublicKey  = "DXsJkoTSkHlG26d75LyHJG+KQsXPr8VKPpmH/78zmko="
	_build             = &api.Build{
		ID: new(int64(1)),
		Repo: &api.Repo{
			ID: new(int64(1)),
			Owner: &api.User{
				ID:     new(int64(1)),
				Name:   new("octocat"),
				Token:  nil,
				Active: new(true),
				Admin:  new(false),
			},
			Org:        new("github"),
			Name:       new("octocat"),
			FullName:   new("github/octocat"),
			Link:       new("https://github.com/github/octocat"),
			Clone:      new("https://github.com/github/octocat.git"),
			Branch:     new("main"),
			Timeout:    new(int32(60)),
			Visibility: new("public"),
			Private:    new(false),
			Trusted:    new(false),
			Active:     new(true),
		},
		Number:       new(int64(2)),
		Parent:       new(int64(1)),
		Event:        new("push"),
		Status:       new("success"),
		Error:        new(""),
		Enqueued:     new(int64(1563474077)),
		Created:      new(int64(1563474076)),
		Started:      new(int64(1563474077)),
		Finished:     new(int64(0)),
		Deploy:       new(""),
		Clone:        new("https://github.com/github/octocat.git"),
		Source:       new("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        new("push received from https://github.com/github/octocat"),
		Message:      new("First commit..."),
		Commit:       new("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       new("OctoKitty"),
		Author:       new("OctoKitty"),
		Branch:       new("main"),
		Ref:          new("refs/heads/main"),
		BaseRef:      new(""),
		Host:         new("example.company.com"),
		Runtime:      new("docker"),
		Distribution: new("linux"),
	}
)

func TestRedis_New(t *testing.T) {
	// setup types
	// create a local fake redis instance
	//
	// https://pkg.go.dev/github.com/alicebob/miniredis/v2#Run
	_redis, err := miniredis.Run()
	if err != nil {
		t.Errorf("unable to create miniredis instance: %v", err)
	}
	defer _redis.Close()

	// setup tests
	tests := []struct {
		failure bool
		address string
	}{
		{
			failure: false,
			address: fmt.Sprintf("redis://%s", _redis.Addr()),
		},
		{
			failure: true,
			address: "",
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(
			context.Background(),
			WithAddress(test.address),
			WithRoutes("foo"),
			WithCluster(false),
			WithTimeout(5*time.Second),
		)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}
	}
}
