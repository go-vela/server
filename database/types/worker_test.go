// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"strconv"
	"testing"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/types/library"
)

func TestTypes_Worker_Nullify(t *testing.T) {
	// setup types
	var w *Worker

	want := &Worker{
		ID:                  sql.NullInt64{Int64: 0, Valid: false},
		Hostname:            sql.NullString{String: "", Valid: false},
		Address:             sql.NullString{String: "", Valid: false},
		Active:              sql.NullBool{Bool: false, Valid: false},
		Status:              sql.NullString{String: "", Valid: false},
		LastStatusUpdateAt:  sql.NullInt64{Int64: 0, Valid: false},
		LastBuildStartedAt:  sql.NullInt64{Int64: 0, Valid: false},
		LastBuildFinishedAt: sql.NullInt64{Int64: 0, Valid: false},
		LastCheckedIn:       sql.NullInt64{Int64: 0, Valid: false},
		BuildLimit:          sql.NullInt64{Int64: 0, Valid: false},
	}

	// setup tests
	tests := []struct {
		repo *Worker
		want *Worker
	}{
		{
			repo: testWorker(),
			want: testWorker(),
		},
		{
			repo: w,
			want: nil,
		},
		{
			repo: new(Worker),
			want: want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.repo.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestTypes_Worker_ToAPI(t *testing.T) {
	// setup types
	b := new(library.Build)
	b.SetID(12345)

	builds := []*library.Build{b}

	want := new(types.Worker)

	want.SetID(1)
	want.SetHostname("worker_0")
	want.SetAddress("http://localhost:8080")
	want.SetRoutes([]string{"vela"})
	want.SetActive(true)
	want.SetStatus("available")
	want.SetLastStatusUpdateAt(1563474077)
	want.SetRunningBuilds(builds)
	want.SetLastBuildStartedAt(1563474077)
	want.SetLastBuildFinishedAt(1563474077)
	want.SetLastCheckedIn(1563474077)
	want.SetBuildLimit(2)

	// run test
	got := testWorker().ToAPI(builds)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestTypes_Worker_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		worker  *Worker
	}{
		{
			failure: false,
			worker:  testWorker(),
		},
		{ // no Hostname set for worker
			failure: true,
			worker: &Worker{
				ID:            sql.NullInt64{Int64: 1, Valid: true},
				Address:       sql.NullString{String: "http://localhost:8080", Valid: true},
				Active:        sql.NullBool{Bool: true, Valid: true},
				LastCheckedIn: sql.NullInt64{Int64: 1563474077, Valid: true},
			},
		},
		{ // no Address set for worker
			failure: true,
			worker: &Worker{
				ID:            sql.NullInt64{Int64: 1, Valid: true},
				Hostname:      sql.NullString{String: "worker_0", Valid: true},
				Active:        sql.NullBool{Bool: true, Valid: true},
				LastCheckedIn: sql.NullInt64{Int64: 1563474077, Valid: true},
			},
		},
		{ // invalid RunningBuildIDs set for worker
			failure: true,
			worker: &Worker{
				ID:              sql.NullInt64{Int64: 1, Valid: true},
				Address:         sql.NullString{String: "http://localhost:8080", Valid: true},
				Hostname:        sql.NullString{String: "worker_0", Valid: true},
				Active:          sql.NullBool{Bool: true, Valid: true},
				RunningBuildIDs: exceededRunningBuildIDs(),
				LastCheckedIn:   sql.NullInt64{Int64: 1563474077, Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.worker.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}

func TestTypes_WorkerFromLibrary(t *testing.T) {
	// setup types
	b := new(library.Build)
	b.SetID(12345)

	builds := []*library.Build{b}

	w := new(types.Worker)

	w.SetID(1)
	w.SetHostname("worker_0")
	w.SetAddress("http://localhost:8080")
	w.SetRoutes([]string{"vela"})
	w.SetActive(true)
	w.SetStatus("available")
	w.SetLastStatusUpdateAt(1563474077)
	w.SetRunningBuilds(builds)
	w.SetLastBuildStartedAt(1563474077)
	w.SetLastBuildFinishedAt(1563474077)
	w.SetLastCheckedIn(1563474077)
	w.SetBuildLimit(2)

	want := testWorker()

	// run test
	got := WorkerFromAPI(w)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("WorkerFromLibrary is %v, want %v", got, want)
	}
}

// testWorker is a test helper function to create a Worker
// type with all fields set to a fake value.
func testWorker() *Worker {
	return &Worker{
		ID:                  sql.NullInt64{Int64: 1, Valid: true},
		Hostname:            sql.NullString{String: "worker_0", Valid: true},
		Address:             sql.NullString{String: "http://localhost:8080", Valid: true},
		Routes:              []string{"vela"},
		Active:              sql.NullBool{Bool: true, Valid: true},
		Status:              sql.NullString{String: "available", Valid: true},
		LastStatusUpdateAt:  sql.NullInt64{Int64: 1563474077, Valid: true},
		RunningBuildIDs:     []string{"12345"},
		LastBuildStartedAt:  sql.NullInt64{Int64: 1563474077, Valid: true},
		LastBuildFinishedAt: sql.NullInt64{Int64: 1563474077, Valid: true},
		LastCheckedIn:       sql.NullInt64{Int64: 1563474077, Valid: true},
		BuildLimit:          sql.NullInt64{Int64: 2, Valid: true},
	}
}

// exceededRunningBuildIDs returns a list of valid running builds that exceed the maximum size.
func exceededRunningBuildIDs() []string {
	// initialize empty runningBuildIDs
	runningBuildIDs := []string{}

	// add enough build ids to exceed the character limit
	for i := 0; i < 50; i++ {
		// construct runningBuildID
		// use i to adhere to unique runningBuildIDs
		runningBuildID := "1234567890-" + strconv.Itoa(i)

		runningBuildIDs = append(runningBuildIDs, runningBuildID)
	}

	return runningBuildIDs
}
