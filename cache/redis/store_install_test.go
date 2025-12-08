// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache/models"
)

func TestRedis_StoreInstall(t *testing.T) {
	// setup types
	// use global variables in redis_test.go
	_item := &models.InstallToken{
		Token:        "test_token",
		Repositories: []string{"octocat/hello-world"},
		Permissions: map[string]string{
			"contents": "read",
		},
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
		failure bool
		redis   *Client
		token   *models.InstallToken
	}{
		{
			failure: false,
			redis:   _redis,
			token:   _item,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.redis.StoreInstallToken(t.Context(), test.token, _repo)

		if test.failure {
			if err == nil {
				t.Errorf("StoreInstallToken should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("StoreInstallToken returned err: %v", err)
		}
	}
}
