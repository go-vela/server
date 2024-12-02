// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"flag"
	"reflect"
	"testing"

	"github.com/urfave/cli/v2"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler/registry/github"
	"github.com/go-vela/server/internal"
)

func TestNative_New(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)
	public, _ := github.New(context.Background(), "", "")
	want := &client{
		Github:        public,
		Compiler:      settings.CompilerMockEmpty(),
		TemplateCache: make(map[string][]byte),
	}
	want.SetCloneImage(defaultCloneImage)

	// run test
	got, err := FromCLIContext(c)

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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)
	public, _ := github.New(context.Background(), "", "")
	private, _ := github.New(context.Background(), url, token)
	want := &client{
		Github:           public,
		PrivateGithub:    private,
		UsePrivateGithub: true,
		TemplateCache:    make(map[string][]byte),
		Compiler:         settings.CompilerMockEmpty(),
	}
	want.SetCloneImage(defaultCloneImage)

	// run test
	got, err := FromCLIContext(c)

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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)
	public, _ := github.New(context.Background(), "", "")
	private, _ := github.New(context.Background(), url, token)
	want := &client{
		Github:           public,
		PrivateGithub:    private,
		UsePrivateGithub: true,
		TemplateCache:    make(map[string][]byte),
		Compiler:         settings.CompilerMockEmpty(),
	}
	want.SetCloneImage(defaultCloneImage)

	// run test
	got, err := FromCLIContext(c)

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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	b := &api.Build{ID: &id}

	want, _ := FromCLIContext(c)

	// run test
	got, err := FromCLIContext(c)
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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	b := &api.Build{ID: &id}

	want, _ := FromCLIContext(c)
	want.build = b

	// run test
	got, err := FromCLIContext(c)
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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	f := []string{"foo"}

	want, _ := FromCLIContext(c)
	want.files = f

	// run test
	got, err := FromCLIContext(c)
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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	comment := "ok to test"
	want, _ := FromCLIContext(c)
	want.comment = comment

	// run test
	got, err := FromCLIContext(c)
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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	local := true
	want, _ := FromCLIContext(c)
	want.local = true

	// run test
	got, err := FromCLIContext(c)
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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	localTemplates := []string{"example:tmpl.yml", "exmpl:template.yml"}
	want, _ := FromCLIContext(c)
	want.localTemplates = []string{"example:tmpl.yml", "exmpl:template.yml"}

	// run test
	got, err := FromCLIContext(c)
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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	m := &internal.Metadata{
		Database: &internal.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &internal.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &internal.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &internal.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	want, _ := FromCLIContext(c)
	want.metadata = m

	// run test
	got, err := FromCLIContext(c)
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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	private, _ := github.New(context.Background(), url, token)

	want, _ := FromCLIContext(c)
	want.PrivateGithub = private

	// run test
	got, err := FromCLIContext(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithPrivateGitHub(context.Background(), url, token), want) {
		t.Errorf("WithRepo is %v, want %v", got, want)
	}
}

func TestNative_WithRepo(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	r := &api.Repo{ID: &id}

	want, _ := FromCLIContext(c)
	want.repo = r

	// run test
	got, err := FromCLIContext(c)
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
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	u := &api.User{ID: &id}

	want, _ := FromCLIContext(c)
	want.user = u

	// run test
	got, err := FromCLIContext(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithUser(u), want) {
		t.Errorf("WithUser is %v, want %v", got, want)
	}
}

func TestNative_WithLabels(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

	labels := []string{"documentation", "enhancement"}
	want, _ := FromCLIContext(c)
	want.labels = []string{"documentation", "enhancement"}

	// run test
	got, err := FromCLIContext(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithLabels(labels), want) {
		t.Errorf("WithLocalTemplates is %v, want %v", got, want)
	}
}
