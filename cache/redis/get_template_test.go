// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/cache/models"
)

func TestRedis_GetTemplateContents(t *testing.T) {
	// setup types
	entry := &models.TemplateEntry{
		ETag:      `"etag-123"`,
		Status:    http.StatusOK,
		Header:    http.Header{"Content-Type": {"application/json"}},
		Body:      []byte("template body"),
		UpdatedAt: time.Now().UTC(),
	}

	// setup redis mock
	_redis, err := NewTest("c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90")
	if err != nil {
		t.Errorf("unable to create cache service: %v", err)
	}

	err = _redis.StoreTemplateContents(t.Context(), "github:contents:abc123", entry)
	if err != nil {
		t.Errorf("unable to store template contents: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		key     string
		want    *models.TemplateEntry
		wantErr bool
	}{
		{
			name:    "existing entry",
			key:     "github:contents:abc123",
			want:    entry,
			wantErr: false,
		},
		{
			name:    "non-existent entry",
			key:     "github:contents:missing",
			want:    nil,
			wantErr: false,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := _redis.GetTemplateContents(t.Context(), test.key)

			if test.wantErr {
				if err == nil {
					t.Errorf("GetTemplateContents should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("GetTemplateContents returned err: %v", err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetTemplateContents() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
