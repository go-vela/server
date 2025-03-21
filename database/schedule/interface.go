// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// ScheduleInterface represents the Vela interface for schedule
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type ScheduleInterface interface {
	// Schedule Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateScheduleIndexes defines a function that creates the indexes for the schedules table.
	CreateScheduleIndexes(context.Context) error
	// CreateScheduleTable defines a function that creates the schedules table.
	CreateScheduleTable(context.Context, string) error

	// Schedule Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountSchedules defines a function that gets the count of all schedules.
	CountSchedules(context.Context) (int64, error)
	// CountSchedulesForRepo defines a function that gets the count of schedules by repo ID.
	CountSchedulesForRepo(context.Context, *api.Repo) (int64, error)
	// CreateSchedule defines a function that creates a new schedule.
	CreateSchedule(context.Context, *api.Schedule) (*api.Schedule, error)
	// DeleteSchedule defines a function that deletes an existing schedule.
	DeleteSchedule(context.Context, *api.Schedule) error
	// GetSchedule defines a function that gets a schedule by ID.
	GetSchedule(context.Context, int64) (*api.Schedule, error)
	// GetScheduleForRepo defines a function that gets a schedule by repo ID and name.
	GetScheduleForRepo(context.Context, *api.Repo, string) (*api.Schedule, error)
	// ListActiveSchedules defines a function that gets a list of all active schedules.
	ListActiveSchedules(context.Context) ([]*api.Schedule, error)
	// ListSchedules defines a function that gets a list of all schedules.
	ListSchedules(context.Context) ([]*api.Schedule, error)
	// ListSchedulesForRepo defines a function that gets a list of schedules by repo ID.
	ListSchedulesForRepo(context.Context, *api.Repo, int, int) ([]*api.Schedule, error)
	// UpdateSchedule defines a function that updates an existing schedule.
	UpdateSchedule(context.Context, *api.Schedule, bool) (*api.Schedule, error)
}
