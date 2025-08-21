// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/testattachment"
)

// TestAttachmentHandlers is a function that extends the provided base router group
// with the API handlers for test attachment functionality.
//
// POST   /api/v1/repos/:org/:repo/builds/:build/reports/testattachments.
func TestAttachmentHandlers(base *gin.RouterGroup) {
	// test attachment endpoints
	_testattachment := base.Group("/reports/testattachment")
	{
		_testattachment.PUT("", testattachment.CreateTestAttachment)
	} // end of test attachment endpoints
}
