package testreport

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

func CreateTestReport(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)

	// capture the test report from the request body
	input := new(types.TestReport)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new testreport: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	l.Debugf("creating new test report")

	tr := new(types.TestReport)

	// update fields in test report object
	tr.SetBuildID((input.GetBuildID()))
	tr.SetCreatedAt(time.Now().UTC().Unix())

	// create the test report in the database
	tr, err = database.FromContext(c).CreateTestReport(c, tr)
	if err != nil {

		retErr := fmt.Errorf("unable to create new test report: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, tr)
}
