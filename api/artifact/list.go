// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/util"
)

func ListArtifactsForBuild(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()
	b := build.Retrieve(c)

	l.Debugf("listing artifacts for build %d", b.GetNumber())

	// retrieve artifacts from the database
	artifacts, err := database.FromContext(c).ListArtifactsByBuildID(ctx, b.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to list artifacts for build %d: %w", b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, artifacts)
}
