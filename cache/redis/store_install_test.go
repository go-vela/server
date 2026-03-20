// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

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

	// setup redis mock
	_redis, err := NewTest("c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90")
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
		err := test.redis.StoreInstallToken(t.Context(), test.token, 1, 30)

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

func TestRedis_StoreInstallStatusToken(t *testing.T) {
	// setup redis mock
	_redis, err := NewTest("c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90")
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
			name:    "valid status token",
			build:   1,
			token:   "ghs_status_token_123",
			wantErr: false,
		},
		{
			name:    "empty status token",
			build:   2,
			token:   "",
			wantErr: false,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := _redis.StoreInstallStatusToken(t.Context(), test.build, test.token)

			if test.wantErr {
				if err == nil {
					t.Errorf("StoreInstallStatusToken should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("StoreInstallStatusToken returned err: %v", err)
			}
		})
	}
}
func TestRedis_StoreInstallToken_TTL(t *testing.T) {
	// setup types
	_item := &models.InstallToken{
		Token:        "ttl_test_token",
		Repositories: []string{"octocat/hello-world"},
		Permissions:  map[string]string{"contents": "read"},
	}

	// setup redis mock
	_redis, err := NewTest("c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// store token with a 30-minute timeout
	err = _redis.StoreInstallToken(t.Context(), _item, 1, 30)
	if err != nil {
		t.Fatalf("StoreInstallToken returned err: %v", err)
	}

	// verify the token key has a TTL set
	tokenTTL := _redis.Redis.TTL(t.Context(), tokenKey(t, "c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90", "ttl_test_token"))
	if tokenTTL.Err() != nil {
		t.Fatalf("unable to get TTL for token key: %v", tokenTTL.Err())
	}

	if tokenTTL.Val() <= 0 {
		t.Errorf("token key TTL = %v, want > 0", tokenTTL.Val())
	}

	if tokenTTL.Val() > 30*time.Minute {
		t.Errorf("token key TTL = %v, want <= 30m", tokenTTL.Val())
	}

	// verify the index key has a TTL set
	idxTTL := _redis.Redis.TTL(t.Context(), "idx:build:1")
	if idxTTL.Err() != nil {
		t.Fatalf("unable to get TTL for index key: %v", idxTTL.Err())
	}

	if idxTTL.Val() <= 0 {
		t.Errorf("index key TTL = %v, want > 0", idxTTL.Val())
	}

	// verify the index key contains the token key
	members := _redis.Redis.SMembers(t.Context(), "idx:build:1")
	if members.Err() != nil {
		t.Fatalf("unable to get members for index key: %v", members.Err())
	}

	if len(members.Val()) != 1 {
		t.Errorf("index key members = %d, want 1", len(members.Val()))
	}
}

func TestRedis_StoreInstallStatusToken_TTL(t *testing.T) {
	// setup redis mock
	_redis, err := NewTest("c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// store a status token
	err = _redis.StoreInstallStatusToken(t.Context(), 1, "ghs_ttl_test")
	if err != nil {
		t.Fatalf("StoreInstallStatusToken returned err: %v", err)
	}

	// verify the status token key has a TTL set
	statusTTL := _redis.Redis.TTL(t.Context(), "install_status_token:1")
	if statusTTL.Err() != nil {
		t.Fatalf("unable to get TTL for status token key: %v", statusTTL.Err())
	}

	if statusTTL.Val() <= 0 {
		t.Errorf("status token key TTL = %v, want > 0", statusTTL.Val())
	}

	if statusTTL.Val() > 59*time.Minute {
		t.Errorf("status token key TTL = %v, want <= 59m", statusTTL.Val())
	}

	// verify we can retrieve the token
	got, err := _redis.GetInstallStatusToken(t.Context(), 1)
	if err != nil {
		t.Fatalf("GetInstallStatusToken returned err: %v", err)
	}

	if got != "ghs_ttl_test" {
		t.Errorf("GetInstallStatusToken() = %v, want ghs_ttl_test", got)
	}
}

// tokenKey is a helper to compute the Redis key for an install token,
// mirroring the HMAC logic in StoreInstallToken.
func tokenKey(t *testing.T, installTokenKey, token string) string {
	t.Helper()

	keyBytes, err := hex.DecodeString(installTokenKey)
	if err != nil {
		t.Fatalf("invalid install token key: %v", err)
	}

	h := hmac.New(sha256.New, keyBytes)
	h.Write([]byte(token))

	return fmt.Sprintf("install_token:%s", hex.EncodeToString(h.Sum(nil)))
}
