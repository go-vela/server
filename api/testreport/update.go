package testreport

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/testreport"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

func UpdateTestReport(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	tr := testreport.Retrieve(c)
	ctx := c.Request.Context()

	// TODO: update this path if necessary
	entry := fmt.Sprintf("%s/%d/testreport", r.GetFullName(), b.GetNumber())

	l.Debugf("updating test report %s", entry)

	// capture body from API request
	input := new(types.TestReport)

	// update test report fields if provided
	if input.GetBuildID() > 0 {
		// update build ID if set
		tr.SetBuildID(b.GetID())
	}

	if input.GetCreatedAt() > 0 {
		// update created_at if set
		tr.SetCreatedAt(input.GetCreatedAt())
	}

	// send API call to update the test report
	tr, err := database.FromContext(c).UpdateTestReport(ctx, tr)
	if err != nil {
		retErr := fmt.Errorf("unable to update test report %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, tr)
}
