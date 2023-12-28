// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// The following functions were taken from
// https://github.com/go-vela/sdk-go/blob/main/vela/go
// which is the only reason go-vela/sdk-go is
// a dependency for go-vela/server
// TODO: consider moving to go-vela/types?

// Bool is a helper routine that allocates a new boolean
// value to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Bytes is a helper routine that allocates a new byte
// array value to store v and returns a pointer to it.
func Bytes(v []byte) *[]byte { return &v }

// Int is a helper routine that allocates a new integer
// value to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new 64 bit
// integer value to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// String is a helper routine that allocates a new string
// value to store v and returns a pointer to it.
func String(v string) *string { return &v }

// Strings is a helper routine that allocates a new string
// array value to store v and returns a pointer to it.
func Strings(v []string) *[]string { return &v }

// setup global variables used for testing.
var (
	_signingPrivateKey = "tCIevHOBq6DdN5SSBtteXUusjjd0fOqzk2eyi0DMq04NewmShNKQeUbbp3vkvIckb4pCxc+vxUo+mYf/vzOaSg=="
	_signingPublicKey  = "DXsJkoTSkHlG26d75LyHJG+KQsXPr8VKPpmH/78zmko="
	_build             = &library.Build{
		ID:           Int64(1),
		Number:       Int(1),
		Parent:       Int(1),
		Event:        String("push"),
		Status:       String("success"),
		Error:        String(""),
		Enqueued:     Int64(1563474077),
		Created:      Int64(1563474076),
		Started:      Int64(1563474077),
		Finished:     Int64(0),
		Deploy:       String(""),
		Clone:        String("https://github.com/github/octocat.git"),
		Source:       String("https://github.com/github/octocat/abcdefghi123456789"),
		Title:        String("push received from https://github.com/github/octocat"),
		Message:      String("First commit..."),
		Commit:       String("48afb5bdc41ad69bf22588491333f7cf71135163"),
		Sender:       String("OctoKitty"),
		Author:       String("OctoKitty"),
		Branch:       String("main"),
		Ref:          String("refs/heads/main"),
		BaseRef:      String(""),
		Host:         String("example.company.com"),
		Runtime:      String("docker"),
		Distribution: String("linux"),
	}

	_repo = &library.Repo{
		ID:          Int64(1),
		Org:         String("github"),
		Name:        String("octocat"),
		FullName:    String("github/octocat"),
		Link:        String("https://github.com/github/octocat"),
		Clone:       String("https://github.com/github/octocat.git"),
		Branch:      String("main"),
		Timeout:     Int64(60),
		Visibility:  String("public"),
		Private:     Bool(false),
		Trusted:     Bool(false),
		Active:      Bool(true),
		AllowPull:   Bool(false),
		AllowPush:   Bool(true),
		AllowDeploy: Bool(false),
		AllowTag:    Bool(false),
	}

	_steps = &pipeline.Build{
		Version: "1",
		ID:      "github_octocat_1",
		Services: pipeline.ContainerSlice{
			{
				ID:          "service_github_octocat_1_postgres",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		Steps: pipeline.ContainerSlice{
			{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			{
				ID:          "step_github_octocat_1_clone",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.5.1",
				Name:        "clone",
				Number:      2,
				Pull:        "always",
			},
			{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo hello"},
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      3,
				Pull:        "always",
			},
		},
	}

	_user = &library.User{
		ID:     Int64(1),
		Name:   String("octocat"),
		Token:  nil,
		Hash:   nil,
		Active: Bool(true),
		Admin:  Bool(false),
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
			WithChannels("foo"),
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
