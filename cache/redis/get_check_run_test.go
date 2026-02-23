// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache/models"
)

func TestRedis_GetCheckRuns(t *testing.T) {
	// setup types
	build := new(api.Build)
	build.SetID(1)

	repo := new(api.Repo)
	repo.SetInstallID(1)
	repo.SetApprovalTimeout(7)
	repo.SetTimeout(30)

	checkRuns := []models.CheckRun{
		{
			ID:          1,
			Context:     "continuous-integration/vela/push",
			Repo:        "octocat/hello-world",
			BuildNumber: 1,
		},
	}

	_redis, err := NewTest("123abc")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	err = _redis.StoreCheckRuns(t.Context(), 1, checkRuns, repo)
	if err != nil {
		t.Errorf("unable to store check runs: %v", err)
	}

	// setup tests
	tests := []struct {
		wantErr bool
		want    []models.CheckRun
	}{
		{
			wantErr: false,
			want:    checkRuns,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _redis.GetCheckRuns(t.Context(), build)

		if test.wantErr {
			if err == nil {
				t.Errorf("GetCheckRuns should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetCheckRuns returned err: %v", err)
		}

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("GetCheckRuns() mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestRedis_GetStepCheckRuns(t *testing.T) {
	// setup types
	step := new(api.Step)
	step.SetID(1)

	repo := new(api.Repo)
	repo.SetInstallID(1)
	repo.SetApprovalTimeout(7)
	repo.SetTimeout(30)

	checkRuns := []models.CheckRun{
		{
			ID:          1,
			Context:     "continuous-integration/vela/push/test",
			Repo:        "octocat/hello-world",
			BuildNumber: 1,
		},
	}

	_redis, err := NewTest("123abc")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	err = _redis.StoreStepCheckRuns(t.Context(), 1, checkRuns, repo)
	if err != nil {
		t.Errorf("unable to store step check runs: %v", err)
	}

	// setup tests
	tests := []struct {
		wantErr bool
		want    []models.CheckRun
	}{
		{
			wantErr: false,
			want:    checkRuns,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _redis.GetStepCheckRuns(t.Context(), step)

		if test.wantErr {
			if err == nil {
				t.Errorf("GetStepCheckRuns should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetStepCheckRuns returned err: %v", err)
		}

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("GetStepCheckRuns() mismatch (-want +got):\n%s", diff)
		}
	}
}
