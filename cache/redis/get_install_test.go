// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"testing"
)

func TestRedis_GetInstallToken(t *testing.T) {
	_redis, err := NewTest("c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	err = _redis.StoreInstallToken(t.Context(), "test_token", 1, 30)
	if err != nil {
		t.Errorf("unable to store install token: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "existing token",
			token:   "test_token",
			wantErr: false,
		},
		{
			name:    "non-existent token",
			token:   "missing_token",
			wantErr: true,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := _redis.GetInstallToken(t.Context(), test.token)

			if test.wantErr {
				if err == nil {
					t.Errorf("GetInstallToken should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("GetInstallToken returned err: %v", err)
			}
		})
	}
}

func TestRedis_GetInstallStatusToken(t *testing.T) {
	// setup redis mock
	_redis, err := NewTest("c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90")
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
