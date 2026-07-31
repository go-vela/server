// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-vela/server/cache/models"
)

func TestRedis_ExtendTemplateExpiry(t *testing.T) {
	// setup redis mock
	_redis, err := NewTest("c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90")
	if err != nil {
		t.Errorf("unable to create cache service: %v", err)
	}

	// store an entry to extend
	entry := &models.TemplateEntry{
		ETag:      `"etag-123"`,
		Status:    http.StatusOK,
		Header:    http.Header{"Content-Type": {"application/json"}},
		Body:      []byte("template body"),
		UpdatedAt: time.Now().UTC(),
	}

	err = _redis.StoreTemplateContents(t.Context(), "github:contents:abc123", entry)
	if err != nil {
		t.Errorf("unable to store template contents: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "existing entry",
			key:     "github:contents:abc123",
			wantErr: false,
		},
		{
			name:    "non-existent entry",
			key:     "github:contents:missing",
			wantErr: false,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := _redis.ExtendTemplateExpiry(t.Context(), test.key)

			if test.wantErr {
				if err == nil {
					t.Errorf("ExtendTemplateExpiry should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("ExtendTemplateExpiry returned err: %v", err)
			}
		})
	}
}
