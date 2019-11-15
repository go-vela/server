// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v26/github"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// ProcessWebhook parses the webhook from a repo
func (c *client) ProcessWebhook(request *http.Request) (*library.Repo, *library.Build, error) {
	payload, err := github.ValidatePayload(request, nil)
	if err != nil {
		return nil, nil, err
	}

	// parse the payload from the webhook
	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		return nil, nil, err
	}

	// process the event from the webhook
	switch event := event.(type) {
	case *github.PushEvent:
		return processPushEvent(event)
	case *github.PullRequestEvent:
		return processPREvent(event)
	}

	return nil, nil, nil
}

// processPushEvent is a helper function to process the push event
func processPushEvent(payload *github.PushEvent) (*library.Repo, *library.Build, error) {
	repo := payload.GetRepo()

	r := new(library.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())

	b := new(library.Build)
	b.SetEvent(constants.EventPush)
	b.SetClone(repo.GetCloneURL())
	b.SetSource(payload.GetHeadCommit().GetURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventPush, repo.GetHTMLURL()))
	b.SetMessage(payload.GetHeadCommit().GetMessage())
	b.SetCommit(payload.GetHeadCommit().GetID())
	b.SetSender(payload.GetSender().GetLogin())
	b.SetAuthor(payload.GetSender().GetLogin())
	b.SetBranch(strings.Replace(payload.GetRef(), "refs/heads/", "", -1))
	b.SetRef(payload.GetRef())
	b.SetBaseRef(payload.GetBaseRef())

	// ensure build author is set
	if len(b.GetAuthor()) == 0 {
		b.SetAuthor(payload.GetHeadCommit().GetAuthor().GetLogin())
		b.SetSender(b.GetAuthor())
	}

	// handle when push event is a tag
	if strings.HasPrefix(b.GetRef(), "refs/tags/") {
		// set the proper event for the build
		b.SetEvent(constants.EventTag)

		// set the proper branch from the base ref
		if strings.HasPrefix(payload.GetBaseRef(), "refs/heads/") {
			b.SetBranch(strings.Replace(payload.GetBaseRef(), "refs/heads/", "", -1))
		}
	}

	return r, b, nil
}

// processPREvent is a helper function to process the pull_request event
func processPREvent(payload *github.PullRequestEvent) (*library.Repo, *library.Build, error) {
	// if the pull request state isn't open we ignore it
	if payload.GetPullRequest().GetState() != "open" {
		return nil, nil, nil
	}

	// skip if the pull request action is not opened or synchronize
	if !strings.EqualFold(payload.GetAction(), "opened") &&
		!strings.EqualFold(payload.GetAction(), "synchronize") {
		return nil, nil, nil
	}

	// capture the repo from the payload
	repo := payload.GetRepo()

	r := new(library.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())

	b := new(library.Build)
	b.SetEvent(constants.EventPull)
	b.SetClone(repo.GetCloneURL())
	b.SetSource(payload.GetPullRequest().GetHTMLURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventPull, repo.GetHTMLURL()))
	b.SetMessage(payload.GetPullRequest().GetTitle())
	b.SetCommit(payload.GetPullRequest().GetHead().GetSHA())
	b.SetSender(payload.GetSender().GetLogin())
	b.SetAuthor(payload.GetSender().GetLogin())
	b.SetBranch(payload.GetPullRequest().GetBase().GetRef())
	b.SetRef(fmt.Sprintf("refs/pull/%d/head", payload.GetNumber()))
	b.SetBaseRef(payload.GetPullRequest().GetBase().GetRef())

	// ensure the build reference is set
	if payload.GetPullRequest().GetMerged() {
		b.SetRef(fmt.Sprintf("refs/pull/%d/merge", payload.GetNumber()))
	}

	// ensure the build author and sender are set
	if len(b.GetAuthor()) == 0 {
		b.SetAuthor(payload.GetPullRequest().GetUser().GetLogin())
		b.SetSender(b.GetAuthor())
	}

	return r, b, nil
}
