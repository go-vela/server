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
		err := _redis.StoreInstallToken(t.Context(), test.token, _repo)
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
