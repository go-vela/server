// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache/models"
)

func TestRedis_StoreCheckRuns(t *testing.T) {
	// setup types
	checkRuns := []models.CheckRun{
		{
			ID:          1,
			Context:     "continuous-integration/vela/push",
			Repo:        "octocat/hello-world",
			BuildNumber: 1,
		},
	}

	repo := new(api.Repo)
	repo.SetInstallID(1)
	repo.SetApprovalTimeout(7)
	repo.SetTimeout(30)

	// setup redis mock
	_redis, err := NewTest("123abc")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		redis    *Client
		buildID  int64
		checkRun []models.CheckRun
	}{
		{
			failure:  false,
			redis:    _redis,
			buildID:  1,
			checkRun: checkRuns,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.redis.StoreCheckRuns(t.Context(), test.buildID, test.checkRun, repo)

		if test.failure {
			if err == nil {
				t.Errorf("StoreCheckRuns should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("StoreCheckRuns returned err: %v", err)
		}
	}
}

func TestRedis_StoreStepCheckRuns(t *testing.T) {
	// setup types
	checkRuns := []models.CheckRun{
		{
			ID:          1,
			Context:     "continuous-integration/vela/push/test",
			Repo:        "octocat/hello-world",
			BuildNumber: 1,
		},
	}

	repo := new(api.Repo)
	repo.SetInstallID(1)
	repo.SetApprovalTimeout(7)
	repo.SetTimeout(30)

	// setup redis mock
	_redis, err := NewTest("123abc")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		redis    *Client
		stepID   int64
		checkRun []models.CheckRun
	}{
		{
			failure:  false,
			redis:    _redis,
			stepID:   1,
			checkRun: checkRuns,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.redis.StoreStepCheckRuns(t.Context(), test.stepID, test.checkRun, repo)

		if test.failure {
			if err == nil {
				t.Errorf("StoreStepCheckRuns should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("StoreStepCheckRuns returned err: %v", err)
		}
	}
}
