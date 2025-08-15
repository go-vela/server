// SPDX-License-Identifier: Apache-2.0

package github

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

func TestGitHub_cloneRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *http.Request
	}{
		{
			name: "basic request",
			request: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "https",
					Path:   "/",
				},
				Header: http.Header{
					"Accept": []string{"application/json"},
				},
			},
		},
		{
			name: "request with body",
			request: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Scheme: "https",
					Path:   "/",
				},
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"key":"value"}`)),
			},
		},
		{
			name: "request with multiple headers",
			request: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "https",
					Path:   "/",
				},
				Header: http.Header{
					"Accept":        []string{"application/json"},
					"Authorization": []string{"Bearer token"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clonedReq := cloneRequest(tt.request)

			if clonedReq == tt.request {
				t.Errorf("cloneRequest() = %v, want different instance", clonedReq)
			}

			if diff := cmp.Diff(clonedReq.Header, tt.request.Header); diff != "" {
				t.Errorf("cloneRequest() headers mismatch (-want +got):\n%s", diff)
			}

			if clonedReq.Method != tt.request.Method {
				t.Errorf("cloneRequest() method = %v, want %v", clonedReq.Method, tt.request.Method)
			}

			if clonedReq.URL.String() != tt.request.URL.String() {
				t.Errorf("cloneRequest() URL = %v, want %v", clonedReq.URL, tt.request.URL)
			}
		})
	}
}

func TestAppsTransport_RoundTrip(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	_url, _ := url.Parse(s.URL)

	tests := []struct {
		name       string
		transport  *AppsTransport
		request    *http.Request
		wantHeader string
		wantErr    bool
	}{
		{
			name:      "valid GET request",
			transport: NewTestAppsTransport(s.URL),
			request: &http.Request{
				Method: "GET",
				URL:    _url,
				Header: http.Header{
					"Accept": []string{"application/json"},
				},
			},
			wantHeader: "Bearer ",
			wantErr:    false,
		},
		{
			name:      "valid POST request",
			transport: NewTestAppsTransport(s.URL),
			request: &http.Request{
				Method: "POST",
				URL:    _url,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"key":"value"}`)),
			},
			wantHeader: "Bearer ",
			wantErr:    false,
		},
		{
			name:      "request with invalid URL",
			transport: NewTestAppsTransport(s.URL),
			request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "://invalid-url"},
				Header: http.Header{},
			},
			wantHeader: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.transport.RoundTrip(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got := tt.request.Header.Get("Authorization"); !strings.HasPrefix(got, tt.wantHeader) {
					t.Errorf("RoundTrip() Authorization header = %v, want prefix %v", got, tt.wantHeader)
				}
			}

			if resp != nil {
				resp.Body.Close()
			}
		})
	}
}
