// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/v75/github"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/raw"
)

// GetDeployment gets a deployment from the GitHub repo.
func (c *Client) GetDeployment(ctx context.Context, u *api.User, r *api.Repo, id int64) (*api.Deployment, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("capturing deployment %d for repo %s", id, r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

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

	createdAt := deployment.CreatedAt.Unix()

	return &api.Deployment{
		ID:          deployment.ID,
		Repo:        r,
		URL:         deployment.URL,
		Commit:      deployment.SHA,
		Ref:         deployment.Ref,
		Task:        deployment.Task,
		Target:      deployment.Environment,
		Description: deployment.Description,
		Payload:     payload,
		CreatedAt:   &createdAt,
		CreatedBy:   deployment.Creator.Login,
	}, nil
}

// GetDeploymentCount counts a list of deployments from the GitHub repo.
func (c *Client) GetDeploymentCount(ctx context.Context, u *api.User, r *api.Repo) (int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("counting deployments for repo %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)
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
func (c *Client) GetDeploymentList(ctx context.Context, u *api.User, r *api.Repo, page, perPage int) ([]*api.Deployment, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("listing deployments for repo %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

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
	deployments := []*api.Deployment{}

	// iterate through all API results
	for _, deployment := range d {
		var payload *raw.StringSliceMap

		err := json.Unmarshal(deployment.Payload, &payload)
		if err != nil {
			c.Logger.Tracef("Unable to unmarshal payload for deployment id %v", deployment.ID)
		}

		createdAt := deployment.CreatedAt.Unix()

		// convert query result to API type
		deployments = append(deployments, &api.Deployment{
			ID:          deployment.ID,
			Repo:        r,
			URL:         deployment.URL,
			Commit:      deployment.SHA,
			Ref:         deployment.Ref,
			Task:        deployment.Task,
			Target:      deployment.Environment,
			Description: deployment.Description,
			Payload:     payload,
			CreatedAt:   &createdAt,
			CreatedBy:   deployment.Creator.Login,
		})
	}

	return deployments, nil
}

// CreateDeployment creates a new deployment for the GitHub repo.
func (c *Client) CreateDeployment(ctx context.Context, u *api.User, r *api.Repo, d *api.Deployment) error {
	c.Logger.WithFields(logrus.Fields{
		"org":     r.GetOrg(),
		"repo":    r.GetName(),
		"user":    u.GetName(),
		"user_id": u.GetID(),
	}).Tracef("creating deployment for repo %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

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
		AutoMerge:        github.Ptr(false),
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

	d.SetNumber(deploy.GetID())
	d.SetRepo(r)
	d.SetURL(deploy.GetURL())
	d.SetCommit(deploy.GetSHA())
	d.SetRef(deploy.GetRef())
	d.SetTask(deploy.GetTask())
	d.SetTarget(deploy.GetEnvironment())
	d.SetDescription(deploy.GetDescription())

	return nil
}
