// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGithub_InstallRateLimit(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/rate_limit", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/rate_limit.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	wantLimit := 5000
	wantRemaining := 4999
	wantReset := int64(1609459200)

	client, _ := NewTest(s.URL)

	// run test
	gotLimit, gotRemaining, gotReset, err := client.InstallRateLimit(context.TODO(), "foobar", 1)

	if resp.Code != http.StatusOK {
		t.Errorf("InstallRateLimit returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("InstallRateLimit returned err: %v", err)
	}

	if gotLimit != wantLimit {
		t.Errorf("InstallRateLimit limit is %v, want %v", gotLimit, wantLimit)
	}

	if gotRemaining != wantRemaining {
		t.Errorf("InstallRateLimit remaining is %v, want %v", gotRemaining, wantRemaining)
	}

	if gotReset != wantReset {
		t.Errorf("InstallRateLimit reset is %v, want %v", gotReset, wantReset)
	}
}
