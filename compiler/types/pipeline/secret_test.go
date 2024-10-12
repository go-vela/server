// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestPipeline_SecretSlice_Purge(t *testing.T) {
	// setup types
	secrets := testSecrets()
	*secrets = (*secrets)[:len(*secrets)-1]

	// setup tests
	tests := []struct {
		secrets *SecretSlice
		want    *SecretSlice
	}{
		{
			secrets: testSecrets(),
			want:    secrets,
		},
		{
			secrets: new(SecretSlice),
			want:    new(SecretSlice),
		},
	}

	// run tests
	for _, test := range tests {
		r := &RuleData{
			Branch: "main",
			Event:  "push",
			Path:   []string{},
			Repo:   "foo/bar",
			Tag:    "refs/heads/main",
		}

		got, _ := test.secrets.Purge(r)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Purge is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_Secret_ParseOrg_success(t *testing.T) {
	// setup tests
	tests := []struct {
		secret *Secret
		org    string
	}{
		{ // success with good data
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/foo",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org: "octocat",
		},
		{ // success with multilevel & special characters
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/👋/🧪/🔑",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org: "octocat",
		},
	}

	// run tests
	for _, test := range tests {
		org, key, err := test.secret.ParseOrg(test.org)
		if err != nil {
			t.Errorf("ParseOrg had an error occur: %+v", err)
		}

		p := strings.SplitN(test.secret.Key, "/", 2)

		if !strings.EqualFold(org, p[0]) {
			t.Errorf("org is %s want %s", org, p[0])
		}

		if !strings.EqualFold(key, p[1]) {
			t.Errorf("key is %s want %s", key, p[1])
		}
	}
}

func TestPipeline_Secret_ParseOrg_failure(t *testing.T) {
	// setup tests
	tests := []struct {
		secret  *Secret
		org     string
		wantErr error
	}{
		{ // failure with bad org
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/foo",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "wrongorg",
			wantErr: ErrInvalidOrg,
		},
		{ // failure with bad key
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidPath,
		},
		{ // failure with bad key
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidPath,
		},
		{ // failure with missing name
			secret: &Secret{
				Value:  "bar",
				Key:    "octocat/foo/bar",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidName,
		},
		{ // failure with bad name
			secret: &Secret{
				Name:   "This is a null char \u0000",
				Value:  "bar",
				Key:    "octocat/foo/bar",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidName,
		},
		{ // failure with bad engine
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/foo",
				Engine: "invalid",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidEngine,
		},
	}

	// run tests
	for _, test := range tests {
		_, _, err := test.secret.ParseOrg(test.org)
		if test.wantErr != nil && err != nil && !errors.Is(err, test.wantErr) {
			t.Errorf("ParseOrg should have failed with error '%s' but got '%s'", test.wantErr, err)
		}

		if err == nil {
			t.Errorf("ParseOrg should have failed")
		}
	}
}

func TestPipeline_Secret_ParseRepo_success(t *testing.T) {
	// setup tests
	tests := []struct {
		secret *Secret
		org    string
		repo   string
	}{
		{ // success with explicit
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/helloworld/foo",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org:  "octocat",
			repo: "helloworld",
		},
		{ // success with multilevel & special characters
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/👋/🧪/🔑",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org:  "octocat",
			repo: "👋",
		},
	}

	// run tests
	for _, test := range tests {
		org, repo, key, err := test.secret.ParseRepo(test.org, test.repo)
		if err != nil {
			t.Errorf("ParseRepo had an error occur: %+v", err)
		}

		// checks for explicit only
		if strings.Contains(test.secret.Key, "/") {
			p := strings.SplitN(test.secret.Key, "/", 3)

			if !strings.EqualFold(org, p[0]) {
				t.Errorf("org is %s want %s", org, p[0])
			}

			if !strings.EqualFold(repo, p[1]) {
				t.Errorf("repo is %s want %s", key, p[1])
			}

			if !strings.EqualFold(key, p[2]) {
				t.Errorf("key is %s want %s", key, p[2])
			}
		}
	}
}

func TestPipeline_Secret_ParseRepo_failure(t *testing.T) {
	// setup tests
	tests := []struct {
		secret  *Secret
		org     string
		repo    string
		wantErr error
	}{
		{ // failure with bad org
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/helloworld/foo",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org:     "wrongorg",
			repo:    "helloworld",
			wantErr: ErrInvalidOrg,
		},
		{ // failure with bad repo
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/helloworld/foo",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org:     "octocat",
			repo:    "badrepo",
			wantErr: ErrInvalidRepo,
		},
		{ // failure with bad key
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidPath,
		},
		{ // failure with bad key
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/helloworld",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			repo:    "helloworld",
			org:     "octocat",
			wantErr: ErrInvalidPath,
		},
		{ // failure with bad key
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/helloworld/",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			repo:    "helloworld",
			org:     "octocat",
			wantErr: ErrInvalidPath,
		},
		{ // failure with missing name
			secret: &Secret{
				Value:  "bar",
				Key:    "octocat/helloworld/foo/bar",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org:     "octocat",
			repo:    "helloworld",
			wantErr: ErrInvalidName,
		},
		{ // failure with bad name
			secret: &Secret{
				Name:   "SOME=PASSWORD",
				Value:  "bar",
				Key:    "octocat/helloworld/foo/bar",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org:     "octocat",
			repo:    "helloworld",
			wantErr: ErrInvalidName,
		},
		{ // failure with bad engine
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat",
				Engine: "invalid",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidEngine,
		},
		{ // failure with deprecated implicit syntax
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "foo",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org:     "octocat",
			repo:    "helloworld",
			wantErr: ErrInvalidPath,
		},
	}

	// run tests
	for _, test := range tests {
		_, _, _, err := test.secret.ParseRepo(test.org, test.repo)
		if test.wantErr != nil && err != nil && !errors.Is(err, test.wantErr) {
			t.Errorf("ParseRepo should have failed with error '%s' but got '%s'", test.wantErr, err)
		}

		if err == nil {
			t.Errorf("ParseRepo should have failed")
		}
	}
}

func TestPipeline_Secret_ParseShared_success(t *testing.T) {
	// setup tests
	tests := []struct {
		secret *Secret
		org    string
	}{
		{ // success with good data
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/helloworld/foo",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org: "octocat",
		},
		{ // success with multilevel & special characters
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/👋/🧪/🔑",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org: "octocat",
		},
	}

	// run tests
	for _, test := range tests {
		org, team, key, err := test.secret.ParseShared()
		if err != nil {
			t.Errorf("ParseShared had an error occur: %+v", err)
		}

		p := strings.SplitN(test.secret.Key, "/", 3)

		if !strings.EqualFold(org, p[0]) {
			t.Errorf("org is %s want %s", org, p[0])
		}

		if !strings.EqualFold(team, p[1]) {
			t.Errorf("repo is %s want %s", key, p[1])
		}

		if !strings.EqualFold(key, p[2]) {
			t.Errorf("key is %s want %s", key, p[2])
		}
	}
}

func TestPipeline_Secret_ParseShared_failure(t *testing.T) {
	// setup tests
	tests := []struct {
		secret  *Secret
		org     string
		wantErr error
	}{
		{ // failure with bad key
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat",
				Engine: "native",
				Type:   "repo",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidPath,
		},
		{ // failure with bad engine
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat",
				Engine: "invalid",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidEngine,
		},
		{ // failure with bad path
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/foo",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidPath,
		},
		{ // failure with bad path
			secret: &Secret{
				Name:   "foo",
				Value:  "bar",
				Key:    "octocat/foo/",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidPath,
		},
		{ // failure with missing name
			secret: &Secret{
				Value:  "bar",
				Key:    "octocat/foo/bar",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidName,
		},
		{ // failure with bad name
			secret: &Secret{
				Name:   "=",
				Value:  "bar",
				Key:    "octocat/foo/bar",
				Engine: "native",
				Type:   "org",
				Pull:   "build_start",
			},
			org:     "octocat",
			wantErr: ErrInvalidName,
		},
	}

	// run tests
	for _, test := range tests {
		_, _, _, err := test.secret.ParseShared()
		if test.wantErr != nil && err != nil && !errors.Is(err, test.wantErr) {
			t.Errorf("ParseShared should have failed with error '%s' but got '%s'", test.wantErr, err)
		}

		if err == nil {
			t.Errorf("ParseShared should have failed")
		}
	}
}

func testSecrets() *SecretSlice {
	return &SecretSlice{
		{
			Engine: "native",
			Key:    "github/octocat/foobar",
			Name:   "foobar",
			Type:   "repo",
			Origin: &Container{},
			Pull:   "build_start",
		},
		{
			Engine: "native",
			Key:    "github/foobar",
			Name:   "foobar",
			Type:   "org",
			Origin: &Container{},
			Pull:   "build_start",
		},
		{
			Engine: "native",
			Key:    "github/octokitties/foobar",
			Name:   "foobar",
			Type:   "shared",
			Origin: &Container{},
			Pull:   "build_start",
		},
		{
			Name: "",
			Origin: &Container{
				ID:          "secret_github octocat._1_vault",
				Directory:   "/vela/src/foo//",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "vault:latest",
				Name:        "vault",
				Number:      1,
				Pull:        "always",
				Ruleset: Ruleset{
					If:       Rules{Event: []string{"push"}},
					Operator: "and",
				},
			},
		},
		{
			Origin: &Container{
				ID:          "secret_github octocat._2_vault",
				Directory:   "/vela/src/foo//",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "vault:latest",
				Name:        "vault",
				Number:      2,
				Pull:        "always",
				Ruleset: Ruleset{
					If:       Rules{Event: []string{"pull_request"}},
					Operator: "and",
				},
			},
		},
	}
}
