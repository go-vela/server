// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/util"
)

func ListTestAttachmentsForBuild(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()
	b := build.Retrieve(c)

	l.Debugf("listing test attachments for build %d", b.GetNumber())

	// retrieve test attachments from the database
	attachments, err := database.FromContext(c).ListTestAttachmentsByBuildID(ctx, b.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to list test attachments for build %d: %w", b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, attachments)
}
