// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"github.com/go-vela/types/library"
)

// LogService represents the Vela interface for log
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type LogService interface {
	// Log Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateLogIndexes defines a function that creates the indexes for the logs table.
	CreateLogIndexes() error
	// CreateLogTable defines a function that creates the logs table.
	CreateLogTable(string) error

	// Log Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountLogs defines a function that gets the count of all logs.
	CountLogs() (int64, error)
	// CountLogsForBuild defines a function that gets the count of logs by build ID.
	CountLogsForBuild(*library.Build) (int64, error)
	// CreateLog defines a function that creates a new log.
	CreateLog(*library.Log) error
	// DeleteLog defines a function that deletes an existing log.
	DeleteLog(*library.Log) error
	// GetLog defines a function that gets a log by ID.
	GetLog(int64) (*library.Log, error)
	// GetLogForService defines a function that gets a log by service ID.
	GetLogForService(*library.Service) (*library.Log, error)
	// GetLogForStep defines a function that gets a log by step ID.
	GetLogForStep(*library.Step) (*library.Log, error)
	// ListLogs defines a function that gets a list of all logs.
	ListLogs() ([]*library.Log, error)
	// ListLogsForBuild defines a function that gets a list of logs by build ID.
	ListLogsForBuild(*library.Build, int, int) ([]*library.Log, int64, error)
	// UpdateLog defines a function that updates an existing log.
	UpdateLog(*library.Log) error
}
