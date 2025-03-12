// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	api "github.com/go-vela/server/api/types"
)

// Ptr is a helper routine that allocates a new T value
// to store v and returns a pointer to it.
func Ptr[T any](v T) *T {
	return &v
}

// setup global variables used for testing.
var (
	_signingPrivateKey = "tCIevHOBq6DdN5SSBtteXUusjjd0fOqzk2eyi0DMq04NewmShNKQeUbbp3vkvIckb4pCxc+vxUo+mYf/vzOaSg=="
	_signingPublicKey  = "DXsJkoTSkHlG26d75LyHJG+KQsXPr8VKPpmH/78zmko="
	_build             = &api.Build{
		ID: Ptr(int64(1)),
		Repo: &api.Repo{
			ID: Ptr(int64(1)),
			Owner: &api.User{
				ID:     Ptr(int64(1)),
				Name:   Ptr("octocat"),
				Token:  nil,
				Active: Ptr(true),
				Admin:  Ptr(false),
			},
			Org:        Ptr("github"),
			Name:       Ptr("octocat"),
			FullName:   Ptr("github/octocat"),
			Link:       Ptr("https://github.com/github/octocat"),
			Clone:      Ptr("https://github.com/github/octocat.git"),
			Branch:     Ptr("main"),
			Timeout:    Ptr(int32(60)),
			Visibility: Ptr("public"),
			Private:    Ptr(false),
			Trusted:    Ptr(false),
			Active:     Ptr(true),
		},
		Number:       Ptr(int64(2)),
		Parent:       Ptr(int64(1)),
		Event:        Ptr("push"),
		Status:       Ptr("success"),
		Error:        Ptr(""),
		Enqueued:     Ptr(int64(1563474077)),
		Created:      Ptr(int64(1563474076)),
		Started:      Ptr(int64(1563474077)),
		Finished:     Ptr(int64(0)),
		Deploy:       Ptr(""),
		Clone:        Ptr("https://github.com/github/octocat.git"),
		Source:       Ptr("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        Ptr("push received from https://github.com/github/octocat"),
		Message:      Ptr("First commit..."),
		Commit:       Ptr("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       Ptr("OctoKitty"),
		Author:       Ptr("OctoKitty"),
		Branch:       Ptr("main"),
		Ref:          Ptr("refs/heads/main"),
		BaseRef:      Ptr(""),
		Host:         Ptr("example.company.com"),
		Runtime:      Ptr("docker"),
		Distribution: Ptr("linux"),
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
