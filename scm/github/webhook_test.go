// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/types/constants"
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

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Event", "push")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("push")
	wantHook.SetBranch("main")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics([]string{"go", "vela"})

	wantBuild := new(api.Build)
	wantBuild.SetEvent("push")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/commit/9c93babf58917cd6f6f6772b5df2b098f507ff95")
	wantBuild.SetTitle("push received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("Update README.md")
	wantBuild.SetCommit("9c93babf58917cd6f6f6772b5df2b098f507ff95")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetSenderSCMID("21031067")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("21031067+Codertocat@users.noreply.github.com")
	wantBuild.SetBranch("main")
	wantBuild.SetRef("refs/heads/main")
	wantBuild.SetBaseRef("")

	want := &internal.Webhook{
		Hook:  wantHook,
		Repo:  wantRepo,
		Build: wantBuild,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

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

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "push")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("push")
	wantHook.SetBranch("main")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics([]string{"go", "vela"})

	wantBuild := new(api.Build)
	wantBuild.SetEvent("push")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/commit/9c93babf58917cd6f6f6772b5df2b098f507ff95")
	wantBuild.SetTitle("push received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("Update README.md")
	wantBuild.SetCommit("9c93babf58917cd6f6f6772b5df2b098f507ff95")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetSenderSCMID("0")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("21031067+Codertocat@users.noreply.github.com")
	wantBuild.SetBranch("main")
	wantBuild.SetRef("refs/heads/main")
	wantBuild.SetBaseRef("")

	want := &internal.Webhook{
		Hook:  wantHook,
		Repo:  wantRepo,
		Build: wantBuild,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ProcessWebhook() mismatch (-want +got):\n%s", diff)
	}
}

func TestGithub_ProcessWebhook_Push_Branch_Delete(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/push_delete_branch.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Event", "push")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("delete")
	wantHook.SetBranch("main")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics([]string{"go", "vela"})

	wantBuild := new(api.Build)
	wantBuild.SetEvent("delete")
	wantBuild.SetEventAction("branch")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/commit/d3d9188fc87a6977343e922c128f162a86018d76")
	wantBuild.SetTitle("delete received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("main branch deleted")
	wantBuild.SetCommit("d3d9188fc87a6977343e922c128f162a86018d76")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetSenderSCMID("21031067")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("21031067+Codertocat@users.noreply.github.com")
	wantBuild.SetBranch("main")
	wantBuild.SetRef("d3d9188fc87a6977343e922c128f162a86018d76")
	wantBuild.SetBaseRef("")

	want := &internal.Webhook{
		Hook:  wantHook,
		Repo:  wantRepo,
		Build: wantBuild,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_ProcessWebhook_Push_Tag_Delete(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/push_delete_tag.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Event", "push")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("delete")
	wantHook.SetBranch("refs/tags/v0.1")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics([]string{"go", "vela"})

	wantBuild := new(api.Build)
	wantBuild.SetEvent("delete")
	wantBuild.SetEventAction("tag")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/commit/d3d9188fc87a6977343e922c128f162a86018d76")
	wantBuild.SetTitle("delete received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("v0.1 tag deleted")
	wantBuild.SetCommit("d3d9188fc87a6977343e922c128f162a86018d76")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetSenderSCMID("21031067")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("21031067+Codertocat@users.noreply.github.com")
	wantBuild.SetBranch("v0.1")
	wantBuild.SetRef("d3d9188fc87a6977343e922c128f162a86018d76")
	wantBuild.SetBaseRef("")

	want := &internal.Webhook{
		Hook:  wantHook,
		Repo:  wantRepo,
		Build: wantBuild,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

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

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("pull_request")
	wantHook.SetBranch("main")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics(nil)

	wantBuild := new(api.Build)
	wantBuild.SetEvent("pull_request")
	wantBuild.SetEventAction("opened")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/pull/1")
	wantBuild.SetTitle("pull_request received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("Update the README with new information")
	wantBuild.SetCommit("34c5c7793cb3b279e22454cb6750c80560547b3a")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetSenderSCMID("21031067")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("")
	wantBuild.SetBranch("main")
	wantBuild.SetRef("refs/pull/1/head")
	wantBuild.SetBaseRef("main")
	wantBuild.SetHeadRef("changes")

	wantBuild2 := new(api.Build)
	wantBuild2.SetEvent("pull_request")
	wantBuild2.SetEventAction("labeled")
	wantBuild2.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild2.SetSource("https://github.com/Codertocat/Hello-World/pull/1")
	wantBuild2.SetTitle("pull_request received from https://github.com/Codertocat/Hello-World")
	wantBuild2.SetMessage("Update the README with new information")
	wantBuild2.SetCommit("34c5c7793cb3b279e22454cb6750c80560547b3a")
	wantBuild2.SetSender("Codertocat")
	wantBuild2.SetSenderSCMID("21031067")
	wantBuild2.SetAuthor("Codertocat")
	wantBuild2.SetEmail("")
	wantBuild2.SetBranch("main")
	wantBuild2.SetRef("refs/pull/1/head")
	wantBuild2.SetBaseRef("main")
	wantBuild2.SetHeadRef("changes")

	wantBuild3 := new(api.Build)
	wantBuild3.SetEvent("pull_request")
	wantBuild3.SetEventAction("unlabeled")
	wantBuild3.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild3.SetSource("https://github.com/Codertocat/Hello-World/pull/1")
	wantBuild3.SetTitle("pull_request received from https://github.com/Codertocat/Hello-World")
	wantBuild3.SetMessage("Update the README with new information")
	wantBuild3.SetCommit("34c5c7793cb3b279e22454cb6750c80560547b3a")
	wantBuild3.SetSender("Codertocat")
	wantBuild3.SetSenderSCMID("21031067")
	wantBuild3.SetAuthor("Codertocat")
	wantBuild3.SetEmail("")
	wantBuild3.SetBranch("main")
	wantBuild3.SetRef("refs/pull/1/head")
	wantBuild3.SetBaseRef("main")
	wantBuild3.SetHeadRef("changes")

	wantBuild4 := new(api.Build)
	wantBuild4.SetEvent("pull_request")
	wantBuild4.SetEventAction("edited")
	wantBuild4.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild4.SetSource("https://github.com/Codertocat/Hello-World/pull/1")
	wantBuild4.SetTitle("pull_request received from https://github.com/Codertocat/Hello-World")
	wantBuild4.SetMessage("Update the README with new information")
	wantBuild4.SetCommit("34c5c7793cb3b279e22454cb6750c80560547b3a")
	wantBuild4.SetSender("Codertocat")
	wantBuild4.SetSenderSCMID("21031067")
	wantBuild4.SetAuthor("Codertocat")
	wantBuild4.SetEmail("")
	wantBuild4.SetBranch("main")
	wantBuild4.SetRef("refs/pull/1/head")
	wantBuild4.SetBaseRef("main")
	wantBuild4.SetHeadRef("changes")

	tests := []struct {
		name     string
		testData string
		want     *internal.Webhook
		wantErr  bool
	}{
		{
			name:     "success",
			testData: "testdata/hooks/pull_request.json",
			want: &internal.Webhook{
				PullRequest: internal.PullRequest{
					Number:     wantHook.GetNumber(),
					IsFromFork: false,
				},
				Hook:  wantHook,
				Repo:  wantRepo,
				Build: wantBuild,
			},
		},
		{
			name:     "fork",
			testData: "testdata/hooks/pull_request_fork.json",
			want: &internal.Webhook{
				PullRequest: internal.PullRequest{
					Number:     wantHook.GetNumber(),
					IsFromFork: true,
				},
				Hook:  wantHook,
				Repo:  wantRepo,
				Build: wantBuild,
			},
		},
		{
			name:     "fork same repo",
			testData: "testdata/hooks/pull_request_fork_same-repo.json",
			want: &internal.Webhook{
				PullRequest: internal.PullRequest{
					Number:     wantHook.GetNumber(),
					IsFromFork: false,
				},
				Hook:  wantHook,
				Repo:  wantRepo,
				Build: wantBuild,
			},
		},
		{
			name:     "closed action",
			testData: "testdata/hooks/pull_request_closed_action.json",
			want: &internal.Webhook{
				Hook:  wantHook,
				Repo:  nil,
				Build: nil,
			},
		},
		{
			name:     "closed state",
			testData: "testdata/hooks/pull_request_closed_state.json",
			want: &internal.Webhook{
				Hook:  wantHook,
				Repo:  nil,
				Build: nil,
			},
		},
		{
			name:     "labeled documentation",
			testData: "testdata/hooks/pull_request_labeled.json",
			want: &internal.Webhook{
				PullRequest: internal.PullRequest{
					Number:     wantHook.GetNumber(),
					IsFromFork: false,
					Labels:     []string{"documentation"},
				},
				Hook:  wantHook,
				Repo:  wantRepo,
				Build: wantBuild2,
			},
		},
		{
			name:     "unlabeled documentation",
			testData: "testdata/hooks/pull_request_unlabeled.json",
			want: &internal.Webhook{
				PullRequest: internal.PullRequest{
					Number:     wantHook.GetNumber(),
					IsFromFork: false,
					Labels:     []string{"documentation"},
				},
				Hook:  wantHook,
				Repo:  wantRepo,
				Build: wantBuild3,
			},
		},
		{
			name:     "edited while labeled documentation",
			testData: "testdata/hooks/pull_request_edited_while_labeled.json",
			want: &internal.Webhook{
				PullRequest: internal.PullRequest{
					Number:     wantHook.GetNumber(),
					IsFromFork: false,
					Labels:     []string{"documentation", "enhancement"},
				},
				Hook:  wantHook,
				Repo:  wantRepo,
				Build: wantBuild4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := os.Open(tt.testData)
			if err != nil {
				t.Errorf("unable to open file: %v", err)
			}

			defer body.Close()

			request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
			request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
			request.Header.Set("X-GitHub-Hook-ID", "123456")
			request.Header.Set("X-GitHub-Host", "github.com")
			request.Header.Set("X-GitHub-Version", "2.16.0")
			request.Header.Set("X-GitHub-Event", "pull_request")

			client, _ := NewTest(s.URL)

			got, err := client.ProcessWebhook(context.TODO(), request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ProcessWebhook() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGithub_ProcessWebhook_Deployment(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetBranch("main")
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")
	wantHook.SetHost("github.com")
	wantHook.SetEvent(constants.EventDeploy)
	wantHook.SetEventAction(constants.ActionCreated)
	wantHook.SetStatus(constants.StatusSuccess)

	wantRepo := new(api.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics(nil)

	wantBuild := new(api.Build)
	wantBuild.SetEvent(constants.EventDeploy)
	wantBuild.SetEventAction(constants.ActionCreated)
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetDeploy("production")
	wantBuild.SetDeployNumber(145988746)
	wantBuild.SetSource("https://api.github.com/repos/Codertocat/Hello-World/deployments/145988746")
	wantBuild.SetTitle("deployment received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("")
	wantBuild.SetCommit("f95f852bd8fca8fcc58a9a2d6c842781e32a215e")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetSenderSCMID("21031067")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("")
	wantBuild.SetBranch("main")
	wantBuild.SetRef("refs/heads/main")

	wantDeployment := new(api.Deployment)
	wantDeployment.SetNumber(145988746)
	wantDeployment.SetURL("https://api.github.com/repos/Codertocat/Hello-World/deployments/145988746")
	wantDeployment.SetCommit("f95f852bd8fca8fcc58a9a2d6c842781e32a215e")
	wantDeployment.SetRef("main")
	wantDeployment.SetTask("deploy")
	wantDeployment.SetTarget("production")
	wantDeployment.SetDescription("")
	wantDeployment.SetCreatedAt(time.Now().UTC().Unix())
	wantDeployment.SetCreatedBy("Codertocat")

	type args struct {
		file              string
		hook              *api.Hook
		repo              *api.Repo
		build             *api.Build
		deploymentPayload raw.StringSliceMap
		deployment        *api.Deployment
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{file: "deployment.json", hook: wantHook, repo: wantRepo, build: wantBuild, deploymentPayload: raw.StringSliceMap{"foo": "test1", "bar": "test2"}, deployment: wantDeployment}, false},
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

			request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
			request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
			request.Header.Set("X-GitHub-Hook-ID", "123456")
			request.Header.Set("X-GitHub-Host", "github.com")
			request.Header.Set("X-GitHub-Version", "2.16.0")
			request.Header.Set("X-GitHub-Event", "deployment")

			client, _ := NewTest(s.URL)
			wantBuild.SetDeployPayload(tt.args.deploymentPayload)

			want := &internal.Webhook{
				Hook:       tt.args.hook,
				Repo:       tt.args.repo,
				Build:      tt.args.build,
				Deployment: tt.args.deployment,
			}

			got, err := client.ProcessWebhook(context.TODO(), request)
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

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "deployment")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetBranch("main")
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")
	wantHook.SetHost("github.com")
	wantHook.SetEvent(constants.EventDeploy)
	wantHook.SetEventAction(constants.ActionCreated)
	wantHook.SetStatus(constants.StatusSuccess)

	wantRepo := new(api.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics(nil)

	wantBuild := new(api.Build)
	wantBuild.SetEvent(constants.EventDeploy)
	wantBuild.SetEventAction(constants.ActionCreated)
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetDeploy("production")
	wantBuild.SetDeployNumber(145988746)
	wantBuild.SetSource("https://api.github.com/repos/Codertocat/Hello-World/deployments/145988746")
	wantBuild.SetTitle("deployment received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("")
	wantBuild.SetCommit("f95f852bd8fca8fcc58a9a2d6c842781e32a215e")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetSenderSCMID("21031067")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("")
	wantBuild.SetBranch("main")
	wantBuild.SetRef("refs/heads/main")

	wantDeployment := new(api.Deployment)
	wantDeployment.SetNumber(145988746)
	wantDeployment.SetURL("https://api.github.com/repos/Codertocat/Hello-World/deployments/145988746")
	wantDeployment.SetCommit("f95f852bd8fca8fcc58a9a2d6c842781e32a215e")
	wantDeployment.SetRef("f95f852bd8fca8fcc58a9a2d6c842781e32a215e")
	wantDeployment.SetTask("deploy")
	wantDeployment.SetTarget("production")
	wantDeployment.SetDescription("")
	//wantDeployment.SetPayload(map[string]string{"foo": "test1"})
	wantDeployment.SetCreatedAt(time.Now().UTC().Unix())
	wantDeployment.SetCreatedBy("Codertocat")

	want := &internal.Webhook{
		Hook:       wantHook,
		Repo:       wantRepo,
		Build:      wantBuild,
		Deployment: wantDeployment,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

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

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "foobar")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("foobar")
	wantHook.SetStatus(constants.StatusSuccess)

	want := &internal.Webhook{
		Hook:       wantHook,
		Repo:       nil,
		Build:      nil,
		Deployment: nil,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

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

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "foobar")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "pull_request")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("pull_request")
	wantHook.SetStatus(constants.StatusSuccess)

	want := &internal.Webhook{
		Hook:       wantHook,
		Repo:       nil,
		Build:      nil,
		Deployment: nil,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

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

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "deployment")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	err = client.VerifyWebhook(context.TODO(), request, new(api.Repo))
	if err != nil {
		t.Errorf("VerifyWebhook should have returned err")
	}
}

func TestGithub_VerifyWebhook_NoSecret(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	r := new(api.Repo)
	r.SetOrg("Codertocat")
	r.SetName("Hello-World")
	r.SetFullName("Codertocat/Hello-World")
	r.SetLink("https://github.com/Codertocat/Hello-World")
	r.SetClone("https://github.com/Codertocat/Hello-World.git")
	r.SetBranch("main")
	r.SetPrivate(false)

	// setup request
	body, err := os.Open("testdata/hooks/push.json")
	if err != nil {
		t.Errorf("Opening file for ProcessWebhook returned err: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "push")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	err = client.VerifyWebhook(context.TODO(), request, r)
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

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "issue_comment")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("comment")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://github.com/Codertocat/Hello-World")
	wantRepo.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics(nil)

	wantBuild := new(api.Build)
	wantBuild.SetEvent("comment")
	wantBuild.SetEventAction("created")
	wantBuild.SetClone("https://github.com/Codertocat/Hello-World.git")
	wantBuild.SetSource("https://github.com/Codertocat/Hello-World/pull/1")
	wantBuild.SetTitle("comment received from https://github.com/Codertocat/Hello-World")
	wantBuild.SetMessage("Update the README with new information")
	wantBuild.SetSender("Codertocat")
	wantBuild.SetSenderSCMID("2172")
	wantBuild.SetAuthor("Codertocat")
	wantBuild.SetEmail("")
	wantBuild.SetRef("refs/pull/1/head")

	want := &internal.Webhook{
		PullRequest: internal.PullRequest{
			Comment: "ok to test",
			Number:  wantHook.GetNumber(),
		},
		Hook:  wantHook,
		Repo:  wantRepo,
		Build: wantBuild,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ProcessWebhook() mismatch (-want +got):\n%s", diff)
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

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "issue_comment")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("comment")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	want := &internal.Webhook{
		Hook: wantHook,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

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

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Host", "github.com")
	request.Header.Set("X-GitHub-Version", "2.16.0")
	request.Header.Set("X-GitHub-Event", "issue_comment")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent("comment")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	want := &internal.Webhook{
		Hook: wantHook,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGitHub_ProcessWebhook_RepositoryRename(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/repository_rename.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Event", "repository")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent(constants.EventRepository)
	wantHook.SetEventAction(constants.ActionRenamed)
	wantHook.SetBranch("main")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetActive(true)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://octocoders.github.io/Codertocat/Hello-World")
	wantRepo.SetClone("https://octocoders.github.io/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics(nil)

	want := &internal.Webhook{
		Hook: wantHook,
		Repo: wantRepo,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGitHub_ProcessWebhook_RepositoryTransfer(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/repository_transferred.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Event", "repository")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent(constants.EventRepository)
	wantHook.SetEventAction(constants.ActionTransferred)
	wantHook.SetBranch("main")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetActive(true)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://octocoders.github.io/Codertocat/Hello-World")
	wantRepo.SetClone("https://octocoders.github.io/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics(nil)

	want := &internal.Webhook{
		Hook: wantHook,
		Repo: wantRepo,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGitHub_ProcessWebhook_RepositoryArchived(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/repository_archived.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Event", "repository")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent(constants.EventRepository)
	wantHook.SetEventAction("archived")
	wantHook.SetBranch("main")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetActive(false)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://octocoders.github.io/Codertocat/Hello-World")
	wantRepo.SetClone("https://octocoders.github.io/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics(nil)

	want := &internal.Webhook{
		Hook: wantHook,
		Repo: wantRepo,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGitHub_ProcessWebhook_RepositoryEdited(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/repository_edited.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Event", "repository")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent(constants.EventRepository)
	wantHook.SetEventAction(constants.ActionEdited)
	wantHook.SetBranch("main")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetActive(true)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://octocoders.github.io/Codertocat/Hello-World")
	wantRepo.SetClone("https://octocoders.github.io/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetTopics([]string{"cloud", "security"})
	wantRepo.SetPrivate(false)

	want := &internal.Webhook{
		Hook: wantHook,
		Repo: wantRepo,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGitHub_ProcessWebhook_Repository(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup request
	body, err := os.Open("testdata/hooks/repository_publicized.json")
	if err != nil {
		t.Errorf("unable to open file: %v", err)
	}

	defer body.Close()

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "GitHub-Hookshot/a22606a")
	request.Header.Set("X-GitHub-Delivery", "7bd477e4-4415-11e9-9359-0d41fdf9567e")
	request.Header.Set("X-GitHub-Hook-ID", "123456")
	request.Header.Set("X-GitHub-Event", "repository")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	wantHook := new(api.Hook)
	wantHook.SetNumber(1)
	wantHook.SetSourceID("7bd477e4-4415-11e9-9359-0d41fdf9567e")
	wantHook.SetWebhookID(123456)
	wantHook.SetCreated(time.Now().UTC().Unix())
	wantHook.SetHost("github.com")
	wantHook.SetEvent(constants.EventRepository)
	wantHook.SetEventAction("publicized")
	wantHook.SetBranch("main")
	wantHook.SetStatus(constants.StatusSuccess)
	wantHook.SetLink("https://github.com/Codertocat/Hello-World/settings/hooks")

	wantRepo := new(api.Repo)
	wantRepo.SetActive(true)
	wantRepo.SetOrg("Codertocat")
	wantRepo.SetName("Hello-World")
	wantRepo.SetFullName("Codertocat/Hello-World")
	wantRepo.SetLink("https://octocoders.github.io/Codertocat/Hello-World")
	wantRepo.SetClone("https://octocoders.github.io/Codertocat/Hello-World.git")
	wantRepo.SetBranch("main")
	wantRepo.SetPrivate(false)
	wantRepo.SetTopics(nil)

	want := &internal.Webhook{
		Hook: wantHook,
		Repo: wantRepo,
	}

	got, err := client.ProcessWebhook(context.TODO(), request)

	if err != nil {
		t.Errorf("ProcessWebhook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessWebhook is %v, want %v", got, want)
	}
}

func TestGithub_Redeliver_Webhook(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/hooks/:repo_id/deliveries/:delivery_id/attempts", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/hooks/push.json")
	})
	engine.GET("/api/v3/repos/:org/:repo/hooks/:hook_id/deliveries", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/delivery_summaries.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(api.User)
	u.SetName("octocat")
	u.SetToken("foo")

	_repo := new(api.Repo)
	_repo.SetID(1)
	_repo.SetOwner(u)
	_repo.SetName("bar")
	_repo.SetOrg("foo")

	_hook := new(api.Hook)
	_hook.SetSourceID("b595f0e0-aee1-11ec-86cf-9418381395c4")
	_hook.SetID(1)
	_hook.SetRepo(_repo)
	_hook.SetNumber(1)
	_hook.SetWebhookID(1234)

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	err := client.RedeliverWebhook(context.TODO(), u, _hook)

	if err != nil {
		t.Errorf("RedeliverWebhook returned err: %v", err)
	}
}

func TestGithub_GetDeliveryID(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	engine.GET("/api/v3/repos/:org/:repo/hooks/:hook_id/deliveries", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/delivery_summaries.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(api.User)
	u.SetName("octocat")
	u.SetToken("foo")

	_repo := new(api.Repo)
	_repo.SetID(1)
	_repo.SetOwner(u)
	_repo.SetName("bar")
	_repo.SetOrg("foo")

	_hook := new(api.Hook)
	_hook.SetSourceID("b595f0e0-aee1-11ec-86cf-9418381395c4")
	_hook.SetID(1)
	_hook.SetRepo(_repo)
	_hook.SetNumber(1)
	_hook.SetWebhookID(1234)

	want := int64(22948188373)

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	ghClient := client.newClientToken(context.Background(), *u.Token)

	// run test
	got, err := client.getDeliveryID(context.TODO(), ghClient, _hook)

	if err != nil {
		t.Errorf("RedeliverWebhook returned err: %v", err)
	}

	if got != want {
		t.Errorf("getDeliveryID returned: %v; want: %v", got, want)
	}
}
