// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestGithub_ProcessWebhook_Push(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/push.json")
	defer body.Close()
	if err != nil {
		t.Errorf("Opening file for ProcessWebhook returned err: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "push")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	rOrg := "Codertocat"
	rName := "Hello-World"
	rFullName := "Codertocat/Hello-World"
	rLink := "https://github.com/Codertocat/Hello-World"
	rClone := "https://github.com/Codertocat/Hello-World.git"
	rBranch := "master"
	zeroBool := false
	wantRepo := &library.Repo{
		Org:      &rOrg,
		Name:     &rName,
		FullName: &rFullName,
		Link:     &rLink,
		Clone:    &rClone,
		Branch:   &rBranch,
		Private:  &zeroBool,
	}

	bEvent := "push"
	bClone := "https://github.com/Codertocat/Hello-World.git"
	bSource := "https://github.com/Codertocat/Hello-World/commit/9c93babf58917cd6f6f6772b5df2b098f507ff95"
	bTitle := "push received from https://github.com/Codertocat/Hello-World"
	bMessage := "Update README.md"
	bCommit := "9c93babf58917cd6f6f6772b5df2b098f507ff95"
	bSender := "Codertocat"
	bAuthor := "Codertocat"
	bBranch := "master"
	bRef := "refs/heads/master"
	bBaseRef := ""
	wantBuild := &library.Build{
		Event:   &bEvent,
		Clone:   &bClone,
		Source:  &bSource,
		Title:   &bTitle,
		Message: &bMessage,
		Commit:  &bCommit,
		Sender:  &bSender,
		Author:  &bAuthor,
		Branch:  &bBranch,
		Ref:     &bRef,
		BaseRef: &bBaseRef,
	}
	gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(gotRepo, wantRepo) {
		t.Errorf("ProcessWebhook repo is %v, want %v", gotRepo, wantRepo)
	}

	if !reflect.DeepEqual(gotBuild, wantBuild) {
		t.Errorf("ProcessWebhook build is %v, want %v", gotBuild, wantBuild)
	}
}

func TestGithub_ProcessWebhook_Push_NoSender(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/pushNoSender.json")
	defer body.Close()
	if err != nil {
		t.Errorf("Opening file for ProcessWebhook returned err: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "push")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	rOrg := "Codertocat"
	rName := "Hello-World"
	rFullName := "Codertocat/Hello-World"
	rLink := "https://github.com/Codertocat/Hello-World"
	rClone := "https://github.com/Codertocat/Hello-World.git"
	rBranch := "master"
	zeroBool := false
	wantRepo := &library.Repo{
		Org:      &rOrg,
		Name:     &rName,
		FullName: &rFullName,
		Link:     &rLink,
		Clone:    &rClone,
		Branch:   &rBranch,
		Private:  &zeroBool,
	}

	bEvent := "push"
	bClone := "https://github.com/Codertocat/Hello-World.git"
	bSource := "https://github.com/Codertocat/Hello-World/commit/9c93babf58917cd6f6f6772b5df2b098f507ff95"
	bTitle := "push received from https://github.com/Codertocat/Hello-World"
	bMessage := "Update README.md"
	bCommit := "9c93babf58917cd6f6f6772b5df2b098f507ff95"
	bSender := "Codertocat"
	bAuthor := "Codertocat"
	bBranch := "master"
	bRef := "refs/heads/master"
	bBaseRef := ""
	wantBuild := &library.Build{
		Event:   &bEvent,
		Clone:   &bClone,
		Source:  &bSource,
		Title:   &bTitle,
		Message: &bMessage,
		Commit:  &bCommit,
		Sender:  &bSender,
		Author:  &bAuthor,
		Branch:  &bBranch,
		Ref:     &bRef,
		BaseRef: &bBaseRef,
	}
	gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(gotRepo, wantRepo) {
		t.Errorf("ProcessWebhook repo is %v, want %v", gotRepo, wantRepo)
	}

	if !reflect.DeepEqual(gotBuild, wantBuild) {
		t.Errorf("ProcessWebhook build is %v, want %v", gotBuild, wantBuild)
	}
}

func TestGithub_ProcessWebhook_PullRequest(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/pull.json")
	defer body.Close()
	if err != nil {
		t.Errorf("Opening file for ProcessWebhook returned err: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "pull_request")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	rOrg := "Codertocat"
	rName := "Hello-World"
	rFullName := "Codertocat/Hello-World"
	rLink := "https://github.com/Codertocat/Hello-World"
	rClone := "https://github.com/Codertocat/Hello-World.git"
	rBranch := "master"
	zeroBool := false
	wantRepo := &library.Repo{
		Org:      &rOrg,
		Name:     &rName,
		FullName: &rFullName,
		Link:     &rLink,
		Clone:    &rClone,
		Branch:   &rBranch,
		Private:  &zeroBool,
	}

	bEvent := "pull_request"
	bClone := "https://github.com/Codertocat/Hello-World.git"
	bSource := "https://github.com/Codertocat/Hello-World/pull/1"
	bTitle := "pull_request received from https://github.com/Codertocat/Hello-World"
	bMessage := "Update the README with new information"
	bCommit := "34c5c7793cb3b279e22454cb6750c80560547b3a"
	bSender := "Codertocat"
	bAuthor := "Codertocat"
	bBranch := "master"
	bRef := "refs/pull/1/head"
	bBaseRef := "master"
	wantBuild := &library.Build{
		Event:   &bEvent,
		Clone:   &bClone,
		Source:  &bSource,
		Title:   &bTitle,
		Message: &bMessage,
		Commit:  &bCommit,
		Sender:  &bSender,
		Author:  &bAuthor,
		Branch:  &bBranch,
		Ref:     &bRef,
		BaseRef: &bBaseRef,
	}
	gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(gotRepo, wantRepo) {
		t.Errorf("ProcessWebhook repo is %v, want %v", gotRepo, wantRepo)
	}

	if !reflect.DeepEqual(gotBuild, wantBuild) {
		t.Errorf("ProcessWebhook build is %v, want %v", gotBuild, wantBuild)
	}
}

func TestGithub_ProcessWebhook_PullRequest_ClosedAction(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/pullClosedAction.json")
	defer body.Close()
	if err != nil {
		t.Errorf("Opening file for ProcessWebhook returned err: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "pull_request")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if gotRepo != nil {
		t.Errorf("ProcessWebhook repo is %v, want nil", gotRepo)
	}

	if gotBuild != nil {
		t.Errorf("ProcessWebhook build is %v, want nil", gotBuild)
	}
}

func TestGithub_ProcessWebhook_PullRequest_ClosedState(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/pullClosedState.json")
	defer body.Close()
	if err != nil {
		t.Errorf("Opening file for ProcessWebhook returned err: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "pull_request")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if gotRepo != nil {
		t.Errorf("ProcessWebhook repo is %v, want nil", gotRepo)
	}

	if gotBuild != nil {
		t.Errorf("ProcessWebhook build is %v, want nil", gotBuild)
	}
}

func TestGithub_ProcessWebhook_BadContentType(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/pull.json")
	defer body.Close()
	if err != nil {
		t.Errorf("Opening file for ProcessWebhook returned err: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "foobar")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "pull_request")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err == nil {
		t.Errorf("ProcessWebhook should have returned err")
	}

	if gotRepo != nil {
		t.Errorf("ProcessWebhook repo is %v, want nil", gotRepo)
	}

	if gotBuild != nil {
		t.Errorf("ProcessWebhook build is %v, want nil", gotBuild)
	}
}

func TestGithub_ProcessWebhook_BadGithubEvent(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/pull.json")
	defer body.Close()
	if err != nil {
		t.Errorf("Opening file for ProcessWebhook returned err: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "foobar")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err == nil {
		t.Errorf("ProcessWebhook should have returned err")
	}

	if gotRepo != nil {
		t.Errorf("ProcessWebhook repo is %v, want nil", gotRepo)
	}

	if gotBuild != nil {
		t.Errorf("ProcessWebhook build is %v, want nil", gotBuild)
	}
}

func TestGithub_ProcessWebhook_UnsupportedGithubEvent(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/pull.json")
	defer body.Close()
	if err != nil {
		t.Errorf("Opening file for ProcessWebhook returned err: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "deployment")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if gotRepo != nil {
		t.Errorf("ProcessWebhook repo is %v, want nil", gotRepo)
	}

	if gotBuild != nil {
		t.Errorf("ProcessWebhook build is %v, want nil", gotBuild)
	}
}
