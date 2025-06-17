// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/testattachment"
)

// TestAttachmentHandlers is a function that extends the provided base router group
// with the API handlers for test attachment functionality.
//
// POST   /api/v1/reports/testreport/:build_id/attachments/:attachment_id
// TODO: this will be modified to follow a pattern similar to test report
func TestAttachmentHandlers(base *gin.RouterGroup) {
	// test attachment endpoints
	_testattachment := base.Group("/reports/testreport/:testreport_id")
	{
		_testattachment.POST("", testattachment.CreateTestAttachment)

	} // end of test attachment endpoints
}
