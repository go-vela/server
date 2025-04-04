// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// DeploymentInterface represents the Vela interface for deployment
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type DeploymentInterface interface {
	// Deployment Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateDeploymentIndexes defines a function that creates the indexes for the deployment table.
	CreateDeploymentIndexes(context.Context) error
	// CreateDeploymentTable defines a function that creates the deployment table.
	CreateDeploymentTable(context.Context, string) error

	// Deployment Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountDeployments defines a function that gets the count of all deployments.
	CountDeployments(context.Context) (int64, error)
	// CountDeploymentsForRepo defines a function that gets the count of deployments by repo ID.
	CountDeploymentsForRepo(context.Context, *api.Repo) (int64, error)
	// CreateDeployment defines a function that creates a new deployment.
	CreateDeployment(context.Context, *api.Deployment) (*api.Deployment, error)
	// DeleteDeployment defines a function that deletes an existing deployment.
	DeleteDeployment(context.Context, *api.Deployment) error
	// GetDeployment defines a function that gets a deployment by ID.
	GetDeployment(context.Context, int64) (*api.Deployment, error)
	// GetDeploymentForRepo defines a function that gets a deployment by repo ID and number.
	GetDeploymentForRepo(context.Context, *api.Repo, int64) (*api.Deployment, error)
	// ListDeployments defines a function that gets a list of all deployments.
	ListDeployments(context.Context) ([]*api.Deployment, error)
	// ListDeploymentsForRepo defines a function that gets a list of deployments by repo ID.
	ListDeploymentsForRepo(context.Context, *api.Repo, int, int) ([]*api.Deployment, error)
	// UpdateDeployment defines a function that updates an existing deployment.
	UpdateDeployment(context.Context, *api.Deployment) (*api.Deployment, error)
}
