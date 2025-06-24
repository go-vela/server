// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/testreport"
)

// TestReportHandlers is a function that extends the provided base router group
// with the API handlers for test report functionality.
//
// POST   /api/v1/repos/:org/:repo/builds/:build/reports/testreport
func TestReportHandlers(base *gin.RouterGroup) {
	// test report endpoints
	_testreport := base.Group("/reports/testreport")
	{
		_testreport.POST("", testreport.CreateTestReport)

	} // end of test report endpoints
}
