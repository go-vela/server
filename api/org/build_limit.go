// SPDX-License-Identifier: Apache-2.0

package org

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	oMiddleware "github.com/go-vela/server/router/middleware/org"
	sMiddleware "github.com/go-vela/server/router/middleware/settings"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/repos/{org}/limit org GetBuildLimit
//
// Get the concurrent build limit for an organization
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the org build limit
//     schema:
//       "$ref": "#/definitions/OrgBuildLimit"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildLimit represents the API handler to get the concurrent
// build limit for an organization. When no override has been set
// for the org, the effective default limit is returned.
func GetBuildLimit(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()
	o := oMiddleware.Retrieve(c)
	defaultOrgBuildLimit := c.Value("defaultOrgBuildLimit").(int32)

	l.Debugf("reading build limit for org %s", o)

	limit, err := database.FromContext(c).GetOrgBuildLimit(ctx, o)
	if err != nil {
		// no override set for the org - return the effective default
		if errors.Is(err, gorm.ErrRecordNotFound) {
			limit = new(types.OrgBuildLimit)
			limit.SetOrg(o)
			limit.SetBuildLimit(defaultOrgBuildLimit)

			c.JSON(http.StatusOK, limit)

			return
		}

		retErr := fmt.Errorf("unable to read build limit for org %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, limit)
}

// swagger:operation PUT /api/v1/repos/{org}/limit org UpdateBuildLimit
//
// Update the concurrent build limit for an organization
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: body
//   name: body
//   description: The org build limit to apply
//   required: true
//   schema:
//     "$ref": "#/definitions/OrgBuildLimit"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the org build limit
//     schema:
//       "$ref": "#/definitions/OrgBuildLimit"
//   '201':
//     description: Successfully created the org build limit
//     schema:
//       "$ref": "#/definitions/OrgBuildLimit"
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '403':
//     description: Org build limits are disabled
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateBuildLimit represents the API handler to set the concurrent
// build limit for an organization. The provided value is clamped to
// the platform-configured bounds and the record is created when no
// override yet exists for the org.
func UpdateBuildLimit(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()
	o := oMiddleware.Retrieve(c)
	u := user.Retrieve(c)
	ps := sMiddleware.FromContext(c)
	defaultOrgBuildLimit := c.Value("defaultOrgBuildLimit").(int32)
	maxOrgBuildLimit := c.Value("maxOrgBuildLimit").(int32)

	l.Debugf("updating build limit for org %s", o)

	// org build limits must be enabled by platform admins
	if ps == nil || !ps.GetEnableOrgBuildLimit() {
		retErr := fmt.Errorf("organization build limits have been disabled by Vela admins")

		util.HandleError(c, http.StatusForbidden, retErr)

		return
	}

	// capture body from API request
	input := new(types.OrgBuildLimit)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for org build limit: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// clamp the requested limit to the platform-configured bounds
	limit := input.GetBuildLimit()

	switch {
	case limit <= 0:
		// apply the default when no usable value is provided
		limit = defaultOrgBuildLimit
	case limit > maxOrgBuildLimit:
		limit = maxOrgBuildLimit
	default:
		limit = max(constants.OrgBuildLimitMin, limit)
	}

	// look up any existing override for the org
	existing, err := database.FromContext(c).GetOrgBuildLimit(ctx, o)

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// create a new org build limit
		obl := new(types.OrgBuildLimit)
		obl.SetOrg(o)
		obl.SetBuildLimit(limit)
		obl.SetCreatedAt(time.Now().UTC().Unix())
		obl.SetUpdatedAt(time.Now().UTC().Unix())
		obl.SetUpdatedBy(u.GetName())

		obl, err = database.FromContext(c).CreateOrgBuildLimit(ctx, obl)
		if err != nil {
			retErr := fmt.Errorf("unable to create build limit for org %s: %w", o, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.Infof("build limit for org %s created with value %d", o, limit)

		c.JSON(http.StatusCreated, obl)

	case err != nil:
		retErr := fmt.Errorf("unable to read build limit for org %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

	default:
		// update the existing org build limit
		existing.SetBuildLimit(limit)
		existing.SetUpdatedAt(time.Now().UTC().Unix())
		existing.SetUpdatedBy(u.GetName())

		existing, err = database.FromContext(c).UpdateOrgBuildLimit(ctx, existing)
		if err != nil {
			retErr := fmt.Errorf("unable to update build limit for org %s: %w", o, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.Infof("build limit for org %s updated to %d", o, limit)

		c.JSON(http.StatusOK, existing)
	}
}
