// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwk"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/actions"
	"github.com/go-vela/server/compiler/types/raw"
)

// API TEST RESOURCES
//
// These are API resources initialized to their zero values for testing.

func APIBuild() *api.Build {
	return &api.Build{
		ID:           new(int64),
		Repo:         APIRepo(),
		PipelineID:   new(int64),
		Number:       new(int),
		Parent:       new(int),
		Event:        new(string),
		EventAction:  new(string),
		Status:       new(string),
		Error:        new(string),
		Enqueued:     new(int64),
		Created:      new(int64),
		Started:      new(int64),
		Finished:     new(int64),
		Deploy:       new(string),
		DeployNumber: new(int64),
		Clone:        new(string),
		Source:       new(string),
		Title:        new(string),
		Message:      new(string),
		Commit:       new(string),
		Sender:       new(string),
		SenderSCMID:  new(string),
		Fork:         new(bool),
		Author:       new(string),
		Email:        new(string),
		Link:         new(string),
		Branch:       new(string),
		Ref:          new(string),
		BaseRef:      new(string),
		HeadRef:      new(string),
		Host:         new(string),
		Route:        new(string),
		Runtime:      new(string),
		Distribution: new(string),
		ApprovedAt:   new(int64),
		ApprovedBy:   new(string),
	}
}

func APIDeployment() *api.Deployment {
	return &api.Deployment{
		ID:          new(int64),
		Repo:        APIRepo(),
		Number:      new(int64),
		URL:         new(string),
		Commit:      new(string),
		Ref:         new(string),
		Task:        new(string),
		Target:      new(string),
		Description: new(string),
		Payload:     new(raw.StringSliceMap),
		CreatedAt:   new(int64),
		CreatedBy:   new(string),
	}
}

func APIEvents() *api.Events {
	return &api.Events{
		Push: &actions.Push{
			Branch:       new(bool),
			Tag:          new(bool),
			DeleteBranch: new(bool),
			DeleteTag:    new(bool),
		},
		PullRequest: &actions.Pull{
			Opened:      new(bool),
			Edited:      new(bool),
			Synchronize: new(bool),
			Reopened:    new(bool),
			Labeled:     new(bool),
			Unlabeled:   new(bool),
		},
		Deployment: &actions.Deploy{
			Created: new(bool),
		},
		Comment: &actions.Comment{
			Created: new(bool),
			Edited:  new(bool),
		},
		Schedule: &actions.Schedule{
			Run: new(bool),
		},
	}
}

func APIRepo() *api.Repo {
	return &api.Repo{
		ID:              new(int64),
		Owner:           APIUser(),
		BuildLimit:      new(int64),
		Timeout:         new(int64),
		Counter:         new(int),
		PipelineType:    new(string),
		Hash:            new(string),
		Org:             new(string),
		Name:            new(string),
		FullName:        new(string),
		Link:            new(string),
		Clone:           new(string),
		Branch:          new(string),
		Visibility:      new(string),
		PreviousName:    new(string),
		Private:         new(bool),
		Trusted:         new(bool),
		Active:          new(bool),
		AllowEvents:     APIEvents(),
		Topics:          new([]string),
		ApproveBuild:    new(string),
		ApprovalTimeout: new(int64),
		InstallID:       new(int64),
	}
}

func APIUser() *api.User {
	return &api.User{
		ID:           new(int64),
		Name:         new(string),
		RefreshToken: new(string),
		Token:        new(string),
		Favorites:    new([]string),
		Dashboards:   new([]string),
		Active:       new(bool),
		Admin:        new(bool),
	}
}

func APIHook() *api.Hook {
	return &api.Hook{
		ID:          new(int64),
		Repo:        APIRepo(),
		Build:       APIBuild(),
		Number:      new(int),
		SourceID:    new(string),
		Created:     new(int64),
		Host:        new(string),
		Event:       new(string),
		EventAction: new(string),
		Branch:      new(string),
		Error:       new(string),
		Status:      new(string),
		Link:        new(string),
		WebhookID:   new(int64),
	}
}

func APILog() *api.Log {
	return &api.Log{
		ID:        new(int64),
		RepoID:    new(int64),
		BuildID:   new(int64),
		ServiceID: new(int64),
		StepID:    new(int64),
		Data:      new([]byte),
	}
}

func APISchedule() *api.Schedule {
	return &api.Schedule{
		ID:          new(int64),
		Repo:        APIRepo(),
		Active:      new(bool),
		Name:        new(string),
		Entry:       new(string),
		CreatedAt:   new(int64),
		CreatedBy:   new(string),
		UpdatedAt:   new(int64),
		UpdatedBy:   new(string),
		ScheduledAt: new(int64),
		Branch:      new(string),
		Error:       new(string),
	}
}

func APISecret() *api.Secret {
	return &api.Secret{
		ID:                new(int64),
		Org:               new(string),
		Repo:              new(string),
		Team:              new(string),
		Name:              new(string),
		Value:             new(string),
		Type:              new(string),
		Images:            new([]string),
		AllowEvents:       APIEvents(),
		AllowCommand:      new(bool),
		AllowSubstitution: new(bool),
		CreatedAt:         new(int64),
		CreatedBy:         new(string),
		UpdatedAt:         new(int64),
		UpdatedBy:         new(string),
	}
}

func APIService() *api.Service {
	return &api.Service{
		ID:           new(int64),
		BuildID:      new(int64),
		RepoID:       new(int64),
		Number:       new(int),
		Name:         new(string),
		Image:        new(string),
		Status:       new(string),
		Error:        new(string),
		ExitCode:     new(int),
		Created:      new(int64),
		Started:      new(int64),
		Finished:     new(int64),
		Host:         new(string),
		Runtime:      new(string),
		Distribution: new(string),
	}
}

func APIStep() *api.Step {
	return &api.Step{
		ID:           new(int64),
		BuildID:      new(int64),
		RepoID:       new(int64),
		Number:       new(int),
		Name:         new(string),
		Image:        new(string),
		Stage:        new(string),
		Status:       new(string),
		Error:        new(string),
		ExitCode:     new(int),
		Created:      new(int64),
		Started:      new(int64),
		Finished:     new(int64),
		Host:         new(string),
		Runtime:      new(string),
		Distribution: new(string),
		ReportAs:     new(string),
	}
}

func APIPipeline() *api.Pipeline {
	return &api.Pipeline{
		ID:              new(int64),
		Repo:            APIRepo(),
		Commit:          new(string),
		Flavor:          new(string),
		Platform:        new(string),
		Ref:             new(string),
		Type:            new(string),
		Version:         new(string),
		ExternalSecrets: new(bool),
		InternalSecrets: new(bool),
		Services:        new(bool),
		Stages:          new(bool),
		Steps:           new(bool),
		Templates:       new(bool),
		Warnings:        new([]string),
		Data:            new([]byte),
	}
}

func APIDashboard() *api.Dashboard {
	return &api.Dashboard{
		ID:        new(string),
		Name:      new(string),
		CreatedAt: new(int64),
		CreatedBy: new(string),
		UpdatedAt: new(int64),
		UpdatedBy: new(string),
		Admins:    &[]*api.User{APIUser()},
		Repos:     &[]*api.DashboardRepo{APIDashboardRepo()},
	}
}

func APIDashboardRepo() *api.DashboardRepo {
	return &api.DashboardRepo{
		ID:       new(int64),
		Branches: new([]string),
		Events:   new([]string),
	}
}

func JWK() jwk.RSAPublicKey {
	privateRSAKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil
	}

	pubJwk, err := jwk.Import(privateRSAKey.PublicKey)
	if err != nil {
		return nil
	}

	switch j := pubJwk.(type) {
	case jwk.RSAPublicKey:
		// assign KID to key pair
		kid, err := uuid.NewV7()
		if err != nil {
			return nil
		}

		err = pubJwk.Set(jwk.KeyIDKey, kid.String())
		if err != nil {
			return nil
		}

		return j
	default:
		return nil
	}
}
