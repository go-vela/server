// SPDX-License-Identifier: Apache-2.0

package limits

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetOrgBuildLimit gets an org build limit by org from the database.
func (e *Engine) GetOrgBuildLimit(ctx context.Context, org string) (*api.OrgBuildLimit, error) {
	e.logger.Tracef("getting org build limit for %s", org)

	o := new(types.OrgBuildLimit)

	err := e.client.
		WithContext(ctx).
		Table(constants.TableOrgBuildLimit).
		Where("org = ?", org).
		Take(o).
		Error
	if err != nil {
		return nil, err
	}

	return o.ToAPI(), nil
}
