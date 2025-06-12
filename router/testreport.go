// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/testreport"
	"github.com/go-vela/server/router/middleware/perm"
)

// TestReportHandlers is a function that extends the provided base router group
// with the API handlers for test report functionality.
//
// POST   /api/v1/...fill this out
func TestReportHandlers(base *gin.RouterGroup) {
	// test report endpoints
	testreports := base.Group("")
	{
		testreports.POST("", perm.MustWrite(), testreport.CreateTestReport)

	} // end of test report endpoints
}
