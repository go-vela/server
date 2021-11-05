// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"reflect"
	"testing"

	"github.com/go-vela/types/raw"
	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/yaml"

	"github.com/urfave/cli/v2"
)

func TestNative_EnvironmentStages(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"

	e := raw.StringSliceMap{
		"HELLO": "Hello, Global Message",
	}

	s := yaml.StageSlice{
		&yaml.Stage{
			Name: str,
			Steps: yaml.StepSlice{
				&yaml.Step{
					Image: "alpine",
					Name:  str,
					Pull:  "always",
				},
			},
		},
	}

	env := environment(nil, nil, nil, nil)
	env["HELLO"] = "Hello, Global Message"

	want := yaml.StageSlice{
		&yaml.Stage{
			Name: str,
			Steps: yaml.StepSlice{
				&yaml.Step{
					Environment: env,
					Image:       "alpine",
					Name:        str,
					Pull:        "always",
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.EnvironmentStages(s, e)
	if err != nil {
		t.Errorf("EnvironmentStages returned err: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("EnvironmentStages mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_EnvironmentSteps(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	e := raw.StringSliceMap{
		"HELLO": "Hello, Global Message",
	}

	str := "foo"
	s := yaml.StepSlice{
		&yaml.Step{
			Image: "alpine",
			Name:  str,
			Pull:  "always",
			Environment: raw.StringSliceMap{
				"BUILD_CHANNEL": "foo",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Image: "alpine",
			Name:  str,
			Pull:  "always",
			Environment: raw.StringSliceMap{
				"BUILD_AUTHOR":             "",
				"BUILD_AUTHOR_EMAIL":       "",
				"BUILD_BASE_REF":           "",
				"BUILD_BRANCH":             "",
				"BUILD_CHANNEL":            "TODO",
				"BUILD_CLONE":              "",
				"BUILD_COMMIT":             "",
				"BUILD_CREATED":            "0",
				"BUILD_ENQUEUED":           "0",
				"BUILD_EVENT":              "",
				"BUILD_HOST":               "",
				"BUILD_LINK":               "",
				"BUILD_MESSAGE":            "",
				"BUILD_NUMBER":             "0",
				"BUILD_PARENT":             "0",
				"BUILD_REF":                "",
				"BUILD_SENDER":             "",
				"BUILD_SOURCE":             "",
				"BUILD_STARTED":            "0",
				"BUILD_STATUS":             "",
				"BUILD_TITLE":              "",
				"BUILD_WORKSPACE":          "/vela/src",
				"CI":                       "vela",
				"REPOSITORY_ACTIVE":        "false",
				"REPOSITORY_ALLOW_COMMENT": "false",
				"REPOSITORY_ALLOW_DEPLOY":  "false",
				"REPOSITORY_ALLOW_PULL":    "false",
				"REPOSITORY_ALLOW_PUSH":    "false",
				"REPOSITORY_ALLOW_TAG":     "false",
				"REPOSITORY_BRANCH":        "",
				"REPOSITORY_CLONE":         "",
				"REPOSITORY_FULL_NAME":     "",
				"REPOSITORY_LINK":          "",
				"REPOSITORY_NAME":          "",
				"REPOSITORY_ORG":           "",
				"REPOSITORY_PRIVATE":       "false",
				"REPOSITORY_TIMEOUT":       "0",
				"REPOSITORY_TRUSTED":       "false",
				"REPOSITORY_VISIBILITY":    "",
				"VELA":                     "true",
				"VELA_ADDR":                "TODO",
				"VELA_BUILD_AUTHOR":        "",
				"VELA_BUILD_AUTHOR_EMAIL":  "",
				"VELA_BUILD_BASE_REF":      "",
				"VELA_BUILD_BRANCH":        "",
				"VELA_BUILD_CHANNEL":       "TODO",
				"VELA_BUILD_CLONE":         "",
				"VELA_BUILD_COMMIT":        "",
				"VELA_BUILD_CREATED":       "0",
				"VELA_BUILD_DISTRIBUTION":  "",
				"VELA_BUILD_ENQUEUED":      "0",
				"VELA_BUILD_EVENT":         "",
				"VELA_BUILD_HOST":          "",
				"VELA_BUILD_LINK":          "",
				"VELA_BUILD_MESSAGE":       "",
				"VELA_BUILD_NUMBER":        "0",
				"VELA_BUILD_PARENT":        "0",
				"VELA_BUILD_REF":           "",
				"VELA_BUILD_RUNTIME":       "",
				"VELA_BUILD_SENDER":        "",
				"VELA_BUILD_SOURCE":        "",
				"VELA_BUILD_STARTED":       "0",
				"VELA_BUILD_STATUS":        "",
				"VELA_BUILD_TITLE":         "",
				"VELA_BUILD_WORKSPACE":     "/vela/src",
				"VELA_CHANNEL":             "TODO",
				"VELA_DATABASE":            "TODO",
				"VELA_DISTRIBUTION":        "TODO",
				"VELA_HOST":                "TODO",
				"VELA_NETRC_MACHINE":       "TODO",
				"VELA_NETRC_PASSWORD":      "",
				"VELA_NETRC_USERNAME":      "x-oauth-basic",
				"VELA_QUEUE":               "TODO",
				"VELA_REPO_ACTIVE":         "false",
				"VELA_REPO_ALLOW_COMMENT":  "false",
				"VELA_REPO_ALLOW_DEPLOY":   "false",
				"VELA_REPO_ALLOW_PULL":     "false",
				"VELA_REPO_ALLOW_PUSH":     "false",
				"VELA_REPO_ALLOW_TAG":      "false",
				"VELA_REPO_BRANCH":         "",
				"VELA_REPO_CLONE":          "",
				"VELA_REPO_FULL_NAME":      "",
				"VELA_REPO_LINK":           "",
				"VELA_REPO_NAME":           "",
				"VELA_REPO_ORG":            "",
				"VELA_REPO_PIPELINE_TYPE":  "",
				"VELA_REPO_PRIVATE":        "false",
				"VELA_REPO_TIMEOUT":        "0",
				"VELA_REPO_TRUSTED":        "false",
				"VELA_REPO_VISIBILITY":     "",
				"VELA_RUNTIME":             "TODO",
				"VELA_SOURCE":              "TODO",
				"VELA_USER_ACTIVE":         "false",
				"VELA_USER_ADMIN":          "false",
				"VELA_USER_FAVORITES":      "[]",
				"VELA_USER_NAME":           "",
				"VELA_VERSION":             "TODO",
				"VELA_WORKSPACE":           "/vela/src",
				"HELLO":                    "Hello, Global Message",
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.EnvironmentSteps(s, e)
	if err != nil {
		t.Errorf("EnvironmentSteps returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("EnvironmentSteps is %v, want %v", got, want)
	}
}

func TestNative_EnvironmentServices(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	e := raw.StringSliceMap{
		"HELLO": "Hello, Global Message",
	}

	str := "foo"
	s := yaml.ServiceSlice{
		&yaml.Service{
			Image: "postgres",
			Name:  str,
			Pull:  "always",
			Environment: raw.StringSliceMap{
				"BUILD_CHANNEL": "foo",
			},
		},
	}

	want := yaml.ServiceSlice{
		&yaml.Service{
			Image: "postgres",
			Name:  str,
			Pull:  "always",
			Environment: raw.StringSliceMap{
				"BUILD_AUTHOR":             "",
				"BUILD_AUTHOR_EMAIL":       "",
				"BUILD_BASE_REF":           "",
				"BUILD_BRANCH":             "",
				"BUILD_CHANNEL":            "TODO",
				"BUILD_CLONE":              "",
				"BUILD_COMMIT":             "",
				"BUILD_CREATED":            "0",
				"BUILD_ENQUEUED":           "0",
				"BUILD_EVENT":              "",
				"BUILD_HOST":               "",
				"BUILD_LINK":               "",
				"BUILD_MESSAGE":            "",
				"BUILD_NUMBER":             "0",
				"BUILD_PARENT":             "0",
				"BUILD_REF":                "",
				"BUILD_SENDER":             "",
				"BUILD_SOURCE":             "",
				"BUILD_STARTED":            "0",
				"BUILD_STATUS":             "",
				"BUILD_TITLE":              "",
				"BUILD_WORKSPACE":          "/vela/src",
				"CI":                       "vela",
				"REPOSITORY_ACTIVE":        "false",
				"REPOSITORY_ALLOW_COMMENT": "false",
				"REPOSITORY_ALLOW_DEPLOY":  "false",
				"REPOSITORY_ALLOW_PULL":    "false",
				"REPOSITORY_ALLOW_PUSH":    "false",
				"REPOSITORY_ALLOW_TAG":     "false",
				"REPOSITORY_BRANCH":        "",
				"REPOSITORY_CLONE":         "",
				"REPOSITORY_FULL_NAME":     "",
				"REPOSITORY_LINK":          "",
				"REPOSITORY_NAME":          "",
				"REPOSITORY_ORG":           "",
				"REPOSITORY_PRIVATE":       "false",
				"REPOSITORY_TIMEOUT":       "0",
				"REPOSITORY_TRUSTED":       "false",
				"REPOSITORY_VISIBILITY":    "",
				"VELA":                     "true",
				"VELA_ADDR":                "TODO",
				"VELA_BUILD_AUTHOR":        "",
				"VELA_BUILD_AUTHOR_EMAIL":  "",
				"VELA_BUILD_BASE_REF":      "",
				"VELA_BUILD_BRANCH":        "",
				"VELA_BUILD_CHANNEL":       "TODO",
				"VELA_BUILD_CLONE":         "",
				"VELA_BUILD_COMMIT":        "",
				"VELA_BUILD_CREATED":       "0",
				"VELA_BUILD_DISTRIBUTION":  "",
				"VELA_BUILD_ENQUEUED":      "0",
				"VELA_BUILD_EVENT":         "",
				"VELA_BUILD_HOST":          "",
				"VELA_BUILD_LINK":          "",
				"VELA_BUILD_MESSAGE":       "",
				"VELA_BUILD_NUMBER":        "0",
				"VELA_BUILD_PARENT":        "0",
				"VELA_BUILD_REF":           "",
				"VELA_BUILD_RUNTIME":       "",
				"VELA_BUILD_SENDER":        "",
				"VELA_BUILD_SOURCE":        "",
				"VELA_BUILD_STARTED":       "0",
				"VELA_BUILD_STATUS":        "",
				"VELA_BUILD_TITLE":         "",
				"VELA_BUILD_WORKSPACE":     "/vela/src",
				"VELA_CHANNEL":             "TODO",
				"VELA_DATABASE":            "TODO",
				"VELA_DISTRIBUTION":        "TODO",
				"VELA_HOST":                "TODO",
				"VELA_NETRC_MACHINE":       "TODO",
				"VELA_NETRC_PASSWORD":      "",
				"VELA_NETRC_USERNAME":      "x-oauth-basic",
				"VELA_QUEUE":               "TODO",
				"VELA_REPO_ACTIVE":         "false",
				"VELA_REPO_ALLOW_COMMENT":  "false",
				"VELA_REPO_ALLOW_DEPLOY":   "false",
				"VELA_REPO_ALLOW_PULL":     "false",
				"VELA_REPO_ALLOW_PUSH":     "false",
				"VELA_REPO_ALLOW_TAG":      "false",
				"VELA_REPO_BRANCH":         "",
				"VELA_REPO_CLONE":          "",
				"VELA_REPO_FULL_NAME":      "",
				"VELA_REPO_LINK":           "",
				"VELA_REPO_NAME":           "",
				"VELA_REPO_ORG":            "",
				"VELA_REPO_PIPELINE_TYPE":  "",
				"VELA_REPO_PRIVATE":        "false",
				"VELA_REPO_TIMEOUT":        "0",
				"VELA_REPO_TRUSTED":        "false",
				"VELA_REPO_VISIBILITY":     "",
				"VELA_RUNTIME":             "TODO",
				"VELA_SOURCE":              "TODO",
				"VELA_USER_ACTIVE":         "false",
				"VELA_USER_ADMIN":          "false",
				"VELA_USER_FAVORITES":      "[]",
				"VELA_USER_NAME":           "",
				"VELA_VERSION":             "TODO",
				"VELA_WORKSPACE":           "/vela/src",
				"HELLO":                    "Hello, Global Message",
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.EnvironmentServices(s, e)
	if err != nil {
		t.Errorf("EnvironmentServices returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("EnvironmentServices is %v, want %v", got, want)
	}
}

func TestNative_EnvironmentSecrets(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	e := raw.StringSliceMap{
		"HELLO": "Hello, Global Message",
	}

	str := "foo"
	s := yaml.SecretSlice{
		&yaml.Secret{
			Name: str,
			Origin: yaml.Origin{
				Image: "vault",
				Name:  str,
				Pull:  "always",
				Parameters: map[string]interface{}{
					"foo": "bar",
				},
				Environment: raw.StringSliceMap{
					"BUILD_CHANNEL": "foo",
				},
			},
		},
	}

	want := yaml.SecretSlice{
		&yaml.Secret{
			Name: str,
			Origin: yaml.Origin{
				Image: "vault",
				Name:  str,
				Pull:  "always",
				Parameters: map[string]interface{}{
					"foo": "bar",
				},
				Environment: raw.StringSliceMap{
					"BUILD_AUTHOR":             "",
					"BUILD_AUTHOR_EMAIL":       "",
					"BUILD_BASE_REF":           "",
					"BUILD_BRANCH":             "",
					"BUILD_CHANNEL":            "TODO",
					"BUILD_CLONE":              "",
					"BUILD_COMMIT":             "",
					"BUILD_CREATED":            "0",
					"BUILD_ENQUEUED":           "0",
					"BUILD_EVENT":              "",
					"BUILD_HOST":               "",
					"BUILD_LINK":               "",
					"BUILD_MESSAGE":            "",
					"BUILD_NUMBER":             "0",
					"BUILD_PARENT":             "0",
					"BUILD_REF":                "",
					"BUILD_SENDER":             "",
					"BUILD_SOURCE":             "",
					"BUILD_STARTED":            "0",
					"BUILD_STATUS":             "",
					"BUILD_TITLE":              "",
					"BUILD_WORKSPACE":          "/vela/src",
					"CI":                       "vela",
					"PARAMETER_FOO":            "bar",
					"REPOSITORY_ACTIVE":        "false",
					"REPOSITORY_ALLOW_COMMENT": "false",
					"REPOSITORY_ALLOW_DEPLOY":  "false",
					"REPOSITORY_ALLOW_PULL":    "false",
					"REPOSITORY_ALLOW_PUSH":    "false",
					"REPOSITORY_ALLOW_TAG":     "false",
					"REPOSITORY_BRANCH":        "",
					"REPOSITORY_CLONE":         "",
					"REPOSITORY_FULL_NAME":     "",
					"REPOSITORY_LINK":          "",
					"REPOSITORY_NAME":          "",
					"REPOSITORY_ORG":           "",
					"REPOSITORY_PRIVATE":       "false",
					"REPOSITORY_TIMEOUT":       "0",
					"REPOSITORY_TRUSTED":       "false",
					"REPOSITORY_VISIBILITY":    "",
					"VELA":                     "true",
					"VELA_ADDR":                "TODO",
					"VELA_BUILD_AUTHOR":        "",
					"VELA_BUILD_AUTHOR_EMAIL":  "",
					"VELA_BUILD_BASE_REF":      "",
					"VELA_BUILD_BRANCH":        "",
					"VELA_BUILD_CHANNEL":       "TODO",
					"VELA_BUILD_CLONE":         "",
					"VELA_BUILD_COMMIT":        "",
					"VELA_BUILD_CREATED":       "0",
					"VELA_BUILD_DISTRIBUTION":  "",
					"VELA_BUILD_ENQUEUED":      "0",
					"VELA_BUILD_EVENT":         "",
					"VELA_BUILD_HOST":          "",
					"VELA_BUILD_LINK":          "",
					"VELA_BUILD_MESSAGE":       "",
					"VELA_BUILD_NUMBER":        "0",
					"VELA_BUILD_PARENT":        "0",
					"VELA_BUILD_REF":           "",
					"VELA_BUILD_RUNTIME":       "",
					"VELA_BUILD_SENDER":        "",
					"VELA_BUILD_SOURCE":        "",
					"VELA_BUILD_STARTED":       "0",
					"VELA_BUILD_STATUS":        "",
					"VELA_BUILD_TITLE":         "",
					"VELA_BUILD_WORKSPACE":     "/vela/src",
					"VELA_CHANNEL":             "TODO",
					"VELA_DATABASE":            "TODO",
					"VELA_DISTRIBUTION":        "TODO",
					"VELA_HOST":                "TODO",
					"VELA_NETRC_MACHINE":       "TODO",
					"VELA_NETRC_PASSWORD":      "",
					"VELA_NETRC_USERNAME":      "x-oauth-basic",
					"VELA_QUEUE":               "TODO",
					"VELA_REPO_ACTIVE":         "false",
					"VELA_REPO_ALLOW_COMMENT":  "false",
					"VELA_REPO_ALLOW_DEPLOY":   "false",
					"VELA_REPO_ALLOW_PULL":     "false",
					"VELA_REPO_ALLOW_PUSH":     "false",
					"VELA_REPO_ALLOW_TAG":      "false",
					"VELA_REPO_BRANCH":         "",
					"VELA_REPO_CLONE":          "",
					"VELA_REPO_FULL_NAME":      "",
					"VELA_REPO_LINK":           "",
					"VELA_REPO_NAME":           "",
					"VELA_REPO_ORG":            "",
					"VELA_REPO_PIPELINE_TYPE":  "",
					"VELA_REPO_PRIVATE":        "false",
					"VELA_REPO_TIMEOUT":        "0",
					"VELA_REPO_TRUSTED":        "false",
					"VELA_REPO_VISIBILITY":     "",
					"VELA_RUNTIME":             "TODO",
					"VELA_SOURCE":              "TODO",
					"VELA_USER_ACTIVE":         "false",
					"VELA_USER_ADMIN":          "false",
					"VELA_USER_FAVORITES":      "[]",
					"VELA_USER_NAME":           "",
					"VELA_VERSION":             "TODO",
					"VELA_WORKSPACE":           "/vela/src",
					"HELLO":                    "Hello, Global Message",
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.EnvironmentSecrets(s, e)
	if err != nil {
		t.Errorf("EnvironmentSecrets returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("EnvironmentSecrets is %v, want %v", got, want)
	}
}

func TestNative_environment(t *testing.T) {
	// setup types
	booL := false
	num := 1
	num64 := int64(num)
	str := "foo"
	workspace := "/vela/src/foo/foo/foo"
	// push
	push := "push"
	// tag
	tag := "tag"
	tagref := "refs/tags/1"
	// pull_request
	pull := "pull_request"
	pullref := "refs/pull/1/head"
	// deployment
	deploy := "deployment"
	target := "production"

	tests := []struct {
		w    string
		b    *library.Build
		m    *types.Metadata
		r    *library.Repo
		u    *library.User
		want map[string]string
	}{
		// push
		{
			w:    workspace,
			b:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &push, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &str, BaseRef: &str},
			m:    &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			r:    &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			u:    &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			want: map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "push", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "foo", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "vela", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_COMMENT": "false", "REPOSITORY_ALLOW_DEPLOY": "false", "REPOSITORY_ALLOW_PULL": "false", "REPOSITORY_ALLOW_PUSH": "false", "REPOSITORY_ALLOW_TAG": "false", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CHANNEL": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "push", "VELA_BUILD_HOST": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "foo", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_COMMENT": "false", "VELA_REPO_ALLOW_DEPLOY": "false", "VELA_REPO_ALLOW_PULL": "false", "VELA_REPO_ALLOW_PUSH": "false", "VELA_REPO_ALLOW_TAG": "false", "VELA_REPO_BRANCH": "foo", "VELA_REPO_CLONE": "foo", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo"},
		},
		// tag
		{
			w:    workspace,
			b:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &tag, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &tagref, BaseRef: &str, PusherName: &str, PusherEmail: &str},
			m:    &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			r:    &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			u:    &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			want: map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "tag", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "refs/tags/1", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TAG": "1", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "vela", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_COMMENT": "false", "REPOSITORY_ALLOW_DEPLOY": "false", "REPOSITORY_ALLOW_PULL": "false", "REPOSITORY_ALLOW_PUSH": "false", "REPOSITORY_ALLOW_TAG": "false", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CHANNEL": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "tag", "VELA_BUILD_HOST": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "refs/tags/1", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TAG": "1", "VELA_BUILD_TAG_AUTHOR": "foo", "VELA_BUILD_TAG_AUTHOR_EMAIL": "foo", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_COMMENT": "false", "VELA_REPO_ALLOW_DEPLOY": "false", "VELA_REPO_ALLOW_PULL": "false", "VELA_REPO_ALLOW_PUSH": "false", "VELA_REPO_ALLOW_TAG": "false", "VELA_REPO_BRANCH": "foo", "VELA_REPO_CLONE": "foo", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo"},
		},
		// pull_request
		{
			w:    workspace,
			b:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &pull, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &pullref, BaseRef: &str},
			m:    &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			r:    &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			u:    &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			want: map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "pull_request", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_PULL_REQUEST_NUMBER": "1", "BUILD_REF": "refs/pull/1/head", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "vela", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_COMMENT": "false", "REPOSITORY_ALLOW_DEPLOY": "false", "REPOSITORY_ALLOW_PULL": "false", "REPOSITORY_ALLOW_PUSH": "false", "REPOSITORY_ALLOW_TAG": "false", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CHANNEL": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "pull_request", "VELA_BUILD_HOST": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_PULL_REQUEST": "1", "VELA_BUILD_REF": "refs/pull/1/head", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_PULL_REQUEST": "1", "VELA_PULL_REQUEST_SOURCE": "", "VELA_PULL_REQUEST_TARGET": "foo", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_COMMENT": "false", "VELA_REPO_ALLOW_DEPLOY": "false", "VELA_REPO_ALLOW_PULL": "false", "VELA_REPO_ALLOW_PUSH": "false", "VELA_REPO_ALLOW_TAG": "false", "VELA_REPO_BRANCH": "foo", "VELA_REPO_CLONE": "foo", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo"},
		},
		// deployment
		{
			w:    workspace,
			b:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &deploy, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &target, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &pullref, BaseRef: &str},
			m:    &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			r:    &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			u:    &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			want: map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "deployment", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "refs/pull/1/head", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TARGET": "production", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "vela", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_COMMENT": "false", "REPOSITORY_ALLOW_DEPLOY": "false", "REPOSITORY_ALLOW_PULL": "false", "REPOSITORY_ALLOW_PUSH": "false", "REPOSITORY_ALLOW_TAG": "false", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CHANNEL": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "deployment", "VELA_BUILD_HOST": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "refs/pull/1/head", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TARGET": "production", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DEPLOYMENT": "production", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_COMMENT": "false", "VELA_REPO_ALLOW_DEPLOY": "false", "VELA_REPO_ALLOW_PULL": "false", "VELA_REPO_ALLOW_PUSH": "false", "VELA_REPO_ALLOW_TAG": "false", "VELA_REPO_BRANCH": "foo", "VELA_REPO_CLONE": "foo", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo"},
		},
	}

	// run test
	for _, test := range tests {
		got := environment(test.b, test.m, test.r, test.u)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("environment is %v, want %v", got, test.want)
		}
	}
}

func Test_mergeMap(t *testing.T) {
	type args struct {
		combinedMap map[string]string
		loopMap     map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"empty", args{
			combinedMap: map[string]string{},
			loopMap:     map[string]string{},
		}, map[string]string{}},
		{"content", args{
			combinedMap: map[string]string{
				"VELA_FOO": "bar",
			},
			loopMap: map[string]string{
				"VELA_TEST": "foo",
			},
		}, map[string]string{
			"VELA_FOO":  "bar",
			"VELA_TEST": "foo",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendMap(tt.args.combinedMap, tt.args.loopMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_EnvironmentBuild(t *testing.T) {
	// setup types
	booL := false
	num := 1
	num64 := int64(num)
	str := "foo"
	//workspace := "/vela/src/foo/foo/foo"
	// push
	push := "push"
	// tag
	tag := "tag"
	tagref := "refs/tags/1"
	// pull_request
	pull := "pull_request"
	pullref := "refs/pull/1/head"
	// deployment
	deploy := "deployment"
	target := "production"
	type fields struct {
		build    *library.Build
		metadata *types.Metadata
		repo     *library.Repo
		user     *library.User
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{"push", fields{
			build:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &push, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &str, BaseRef: &str},
			metadata: &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			repo:     &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			user:     &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
		}, map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "push", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "foo", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "vela", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_COMMENT": "false", "REPOSITORY_ALLOW_DEPLOY": "false", "REPOSITORY_ALLOW_PULL": "false", "REPOSITORY_ALLOW_PUSH": "false", "REPOSITORY_ALLOW_TAG": "false", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CHANNEL": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "push", "VELA_BUILD_HOST": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "foo", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_COMMENT": "false", "VELA_REPO_ALLOW_DEPLOY": "false", "VELA_REPO_ALLOW_PULL": "false", "VELA_REPO_ALLOW_PUSH": "false", "VELA_REPO_ALLOW_TAG": "false", "VELA_REPO_BRANCH": "foo", "VELA_REPO_CLONE": "foo", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo"}},
		{"tag", fields{
			build:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &tag, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &tagref, BaseRef: &str, PusherName: &str, PusherEmail: &str},
			metadata: &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			repo:     &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			user:     &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
		}, map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "tag", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "refs/tags/1", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TAG": "1", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "vela", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_COMMENT": "false", "REPOSITORY_ALLOW_DEPLOY": "false", "REPOSITORY_ALLOW_PULL": "false", "REPOSITORY_ALLOW_PUSH": "false", "REPOSITORY_ALLOW_TAG": "false", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CHANNEL": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "tag", "VELA_BUILD_HOST": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "refs/tags/1", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TAG": "1", "VELA_BUILD_TAG_AUTHOR": "foo", "VELA_BUILD_TAG_AUTHOR_EMAIL": "foo", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_COMMENT": "false", "VELA_REPO_ALLOW_DEPLOY": "false", "VELA_REPO_ALLOW_PULL": "false", "VELA_REPO_ALLOW_PUSH": "false", "VELA_REPO_ALLOW_TAG": "false", "VELA_REPO_BRANCH": "foo", "VELA_REPO_CLONE": "foo", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo"},
		},
		{"pull_request", fields{
			build:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &pull, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &pullref, BaseRef: &str},
			metadata: &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			repo:     &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			user:     &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
		}, map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "pull_request", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_PULL_REQUEST_NUMBER": "1", "BUILD_REF": "refs/pull/1/head", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "vela", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_COMMENT": "false", "REPOSITORY_ALLOW_DEPLOY": "false", "REPOSITORY_ALLOW_PULL": "false", "REPOSITORY_ALLOW_PUSH": "false", "REPOSITORY_ALLOW_TAG": "false", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CHANNEL": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "pull_request", "VELA_BUILD_HOST": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_PULL_REQUEST": "1", "VELA_BUILD_REF": "refs/pull/1/head", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_PULL_REQUEST": "1", "VELA_PULL_REQUEST_SOURCE": "", "VELA_PULL_REQUEST_TARGET": "foo", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_COMMENT": "false", "VELA_REPO_ALLOW_DEPLOY": "false", "VELA_REPO_ALLOW_PULL": "false", "VELA_REPO_ALLOW_PUSH": "false", "VELA_REPO_ALLOW_TAG": "false", "VELA_REPO_BRANCH": "foo", "VELA_REPO_CLONE": "foo", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo"},
		},
		{"deployment", fields{
			build:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &deploy, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &target, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &pullref, BaseRef: &str},
			metadata: &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			repo:     &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			user:     &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
		}, map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "deployment", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "refs/pull/1/head", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TARGET": "production", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "vela", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_COMMENT": "false", "REPOSITORY_ALLOW_DEPLOY": "false", "REPOSITORY_ALLOW_PULL": "false", "REPOSITORY_ALLOW_PUSH": "false", "REPOSITORY_ALLOW_TAG": "false", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CHANNEL": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "deployment", "VELA_BUILD_HOST": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "refs/pull/1/head", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TARGET": "production", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DEPLOYMENT": "production", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_COMMENT": "false", "VELA_REPO_ALLOW_DEPLOY": "false", "VELA_REPO_ALLOW_PULL": "false", "VELA_REPO_ALLOW_PUSH": "false", "VELA_REPO_ALLOW_TAG": "false", "VELA_REPO_BRANCH": "foo", "VELA_REPO_CLONE": "foo", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{
				build:    tt.fields.build,
				metadata: tt.fields.metadata,
				repo:     tt.fields.repo,
				user:     tt.fields.user,
			}
			if got := c.EnvironmentBuild(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvironmentBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
