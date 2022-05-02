// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/raw"
	"github.com/google/go-github/v44/github"
)

// GetDeployment gets a deployment from the GitHub repo.
func (c *client) GetDeployment(u *library.User, r *library.Repo, id int64) (*library.Deployment, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("capturing deployment %d for repo %s", id, r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	// send API call to capture the deployment
	deployment, _, err := client.Repositories.GetDeployment(ctx, r.GetOrg(), r.GetName(), id)
	if err != nil {
		return nil, err
	}

	var payload *raw.StringSliceMap

	err = json.Unmarshal(deployment.Payload, &payload)
	if err != nil {
		c.Logger.Tracef("Unable to unmarshal payload for deployment id %v", deployment.ID)
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
		Payload:     payload,
	}, nil
}

// GetDeploymentCount counts a list of deployments from the GitHub repo.
func (c *client) GetDeploymentCount(u *library.User, r *library.Repo) (int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("counting deployments for repo %s", r.GetFullName())

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
func (c *client) GetDeploymentList(u *library.User, r *library.Repo, page, perPage int) ([]*library.Deployment, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("listing deployments for repo %s", r.GetFullName())

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
		var payload *raw.StringSliceMap

		err := json.Unmarshal(deployment.Payload, &payload)
		if err != nil {
			c.Logger.Tracef("Unable to unmarshal payload for deployment id %v", deployment.ID)
		}
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
			Payload:     payload,
		})
	}

	return deployments, nil
}

// CreateDeployment creates a new deployment for the GitHub repo.
func (c *client) CreateDeployment(u *library.User, r *library.Repo, d *library.Deployment) error {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("creating deployment for repo %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	var payload interface{}
	if d.Payload == nil {
		payload = make(map[string]string)
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
