// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"encoding/json"
	"strconv"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/raw"
	"github.com/jenkins-x/go-scm/scm"

	"github.com/sirupsen/logrus"
)

// GetDeployment gets a deployment from the GitHub repo.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) GetDeployment(u *library.User, r *library.Repo, id int64) (*library.Deployment, error) {
	logrus.Tracef("capturing deployment %d for %s", id, r.GetFullName())

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return nil, err
	}

	// send API call to capture the deployment
	deployment, _, err := client.Deployments.Find(ctx, r.GetFullName(), strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}

	payload, err := converPayload(deployment.Payload)
	if err != nil {
		return nil, err
	}

	deploymentID, err := strconv.ParseInt(deployment.ID, 10, 64)
	if err == nil {
		logrus.Tracef("Unable to convert deployment id %s to int64: %v", deployment.ID, err)
	}

	return &library.Deployment{
		ID:          &deploymentID,
		RepoID:      r.ID,
		URL:         &deployment.Link,
		User:        &deployment.Author.Login,
		Commit:      &deployment.Sha,
		Ref:         &deployment.Ref,
		Task:        &deployment.Task,
		Target:      &deployment.Environment,
		Description: &deployment.Description,
		Payload:     payload,
	}, nil
}

// GetDeploymentCount counts a list of deployments from the GitHub repo.
func (c *client) GetDeploymentCount(u *library.User, r *library.Repo) (int64, error) {
	logrus.Tracef("counting deployments for %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return 0, err
	}

	deployments := []*scm.Deployment{}
	opts := scm.ListOptions{
		Size: 100,
	}

	for {
		// send API call to capture the list of deployments
		d, resp, err := client.Deployments.List(ctx, r.GetFullName(), opts)
		if err != nil {
			return 0, err
		}

		deployments = append(deployments, d...)

		// break the loop if there is no more results to page through
		if resp.Page.Next == 0 {
			break
		}

		opts.Page = resp.Page.Next
	}

	return int64(len(deployments)), nil
}

// GetDeploymentList gets a list of deployments from the GitHub repo.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) GetDeploymentList(u *library.User, r *library.Repo, page, perPage int) ([]*library.Deployment, error) {
	logrus.Tracef("capturing deployments for %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return nil, err
	}

	// set pagination options for listing deployments
	opts := scm.ListOptions{
		Size: 100,
	}

	// send API call to capture the list of deployments
	d, _, err := client.Deployments.List(ctx, r.GetFullName(), opts)
	if err != nil {
		return nil, err
	}

	// variable we want to return
	deployments := []*library.Deployment{}

	// iterate through all API results
	for _, deployment := range d {
		payload, err := converPayload(deployment.Payload)
		if err != nil {
			return nil, err
		}

		deploymentID, err := strconv.ParseInt(deployment.ID, 10, 64)
		if err == nil {
			logrus.Tracef("Unable to convert deployment id %s to int64: %v", deployment.ID, err)
		}

		// convert query result to library type
		deployments = append(deployments, &library.Deployment{
			ID:          &deploymentID,
			RepoID:      r.ID,
			URL:         &deployment.Link,
			User:        &deployment.Author.Login,
			Commit:      &deployment.Sha,
			Ref:         &deployment.Ref,
			Task:        &deployment.Task,
			Target:      &deployment.Environment,
			Description: &deployment.Description,
			Payload:     payload,
		})
	}

	return deployments, nil
}

// CreateDeployment creates a new deployment for the GitHub repo.
func (c *client) CreateDeployment(u *library.User, r *library.Repo, d *library.Deployment) error {
	logrus.Tracef("creating deployment for %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return err
	}
	var payload interface{}
	if d.Payload == nil {
		payload = make(map[string]string)
	} else {
		payload = d.Payload
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// create the hook object to make the API call
	deployment := &scm.DeploymentInput{
		Ref:              d.GetRef(),
		Task:             d.GetTask(),
		AutoMerge:        false,
		RequiredContexts: []string{},
		Payload:          string(bytes),
		Environment:      d.GetTarget(),
		Description:      d.GetDescription(),
	}

	// send API call to create the deployment
	deploy, _, err := client.Deployments.Create(ctx, r.GetFullName(), deployment)
	if err != nil {
		return err
	}

	deploymentID, err := strconv.ParseInt(deploy.ID, 10, 64)
	if err == nil {
		logrus.Tracef("Unable to convert deployment id %s to int64: %v", deploy.ID, err)
	}

	d.SetID(deploymentID)
	d.SetRepoID(r.GetID())
	d.SetURL(deploy.Link)
	d.SetUser(deploy.Author.Login)
	d.SetCommit(deploy.Sha)
	d.SetRef(deploy.Ref)
	d.SetTask(deploy.Task)
	d.SetTarget(deploy.Environment)
	d.SetDescription(deploy.Description)

	return nil
}

// helper function to handle converting generic payloads into known StringSliceMap type
func converPayload(i interface{}) (*raw.StringSliceMap, error) {
	var payload *raw.StringSliceMap

	bytes, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &payload)
	if err != nil {
		logrus.Tracef("Unable to unmarshal deployment for payload: %v", payload)
	}

	return payload, nil
}
