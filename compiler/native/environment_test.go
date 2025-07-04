// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
	"github.com/go-vela/server/internal"
)

func TestNative_EnvironmentStages(t *testing.T) {
	// setup types
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

	env := environment(nil, nil, nil, nil, nil)
	env["HELLO"] = "Hello, Global Message"

	want := yaml.StageSlice{
		&yaml.Stage{
			Name:        str,
			Environment: env,
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
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
	e := raw.StringSliceMap{
		"HELLO": "Hello, Stage Message",
	}

	str := "foo"
	s := yaml.StepSlice{
		&yaml.Step{
			Image: "alpine",
			Name:  str,
			Pull:  "always",
			Environment: raw.StringSliceMap{
				"BUILD_WORKSPACE": "foo",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Image: "alpine",
			Name:  str,
			Pull:  "always",
			Environment: raw.StringSliceMap{
				"BUILD_AUTHOR":               "",
				"BUILD_AUTHOR_EMAIL":         "",
				"BUILD_BASE_REF":             "",
				"BUILD_BRANCH":               "",
				"BUILD_CLONE":                "",
				"BUILD_COMMIT":               "",
				"BUILD_CREATED":              "0",
				"BUILD_ENQUEUED":             "0",
				"BUILD_EVENT":                "",
				"BUILD_HOST":                 "",
				"BUILD_LINK":                 "",
				"BUILD_MESSAGE":              "",
				"BUILD_NUMBER":               "0",
				"BUILD_PARENT":               "0",
				"BUILD_REF":                  "",
				"BUILD_SENDER":               "",
				"BUILD_SOURCE":               "",
				"BUILD_STARTED":              "0",
				"BUILD_STATUS":               "",
				"BUILD_TITLE":                "",
				"BUILD_WORKSPACE":            "/vela/src",
				"CI":                         "true",
				"REPOSITORY_ACTIVE":          "false",
				"REPOSITORY_ALLOW_EVENTS":    "",
				"REPOSITORY_BRANCH":          "",
				"REPOSITORY_CLONE":           "",
				"REPOSITORY_FULL_NAME":       "",
				"REPOSITORY_LINK":            "",
				"REPOSITORY_NAME":            "",
				"REPOSITORY_ORG":             "",
				"REPOSITORY_PRIVATE":         "false",
				"REPOSITORY_TIMEOUT":         "0",
				"REPOSITORY_TRUSTED":         "false",
				"REPOSITORY_VISIBILITY":      "",
				"VELA":                       "true",
				"VELA_ADDR":                  "TODO",
				"VELA_BUILD_APPROVED_AT":     "0",
				"VELA_BUILD_APPROVED_BY":     "",
				"VELA_BUILD_AUTHOR":          "",
				"VELA_BUILD_AUTHOR_EMAIL":    "",
				"VELA_BUILD_BASE_REF":        "",
				"VELA_BUILD_BRANCH":          "",
				"VELA_BUILD_CLONE":           "",
				"VELA_BUILD_COMMIT":          "",
				"VELA_BUILD_CREATED":         "0",
				"VELA_BUILD_DISTRIBUTION":    "",
				"VELA_BUILD_ENQUEUED":        "0",
				"VELA_BUILD_EVENT":           "",
				"VELA_BUILD_EVENT_ACTION":    "",
				"VELA_BUILD_HOST":            "",
				"VELA_BUILD_ROUTE":           "",
				"VELA_BUILD_LINK":            "",
				"VELA_BUILD_MESSAGE":         "",
				"VELA_BUILD_NUMBER":          "0",
				"VELA_BUILD_PARENT":          "0",
				"VELA_BUILD_REF":             "",
				"VELA_BUILD_RUNTIME":         "",
				"VELA_BUILD_SENDER":          "",
				"VELA_BUILD_SENDER_SCM_ID":   "",
				"VELA_BUILD_SOURCE":          "",
				"VELA_BUILD_STARTED":         "0",
				"VELA_BUILD_STATUS":          "",
				"VELA_BUILD_TITLE":           "",
				"VELA_BUILD_WORKSPACE":       "/vela/src",
				"VELA_DATABASE":              "TODO",
				"VELA_DISTRIBUTION":          "TODO",
				"VELA_HOST":                  "TODO",
				"VELA_NETRC_MACHINE":         "TODO",
				"VELA_NETRC_PASSWORD":        "TODO",
				"VELA_NETRC_USERNAME":        "x-oauth-basic",
				"VELA_QUEUE":                 "TODO",
				"VELA_REPO_ACTIVE":           "false",
				"VELA_REPO_ALLOW_EVENTS":     "",
				"VELA_REPO_APPROVAL_TIMEOUT": "0",
				"VELA_REPO_APPROVE_BUILD":    "",
				"VELA_REPO_BRANCH":           "",
				"VELA_REPO_TOPICS":           "",
				"VELA_REPO_BUILD_LIMIT":      "0",
				"VELA_REPO_CLONE":            "",
				"VELA_REPO_CUSTOM_PROPS":     "{}",
				"VELA_REPO_INSTALL_ID":       "0",
				"VELA_REPO_FULL_NAME":        "",
				"VELA_REPO_LINK":             "",
				"VELA_REPO_NAME":             "",
				"VELA_REPO_ORG":              "",
				"VELA_REPO_OWNER":            "",
				"VELA_REPO_PIPELINE_TYPE":    "",
				"VELA_REPO_PRIVATE":          "false",
				"VELA_REPO_TIMEOUT":          "0",
				"VELA_REPO_TRUSTED":          "false",
				"VELA_REPO_VISIBILITY":       "",
				"VELA_RUNTIME":               "TODO",
				"VELA_SOURCE":                "TODO",
				"VELA_USER_ACTIVE":           "false",
				"VELA_USER_ADMIN":            "false",
				"VELA_USER_FAVORITES":        "[]",
				"VELA_USER_NAME":             "",
				"VELA_VERSION":               "TODO",
				"VELA_WORKSPACE":             "/vela/src",
				"HELLO":                      "Hello, Stage Message",
			},
		},
	}

	// run test non-local
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	// run test local
	compiler.WithLocal(true)

	t.Setenv("VELA_BUILD_COMMIT", "123abc")

	got, err := compiler.EnvironmentSteps(s, e)
	if err != nil {
		t.Errorf("EnvironmentSteps returned err: %v", err)
	}

	// cannot use complete diff since local compiler pulls from OS env
	if !strings.EqualFold(got[0].Environment["VELA_BUILD_COMMIT"], "123abc") {
		t.Errorf("EnvironmentSteps with local compiler should have set VELA_BUILD_COMMIT to 123abc, got %s", got[0].Environment["VELA_BUILD_COMMIT"])
	}

	// test without local
	compiler.WithLocal(false)

	// reset s
	s = yaml.StepSlice{
		&yaml.Step{
			Image: "alpine",
			Name:  str,
			Pull:  "always",
			Environment: raw.StringSliceMap{
				"BUILD_WORKSPACE": "foo",
			},
		},
	}

	got, err = compiler.EnvironmentSteps(s, e)
	if err != nil {
		t.Errorf("EnvironmentSteps returned err: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("EnvironmentSteps mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_EnvironmentServices(t *testing.T) {
	// setup types
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
				"BUILD_WORKSPACE": "foo",
			},
		},
	}

	want := yaml.ServiceSlice{
		&yaml.Service{
			Image: "postgres",
			Name:  str,
			Pull:  "always",
			Environment: raw.StringSliceMap{
				"BUILD_AUTHOR":               "",
				"BUILD_AUTHOR_EMAIL":         "",
				"BUILD_BASE_REF":             "",
				"BUILD_BRANCH":               "",
				"BUILD_CLONE":                "",
				"BUILD_COMMIT":               "",
				"BUILD_CREATED":              "0",
				"BUILD_ENQUEUED":             "0",
				"BUILD_EVENT":                "",
				"BUILD_HOST":                 "",
				"BUILD_LINK":                 "",
				"BUILD_MESSAGE":              "",
				"BUILD_NUMBER":               "0",
				"BUILD_PARENT":               "0",
				"BUILD_REF":                  "",
				"BUILD_SENDER":               "",
				"BUILD_SOURCE":               "",
				"BUILD_STARTED":              "0",
				"BUILD_STATUS":               "",
				"BUILD_TITLE":                "",
				"BUILD_WORKSPACE":            "/vela/src",
				"CI":                         "true",
				"REPOSITORY_ACTIVE":          "false",
				"REPOSITORY_ALLOW_EVENTS":    "",
				"REPOSITORY_BRANCH":          "",
				"REPOSITORY_CLONE":           "",
				"REPOSITORY_FULL_NAME":       "",
				"REPOSITORY_LINK":            "",
				"REPOSITORY_NAME":            "",
				"REPOSITORY_ORG":             "",
				"REPOSITORY_PRIVATE":         "false",
				"REPOSITORY_TIMEOUT":         "0",
				"REPOSITORY_TRUSTED":         "false",
				"REPOSITORY_VISIBILITY":      "",
				"VELA":                       "true",
				"VELA_ADDR":                  "TODO",
				"VELA_BUILD_APPROVED_AT":     "0",
				"VELA_BUILD_APPROVED_BY":     "",
				"VELA_BUILD_AUTHOR":          "",
				"VELA_BUILD_AUTHOR_EMAIL":    "",
				"VELA_BUILD_BASE_REF":        "",
				"VELA_BUILD_BRANCH":          "",
				"VELA_BUILD_CLONE":           "",
				"VELA_BUILD_COMMIT":          "",
				"VELA_BUILD_CREATED":         "0",
				"VELA_BUILD_DISTRIBUTION":    "",
				"VELA_BUILD_ENQUEUED":        "0",
				"VELA_BUILD_EVENT":           "",
				"VELA_BUILD_EVENT_ACTION":    "",
				"VELA_BUILD_HOST":            "",
				"VELA_BUILD_ROUTE":           "",
				"VELA_BUILD_LINK":            "",
				"VELA_BUILD_MESSAGE":         "",
				"VELA_BUILD_NUMBER":          "0",
				"VELA_BUILD_PARENT":          "0",
				"VELA_BUILD_REF":             "",
				"VELA_BUILD_RUNTIME":         "",
				"VELA_BUILD_SENDER":          "",
				"VELA_BUILD_SENDER_SCM_ID":   "",
				"VELA_BUILD_SOURCE":          "",
				"VELA_BUILD_STARTED":         "0",
				"VELA_BUILD_STATUS":          "",
				"VELA_BUILD_TITLE":           "",
				"VELA_BUILD_WORKSPACE":       "/vela/src",
				"VELA_DATABASE":              "TODO",
				"VELA_DISTRIBUTION":          "TODO",
				"VELA_HOST":                  "TODO",
				"VELA_NETRC_MACHINE":         "TODO",
				"VELA_NETRC_PASSWORD":        "TODO",
				"VELA_NETRC_USERNAME":        "x-oauth-basic",
				"VELA_QUEUE":                 "TODO",
				"VELA_REPO_ACTIVE":           "false",
				"VELA_REPO_ALLOW_EVENTS":     "",
				"VELA_REPO_APPROVAL_TIMEOUT": "0",
				"VELA_REPO_APPROVE_BUILD":    "",
				"VELA_REPO_BRANCH":           "",
				"VELA_REPO_TOPICS":           "",
				"VELA_REPO_BUILD_LIMIT":      "0",
				"VELA_REPO_CLONE":            "",
				"VELA_REPO_CUSTOM_PROPS":     "{}",
				"VELA_REPO_INSTALL_ID":       "0",
				"VELA_REPO_FULL_NAME":        "",
				"VELA_REPO_LINK":             "",
				"VELA_REPO_NAME":             "",
				"VELA_REPO_ORG":              "",
				"VELA_REPO_OWNER":            "",
				"VELA_REPO_PIPELINE_TYPE":    "",
				"VELA_REPO_PRIVATE":          "false",
				"VELA_REPO_TIMEOUT":          "0",
				"VELA_REPO_TRUSTED":          "false",
				"VELA_REPO_VISIBILITY":       "",
				"VELA_RUNTIME":               "TODO",
				"VELA_SOURCE":                "TODO",
				"VELA_USER_ACTIVE":           "false",
				"VELA_USER_ADMIN":            "false",
				"VELA_USER_FAVORITES":        "[]",
				"VELA_USER_NAME":             "",
				"VELA_VERSION":               "TODO",
				"VELA_WORKSPACE":             "/vela/src",
				"HELLO":                      "Hello, Global Message",
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.EnvironmentServices(s, e)
	if err != nil {
		t.Errorf("EnvironmentServices returned err: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("EnvironmentServices mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_EnvironmentSecrets(t *testing.T) {
	// setup types
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
					"BUILD_WORKSPACE": "foo",
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
					"BUILD_AUTHOR":               "",
					"BUILD_AUTHOR_EMAIL":         "",
					"BUILD_BASE_REF":             "",
					"BUILD_BRANCH":               "",
					"BUILD_CLONE":                "",
					"BUILD_COMMIT":               "",
					"BUILD_CREATED":              "0",
					"BUILD_ENQUEUED":             "0",
					"BUILD_EVENT":                "",
					"BUILD_HOST":                 "",
					"BUILD_LINK":                 "",
					"BUILD_MESSAGE":              "",
					"BUILD_NUMBER":               "0",
					"BUILD_PARENT":               "0",
					"BUILD_REF":                  "",
					"BUILD_SENDER":               "",
					"BUILD_SOURCE":               "",
					"BUILD_STARTED":              "0",
					"BUILD_STATUS":               "",
					"BUILD_TITLE":                "",
					"BUILD_WORKSPACE":            "/vela/src",
					"CI":                         "true",
					"PARAMETER_FOO":              "bar",
					"REPOSITORY_ACTIVE":          "false",
					"REPOSITORY_ALLOW_EVENTS":    "",
					"REPOSITORY_BRANCH":          "",
					"REPOSITORY_CLONE":           "",
					"REPOSITORY_FULL_NAME":       "",
					"REPOSITORY_LINK":            "",
					"REPOSITORY_NAME":            "",
					"REPOSITORY_ORG":             "",
					"REPOSITORY_PRIVATE":         "false",
					"REPOSITORY_TIMEOUT":         "0",
					"REPOSITORY_TRUSTED":         "false",
					"REPOSITORY_VISIBILITY":      "",
					"VELA":                       "true",
					"VELA_ADDR":                  "TODO",
					"VELA_BUILD_APPROVED_AT":     "0",
					"VELA_BUILD_APPROVED_BY":     "",
					"VELA_BUILD_AUTHOR":          "",
					"VELA_BUILD_AUTHOR_EMAIL":    "",
					"VELA_BUILD_BASE_REF":        "",
					"VELA_BUILD_BRANCH":          "",
					"VELA_BUILD_CLONE":           "",
					"VELA_BUILD_COMMIT":          "",
					"VELA_BUILD_CREATED":         "0",
					"VELA_BUILD_DISTRIBUTION":    "",
					"VELA_BUILD_ENQUEUED":        "0",
					"VELA_BUILD_EVENT":           "",
					"VELA_BUILD_EVENT_ACTION":    "",
					"VELA_BUILD_HOST":            "",
					"VELA_BUILD_ROUTE":           "",
					"VELA_BUILD_LINK":            "",
					"VELA_BUILD_MESSAGE":         "",
					"VELA_BUILD_NUMBER":          "0",
					"VELA_BUILD_PARENT":          "0",
					"VELA_BUILD_REF":             "",
					"VELA_BUILD_RUNTIME":         "",
					"VELA_BUILD_SENDER":          "",
					"VELA_BUILD_SENDER_SCM_ID":   "",
					"VELA_BUILD_SOURCE":          "",
					"VELA_BUILD_STARTED":         "0",
					"VELA_BUILD_STATUS":          "",
					"VELA_BUILD_TITLE":           "",
					"VELA_BUILD_WORKSPACE":       "/vela/src",
					"VELA_DATABASE":              "TODO",
					"VELA_DISTRIBUTION":          "TODO",
					"VELA_HOST":                  "TODO",
					"VELA_NETRC_MACHINE":         "TODO",
					"VELA_NETRC_PASSWORD":        "TODO",
					"VELA_NETRC_USERNAME":        "x-oauth-basic",
					"VELA_QUEUE":                 "TODO",
					"VELA_REPO_ACTIVE":           "false",
					"VELA_REPO_ALLOW_EVENTS":     "",
					"VELA_REPO_APPROVAL_TIMEOUT": "0",
					"VELA_REPO_APPROVE_BUILD":    "",
					"VELA_REPO_BRANCH":           "",
					"VELA_REPO_TOPICS":           "",
					"VELA_REPO_BUILD_LIMIT":      "0",
					"VELA_REPO_CLONE":            "",
					"VELA_REPO_CUSTOM_PROPS":     "{}",
					"VELA_REPO_INSTALL_ID":       "0",
					"VELA_REPO_FULL_NAME":        "",
					"VELA_REPO_LINK":             "",
					"VELA_REPO_NAME":             "",
					"VELA_REPO_ORG":              "",
					"VELA_REPO_OWNER":            "",
					"VELA_REPO_PIPELINE_TYPE":    "",
					"VELA_REPO_PRIVATE":          "false",
					"VELA_REPO_TIMEOUT":          "0",
					"VELA_REPO_TRUSTED":          "false",
					"VELA_REPO_VISIBILITY":       "",
					"VELA_RUNTIME":               "TODO",
					"VELA_SOURCE":                "TODO",
					"VELA_USER_ACTIVE":           "false",
					"VELA_USER_ADMIN":            "false",
					"VELA_USER_FAVORITES":        "[]",
					"VELA_USER_NAME":             "",
					"VELA_VERSION":               "TODO",
					"VELA_WORKSPACE":             "/vela/src",
					"HELLO":                      "Hello, Global Message",
				},
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.EnvironmentSecrets(s, e)
	if err != nil {
		t.Errorf("EnvironmentSecrets returned err: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("EnvironmentSecrets mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_environment(t *testing.T) {
	// setup types
	booL := false
	num32 := int32(1)
	num64 := int64(1)
	str := "foo"
	workspace := "/vela/src/foo/foo/foo"
	topics := []string{"cloud", "security"}
	props := map[string]any{"foo": "bar"}
	// push
	push := "push"
	// tag
	tag := "tag"
	tagref := "refs/tags/1"
	// pull_request
	pull := "pull_request"
	pullact := "opened"
	pullref := "refs/pull/1/head"
	// deployment
	deploy := "deployment"
	target := "production"
	// netrc
	netrc := "foo"

	tests := []struct {
		w     string
		b     *api.Build
		m     *internal.Metadata
		r     *api.Repo
		u     *api.User
		netrc *string
		want  map[string]string
	}{
		// push
		{
			w:     workspace,
			b:     &api.Build{ID: &num64, Repo: &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL}, Number: &num64, Parent: &num64, Event: &push, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, SenderSCMID: &str, Author: &str, Branch: &str, Ref: &str, BaseRef: &str},
			m:     &internal.Metadata{Database: &internal.Database{Driver: str, Host: str}, Queue: &internal.Queue{Driver: str, Host: str}, Source: &internal.Source{Driver: str, Host: str}, Vela: &internal.Vela{Address: str, WebAddress: str, OpenIDIssuer: str}},
			r:     &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL},
			u:     &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			netrc: &netrc,
			want:  map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "push", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "foo", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "true", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_EVENTS": "", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_SERVER_ADDR": "foo", "VELA_OPEN_ID_ISSUER": "foo", "VELA_BUILD_APPROVED_AT": "0", "VELA_BUILD_APPROVED_BY": "", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "push", "VELA_BUILD_EVENT_ACTION": "", "VELA_BUILD_HOST": "", "VELA_BUILD_ROUTE": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "foo", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SENDER_SCM_ID": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_EVENTS": "", "VELA_REPO_APPROVAL_TIMEOUT": "0", "VELA_REPO_APPROVE_BUILD": "", "VELA_REPO_BRANCH": "foo", "VELA_REPO_TOPICS": "cloud,security", "VELA_REPO_BUILD_LIMIT": "1", "VELA_REPO_CLONE": "foo", "VELA_REPO_CUSTOM_PROPS": `{"foo":"bar"}`, "VELA_REPO_INSTALL_ID": "0", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_OWNER": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_ID_TOKEN_REQUEST_URL": "foo/api/v1/repos/foo/builds/1/id_token"},
		},
		// tag
		{
			w:     workspace,
			b:     &api.Build{ID: &num64, Repo: &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL}, Number: &num64, Parent: &num64, Event: &tag, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, SenderSCMID: &str, Author: &str, Branch: &str, Ref: &tagref, BaseRef: &str},
			m:     &internal.Metadata{Database: &internal.Database{Driver: str, Host: str}, Queue: &internal.Queue{Driver: str, Host: str}, Source: &internal.Source{Driver: str, Host: str}, Vela: &internal.Vela{Address: str, WebAddress: str, OpenIDIssuer: str}},
			r:     &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL},
			u:     &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			netrc: &netrc,
			want:  map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "tag", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "refs/tags/1", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TAG": "1", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "true", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_EVENTS": "", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_SERVER_ADDR": "foo", "VELA_OPEN_ID_ISSUER": "foo", "VELA_BUILD_APPROVED_AT": "0", "VELA_BUILD_APPROVED_BY": "", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "tag", "VELA_BUILD_EVENT_ACTION": "", "VELA_BUILD_HOST": "", "VELA_BUILD_ROUTE": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "refs/tags/1", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SENDER_SCM_ID": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TAG": "1", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_EVENTS": "", "VELA_REPO_APPROVAL_TIMEOUT": "0", "VELA_REPO_APPROVE_BUILD": "", "VELA_REPO_BRANCH": "foo", "VELA_REPO_TOPICS": "cloud,security", "VELA_REPO_BUILD_LIMIT": "1", "VELA_REPO_CLONE": "foo", "VELA_REPO_CUSTOM_PROPS": `{"foo":"bar"}`, "VELA_REPO_INSTALL_ID": "0", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_OWNER": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_ID_TOKEN_REQUEST_URL": "foo/api/v1/repos/foo/builds/1/id_token"},
		},
		// pull_request
		{
			w:     workspace,
			b:     &api.Build{ID: &num64, Repo: &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL}, Number: &num64, Parent: &num64, Event: &pull, EventAction: &pullact, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, SenderSCMID: &str, Author: &str, Branch: &str, Ref: &pullref, BaseRef: &str},
			m:     &internal.Metadata{Database: &internal.Database{Driver: str, Host: str}, Queue: &internal.Queue{Driver: str, Host: str}, Source: &internal.Source{Driver: str, Host: str}, Vela: &internal.Vela{Address: str, WebAddress: str, OpenIDIssuer: str}},
			r:     &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL},
			u:     &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			netrc: &netrc,
			want:  map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "pull_request", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_PULL_REQUEST_NUMBER": "1", "BUILD_REF": "refs/pull/1/head", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "true", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_EVENTS": "", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_SERVER_ADDR": "foo", "VELA_OPEN_ID_ISSUER": "foo", "VELA_BUILD_APPROVED_AT": "0", "VELA_BUILD_APPROVED_BY": "", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "pull_request", "VELA_BUILD_EVENT_ACTION": "opened", "VELA_BUILD_HOST": "", "VELA_BUILD_ROUTE": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_PULL_REQUEST": "1", "VELA_BUILD_REF": "refs/pull/1/head", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SENDER_SCM_ID": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_PULL_REQUEST": "1", "VELA_PULL_REQUEST_FORK": "false", "VELA_PULL_REQUEST_SOURCE": "", "VELA_PULL_REQUEST_TARGET": "foo", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_EVENTS": "", "VELA_REPO_APPROVAL_TIMEOUT": "0", "VELA_REPO_APPROVE_BUILD": "", "VELA_REPO_BRANCH": "foo", "VELA_REPO_TOPICS": "cloud,security", "VELA_REPO_BUILD_LIMIT": "1", "VELA_REPO_CLONE": "foo", "VELA_REPO_CUSTOM_PROPS": `{"foo":"bar"}`, "VELA_REPO_INSTALL_ID": "0", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_OWNER": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_ID_TOKEN_REQUEST_URL": "foo/api/v1/repos/foo/builds/1/id_token"},
		},
		// deployment
		{
			w:     workspace,
			b:     &api.Build{ID: &num64, Repo: &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL}, Number: &num64, Parent: &num64, Event: &deploy, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &target, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, SenderSCMID: &str, Author: &str, Branch: &str, Ref: &pullref, BaseRef: &str},
			m:     &internal.Metadata{Database: &internal.Database{Driver: str, Host: str}, Queue: &internal.Queue{Driver: str, Host: str}, Source: &internal.Source{Driver: str, Host: str}, Vela: &internal.Vela{Address: str, WebAddress: str, OpenIDIssuer: str}},
			r:     &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL},
			u:     &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			netrc: &netrc,
			want:  map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "deployment", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "refs/pull/1/head", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TARGET": "production", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "true", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_EVENTS": "", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_SERVER_ADDR": "foo", "VELA_OPEN_ID_ISSUER": "foo", "VELA_BUILD_APPROVED_AT": "0", "VELA_BUILD_APPROVED_BY": "", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "deployment", "VELA_BUILD_EVENT_ACTION": "", "VELA_BUILD_HOST": "", "VELA_BUILD_ROUTE": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "refs/pull/1/head", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SENDER_SCM_ID": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TARGET": "production", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_DATABASE": "foo", "VELA_DEPLOYMENT": "production", "VELA_DEPLOYMENT_NUMBER": "0", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_EVENTS": "", "VELA_REPO_APPROVAL_TIMEOUT": "0", "VELA_REPO_APPROVE_BUILD": "", "VELA_REPO_BRANCH": "foo", "VELA_REPO_TOPICS": "cloud,security", "VELA_REPO_BUILD_LIMIT": "1", "VELA_REPO_CLONE": "foo", "VELA_REPO_CUSTOM_PROPS": `{"foo":"bar"}`, "VELA_REPO_INSTALL_ID": "0", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_OWNER": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_ID_TOKEN_REQUEST_URL": "foo/api/v1/repos/foo/builds/1/id_token"},
		},
		// netrc
		{
			w:     workspace,
			b:     &api.Build{ID: &num64, Repo: &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL}, Number: &num64, Parent: &num64, Event: &deploy, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &target, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, SenderSCMID: &str, Author: &str, Branch: &str, Ref: &pullref, BaseRef: &str},
			m:     &internal.Metadata{Database: &internal.Database{Driver: str, Host: str}, Queue: &internal.Queue{Driver: str, Host: str}, Source: &internal.Source{Driver: str, Host: str}, Vela: &internal.Vela{Address: str, WebAddress: str, OpenIDIssuer: str}},
			r:     &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL},
			u:     &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			netrc: nil,
			want:  map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "deployment", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "refs/pull/1/head", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TARGET": "production", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "true", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_EVENTS": "", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_SERVER_ADDR": "foo", "VELA_OPEN_ID_ISSUER": "foo", "VELA_BUILD_APPROVED_AT": "0", "VELA_BUILD_APPROVED_BY": "", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "deployment", "VELA_BUILD_EVENT_ACTION": "", "VELA_BUILD_HOST": "", "VELA_BUILD_ROUTE": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "refs/pull/1/head", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SENDER_SCM_ID": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TARGET": "production", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_DATABASE": "foo", "VELA_DEPLOYMENT": "production", "VELA_DEPLOYMENT_NUMBER": "0", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "TODO", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_EVENTS": "", "VELA_REPO_APPROVAL_TIMEOUT": "0", "VELA_REPO_APPROVE_BUILD": "", "VELA_REPO_BRANCH": "foo", "VELA_REPO_TOPICS": "cloud,security", "VELA_REPO_BUILD_LIMIT": "1", "VELA_REPO_CLONE": "foo", "VELA_REPO_CUSTOM_PROPS": `{"foo":"bar"}`, "VELA_REPO_INSTALL_ID": "0", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_OWNER": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_ID_TOKEN_REQUEST_URL": "foo/api/v1/repos/foo/builds/1/id_token"},
		},
	}

	// run test
	for _, test := range tests {
		got := environment(test.b, test.m, test.r, test.u, test.netrc)

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("environment mismatch (-want +got):\n%s", diff)
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
	num32 := int32(1)
	num64 := int64(1)
	str := "foo"
	topics := []string{"cloud", "security"}
	props := map[string]any{"foo": "bar"}
	//workspace := "/vela/src/foo/foo/foo"
	// push
	push := "push"
	// tag
	tag := "tag"
	tagref := "refs/tags/1"
	// pull_request
	pull := "pull_request"
	pullact := "opened"
	pullref := "refs/pull/1/head"
	// deployment
	deploy := "deployment"
	target := "production"
	// netrc
	netrc := "foo"

	type fields struct {
		build    *api.Build
		metadata *internal.Metadata
		repo     *api.Repo
		user     *api.User
		netrc    *string
	}

	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{"push", fields{
			build:    &api.Build{ID: &num64, Repo: &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL}, Number: &num64, Parent: &num64, Event: &push, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, SenderSCMID: &str, Author: &str, Branch: &str, Ref: &str, BaseRef: &str},
			metadata: &internal.Metadata{Database: &internal.Database{Driver: str, Host: str}, Queue: &internal.Queue{Driver: str, Host: str}, Source: &internal.Source{Driver: str, Host: str}, Vela: &internal.Vela{Address: str, WebAddress: str, OpenIDIssuer: str}},
			repo:     &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL},
			user:     &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			netrc:    &netrc,
		}, map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "push", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "foo", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "true", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_EVENTS": "", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_SERVER_ADDR": "foo", "VELA_OPEN_ID_ISSUER": "foo", "VELA_BUILD_APPROVED_AT": "0", "VELA_BUILD_APPROVED_BY": "", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "push", "VELA_BUILD_EVENT_ACTION": "", "VELA_BUILD_HOST": "", "VELA_BUILD_ROUTE": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "foo", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SENDER_SCM_ID": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_EVENTS": "", "VELA_REPO_APPROVAL_TIMEOUT": "0", "VELA_REPO_APPROVE_BUILD": "", "VELA_REPO_OWNER": "foo", "VELA_REPO_BRANCH": "foo", "VELA_REPO_BUILD_LIMIT": "1", "VELA_REPO_CLONE": "foo", "VELA_REPO_CUSTOM_PROPS": `{"foo":"bar"}`, "VELA_REPO_INSTALL_ID": "0", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TOPICS": "cloud,security", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_ID_TOKEN_REQUEST_URL": "foo/api/v1/repos/foo/builds/1/id_token"}},
		{"tag", fields{
			build:    &api.Build{ID: &num64, Repo: &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL}, Number: &num64, Parent: &num64, Event: &tag, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, SenderSCMID: &str, Author: &str, Branch: &str, Ref: &tagref, BaseRef: &str},
			metadata: &internal.Metadata{Database: &internal.Database{Driver: str, Host: str}, Queue: &internal.Queue{Driver: str, Host: str}, Source: &internal.Source{Driver: str, Host: str}, Vela: &internal.Vela{Address: str, WebAddress: str, OpenIDIssuer: str}},
			repo:     &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL},
			user:     &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			netrc:    &netrc,
		}, map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "tag", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "refs/tags/1", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TAG": "1", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "true", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_EVENTS": "", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_SERVER_ADDR": "foo", "VELA_OPEN_ID_ISSUER": "foo", "VELA_BUILD_APPROVED_AT": "0", "VELA_BUILD_APPROVED_BY": "", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "tag", "VELA_BUILD_EVENT_ACTION": "", "VELA_BUILD_HOST": "", "VELA_BUILD_ROUTE": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "refs/tags/1", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SENDER_SCM_ID": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TAG": "1", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_EVENTS": "", "VELA_REPO_APPROVAL_TIMEOUT": "0", "VELA_REPO_APPROVE_BUILD": "", "VELA_REPO_OWNER": "foo", "VELA_REPO_BRANCH": "foo", "VELA_REPO_BUILD_LIMIT": "1", "VELA_REPO_CLONE": "foo", "VELA_REPO_CUSTOM_PROPS": `{"foo":"bar"}`, "VELA_REPO_INSTALL_ID": "0", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TOPICS": "cloud,security", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_ID_TOKEN_REQUEST_URL": "foo/api/v1/repos/foo/builds/1/id_token"},
		},
		{"pull_request", fields{
			build:    &api.Build{ID: &num64, Repo: &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL}, Number: &num64, Parent: &num64, Event: &pull, EventAction: &pullact, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, SenderSCMID: &str, Fork: &booL, Author: &str, Branch: &str, Ref: &pullref, BaseRef: &str},
			metadata: &internal.Metadata{Database: &internal.Database{Driver: str, Host: str}, Queue: &internal.Queue{Driver: str, Host: str}, Source: &internal.Source{Driver: str, Host: str}, Vela: &internal.Vela{Address: str, WebAddress: str, OpenIDIssuer: str}},
			repo:     &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL},
			user:     &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			netrc:    &netrc,
		}, map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "pull_request", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_PULL_REQUEST_NUMBER": "1", "BUILD_REF": "refs/pull/1/head", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "true", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_EVENTS": "", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_SERVER_ADDR": "foo", "VELA_OPEN_ID_ISSUER": "foo", "VELA_BUILD_APPROVED_AT": "0", "VELA_BUILD_APPROVED_BY": "", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "pull_request", "VELA_BUILD_EVENT_ACTION": "opened", "VELA_BUILD_HOST": "", "VELA_BUILD_ROUTE": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_PULL_REQUEST": "1", "VELA_PULL_REQUEST_FORK": "false", "VELA_BUILD_REF": "refs/pull/1/head", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SENDER_SCM_ID": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_PULL_REQUEST": "1", "VELA_PULL_REQUEST_SOURCE": "", "VELA_PULL_REQUEST_TARGET": "foo", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_EVENTS": "", "VELA_REPO_APPROVAL_TIMEOUT": "0", "VELA_REPO_APPROVE_BUILD": "", "VELA_REPO_OWNER": "foo", "VELA_REPO_BRANCH": "foo", "VELA_REPO_BUILD_LIMIT": "1", "VELA_REPO_CLONE": "foo", "VELA_REPO_CUSTOM_PROPS": `{"foo":"bar"}`, "VELA_REPO_INSTALL_ID": "0", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TOPICS": "cloud,security", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_ID_TOKEN_REQUEST_URL": "foo/api/v1/repos/foo/builds/1/id_token"},
		},
		{"deployment", fields{
			build:    &api.Build{ID: &num64, Repo: &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, CustomProps: &props, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL}, Number: &num64, Parent: &num64, Event: &deploy, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &target, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, SenderSCMID: &str, Author: &str, Branch: &str, Ref: &pullref, BaseRef: &str},
			metadata: &internal.Metadata{Database: &internal.Database{Driver: str, Host: str}, Queue: &internal.Queue{Driver: str, Host: str}, Source: &internal.Source{Driver: str, Host: str}, Vela: &internal.Vela{Address: str, WebAddress: str, OpenIDIssuer: str}},
			repo:     &api.Repo{ID: &num64, Owner: &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL}, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Topics: &topics, BuildLimit: &num32, Timeout: &num32, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL},
			user:     &api.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			netrc:    &netrc,
		}, map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BASE_REF": "foo", "BUILD_BRANCH": "foo", "BUILD_CLONE": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "deployment", "BUILD_HOST": "", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "refs/pull/1/head", "BUILD_SENDER": "foo", "BUILD_SOURCE": "foo", "BUILD_STARTED": "1", "BUILD_STATUS": "foo", "BUILD_TARGET": "production", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "CI": "true", "REPOSITORY_ACTIVE": "false", "REPOSITORY_ALLOW_EVENTS": "", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false", "REPOSITORY_VISIBILITY": "foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_SERVER_ADDR": "foo", "VELA_OPEN_ID_ISSUER": "foo", "VELA_BUILD_APPROVED_AT": "0", "VELA_BUILD_APPROVED_BY": "", "VELA_BUILD_AUTHOR": "foo", "VELA_BUILD_AUTHOR_EMAIL": "", "VELA_BUILD_BASE_REF": "foo", "VELA_BUILD_BRANCH": "foo", "VELA_BUILD_CLONE": "foo", "VELA_BUILD_COMMIT": "foo", "VELA_BUILD_CREATED": "1", "VELA_BUILD_DISTRIBUTION": "", "VELA_BUILD_ENQUEUED": "1", "VELA_BUILD_EVENT": "deployment", "VELA_BUILD_EVENT_ACTION": "", "VELA_BUILD_HOST": "", "VELA_BUILD_ROUTE": "", "VELA_BUILD_LINK": "", "VELA_BUILD_MESSAGE": "foo", "VELA_BUILD_NUMBER": "1", "VELA_BUILD_PARENT": "1", "VELA_BUILD_REF": "refs/pull/1/head", "VELA_BUILD_RUNTIME": "", "VELA_BUILD_SENDER": "foo", "VELA_BUILD_SENDER_SCM_ID": "foo", "VELA_BUILD_SOURCE": "foo", "VELA_BUILD_STARTED": "1", "VELA_BUILD_STATUS": "foo", "VELA_BUILD_TARGET": "production", "VELA_BUILD_TITLE": "foo", "VELA_BUILD_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_DATABASE": "foo", "VELA_DEPLOYMENT": "production", "VELA_DEPLOYMENT_NUMBER": "0", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_REPO_ACTIVE": "false", "VELA_REPO_ALLOW_EVENTS": "", "VELA_REPO_APPROVAL_TIMEOUT": "0", "VELA_REPO_APPROVE_BUILD": "", "VELA_REPO_OWNER": "foo", "VELA_REPO_BRANCH": "foo", "VELA_REPO_BUILD_LIMIT": "1", "VELA_REPO_CLONE": "foo", "VELA_REPO_CUSTOM_PROPS": "{}", "VELA_REPO_INSTALL_ID": "0", "VELA_REPO_FULL_NAME": "foo", "VELA_REPO_LINK": "foo", "VELA_REPO_NAME": "foo", "VELA_REPO_ORG": "foo", "VELA_REPO_PIPELINE_TYPE": "", "VELA_REPO_PRIVATE": "false", "VELA_REPO_TIMEOUT": "1", "VELA_REPO_TOPICS": "cloud,security", "VELA_REPO_TRUSTED": "false", "VELA_REPO_VISIBILITY": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_USER_ACTIVE": "false", "VELA_USER_ADMIN": "false", "VELA_USER_FAVORITES": "[]", "VELA_USER_NAME": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/vela/src/foo/foo/foo", "VELA_ID_TOKEN_REQUEST_URL": "foo/api/v1/repos/foo/builds/1/id_token"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				build:    tt.fields.build,
				metadata: tt.fields.metadata,
				repo:     tt.fields.repo,
				user:     tt.fields.user,
				netrc:    tt.fields.netrc,
			}
			got := c.EnvironmentBuild()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("EnvironmentBuild mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
