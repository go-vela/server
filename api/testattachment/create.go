package testattachment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

func CreateTestAttachment(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)

	// capture the test attachment from the request body
	input := new(types.TestAttachment)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new TestAttachment: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	l.Debugf("creating new test attachment")

	// TODO this isn't working yet, but will compile
	ta := new(types.TestAttachment)

	// create the test attachment in the database
	ta, err = database.FromContext(c).CreateTestAttachment(c, ta)
	if err != nil {

		retErr := fmt.Errorf("unable to create new test report: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, ta)
}
