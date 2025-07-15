package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"net/http"
)

// TestReportResp is the mock response for adding a test report.
const TestReportResp = `{
	"id": 2,
	"build_id": 8,
	"created_at": 1750710551
}`

func addTestReport(c *gin.Context) {
	data := []byte(TestReportResp)
	var body api.TestReport
	_ = json.Unmarshal(data, &body)
	c.JSON(http.StatusCreated, body)
}
