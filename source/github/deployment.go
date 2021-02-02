// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"github.com/google/go-github/v29/github"

	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetDeployment gets a deployment from the GitHub repo.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) GetDeployment(u *library.User, r *library.Repo, id int64) (*library.Deployment, error) {
	logrus.Tracef("capturing deployment %d for %s", id, r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	// send API call to capture the deployment
	deployment, _, err := client.Repositories.GetDeployment(ctx, r.GetOrg(), r.GetName(), id)
	if err != nil {
		return nil, err
	}

	return &library.Deployment{
		ID:          deployment.ID,
		RepoID:      r.ID,
		URL:         deployment.URL,
		User:        deployment.Creator.Login,
		Commit:      deployment.SHA,
		Ref:         deployment.Ref,
		Task:        deployment.Task,
		Target:      deployment.Environment,
		Description: deployment.Description,
	}, nil
}

// GetDeploymentCount counts a list of deployments from the GitHub repo.
func (c *client) GetDeploymentCount(u *library.User, r *library.Repo) (int64, error) {
	logrus.Tracef("counting deployments for %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)
	// create variable to track the deployments
	deployments := []*github.Deployment{}

	// set pagination options for listing deployments
	opts := &github.DeploymentsListOptions{
		// set the max per page for the options
		// to capture the list of deployments
		ListOptions: github.ListOptions{
			PerPage: 100, // 100 is max
		},
	}

	for {
		// send API call to capture the list of deployments
		d, resp, err := client.Repositories.ListDeployments(ctx, r.GetOrg(), r.GetName(), opts)
		if err != nil {
			return 0, err
		}

		deployments = append(deployments, d...)

		// break the loop if there is no more results to page through
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return int64(len(deployments)), nil
}

// GetDeploymentList gets a list of deployments from the GitHub repo.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) GetDeploymentList(u *library.User, r *library.Repo, page, perPage int) ([]*library.Deployment, error) {
	logrus.Tracef("capturing deployments for %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	// set pagination options for listing deployments
	opts := &github.DeploymentsListOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}

	// send API call to capture the list of deployments
	d, _, err := client.Repositories.ListDeployments(ctx, r.GetOrg(), r.GetName(), opts)
	if err != nil {
		return nil, err
	}

	// variable we want to return
	deployments := []*library.Deployment{}

	// iterate through all API results
	for _, deployment := range d {
		// convert query result to library type
		deployments = append(deployments, &library.Deployment{
			ID:          deployment.ID,
			RepoID:      r.ID,
			URL:         deployment.URL,
			User:        deployment.Creator.Login,
			Commit:      deployment.SHA,
			Ref:         deployment.Ref,
			Task:        deployment.Task,
			Target:      deployment.Environment,
			Description: deployment.Description,
		})
	}

	return deployments, nil
}

// CreateDeployment creates a new deployment for the GitHub repo.
func (c *client) CreateDeployment(u *library.User, r *library.Repo, d *library.Deployment) error {
	logrus.Tracef("creating deployment for %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	var payload interface{}
	if d.Payload == nil {
		payload = github.String("")
	} else {
		payload = d.Payload
	}

	// create the hook object to make the API call
	deployment := &github.DeploymentRequest{
		Ref:              d.Ref,
		Task:             d.Task,
		AutoMerge:        github.Bool(false),
		RequiredContexts: &[]string{},
		Payload:          payload,
		Environment:      d.Target,
		Description:      d.Description,
	}

	// send API call to create the deployment
	deploy, _, err := client.Repositories.CreateDeployment(ctx, r.GetOrg(), r.GetName(), deployment)
	if err != nil {
		return err
	}

	d.SetID(deploy.GetID())
	d.SetRepoID(r.GetID())
	d.SetURL(deploy.GetURL())
	d.SetUser(deploy.GetCreator().GetLogin())
	d.SetCommit(deploy.GetSHA())
	d.SetRef(deploy.GetRef())
	d.SetTask(deploy.GetTask())
	d.SetTarget(deploy.GetEnvironment())
	d.SetDescription(deploy.GetDescription())

	return nil
}
