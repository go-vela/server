// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v26/github"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// ProcessWebhook parses the webhook from a repo
func (c *client) ProcessWebhook(request *http.Request) (*library.Hook, *library.Repo, *library.Build, error) {
	h := new(library.Hook)
	h.SetNumber(1)
	h.SetSourceID(request.Header.Get("X-GitHub-Delivery"))
	h.SetCreated(time.Now().UTC().Unix())
	h.SetHost("github.com")
	h.SetEvent(request.Header.Get("X-GitHub-Event"))
	h.SetStatus(constants.StatusSuccess)

	if len(request.Header.Get("X-GitHub-Enterprise-Host")) > 0 {
		h.SetHost(request.Header.Get("X-GitHub-Enterprise-Host"))
	}

	payload, err := github.ValidatePayload(request, nil)
	if err != nil {
		return h, nil, nil, err
	}

	// parse the payload from the webhook
	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		return h, nil, nil, err
	}

	// process the event from the webhook
	switch event := event.(type) {
	case *github.PushEvent:
		return processPushEvent(h, event)
	case *github.PullRequestEvent:
		return processPREvent(h, event)
	}

	return h, nil, nil, nil
}

// processPushEvent is a helper function to process the push event
func processPushEvent(h *library.Hook, payload *github.PushEvent) (*library.Hook, *library.Repo, *library.Build, error) {
	repo := payload.GetRepo()

	// convert payload to library repo
	r := new(library.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())

	// convert payload to library build
	b := new(library.Build)
	b.SetEvent(constants.EventPush)
	b.SetClone(repo.GetCloneURL())
	b.SetSource(payload.GetHeadCommit().GetURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventPush, repo.GetHTMLURL()))
	b.SetMessage(payload.GetHeadCommit().GetMessage())
	b.SetCommit(payload.GetHeadCommit().GetID())
	b.SetSender(payload.GetSender().GetLogin())
	b.SetAuthor(payload.GetHeadCommit().GetAuthor().GetLogin())
	b.SetEmail(payload.GetHeadCommit().GetAuthor().GetEmail())
	b.SetBranch(strings.Replace(payload.GetRef(), "refs/heads/", "", -1))
	b.SetRef(payload.GetRef())
	b.SetBaseRef(payload.GetBaseRef())

	// update the hook object
	h.SetBranch(b.GetBranch())
	h.SetEvent(constants.EventPush)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), r.GetFullName()),
	)

	// ensure the build author is set
	if len(b.GetAuthor()) == 0 {
		b.SetAuthor(payload.GetHeadCommit().GetCommitter().GetName())
	}

	// ensure the build sender is set
	if len(b.GetSender()) == 0 {
		b.SetSender(payload.GetPusher().GetName())
	}

	// ensure the build email is set
	if len(b.GetEmail()) == 0 {
		b.SetEmail(payload.GetHeadCommit().GetCommitter().GetEmail())
	}

	// handle when push event is a tag
	if strings.HasPrefix(b.GetRef(), "refs/tags/") {
		// set the proper event for the hook
		h.SetEvent(constants.EventTag)
		// set the proper event for the build
		b.SetEvent(constants.EventTag)

		// set the proper branch from the base ref
		if strings.HasPrefix(payload.GetBaseRef(), "refs/heads/") {
			b.SetBranch(strings.Replace(payload.GetBaseRef(), "refs/heads/", "", -1))
		}
	}

	return h, r, b, nil
}

// processPREvent is a helper function to process the pull_request event
func processPREvent(h *library.Hook, payload *github.PullRequestEvent) (*library.Hook, *library.Repo, *library.Build, error) {
	// update the hook object
	h.SetBranch(payload.GetPullRequest().GetBase().GetRef())
	h.SetEvent(constants.EventPull)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), payload.GetRepo().GetFullName()),
	)

	// if the pull request state isn't open we ignore it
	if payload.GetPullRequest().GetState() != "open" {
		return h, nil, nil, nil
	}

	// skip if the pull request action is not opened or synchronize
	if !strings.EqualFold(payload.GetAction(), "opened") &&
		!strings.EqualFold(payload.GetAction(), "synchronize") {
		return h, nil, nil, nil
	}

	// capture the repo from the payload
	repo := payload.GetRepo()

	// convert payload to library repo
	r := new(library.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())

	// convert payload to library build
	b := new(library.Build)
	b.SetEvent(constants.EventPull)
	b.SetClone(repo.GetCloneURL())
	b.SetSource(payload.GetPullRequest().GetHTMLURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventPull, repo.GetHTMLURL()))
	b.SetMessage(payload.GetPullRequest().GetTitle())
	b.SetCommit(payload.GetPullRequest().GetHead().GetSHA())
	b.SetSender(payload.GetSender().GetLogin())
	b.SetAuthor(payload.GetPullRequest().GetUser().GetLogin())
	b.SetEmail(payload.GetPullRequest().GetUser().GetEmail())
	b.SetBranch(payload.GetPullRequest().GetBase().GetRef())
	b.SetRef(fmt.Sprintf("refs/pull/%d/head", payload.GetNumber()))
	b.SetBaseRef(payload.GetPullRequest().GetBase().GetRef())

	// ensure the build reference is set
	if payload.GetPullRequest().GetMerged() {
		b.SetRef(fmt.Sprintf("refs/pull/%d/merge", payload.GetNumber()))
	}

	// ensure the build author is set
	if len(b.GetAuthor()) == 0 {
		b.SetAuthor(payload.GetPullRequest().GetHead().GetUser().GetLogin())
	}

	// ensure the build sender is set
	if len(b.GetSender()) == 0 {
		b.SetSender(payload.GetPullRequest().GetUser().GetLogin())
	}

	// ensure the build email is set
	if len(b.GetEmail()) == 0 {
		b.SetEmail(payload.GetPullRequest().GetHead().GetUser().GetEmail())
	}

	return h, r, b, nil
}
