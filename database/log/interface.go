// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// LogInterface represents the Vela interface for log
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type LogInterface interface {
	// Log Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateLogIndexes defines a function that creates the indexes for the logs table.
	CreateLogIndexes(context.Context) error
	// CreateLogTable defines a function that creates the logs table.
	CreateLogTable(context.Context, string) error

	// Log Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountLogs defines a function that gets the count of all logs.
	CountLogs(context.Context) (int64, error)
	// CountLogsForBuild defines a function that gets the count of logs by build ID.
	CountLogsForBuild(context.Context, *api.Build) (int64, error)
	// CreateLog defines a function that creates a new log.
	CreateLog(context.Context, *api.Log) error
	// DeleteLog defines a function that deletes an existing log.
	DeleteLog(context.Context, *api.Log) error
	// GetLog defines a function that gets a log by ID.
	GetLog(context.Context, int64) (*api.Log, error)
	// GetLogForService defines a function that gets a log by service ID.
	GetLogForService(context.Context, *api.Service) (*api.Log, error)
	// GetLogForStep defines a function that gets a log by step ID.
	GetLogForStep(context.Context, *api.Step) (*api.Log, error)
	// ListLogs defines a function that gets a list of all logs.
	ListLogs(context.Context) ([]*api.Log, error)
	// ListLogsForBuild defines a function that gets a list of logs by build ID.
	ListLogsForBuild(context.Context, *api.Build, int, int) ([]*api.Log, error)
	// UpdateLog defines a function that updates an existing log.
	UpdateLog(context.Context, *api.Log) error
	// CleanLogs defines a function that deletes logs older than a specified timestamp in batches.
	CleanLogs(context.Context, int64, int, bool, string) (*CleanupResult, error)
}
