// SPDX-License-Identifier: Apache-2.0

package testreport

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
)

func CreateTestReport(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	l.Debugf("creating new test report for build %s", entry)

	// capture the test report from the request body
	input := new(types.TestReport)

	input.SetBuildID(b.GetID())
	input.SetCreatedAt(time.Now().UTC().Unix())

	// create the test report in the database
	tr, err := database.FromContext(c).CreateTestReport(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create new test report: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, tr)
}
