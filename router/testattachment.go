// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/testattachment"
	tamiddleware "github.com/go-vela/server/router/middleware/testattachment"
)

// TestAttachmentHandlers is a function that extends the provided base router group
// with the API handlers for test attachment functionality.
//
// GET    /api/v1/repos/:org/:repo/builds/:build/reports/testattachment
// GET    /api/v1/repos/:org/:repo/builds/:build/reports/testattachment/:attachment
// PUT    /api/v1/repos/:org/:repo/builds/:build/reports/testattachment
func TestAttachmentHandlers(base *gin.RouterGroup) {
	// test attachment endpoints
	_testattachment := base.Group("/reports/testattachment")
	{
		_testattachment.GET("", testattachment.ListTestAttachmentsForBuild)
		_testattachment.PUT("", testattachment.CreateTestAttachment)

		// Individual test attachment endpoints
		ta := _testattachment.Group("/:attachment", tamiddleware.Establish())
		{
			ta.GET("", testattachment.GetTestAttachment)
		} // end of individual test attachment endpoints
	} // end of test attachment endpoints
}
