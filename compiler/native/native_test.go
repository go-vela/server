// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"reflect"
	"testing"

	"github.com/go-vela/server/compiler/registry/github"

	"github.com/go-vela/types"
	"github.com/go-vela/types/library"

	"github.com/urfave/cli/v2"
)

func TestNative_New(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	public, _ := github.New("", "")
	want := &client{
		Github: public,
	}

	// run test
	got, err := New(c)

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestNative_New_PrivateGithub(t *testing.T) {
	// setup types
	url := "http://foo.example.com"
	token := "someToken"
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-url", url, "doc")
	set.String("github-token", token, "doc")
	c := cli.NewContext(nil, set, nil)
	public, _ := github.New("", "")
	private, _ := github.New(url, token)
	want := &client{
		Github:           public,
		PrivateGithub:    private,
		UsePrivateGithub: true,
	}

	// run test
	got, err := New(c)

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestNative_DuplicateRetainSettings(t *testing.T) {
	// setup types
	url := "http://foo.example.com"
	token := "someToken"
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-url", url, "doc")
	set.String("github-token", token, "doc")
	c := cli.NewContext(nil, set, nil)
	public, _ := github.New("", "")
	private, _ := github.New(url, token)
	want := &client{
		Github:           public,
		PrivateGithub:    private,
		UsePrivateGithub: true,
	}

	// run test
	got, err := New(c)

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got.Duplicate(), want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestNative_DuplicateStripBuild(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	b := &library.Build{ID: &id}

	want, _ := New(c)

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	// modify engine with WithBuild and then call Duplicate
	// to get a copy of the Engine without build attached.
	if !reflect.DeepEqual(got.WithBuild(b).Duplicate(), want) {
		t.Errorf("WithBuild is %v, want %v", got, want)
	}
}

func TestNative_WithBuild(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	b := &library.Build{ID: &id}

	want, _ := New(c)
	want.build = b

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithBuild(b), want) {
		t.Errorf("WithBuild is %v, want %v", got, want)
	}
}

func TestNative_WithFiles(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	f := []string{"foo"}

	want, _ := New(c)
	want.files = f

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithFiles(f), want) {
		t.Errorf("WithFiles is %v, want %v", got, want)
	}
}

func TestNative_WithComment(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	comment := "ok to test"
	want, _ := New(c)
	want.comment = comment

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithComment(comment), want) {
		t.Errorf("WithComment is %v, want %v", got, want)
	}
}

func TestNative_WithLocal(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	local := true
	want, _ := New(c)
	want.local = true

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithLocal(local), want) {
		t.Errorf("WithLocal is %v, want %v", got, want)
	}
}

func TestNative_WithLocalTemplates(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	localTemplates := []string{"example:tmpl.yml", "exmpl:template.yml"}
	want, _ := New(c)
	want.localTemplates = []string{"example:tmpl.yml", "exmpl:template.yml"}

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithLocalTemplates(localTemplates), want) {
		t.Errorf("WithLocalTemplates is %v, want %v", got, want)
	}
}

func TestNative_WithMetadata(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	want, _ := New(c)
	want.metadata = m

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithMetadata(m), want) {
		t.Errorf("WithMetadata is %v, want %v", got, want)
	}
}

func TestNative_WithPrivateGitHub(t *testing.T) {
	// setup types
	url := "http://foo.example.com"
	token := "someToken"
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-url", url, "doc")
	set.String("github-token", token, "doc")
	c := cli.NewContext(nil, set, nil)

	private, _ := github.New(url, token)

	want, _ := New(c)
	want.PrivateGithub = private

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithPrivateGitHub(url, token), want) {
		t.Errorf("WithRepo is %v, want %v", got, want)
	}
}

func TestNative_WithRepo(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	r := &library.Repo{ID: &id}

	want, _ := New(c)
	want.repo = r

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithRepo(r), want) {
		t.Errorf("WithRepo is %v, want %v", got, want)
	}
}

func TestNative_WithUser(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	u := &library.User{ID: &id}

	want, _ := New(c)
	want.user = u

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithUser(u), want) {
		t.Errorf("WithUser is %v, want %v", got, want)
	}
}
