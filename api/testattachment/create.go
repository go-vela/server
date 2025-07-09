package testattachment

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

func CreateTestAttachment(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	// capture the test attachment from the request body
	input := new(types.TestAttachment)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new TestAttachment: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure test_report_id is defined
	if input.GetTestReportID() <= 0 {
		util.HandleError(c, http.StatusBadRequest, fmt.Errorf("test_report_id must set and greater than 0"))
		return
	}

	input.SetCreatedAt(time.Now().UTC().Unix())

	l.Debugf("creating new test attachment")

	// create the test attachment in the database using the input from request
	ta, err := database.FromContext(c).CreateTestAttachment(ctx, input)

	if err != nil {

		retErr := fmt.Errorf("unable to create new test attachment: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, ta)
}
