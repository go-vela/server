// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/lib/pq"
)

var (
	// ErrEmptyWorkerHost defines the error type when a
	// Worker type has an empty Host field provided.
	ErrEmptyWorkerHost = errors.New("empty worker address provided")

	// ErrEmptyWorkerAddress defines the error type when a
	// Worker type has an empty Address field provided.
	ErrEmptyWorkerAddress = errors.New("empty worker address provided")

	// ErrExceededRunningBuildIDsLimit defines the error type when a
	// Worker type has RunningBuildIDs field provided that exceeds the database limit.
	ErrExceededRunningBuildIDsLimit = errors.New("exceeded running build ids limit")
)

// Worker is the database representation of a worker.
type Worker struct {
	ID                  sql.NullInt64  `sql:"id"`
	Hostname            sql.NullString `sql:"hostname"`
	Address             sql.NullString `sql:"address"`
	Routes              pq.StringArray `sql:"routes" gorm:"type:varchar(1000)"`
	Active              sql.NullBool   `sql:"active"`
	Status              sql.NullString `sql:"status"`
	LastStatusUpdateAt  sql.NullInt64  `sql:"last_status_update_at"`
	RunningBuildIDs     pq.StringArray `sql:"running_build_ids" gorm:"type:varchar(500)"`
	LastBuildStartedAt  sql.NullInt64  `sql:"last_build_started_at"`
	LastBuildFinishedAt sql.NullInt64  `sql:"last_build_finished_at"`
	LastCheckedIn       sql.NullInt64  `sql:"last_checked_in"`
	BuildLimit          sql.NullInt64  `sql:"build_limit"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Build type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (w *Worker) Nullify() *Worker {
	if w == nil {
		return nil
	}

	// check if the ID field should be false
	if w.ID.Int64 == 0 {
		w.ID.Valid = false
	}

	// check if the Hostname field should be false
	if len(w.Hostname.String) == 0 {
		w.Hostname.Valid = false
	}

	// check if the Address field should be false
	if len(w.Address.String) == 0 {
		w.Address.Valid = false
	}

	// check if the Status field should be false
	if len(w.Status.String) == 0 {
		w.Status.Valid = false
	}

	// check if the LastStatusUpdateAt field should be false
	if w.LastStatusUpdateAt.Int64 == 0 {
		w.LastStatusUpdateAt.Valid = false
	}

	// check if the LastBuildStartedAt field should be false
	if w.LastBuildStartedAt.Int64 == 0 {
		w.LastBuildStartedAt.Valid = false
	}

	// check if the LastBuildFinishedAt field should be false
	if w.LastBuildFinishedAt.Int64 == 0 {
		w.LastBuildFinishedAt.Valid = false
	}

	// check if the LastCheckedIn field should be false
	if w.LastCheckedIn.Int64 == 0 {
		w.LastCheckedIn.Valid = false
	}

	if w.BuildLimit.Int64 == 0 {
		w.BuildLimit.Valid = false
	}

	return w
}

// ToAPI converts the Worker type
// to an API Worker type.
func (w *Worker) ToAPI(builds []*library.Build) *types.Worker {
	worker := new(types.Worker)

	worker.SetID(w.ID.Int64)
	worker.SetHostname(w.Hostname.String)
	worker.SetAddress(w.Address.String)
	worker.SetRoutes(w.Routes)
	worker.SetActive(w.Active.Bool)
	worker.SetStatus(w.Status.String)
	worker.SetLastStatusUpdateAt(w.LastStatusUpdateAt.Int64)
	worker.SetRunningBuilds(builds)
	worker.SetLastBuildStartedAt(w.LastBuildStartedAt.Int64)
	worker.SetLastBuildFinishedAt(w.LastBuildFinishedAt.Int64)
	worker.SetLastCheckedIn(w.LastCheckedIn.Int64)
	worker.SetBuildLimit(w.BuildLimit.Int64)

	return worker
}

// Validate verifies the necessary fields for
// the Worker type are populated correctly.
func (w *Worker) Validate() error {
	// verify the Host field is populated
	if len(w.Hostname.String) == 0 {
		return ErrEmptyWorkerHost
	}

	// verify the Address field is populated
	if len(w.Address.String) == 0 {
		return ErrEmptyWorkerAddress
	}

	// calculate total size of RunningBuildIds
	total := 0
	for _, f := range w.RunningBuildIDs {
		total += len(f)
	}

	// verify the RunningBuildIds field is within the database constraints
	// len is to factor in number of comma separators included in the database field,
	// removing 1 due to the last item not having an appended comma
	if (total + len(w.RunningBuildIDs) - 1) > constants.RunningBuildIDsMaxSize {
		return ErrExceededRunningBuildIDsLimit
	}

	// ensure that all Worker string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	w.Hostname = sql.NullString{String: sanitize(w.Hostname.String), Valid: w.Hostname.Valid}
	w.Address = sql.NullString{String: sanitize(w.Address.String), Valid: w.Address.Valid}

	// ensure that all Routes are sanitized
	// to avoid unsafe HTML content
	for i, v := range w.Routes {
		w.Routes[i] = sanitize(v)
	}

	return nil
}

// WorkerFromAPI converts the library worker type
// to a database worker type.
func WorkerFromAPI(w *types.Worker) *Worker {
	var rBs []string

	for _, b := range w.GetRunningBuilds() {
		rBs = append(rBs, fmt.Sprint(b.GetID()))
	}

	worker := &Worker{
		ID:                  sql.NullInt64{Int64: w.GetID(), Valid: true},
		Hostname:            sql.NullString{String: w.GetHostname(), Valid: true},
		Address:             sql.NullString{String: w.GetAddress(), Valid: true},
		Routes:              pq.StringArray(w.GetRoutes()),
		Active:              sql.NullBool{Bool: w.GetActive(), Valid: true},
		Status:              sql.NullString{String: w.GetStatus(), Valid: true},
		LastStatusUpdateAt:  sql.NullInt64{Int64: w.GetLastStatusUpdateAt(), Valid: true},
		RunningBuildIDs:     pq.StringArray(rBs),
		LastBuildStartedAt:  sql.NullInt64{Int64: w.GetLastBuildStartedAt(), Valid: true},
		LastBuildFinishedAt: sql.NullInt64{Int64: w.GetLastBuildFinishedAt(), Valid: true},
		LastCheckedIn:       sql.NullInt64{Int64: w.GetLastCheckedIn(), Valid: true},
		BuildLimit:          sql.NullInt64{Int64: w.GetBuildLimit(), Valid: true},
	}

	return worker.Nullify()
}
