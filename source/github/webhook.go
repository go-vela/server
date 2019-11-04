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

	r := &library.Repo{
		Org:      repo.GetOwner().Login,
		Name:     repo.Name,
		FullName: repo.FullName,
		Link:     repo.HTMLURL,
		Clone:    repo.CloneURL,
		Branch:   repo.DefaultBranch,
		Private:  repo.Private,
	}

	title := fmt.Sprintf("%s received from %s", constants.EventPush, repo.GetHTMLURL())
	event := constants.EventPush
	branch := strings.Replace(payload.GetRef(), "refs/heads/", "", -1)

	bClone := repo.GetCloneURL()
	bSource := payload.GetHeadCommit().GetURL()
	bMessage := payload.GetHeadCommit().GetMessage()
	bCommit := payload.GetHeadCommit().GetID()
	bSender := payload.GetSender().GetLogin()
	bAuthor := payload.GetSender().GetLogin()
	bRef := payload.GetRef()
	bBaseRef := payload.GetBaseRef()

	b := &library.Build{
		Event:   &event,
		Clone:   &bClone,
		Source:  &bSource,
		Title:   &title,
		Message: &bMessage,
		Commit:  &bCommit,
		Sender:  &bSender,
		Author:  &bAuthor,
		Branch:  &branch,
		Ref:     &bRef,
		BaseRef: &bBaseRef,
	}

	// ensure build author is set
	if len(b.GetAuthor()) == 0 {
		b.Author = payload.GetHeadCommit().GetAuthor().Login
		b.Sender = b.Author
	}

	// handle when push event is a tag
	if strings.HasPrefix(b.GetRef(), "refs/tags/") {
		// set the proper event for the build
		*b.Event = constants.EventTag

		// set the proper branch from the base ref
		if strings.HasPrefix(payload.GetBaseRef(), "refs/heads/") {
			*b.Branch = strings.Replace(payload.GetBaseRef(), "refs/heads/", "", -1)
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

	r := &library.Repo{
		Org:      repo.GetOwner().Login,
		Name:     repo.Name,
		FullName: repo.FullName,
		Link:     repo.HTMLURL,
		Clone:    repo.CloneURL,
		Branch:   repo.DefaultBranch,
		Private:  repo.Private,
	}

	event := constants.EventPull
	clone := repo.GetCloneURL()
	source := payload.GetPullRequest().GetHTMLURL()
	title := fmt.Sprintf("%s received from %s", constants.EventPull, repo.GetHTMLURL())
	ref := fmt.Sprintf("refs/pull/%d/head", payload.GetNumber())
	b := &library.Build{
		Event:   &event,
		Clone:   &clone,
		Source:  &source,
		Title:   &title,
		Message: payload.GetPullRequest().Title,
		Commit:  payload.GetPullRequest().GetHead().SHA,
		Sender:  payload.GetSender().Login,
		Author:  payload.GetSender().Login,
		Branch:  payload.GetPullRequest().GetBase().Ref,
		Ref:     &ref,
		BaseRef: payload.GetPullRequest().GetBase().Ref,
	}

	// ensure the build reference is set
	if payload.GetPullRequest().GetMerged() {
		*b.Ref = fmt.Sprintf("refs/pull/%d/merge", payload.GetNumber())
	}

	// ensure the build author and sender are set
	if len(b.GetAuthor()) == 0 {
		*b.Author = payload.GetPullRequest().GetUser().GetLogin()
		*b.Sender = b.GetAuthor()
	}

	return r, b, nil
}
