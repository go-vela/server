// SPDX-License-Identifier: Apache-2.0

package github

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/cache/models"
	cacheredis "github.com/go-vela/server/cache/redis"
)

// mockTransport is a configurable http.RoundTripper for testing.
type mockTransport struct {
	handler func(*http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.handler(req)
}

func testCacheKey(rawURL string) string {
	sum := sha256.Sum256([]byte(rawURL))
	return "github:contents:" + hex.EncodeToString(sum[:])
}

func TestContentsTransport_RoundTrip(t *testing.T) {
	testURL := "http://example.com/repos/org/repo/contents/file.yml"
	key := testCacheKey(testURL)

	// setup tests
	tests := []struct {
		name            string
		cachedEntry     *models.TemplateEntry
		baseStatus      int
		baseHeaders     http.Header
		baseBody        string
		wantStatus      int
		wantBody        string
		wantStored      bool
		wantIfNoneMatch string
	}{
		{
			name:       "cache miss with 200 and etag stores response",
			baseStatus: http.StatusOK,
			baseHeaders: http.Header{
				"Etag":         {`"etag-123"`},
				"Content-Type": {"application/json"},
			},
			baseBody:   "template content",
			wantStatus: http.StatusOK,
			wantBody:   "template content",
			wantStored: true,
		},
		{
			name:       "cache miss with 200 and no etag skips store",
			baseStatus: http.StatusOK,
			baseHeaders: http.Header{
				"Content-Type": {"application/json"},
			},
			baseBody:   "template content",
			wantStatus: http.StatusOK,
			wantBody:   "template content",
		},
		{
			name: "cache hit with 304 returns cached response",
			cachedEntry: &models.TemplateEntry{
				ETag:   `"etag-123"`,
				Status: http.StatusOK,
				Header: http.Header{"Content-Type": {"application/json"}},
				Body:   []byte("cached body"),
			},
			baseStatus:      http.StatusNotModified,
			baseHeaders:     http.Header{},
			wantStatus:      http.StatusOK,
			wantBody:        "cached body",
			wantIfNoneMatch: `"etag-123"`,
		},
		{
			name:        "cache miss with 404 passes through",
			baseStatus:  http.StatusNotFound,
			baseHeaders: http.Header{},
			baseBody:    "not found",
			wantStatus:  http.StatusNotFound,
			wantBody:    "not found",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup redis mock
			_redis, err := cacheredis.NewTest("c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90")
			if err != nil {
				t.Errorf("unable to create cache service: %v", err)
			}

			// seed the cache if a cached entry is provided
			if test.cachedEntry != nil {
				err = _redis.StoreTemplateContents(t.Context(), key, test.cachedEntry)
				if err != nil {
					t.Errorf("unable to store template contents: %v", err)
				}
			}

			var gotIfNoneMatch string

			base := &mockTransport{
				handler: func(req *http.Request) (*http.Response, error) {
					gotIfNoneMatch = req.Header.Get("If-None-Match")

					return &http.Response{
						StatusCode: test.baseStatus,
						Header:     test.baseHeaders.Clone(),
						Body:       io.NopCloser(bytes.NewBufferString(test.baseBody)),
					}, nil
				},
			}

			transport := NewContentsTransport(_redis, base)

			req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, testURL, nil)

			resp, err := transport.RoundTrip(req)
			if err != nil {
				t.Errorf("RoundTrip returned err: %v", err)

				return
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode != test.wantStatus {
				t.Errorf("RoundTrip status = %d, want %d", resp.StatusCode, test.wantStatus)
			}

			if diff := cmp.Diff(test.wantBody, string(body)); diff != "" {
				t.Errorf("RoundTrip body mismatch (-want +got):\n%s", diff)
			}

			// verify entry was stored in redis
			got, err := _redis.GetTemplateContents(t.Context(), key)
			if err != nil {
				t.Errorf("GetTemplateContents returned err: %v", err)
			}

			if (got != nil) != test.wantStored && test.cachedEntry == nil {
				t.Errorf("stored = %v, want %v", got != nil, test.wantStored)
			}

			if diff := cmp.Diff(test.wantIfNoneMatch, gotIfNoneMatch); diff != "" {
				t.Errorf("If-None-Match mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestContentsTransport_RoundTrip_NilStore(t *testing.T) {
	wantBody := "passthrough response"

	base := &mockTransport{
		handler: func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(wantBody)),
			}, nil
		},
	}

	transport := NewContentsTransport(nil, base)

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/test", nil)

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Errorf("RoundTrip returned err: %v", err)

		return
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("RoundTrip status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if diff := cmp.Diff(wantBody, string(body)); diff != "" {
		t.Errorf("RoundTrip body mismatch (-want +got):\n%s", diff)
	}
}
