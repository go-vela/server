// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetDeployment gets a deployment from the GitHub repo.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) GetDeployment(u *library.User, r *library.Repo, id int64) (*library.Deployment, error) {
	logrus.Tracef("capturing deployment %d for %s", id, r.GetFullName())
	return nil, nil
}

// GetDeploymentCount counts a list of deployments from the GitHub repo.
func (c *client) GetDeploymentCount(u *library.User, r *library.Repo) (int64, error) {
	logrus.Tracef("counting deployments for %s", r.GetFullName())
	return 0, nil
}

// GetDeploymentList gets a list of deployments from the GitHub repo.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) GetDeploymentList(u *library.User, r *library.Repo, page, perPage int) ([]*library.Deployment, error) {
	logrus.Tracef("capturing deployments for %s", r.GetFullName())
	return nil, nil
}

// CreateDeployment creates a new deployment for the GitHub repo.
func (c *client) CreateDeployment(u *library.User, r *library.Repo, d *library.Deployment) error {
	logrus.Tracef("creating deployment for %s", r.GetFullName())
	return nil
}
