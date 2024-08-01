// SPDX-License-Identifier: Apache-2.0

package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/go-vela/server/compiler/registry"
)

func TestGithub_Parse(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup types
	want := &registry.Source{
		Host: "github.example.com",
		Org:  "github",
		Repo: "octocat",
		Name: "template.yml",
		Ref:  "",
	}

	// run test
	c, err := New(s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	path := "github.example.com/github/octocat/template.yml"

	got, err := c.Parse(path)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestGithub_ParseWithBranch(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup types
	want := &registry.Source{
		Host: "github.example.com",
		Org:  "github",
		Repo: "octocat",
		Name: "template.yml",
		Ref:  "dev",
	}

	// run test
	c, err := New(s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	path := "github.example.com/github/octocat/template.yml@dev"

	got, err := c.Parse(path)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}
func TestGithub_Parse_Custom(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup types
	want := &registry.Source{
		Host: "github.example.com",
		Org:  "github",
		Repo: "octocat",
		Name: "path/to/template.yml",
		Ref:  "test",
	}

	// run test
	c, err := New(s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	path := "github.example.com/github/octocat/path/to/template.yml@test"

	got, err := c.Parse(path)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestGithub_Parse_Full(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup types
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Errorf("Parsing url returned err: %v", err)
	}

	want := &registry.Source{
		Host: u.Host,
		Org:  "github",
		Repo: "octocat",
		Name: "template.yml",
		Ref:  "test",
	}

	// run test
	c, err := New(s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	path := fmt.Sprintf("%s/%s", s.URL, "github/octocat/template.yml@test")

	got, err := c.Parse(path)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestGithub_Parse_Invalid(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// run test
	c, err := New(s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Parse("!@#$%^&*()")
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestGithub_Parse_Hostname(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup types
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Errorf("Parsing url returned err: %v", err)
	}

	want := &registry.Source{
		Host: u.Hostname(),
		Org:  "github",
		Repo: "octocat",
		Name: "template.yml",
		Ref:  "",
	}

	// run test
	c, err := New(s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	path := fmt.Sprintf("%s/%s", u.Hostname(), "github/octocat/template.yml")

	got, err := c.Parse(path)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestGithub_Parse_Path(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup types
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Errorf("Parsing url returned err: %v", err)
	}

	want := &registry.Source{
		Host: u.Host,
		Org:  "github",
		Repo: "octocat",
		Name: "path/to/template.yml",
		Ref:  "",
	}

	// run test
	c, err := New(s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	path := fmt.Sprintf("%s/%s", s.URL, "github/octocat/path/to/template.yml")

	got, err := c.Parse(path)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}
