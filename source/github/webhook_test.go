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
	"time"

	"github.com/go-vela/types/constants"
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
	request.Header.Set("X-GitHub-Event", "push")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("push")
	wantHook.SetBranch("master")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(library.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("master")
	wantRepo.SetPrivate(false)

	wantBuild := new(library.Build)
	wantBuild.SetEvent("push")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/commit/9c93babf58917cd6f6f6772b5df2b098f507ff95")
	wantBuild.SetTitle("push received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("Update README.md")
	wantBuild.SetCommit("9c93babf58917cd6f6f6772b5df2b098f507ff95")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetBranch("master")
	wantBuild.SetRef("refs/heads/master")
	wantBuild.SetBaseRef("")

	gotHook, gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(gotHook, wantHook) {
		t.Errorf("ProcessWebhook webhook is %v, want %v", gotHook, wantHook)
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
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("push")
	wantHook.SetBranch("master")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(library.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("master")
	wantRepo.SetPrivate(false)

	wantBuild := new(library.Build)
	wantBuild.SetEvent("push")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/commit/9c93babf58917cd6f6f6772b5df2b098f507ff95")
	wantBuild.SetTitle("push received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("Update README.md")
	wantBuild.SetCommit("9c93babf58917cd6f6f6772b5df2b098f507ff95")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetBranch("master")
	wantBuild.SetRef("refs/heads/master")
	wantBuild.SetBaseRef("")

	gotHook, gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(gotHook, wantHook) {
		t.Errorf("ProcessWebhook webhook is %v, want %v", gotHook, wantHook)
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
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("pull_request")
	wantHook.SetBranch("master")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(library.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("master")
	wantRepo.SetPrivate(false)

	wantBuild := new(library.Build)
	wantBuild.SetEvent("pull_request")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/pull/1")
	wantBuild.SetTitle("pull_request received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("Update the README with new information")
	wantBuild.SetCommit("34c5c7793cb3b279e22454cb6750c80560547b3a")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetBranch("master")
	wantBuild.SetRef("refs/pull/1/head")
	wantBuild.SetBaseRef("master")

	gotHook, gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(gotHook, wantHook) {
		t.Errorf("ProcessWebhook webhook is %v, want %v", gotHook, wantHook)
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
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("pull_request")
	wantHook.SetBranch("master")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	gotHook, gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(gotHook, wantHook) {
		t.Errorf("ProcessWebhook webhook is %v, want %v", gotHook, wantHook)
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
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("pull_request")
	wantHook.SetBranch("master")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	gotHook, gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(gotHook, wantHook) {
		t.Errorf("ProcessWebhook webhook is %v, want %v", gotHook, wantHook)
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
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("pull_request")
	wantHook.SetStatus(constants.StatusSuccess)

	gotHook, gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err == nil {
		t.Errorf("ProcessWebhook should have returned err")
	}

	if !reflect.DeepEqual(gotHook, wantHook) {
		t.Errorf("ProcessWebhook webhook is %v, want %v", gotHook, wantHook)
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
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("foobar")
	wantHook.SetStatus(constants.StatusSuccess)

	gotHook, gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err == nil {
		t.Errorf("ProcessWebhook should have returned err")
	}

	if !reflect.DeepEqual(gotHook, wantHook) {
		t.Errorf("ProcessWebhook webhook is %v, want %v", gotHook, wantHook)
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
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("deployment")
	wantHook.SetStatus(constants.StatusSuccess)

	gotHook, gotRepo, gotBuild, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(gotHook, wantHook) {
		t.Errorf("ProcessWebhook webhook is %v, want %v", gotHook, wantHook)
	}

	if gotRepo != nil {
		t.Errorf("ProcessWebhook repo is %v, want nil", gotRepo)
	}

	if gotBuild != nil {
		t.Errorf("ProcessWebhook build is %v, want nil", gotBuild)
	}
}
