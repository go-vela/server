// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"
	"mime"
	"net/http"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/google/go-github/v63/github"
)

// ProcessGitHubAppWebhook parses the webhook from a GitHub App.
//
//nolint:nilerr // ignore webhook returning nil
func (c *client) ProcessGitHubAppWebhook(ctx context.Context, request *http.Request) error {
	c.Logger.Tracef("processing GitHub App webhook")

	// get content type
	contentType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	payload, err := github.ValidatePayloadFromBody(contentType, request.Body, "", nil)
	if err != nil {
		return err
	}

	// parse the payload from the webhook
	event, err := github.ParseWebHook(github.WebHookType(request), payload)

	if err != nil {
		return err
	}

	// TODO: Implement
	// process the event from the webhook
	switch event := event.(type) {
	case *github.InstallationEvent:
		return c.processInstallationEvent(ctx, event)
	case *github.InstallationRepositoriesEvent:
		return c.processInstallationRepositoriesEvent(ctx, event)
	}

	return nil
}

func (c *client) GetInstallations() ([]*github.Installation, error) {
	client, err := c.newClientGitHubApp()
	if err != nil {
		return nil, err
	}

	// list all installations (a.k.a. orgs) where the GitHub App is installed
	installations, _, err := client.Apps.ListInstallations(context.Background(), &github.ListOptions{})
	if err != nil {
		return nil, err
	}
	return installations, nil
}

func (c *client) GetInstallationRepos(installation *github.Installation) ([]*github.Repository, error) {
	client, err := c.newClientGitHubAppInstallation(installation)
	if err != nil {
		return nil, err
	}

	// lists the repositories that are accessible to the authenticated installation
	repos, _, err := client.Apps.ListRepos(context.Background(), &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return repos.Repositories, nil
}

func (c *client) GetInstallationAccessToken(r *api.Repo) (string, error) {
	client, err := c.newClientGitHubApp()
	if err != nil {
		return "", err
	}

	// if repo has an install ID, use it to create an installation token
	if r.GetInstallID() != 0 {
		// create installation token for the repo
		t, _, err := client.Apps.CreateInstallationToken(context.Background(), r.GetInstallID(), &github.InstallationTokenOptions{})
		if err != nil {
			panic(err)
		}

		return t.GetToken(), nil
	}
	return "", nil
}

// CreateChecks defines a function that does stuff...
func (c *client) CreateChecks(ctx context.Context, r *api.Repo, commit, step, event string) (int64, error) {
	installation, err := c.getInstallationFromRepo(r)

	if err != nil {
		return 0, err
	}

	client, err := c.newClientGitHubAppInstallation(installation)

	if err != nil {
		return 0, err
	}

	opts := github.CreateCheckRunOptions{
		Name:    fmt.Sprintf("vela-%s-%s", event, step),
		HeadSHA: commit,
	}

	check, _, err := client.Checks.CreateCheckRun(ctx, r.GetOrg(), r.GetName(), opts)
	if err != nil {
		return 0, err
	}

	return check.GetID(), nil
}

// UpdateChecks defines a function that does stuff...
func (c *client) UpdateChecks(ctx context.Context, r *api.Repo, s *library.Step, commit, event string) error {
	installation, err := c.getInstallationFromRepo(r)

	if err != nil {
		return err
	}

	client, err := c.newClientGitHubAppInstallation(installation)

	if err != nil {
		return err
	}

	var (
		conclusion string
		status     string
	)
	// set the conclusion and status for the step check depending on what the status of the step is
	switch s.GetStatus() {
	case constants.StatusPending:
		conclusion = "neutral"
		status = "queued"
	case constants.StatusPendingApproval:
		conclusion = "action_required"
		status = "queued"
	case constants.StatusRunning:
		conclusion = "neutral"
		status = "in_progress"
	case constants.StatusSuccess:
		conclusion = "success"
		status = "completed"
	case constants.StatusFailure:
		conclusion = "failure"
		status = "completed"
	case constants.StatusCanceled:
		conclusion = "cancelled"
		status = "completed"
	case constants.StatusKilled:
		conclusion = "cancelled"
		status = "completed"
	case constants.StatusSkipped:
		conclusion = "skipped"
		status = "completed"
	default:
		conclusion = "neutral"
		status = "completed"
	}

	var annotations []*github.CheckRunAnnotation

	for _, reportAnnotation := range s.GetReport().GetAnnotations() {
		annotation := &github.CheckRunAnnotation{
			Path:            github.String(reportAnnotation.GetPath()),
			StartLine:       github.Int(reportAnnotation.GetStartLine()),
			EndLine:         github.Int(reportAnnotation.GetEndLine()),
			StartColumn:     github.Int(reportAnnotation.GetStartColumn()),
			EndColumn:       github.Int(reportAnnotation.GetEndColumn()),
			AnnotationLevel: github.String(reportAnnotation.GetAnnotationLevel()),
			Message:         github.String(reportAnnotation.GetMessage()),
			Title:           github.String(reportAnnotation.GetTitle()),
			RawDetails:      github.String(reportAnnotation.GetRawDetails()),
		}

		annotations = append(annotations, annotation)
	}

	output := &github.CheckRunOutput{
		Title:            github.String(s.GetReport().GetTitle()),
		Summary:          github.String(s.GetReport().GetSummary()),
		Text:             github.String(s.GetReport().GetText()),
		AnnotationsCount: github.Int(s.GetReport().GetAnnotationsCount()),
		AnnotationsURL:   github.String(s.GetReport().GetAnnotationsURL()),
		Annotations:      annotations,
	}

	opts := github.UpdateCheckRunOptions{
		Name:       fmt.Sprintf("vela-%s-%s", event, s.GetName()),
		Conclusion: github.String(conclusion),
		Status:     github.String(status),
		Output:     output,
	}

	_, _, err = client.Checks.UpdateCheckRun(ctx, r.GetOrg(), r.GetName(), s.GetCheckID(), opts)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) getInstallationFromRepo(r *api.Repo) (*github.Installation, error) {
	client, err := c.newClientGitHubApp()
	if err != nil {
		return nil, err
	}

	// if repo has an install ID, use it to find installation
	if r.GetInstallID() != 0 {
		// create installation token for the repo
		installation, _, err := client.Apps.GetInstallation(context.Background(), r.GetInstallID())
		if err != nil {
			return nil, err
		}

		return installation, nil
	}

	// TODO: Handle case where repo doesn't have an install ID, but app is installed on repo

	return nil, nil
}

// processPushEvent is a helper function to process an installation event.
func (c *client) processInstallationEvent(ctx context.Context, payload *github.InstallationEvent) error {
	// TODO: Implement

	return nil
}

// processPushEvent is a helper function to process an installation repositories event.
func (c *client) processInstallationRepositoriesEvent(ctx context.Context, payload *github.InstallationRepositoriesEvent) error {
	// TODO: Implement

	return nil
}
