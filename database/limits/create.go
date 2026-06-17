// SPDX-License-Identifier: Apache-2.0

package limits

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateOrgBuildLimit creates a new org build limit in the database.
func (e *Engine) CreateOrgBuildLimit(ctx context.Context, o *api.OrgBuildLimit) (*api.OrgBuildLimit, error) {
	e.logger.WithFields(logrus.Fields{
		"org": o.GetOrg(),
	}).Tracef("creating org build limit for %s", o.GetOrg())

	orgBuildLimit := types.OrgBuildLimitFromAPI(o)

	err := orgBuildLimit.Validate()
	if err != nil {
		return nil, err
	}

	result := e.client.
		WithContext(ctx).
		Table(constants.TableOrgBuildLimit).
		Create(orgBuildLimit)

	return orgBuildLimit.ToAPI(), result.Error
}
