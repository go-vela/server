// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"reflect"
	"testing"

	"github.com/urfave/cli/v3"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler/registry/github"
	"github.com/go-vela/server/internal"
)

func TestNative_New(t *testing.T) {
	// setup types
	want := &Client{
		Compiler:      settings.CompilerMockEmpty(),
		TemplateCache: make(map[string][]byte),
	}
	want.SetCloneImage(defaultCloneImage)
	want.SetTemplateDepth(5)

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if got.Github == nil {
		t.Errorf("New returned nil Github client")
	}

	got.Github = nil

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestNative_New_PrivateGithub(t *testing.T) {
	// setup types
	want := &Client{
		UsePrivateGithub: true,
		TemplateCache:    make(map[string][]byte),
		Compiler:         settings.CompilerMockEmpty(),
	}
	want.SetCloneImage(defaultCloneImage)
	want.SetTemplateDepth(5)

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if got.Github == nil {
		t.Errorf("New returned nil Github client")
	}

	if got.PrivateGithub == nil {
		t.Errorf("New returned nil Private Github client")
	}

	got.Github = nil
	got.PrivateGithub = nil

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestNative_DuplicateRetainSettings(t *testing.T) {
	// setup types
	want := &Client{
		UsePrivateGithub: true,
		TemplateCache:    make(map[string][]byte),
		Compiler:         settings.CompilerMockEmpty(),
	}
	want.SetCloneImage(defaultCloneImage)
	want.SetTemplateDepth(5)

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if got.Github == nil {
		t.Errorf("New returned nil Github client")
	}

	if got.PrivateGithub == nil {
		t.Errorf("New returned nil Private Github client")
	}

	got.Github = nil
	got.PrivateGithub = nil

	if !reflect.DeepEqual(got.Duplicate(), want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestNative_DuplicateStripBuild(t *testing.T) {
	// setup types
	id := int64(1)
	b := &api.Build{ID: &id}

	want, _ := FromCLICommand(context.Background(), testCommand(t, ""))

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.Github = want.Github

	// modify engine with WithBuild and then call Duplicate
	// to get a copy of the Engine without build attached.
	if !reflect.DeepEqual(got.WithBuild(b).Duplicate(), want) {
		t.Errorf("WithBuild is %v, want %v", got, want)
	}
}

func TestNative_WithBuild(t *testing.T) {
	// setup types
	id := int64(1)
	b := &api.Build{ID: &id}

	want, _ := FromCLICommand(context.Background(), testCommand(t, ""))
	want.build = b

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.Github = want.Github

	if !reflect.DeepEqual(got.WithBuild(b), want) {
		t.Errorf("WithBuild is %v, want %v", got, want)
	}
}

func TestNative_WithFiles(t *testing.T) {
	// setup types
	f := []string{"foo"}

	want, _ := FromCLICommand(context.Background(), testCommand(t, ""))
	want.files = f

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.Github = want.Github

	if !reflect.DeepEqual(got.WithFiles(f), want) {
		t.Errorf("WithFiles is %v, want %v", got, want)
	}
}

func TestNative_WithComment(t *testing.T) {
	// setup types
	comment := "ok to test"
	want, _ := FromCLICommand(context.Background(), testCommand(t, ""))
	want.comment = comment

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.Github = want.Github

	if !reflect.DeepEqual(got.WithComment(comment), want) {
		t.Errorf("WithComment is %v, want %v", got, want)
	}
}

func TestNative_WithLocal(t *testing.T) {
	// setup types
	local := true
	want, _ := FromCLICommand(context.Background(), testCommand(t, ""))
	want.local = true

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.Github = want.Github

	if !reflect.DeepEqual(got.WithLocal(local), want) {
		t.Errorf("WithLocal is %v, want %v", got, want)
	}
}

func TestNative_WithLocalTemplates(t *testing.T) {
	// setup types
	localTemplates := []string{"example:tmpl.yml", "exmpl:template.yml"}
	want, _ := FromCLICommand(context.Background(), testCommand(t, ""))
	want.localTemplates = []string{"example:tmpl.yml", "exmpl:template.yml"}

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.Github = want.Github

	if !reflect.DeepEqual(got.WithLocalTemplates(localTemplates), want) {
		t.Errorf("WithLocalTemplates is %v, want %v", got, want)
	}
}

func TestNative_WithMetadata(t *testing.T) {
	// setup types
	m := &internal.Metadata{
		Database: &internal.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &internal.Queue{
			Driver: "foo",
			Host:   "foo",
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

	want, _ := FromCLICommand(context.Background(), testCommand(t, ""))
	want.metadata = m

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.Github = want.Github

	if !reflect.DeepEqual(got.WithMetadata(m), want) {
		t.Errorf("WithMetadata is %v, want %v", got, want)
	}
}

func TestNative_WithPrivateGitHub(t *testing.T) {
	// setup types
	private, _ := github.New(context.Background(), "http://foo.example.com", "someToken")

	want, _ := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	want.PrivateGithub = private

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.WithPrivateGitHub(context.Background(), "http://foo.example.com", "someToken")

	got.Github = want.Github
	got.PrivateGithub = want.PrivateGithub

	if !reflect.DeepEqual(got, want) {
		t.Errorf("WithPrivateGitHub is %v, want %v", got, want)
	}
}

func TestNative_WithRepo(t *testing.T) {
	// setup types
	id := int64(1)
	r := &api.Repo{ID: &id}

	want, _ := FromCLICommand(context.Background(), testCommand(t, ""))
	want.repo = r

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.Github = want.Github

	if !reflect.DeepEqual(got.WithRepo(r), want) {
		t.Errorf("WithRepo is %v, want %v", got, want)
	}
}

func TestNative_WithUser(t *testing.T) {
	// setup types
	id := int64(1)
	u := &api.User{ID: &id}

	want, _ := FromCLICommand(context.Background(), testCommand(t, ""))
	want.user = u

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.Github = want.Github

	if !reflect.DeepEqual(got.WithUser(u), want) {
		t.Errorf("WithUser is %v, want %v", got, want)
	}
}

func TestNative_WithLabels(t *testing.T) {
	// setup types
	labels := []string{"documentation", "enhancement"}
	want, _ := FromCLICommand(context.Background(), testCommand(t, ""))
	want.labels = []string{"documentation", "enhancement"}

	// run test
	got, err := FromCLICommand(context.Background(), testCommand(t, ""))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got.Github = want.Github

	if !reflect.DeepEqual(got.WithLabels(labels), want) {
		t.Errorf("WithLocalTemplates is %v, want %v", got, want)
	}
}

func testCommand(t *testing.T, url string) *cli.Command {
	t.Helper()

	c := new(cli.Command)

	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: defaultCloneImage,
		},
		&cli.IntFlag{
			Name:  "max-template-depth",
			Value: 5,
		},
		&cli.BoolFlag{
			Name:  "github-driver",
			Value: len(url) > 0,
		},
		&cli.StringFlag{
			Name:  "github-url",
			Value: url,
		},
		&cli.StringFlag{
			Name:  "github-token",
			Value: "someToken",
		},
	}

	return c
}
