// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types/pipeline"
	"github.com/google/go-cmp/cmp"
)

// TestPrepareRuleData tests the prepareRuleData function.
func TestPrepareRuleData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		parameters map[string]string
		want       *pipeline.RuleData
	}{
		{
			name: "all params provided",
			parameters: map[string]string{
				"branch":  "main",
				"comment": "Test comment",
				"event":   "push",
				"repo":    "my-repo",
				"status":  "success",
				"tag":     "v1.0.0",
				"target":  "production",
				"path":    "README.md",
			},
			want: &pipeline.RuleData{
				Branch:  "main",
				Comment: "Test comment",
				Event:   "push",
				Repo:    "my-repo",
				Status:  "success",
				Tag:     "v1.0.0",
				Target:  "production",
				Path:    []string{"README.md"},
			},
		},
		{
			name: "multiple path params",
			parameters: map[string]string{
				"path":   "README.md",
				"branch": "main",
			},
			want: &pipeline.RuleData{
				Branch: "main",
				Path:   []string{"README.md", "src/main.go"},
			},
		},
		{
			name:       "no params provided",
			parameters: map[string]string{},
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)

			q := req.URL.Query()
			for key, value := range tt.parameters {
				q.Add(key, value)
			}

			// add additional path parameter for multiple path test
			if strings.EqualFold(tt.name, "multiple path params") {
				q.Add("path", "src/main.go")
			}

			req.URL.RawQuery = q.Encode()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			got := prepareRuleData(ctx)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("prepareRuleData() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
