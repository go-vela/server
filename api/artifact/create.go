// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
)

func CreateArtifact(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	// capture the artifact from the request body
	input := new(types.Artifact)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new Artifact: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure build_id is defined
	if input.GetBuildID() <= 0 {
		util.HandleError(c, http.StatusBadRequest, fmt.Errorf("build_id must set and greater than 0"))
		return
	}

	input.SetCreatedAt(time.Now().UTC().Unix())

	l.Debugf("creating new artifact")
	// create the artifact in the database using the input from request
	a, err := database.FromContext(c).CreateArtifact(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create new artifact: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, a)
}
