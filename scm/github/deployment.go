// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/library"
	"github.com/google/go-github/v56/github"
)

// CreateDeployment creates a new deployment for the GitHub repo.
func (c *client) CreateDeployment(ctx context.Context, u *library.User, r *library.Repo, d *library.Deployment) error {
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

	d.SetNumber(deploy.GetID())
	d.SetRepoID(r.GetID())
	d.SetURL(deploy.GetURL())
	d.SetCommit(deploy.GetSHA())
	d.SetRef(deploy.GetRef())
	d.SetTask(deploy.GetTask())
	d.SetTarget(deploy.GetEnvironment())
	d.SetDescription(deploy.GetDescription())

	return nil
}
