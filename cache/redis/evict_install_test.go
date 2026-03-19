// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache/models"
)

func TestRedis_EvictInstall(t *testing.T) {
	// setup types
	// use global variables in redis_test.go
	_item := &models.InstallToken{
		Token:        "ghs_123abc",
		Repositories: []string{"vela"},
		Permissions:  map[string]string{"contents": "read"},
	}

	_repo := new(api.Repo)
	_repo.SetTimeout(30)

	// setup redis mock
	_redis, err := NewTest("installKey")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// setup tests
	tests := []struct {
		token   *models.InstallToken
		evict   string
		wantErr bool
	}{
		{
			token:   _item,
			evict:   "ghs_123abc",
			wantErr: false,
		},
		{
			token:   _item,
			evict:   "not_found",
			wantErr: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := _redis.StoreInstallToken(t.Context(), test.token, 1, 30)
		if err != nil {
			t.Errorf("unable to store install token in cache: %v", err)
		}

		err = _redis.EvictInstallToken(t.Context(), test.evict)
		if test.wantErr && err == nil {
			t.Errorf("EvictInstallToken for %s returned err nil, want err", test.evict)
		}

		if !test.wantErr && err != nil {
			t.Errorf("EvictInstallToken for %s returned err %v, want nil", test.evict, err)
		}
	}
}

func TestRedis_EvictBuildInstallTokens(t *testing.T) {
	// setup types
	_item := &models.InstallToken{
		Token:        "ghs_123abc",
		Repositories: []string{"vela"},
		Permissions:  map[string]string{"contents": "read"},
	}

	_item2 := &models.InstallToken{
		Token:        "ghs_456def",
		Repositories: []string{"vela", "other-repo"},
		Permissions:  map[string]string{"contents": "read"},
	}

	// setup redis mock
	_redis, err := NewTest("installKey")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		build   int64
		tokens  []*models.InstallToken
		wantErr bool
	}{
		{
			name:    "evict single token for build",
			build:   1,
			tokens:  []*models.InstallToken{_item},
			wantErr: false,
		},
		{
			name:    "evict multiple tokens for build",
			build:   2,
			tokens:  []*models.InstallToken{_item, _item2},
			wantErr: false,
		},
		{
			name:    "evict tokens for build with no tokens",
			build:   999,
			tokens:  nil,
			wantErr: false,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// store tokens for the build if provided
			for _, token := range test.tokens {
				err := _redis.StoreInstallToken(t.Context(), token, test.build, 30)
				if err != nil {
					t.Errorf("unable to store install token in cache: %v", err)
				}
			}

			err := _redis.EvictBuildInstallTokens(t.Context(), test.build)
			if test.wantErr && err == nil {
				t.Errorf("EvictBuildInstallTokens for build %d returned err nil, want err", test.build)
			}

			if !test.wantErr && err != nil {
				t.Errorf("EvictBuildInstallTokens for build %d returned err %v, want nil", test.build, err)
			}
		})
	}
}

func TestRedis_EvictInstallStatusToken(t *testing.T) {
	// setup redis mock
	_redis, err := NewTest("installKey")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// setup tests
	tests := []struct {
		name    string
		build   int64
		token   string
		wantErr bool
	}{
		{
			name:    "evict existing status token",
			build:   1,
			token:   "ghs_status_token_123",
			wantErr: false,
		},
		{
			name:    "evict non-existent status token",
			build:   999,
			token:   "",
			wantErr: false,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// store a status token if provided
			if test.token != "" {
				err := _redis.StoreInstallStatusToken(t.Context(), test.build, test.token)
				if err != nil {
					t.Errorf("unable to store install status token in cache: %v", err)
				}
			}

			err := _redis.EvictInstallStatusToken(t.Context(), test.build)
			if test.wantErr && err == nil {
				t.Errorf("EvictInstallStatusToken for build %d returned err nil, want err", test.build)
			}

			if !test.wantErr && err != nil {
				t.Errorf("EvictInstallStatusToken for build %d returned err %v, want nil", test.build, err)
			}
		})
	}
}
