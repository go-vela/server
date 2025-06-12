// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/testattachment"
	"github.com/go-vela/server/router/middleware/perm"
)

// TestAttachmentHandlers is a function that extends the provided base router group
// with the API handlers for test attachment functionality.
//
// POST   /api/v1/...fill this out
func TestAttachmentHandlers(base *gin.RouterGroup) {
	// test attachment endpoints
	testattachments := base.Group("")
	{
		testattachments.POST("", perm.MustWrite(), testattachment.CreateTestAttachment)

	} // end of test attachment endpoints
}
