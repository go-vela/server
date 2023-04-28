// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"github.com/go-vela/server/api/types"
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
	CreateSchedule(*types.Schedule) error
	// DeleteSchedule defines a function that deletes an existing schedule.
	DeleteSchedule(*types.Schedule) error
	// GetSchedule defines a function that gets a schedule by ID.
	GetSchedule(int64) (*types.Schedule, error)
	// GetScheduleForRepo defines a function that gets a schedule by repo ID and number.
	GetScheduleForRepo(*library.Repo, int) (*types.Schedule, error)
	// ListActiveSchedules defines a function that gets a list of all active schedules.
	ListActiveSchedules() ([]*types.Schedule, error)
	// ListSchedules defines a function that gets a list of all schedules.
	ListSchedules() ([]*types.Schedule, error)
	// ListSchedulesForRepo defines a function that gets a list of schedules by repo ID.
	ListSchedulesForRepo(*library.Repo, int, int) ([]*types.Schedule, int64, error)
	// UpdateSchedule defines a function that updates an existing schedule.
	UpdateSchedule(*types.Schedule) error
}
