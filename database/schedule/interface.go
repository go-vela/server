// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"context"

	"github.com/go-vela/types/library"
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
	CreateScheduleIndexes() error
	// CreateScheduleTable defines a function that creates the schedules table.
	CreateScheduleTable(string) error

	// Schedule Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountSchedules defines a function that gets the count of all schedules.
	CountSchedules() (int64, error)
	// CountSchedulesForRepo defines a function that gets the count of schedules by repo ID.
	CountSchedulesForRepo(*library.Repo) (int64, error)
	// CreateSchedule defines a function that creates a new schedule.
	CreateSchedule(context.Context, *library.Schedule) error
	// DeleteSchedule defines a function that deletes an existing schedule.
	DeleteSchedule(*library.Schedule) error
	// GetSchedule defines a function that gets a schedule by ID.
	GetSchedule(int64) (*library.Schedule, error)
	// GetScheduleForRepo defines a function that gets a schedule by repo ID and name.
	GetScheduleForRepo(*library.Repo, string) (*library.Schedule, error)
	// ListActiveSchedules defines a function that gets a list of all active schedules.
	ListActiveSchedules() ([]*library.Schedule, error)
	// ListSchedules defines a function that gets a list of all schedules.
	ListSchedules() ([]*library.Schedule, error)
	// ListSchedulesForRepo defines a function that gets a list of schedules by repo ID.
	ListSchedulesForRepo(*library.Repo, int, int) ([]*library.Schedule, int64, error)
	// UpdateSchedule defines a function that updates an existing schedule.
	UpdateSchedule(*library.Schedule, bool) error
}
