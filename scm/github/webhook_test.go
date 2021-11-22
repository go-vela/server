// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-vela/types/raw"
	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

func TestGithub_ProcessWebhook_Push(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/push.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

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
	wantBuild.SetEmail("21031067+Codertocat@users.noreply.github.com")
	wantBuild.SetBranch("master")
	wantBuild.SetRef("refs/heads/master")
	wantBuild.SetBaseRef("")

	want := &types.Webhook{
		Comment: "",
		Hook:    wantHook,
		Repo:    wantRepo,
		Build:   wantBuild,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_ProcessWebhook_Push_NoSender(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/push_no_sender.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

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
	wantBuild.SetEmail("21031067+Codertocat@users.noreply.github.com")
	wantBuild.SetBranch("master")
	wantBuild.SetRef("refs/heads/master")
	wantBuild.SetBaseRef("")

	want := &types.Webhook{
		Comment: "",
		Hook:    wantHook,
		Repo:    wantRepo,
		Build:   wantBuild,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_ProcessWebhook_PullRequest(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/pull_request.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

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
	wantBuild.SetEmail("")
	wantBuild.SetBranch("master")
	wantBuild.SetRef("refs/pull/1/head")
	wantBuild.SetBaseRef("master")
	wantBuild.SetHeadRef("changes")

	want := &types.Webhook{
		Comment:  "",
		PRNumber: wantHook.GetNumber(),
		Hook:     wantHook,
		Repo:     wantRepo,
		Build:    wantBuild,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_ProcessWebhook_PullRequest_ClosedAction(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/pull_request_closed_action.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

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

	want := &types.Webhook{
		Comment: "",
		Hook:    wantHook,
		Repo:    nil,
		Build:   nil,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_ProcessWebhook_PullRequest_ClosedState(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/pull_request_closed_state.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

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

	want := &types.Webhook{
		Comment: "",
		Hook:    wantHook,
		Repo:    nil,
		Build:   nil,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_ProcessWebhook_Deployment(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetBranch("master")
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")
	wantHook.SetHost("github.com")
	wantHook.SetEvent("deployment")
	wantHook.SetStatus(constants.StatusSuccess)

	wantRepo := new(library.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("master")
	wantRepo.SetPrivate(false)

	wantBuild := new(library.Build)
	wantBuild.SetEvent("deployment")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetDeploy("production")
	wantBuild.SetSource("https://api.github.com/repos/Codertocat/Hello-World/deployments/145988746")
	wantBuild.SetTitle("deployment received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("")
	wantBuild.SetCommit("f95f852bd8fca8fcc58a9a2d6c842781e32a215e")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("")
	wantBuild.SetBranch("master")
	wantBuild.SetRef("refs/heads/master")

	type args struct {
		file              string
		hook              *library.Hook
		repo              *library.Repo
		build             *library.Build
		deploymentPayload raw.StringSliceMap
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{file: "deployment.json", hook: wantHook, repo: wantRepo, build: wantBuild, deploymentPayload: raw.StringSliceMap{"foo": "test1", "bar": "test2"}}, false},
		{"unexpected json payload", args{file: "deployment_unexpected_json_payload.json", deploymentPayload: raw.StringSliceMap{}}, true},
		{"unexpected text payload", args{file: "deployment_unexpected_text_payload.json", deploymentPayload: raw.StringSliceMap{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := os.Open(fmt.Sprintf("testdata/hooks/%s", tt.args.file))
			if err != nil {
				t.Errorf("unable to open file: %v", err)
			}

			defer body.Close()

			request, _ := http.NewRequest(http.MethodGet, "/test", body)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
			request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
			request.Header.Set("X-GitHub-Host", "github.com")
			request.Header.Set("X-GitHub-Version", "2.16.0")
			request.Header.Set("X-GitHub-Event", "deployment")

			client, _ := NewTest(s.URL)
			wantBuild.SetDeployPayload(tt.args.deploymentPayload)

			want := &types.Webhook{
				Comment: "",
				Hook:    tt.args.hook,
				Repo:    tt.args.repo,
				Build:   tt.args.build,
			}

			got, err := client.ProcessWebhook(request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("ProcessWebhook() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGithub_ProcessWebhook_Deployment_Commit(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/deployment_commit.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

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
	wantHook.SetBranch("master")
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")
	wantHook.SetHost("github.com")
	wantHook.SetEvent("deployment")
	wantHook.SetStatus(constants.StatusSuccess)

	wantRepo := new(library.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("master")
	wantRepo.SetPrivate(false)

	wantBuild := new(library.Build)
	wantBuild.SetEvent("deployment")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetDeploy("production")
	wantBuild.SetSource("https://api.github.com/repos/Codertocat/Hello-World/deployments/145988746")
	wantBuild.SetTitle("deployment received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("")
	wantBuild.SetCommit("f95f852bd8fca8fcc58a9a2d6c842781e32a215e")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("")
	wantBuild.SetBranch("master")
	wantBuild.SetRef("refs/heads/master")

	want := &types.Webhook{
		Comment: "",
		Hook:    wantHook,
		Repo:    wantRepo,
		Build:   wantBuild,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_ProcessWebhook_BadGithubEvent(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/pull_request.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

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

	want := &types.Webhook{
		Comment: "",
		Hook:    wantHook,
		Repo:    nil,
		Build:   nil,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_ProcessWebhook_BadContentType(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/pull_request.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

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

	want := &types.Webhook{
		Comment: "",
		Hook:    wantHook,
		Repo:    nil,
		Build:   nil,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_VerifyWebhook_EmptyRepo(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/push.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

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
	err = client.VerifyWebhook(request, new(library.Repo))
	if err != nil {
		t.Errorf("VerifyWebhook should have returned err")
	}
}

func TestGithub_VerifyWebhook_NoSecret(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	r := new(library.Repo)
	r.SetOrg("Codertocat")
	r.SetName("Hello-World")
	r.SetFullName("Codertocat/Hello-World")
	r.SetLink("https://github.com/Codertocat/Hello-World")
	r.SetClone("https://github.com/Codertocat/Hello-World.git")
	r.SetBranch("master")
	r.SetPrivate(false)

	// setup request
	body, err := os.Open("testdata/hooks/push.json")
	if err != nil {
		t.Errorf("Opening file for ProcessWebhook returned err: %v", err)
	}

	defer body.Close()

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
	err = client.VerifyWebhook(request, r)
	if err != nil {
		t.Errorf("VerifyWebhook should have returned err")
	}
}

func TestGithub_ProcessWebhook_IssueComment_PR(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/issue_comment_pr.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "issue_comment")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("comment")
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
	wantBuild.SetEvent("comment")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/pull/1")
	wantBuild.SetTitle("comment received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("Update the README with new information")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("")
	wantBuild.SetRef("refs/pull/1/head")

	want := &types.Webhook{
		Comment:  "ok to test",
		PRNumber: wantHook.GetNumber(),
		Hook:     wantHook,
		Repo:     wantRepo,
		Build:    wantBuild,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_ProcessWebhook_IssueComment_Created(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/issue_comment_created.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "issue_comment")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("comment")
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
	wantBuild.SetEvent("comment")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/issues/1")
	wantBuild.SetTitle("comment received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("Update the README with new information")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("")
	wantBuild.SetRef("refs/heads/master")

	want := &types.Webhook{
		Comment:  "ok to test",
		PRNumber: 0,
		Hook:     wantHook,
		Repo:     wantRepo,
		Build:    wantBuild,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_ProcessWebhook_IssueComment_Deleted(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/issue_comment_deleted.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequest(http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "issue_comment")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(library.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("comment")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	want := &types.Webhook{
		Comment:  "ok to test",
		PRNumber: 0,
		Hook:     wantHook,
		Repo:     nil,
		Build:    nil,
	}

	got, err := client.ProcessWebhook(request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}
