// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/cache/models"
)

func TestRedis_GetInstallToken(t *testing.T) {
	// setup types
	installToken := &models.InstallToken{
		Token:        "test_token",
		InstallID:    1,
		Repositories: []string{"octocat/hello-world"},
		Permissions: map[string]string{
			"contents": "read",
		},
		Timeout: 30,
	}

	_redis, err := NewTest("123abc")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	err = _redis.StoreInstallToken(t.Context(), installToken, 1, 30)
	if err != nil {
		t.Errorf("unable to store install token: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		token   string
		want    *models.InstallToken
		wantErr bool
	}{
		{
			name:    "existing token",
			token:   "test_token",
			want:    installToken,
			wantErr: false,
		},
		{
			name:    "non-existent token",
			token:   "missing_token",
			want:    nil,
			wantErr: true,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := _redis.GetInstallToken(t.Context(), test.token)

			if test.wantErr {
				if err == nil {
					t.Errorf("GetInstallToken should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("GetInstallToken returned err: %v", err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetInstallToken() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRedis_GetInstallStatusToken(t *testing.T) {
	// setup redis mock
	_redis, err := NewTest("123abc")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// store a status token
	err = _redis.StoreInstallStatusToken(t.Context(), 1, "ghs_status_token_123")
	if err != nil {
		t.Errorf("unable to store install status token: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		build   int64
		want    string
		wantErr bool
	}{
		{
			name:    "existing status token",
			build:   1,
			want:    "ghs_status_token_123",
			wantErr: false,
		},
		{
			name:    "non-existent status token",
			build:   999,
			want:    "",
			wantErr: true,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := _redis.GetInstallStatusToken(t.Context(), test.build)

			if test.wantErr {
				if err == nil {
					t.Errorf("GetInstallStatusToken should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("GetInstallStatusToken returned err: %v", err)
			}

			if got != test.want {
				t.Errorf("GetInstallStatusToken() = %v, want %v", got, test.want)
			}
		})
	}
}
