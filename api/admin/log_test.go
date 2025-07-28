// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
)

func TestAdmin_CleanLogs(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	// setup mock database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	// setup request body
	body := types.Error{
		Message: String("Test log cleanup"),
	}
	data, _ := json.Marshal(body)

	cutoffTime := time.Now().Add(-48 * time.Hour).Unix()

	// setup tests
	tests := []struct {
		name        string
		queryParams string
		body        *bytes.Buffer
		wantStatus  int
		wantError   bool
	}{
		{
			name:        "successful cleanup",
			queryParams: fmt.Sprintf("before=%d&batch_size=1000&vacuum=true", cutoffTime),
			body:        bytes.NewBuffer(data),
			wantStatus:  http.StatusOK,
			wantError:   false,
		},
		{
			name:        "missing before parameter",
			queryParams: "batch_size=1000",
			body:        bytes.NewBuffer(data),
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
		{
			name:        "invalid before parameter",
			queryParams: "before=invalid",
			body:        bytes.NewBuffer(data),
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
		{
			name:        "before timestamp too recent",
			queryParams: fmt.Sprintf("before=%d", time.Now().Unix()),
			body:        bytes.NewBuffer(data),
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
		{
			name:        "invalid batch_size parameter",
			queryParams: fmt.Sprintf("before=%d&batch_size=invalid", cutoffTime),
			body:        bytes.NewBuffer(data),
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
		{
			name:        "batch_size too large",
			queryParams: fmt.Sprintf("before=%d&batch_size=20000", cutoffTime),
			body:        bytes.NewBuffer(data),
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
		{
			name:        "invalid vacuum parameter",
			queryParams: fmt.Sprintf("before=%d&vacuum=invalid", cutoffTime),
			body:        bytes.NewBuffer(data),
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup context
			resp := httptest.NewRecorder()
			context, engine := gin.CreateTestContext(resp)

			// setup request
			context.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodDelete, fmt.Sprintf("/api/v1/admin/log/cleanup?%s", test.queryParams), test.body)

			// setup vela mock server
			engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.New())) })
			engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
			engine.DELETE("/api/v1/admin/log/cleanup", CleanLogs)

			// run test
			engine.ServeHTTP(context.Writer, context.Request)

			if resp.Code != test.wantStatus {
				t.Errorf("CleanLogs returned %v, want %v", resp.Code, test.wantStatus)
			}

			if !test.wantError && resp.Code == http.StatusOK {
				// verify response structure for successful requests
				var response types.LogCleanupResponse
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("unable to unmarshal response: %v", err)
				}

				// basic validation of response fields
				if response.DeletedCount < 0 {
					t.Errorf("DeletedCount should be >= 0, got %d", response.DeletedCount)
				}
				if response.BatchesProcessed < 0 {
					t.Errorf("BatchesProcessed should be >= 0, got %d", response.BatchesProcessed)
				}
				if response.DurationSeconds < 0 {
					t.Errorf("DurationSeconds should be >= 0, got %f", response.DurationSeconds)
				}
				if response.Message == "" {
					t.Errorf("Message should not be empty")
				}
			}
		})
	}
}

func TestAdmin_CleanLogs_ValidationEdgeCases(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	// setup mock database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	cutoffTime := time.Now().Add(-48 * time.Hour).Unix()

	// test edge cases for validation
	tests := []struct {
		name        string
		queryParams string
		wantStatus  int
	}{
		{
			name:        "minimum valid batch_size",
			queryParams: fmt.Sprintf("before=%d&batch_size=1", cutoffTime),
			wantStatus:  http.StatusOK,
		},
		{
			name:        "maximum valid batch_size",
			queryParams: fmt.Sprintf("before=%d&batch_size=10000", cutoffTime),
			wantStatus:  http.StatusOK,
		},
		{
			name:        "batch_size zero",
			queryParams: fmt.Sprintf("before=%d&batch_size=0", cutoffTime),
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "vacuum true",
			queryParams: fmt.Sprintf("before=%d&vacuum=true", cutoffTime),
			wantStatus:  http.StatusOK,
		},
		{
			name:        "vacuum false",
			queryParams: fmt.Sprintf("before=%d&vacuum=false", cutoffTime),
			wantStatus:  http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup context
			resp := httptest.NewRecorder()
			context, engine := gin.CreateTestContext(resp)

			// setup request
			context.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodDelete, fmt.Sprintf("/api/v1/admin/log/cleanup?%s", test.queryParams), nil)

			// setup vela mock server
			engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.New())) })
			engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
			engine.DELETE("/api/v1/admin/log/cleanup", CleanLogs)

			// run test
			engine.ServeHTTP(context.Writer, context.Request)

			if resp.Code != test.wantStatus {
				t.Errorf("CleanLogs returned %v, want %v", resp.Code, test.wantStatus)
			}
		})
	}
}

// String is a helper routine that returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}
