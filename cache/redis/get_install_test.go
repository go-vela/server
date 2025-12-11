// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache/models"
)

func TestRedis_Pop(t *testing.T) {
	// setup types
	installToken := &models.InstallToken{
		Token:        "test_token",
		Repositories: []string{"octocat/hello-world"},
		Permissions: map[string]string{
			"contents": "read",
		},
	}

	repo := new(api.Repo)
	repo.SetTimeout(30)

	_redis, err := NewTest("123abc")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	err = _redis.StoreInstallToken(t.Context(), installToken, repo)
	if err != nil {
		t.Errorf("unable to store install token: %v", err)
	}

	// setup tests
	tests := []struct {
		wantErr bool
		want    *models.InstallToken
	}{
		{
			wantErr: false,
			want:    installToken,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _redis.GetInstallToken(t.Context(), "test_token")

		if test.wantErr {
			if err == nil {
				t.Errorf("GetInstallToken should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Pop returned err: %v", err)
		}

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("GetInstallToken() mismatch (-want +got):\n%s", diff)
		}
	}
}
