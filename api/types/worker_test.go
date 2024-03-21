// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-vela/types/library"
)

func TestTypes_Worker_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		worker *Worker
		want   *Worker
	}{
		{
			worker: testWorker(),
			want:   testWorker(),
		},
		{
			worker: new(Worker),
			want:   new(Worker),
		},
	}

	// run tests
	for _, test := range tests {
		if test.worker.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.worker.GetID(), test.want.GetID())
		}

		if test.worker.GetHostname() != test.want.GetHostname() {
			t.Errorf("GetHostname is %v, want %v", test.worker.GetHostname(), test.want.GetHostname())
		}

		if test.worker.GetAddress() != test.want.GetAddress() {
			t.Errorf("Getaddress is %v, want %v", test.worker.GetAddress(), test.want.GetAddress())
		}

		if !reflect.DeepEqual(test.worker.GetRoutes(), test.want.GetRoutes()) {
			t.Errorf("GetRoutes is %v, want %v", test.worker.GetRoutes(), test.want.GetRoutes())
		}

		if test.worker.GetActive() != test.want.GetActive() {
			t.Errorf("GetActive is %v, want %v", test.worker.GetActive(), test.want.GetActive())
		}

		if test.worker.GetStatus() != test.want.GetStatus() {
			t.Errorf("GetStatus is %v, want %v", test.worker.GetStatus(), test.want.GetStatus())
		}

		if test.worker.GetLastStatusUpdateAt() != test.want.GetLastStatusUpdateAt() {
			t.Errorf("GetLastStatusUpdateAt is %v, want %v", test.worker.GetLastStatusUpdateAt(), test.want.GetLastStatusUpdateAt())
		}

		if !reflect.DeepEqual(test.worker.GetRunningBuilds(), test.want.GetRunningBuilds()) {
			t.Errorf("GetRunningBuildIDs is %v, want %v", test.worker.GetRunningBuilds(), test.want.GetRunningBuilds())
		}

		if test.worker.GetLastBuildStartedAt() != test.want.GetLastBuildStartedAt() {
			t.Errorf("GetLastBuildStartedAt is %v, want %v", test.worker.GetLastBuildStartedAt(), test.want.GetLastBuildStartedAt())
		}

		if test.worker.GetLastBuildFinishedAt() != test.want.GetLastBuildFinishedAt() {
			t.Errorf("GetLastBuildFinishedAt is %v, want %v", test.worker.GetLastBuildFinishedAt(), test.want.GetLastBuildFinishedAt())
		}

		if test.worker.GetLastCheckedIn() != test.want.GetLastCheckedIn() {
			t.Errorf("GetLastCheckedIn is %v, want %v", test.worker.GetLastCheckedIn(), test.want.GetLastCheckedIn())
		}

		if test.worker.GetBuildLimit() != test.want.GetBuildLimit() {
			t.Errorf("GetBuildLimit is %v, want %v", test.worker.GetBuildLimit(), test.want.GetBuildLimit())
		}
	}
}

func TestTypes_Worker_Setters(t *testing.T) {
	// setup types
	var w *Worker

	// setup tests
	tests := []struct {
		worker *Worker
		want   *Worker
	}{
		{
			worker: testWorker(),
			want:   testWorker(),
		},
		{
			worker: w,
			want:   new(Worker),
		},
	}

	// run tests
	for _, test := range tests {
		test.worker.SetID(test.want.GetID())
		test.worker.SetHostname(test.want.GetHostname())
		test.worker.SetAddress(test.want.GetAddress())
		test.worker.SetRoutes(test.want.GetRoutes())
		test.worker.SetActive(test.want.GetActive())
		test.worker.SetStatus(test.want.GetStatus())
		test.worker.SetLastStatusUpdateAt(test.want.GetLastStatusUpdateAt())
		test.worker.SetRunningBuilds(test.want.GetRunningBuilds())
		test.worker.SetLastBuildStartedAt(test.want.GetLastBuildStartedAt())
		test.worker.SetLastBuildFinishedAt(test.want.GetLastBuildFinishedAt())
		test.worker.SetLastCheckedIn(test.want.GetLastCheckedIn())
		test.worker.SetBuildLimit(test.want.GetBuildLimit())

		if test.worker.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.worker.GetID(), test.want.GetID())
		}

		if test.worker.GetHostname() != test.want.GetHostname() {
			t.Errorf("SetHostname is %v, want %v", test.worker.GetHostname(), test.want.GetHostname())
		}

		if test.worker.GetAddress() != test.want.GetAddress() {
			t.Errorf("SetAddress is %v, want %v", test.worker.GetAddress(), test.want.GetAddress())
		}

		if !reflect.DeepEqual(test.worker.GetRoutes(), test.want.GetRoutes()) {
			t.Errorf("SetRoutes is %v, want %v", test.worker.GetRoutes(), test.want.GetRoutes())
		}

		if test.worker.GetActive() != test.want.GetActive() {
			t.Errorf("SetActive is %v, want %v", test.worker.GetActive(), test.want.GetActive())
		}

		if test.worker.GetStatus() != test.want.GetStatus() {
			t.Errorf("SetStatus is %v, want %v", test.worker.GetStatus(), test.want.GetStatus())
		}

		if test.worker.GetLastStatusUpdateAt() != test.want.GetLastStatusUpdateAt() {
			t.Errorf("SetLastStatusUpdateAt is %v, want %v", test.worker.GetLastStatusUpdateAt(), test.want.GetLastStatusUpdateAt())
		}

		if test.worker.GetLastBuildStartedAt() != test.want.GetLastBuildStartedAt() {
			t.Errorf("SetLastBuildStartedAt is %v, want %v", test.worker.GetLastBuildStartedAt(), test.want.GetLastBuildStartedAt())
		}

		if test.worker.GetLastBuildFinishedAt() != test.want.GetLastBuildFinishedAt() {
			t.Errorf("SetLastBuildFinishedAt is %v, want %v", test.worker.GetLastBuildFinishedAt(), test.want.GetLastBuildFinishedAt())
		}

		if test.worker.GetLastCheckedIn() != test.want.GetLastCheckedIn() {
			t.Errorf("SetLastCheckedIn is %v, want %v", test.worker.GetLastCheckedIn(), test.want.GetLastCheckedIn())
		}

		if test.worker.GetBuildLimit() != test.want.GetBuildLimit() {
			t.Errorf("SetBuildLimit is %v, want %v", test.worker.GetBuildLimit(), test.want.GetBuildLimit())
		}
	}
}

func TestTypes_Worker_String(t *testing.T) {
	// setup types
	w := testWorker()

	want := fmt.Sprintf(`{
  ID: %d,
  Hostname: %s,
  Address: %s,
  Routes: %s,
  Active: %t,
  Status: %s,
  LastStatusUpdateAt: %v,
  LastBuildStartedAt: %v,
  LastBuildFinishedAt: %v,
  LastCheckedIn: %v,
  BuildLimit: %v,
  RunningBuilds: %v,
}`,
		w.GetID(),
		w.GetHostname(),
		w.GetAddress(),
		w.GetRoutes(),
		w.GetActive(),
		w.GetStatus(),
		w.GetLastStatusUpdateAt(),
		w.GetLastBuildStartedAt(),
		w.GetLastBuildFinishedAt(),
		w.GetLastCheckedIn(),
		w.GetBuildLimit(),
		w.GetRunningBuilds(),
	)

	// run test
	got := w.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testWorker is a test helper function to create a Worker
// type with all fields set to a fake value.
func testWorker() *Worker {
	b := new(library.Build)
	b.SetID(1)

	w := new(Worker)

	w.SetID(1)
	w.SetHostname("worker_0")
	w.SetAddress("http://localhost:8080")
	w.SetRoutes([]string{"vela"})
	w.SetActive(true)
	w.SetStatus("available")
	w.SetLastStatusUpdateAt(time.Time{}.UTC().Unix())
	w.SetRunningBuilds([]*library.Build{b})
	w.SetLastBuildStartedAt(time.Time{}.UTC().Unix())
	w.SetLastBuildFinishedAt(time.Time{}.UTC().Unix())
	w.SetLastCheckedIn(time.Time{}.UTC().Unix())
	w.SetBuildLimit(2)

	return w
}
