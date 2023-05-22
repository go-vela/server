// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/server/database/build"
	"github.com/go-vela/server/database/hook"
	"github.com/go-vela/server/database/log"
	"github.com/go-vela/server/database/pipeline"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/server/database/schedule"
	"github.com/go-vela/server/database/secret"
	"github.com/go-vela/server/database/service"
	"github.com/go-vela/server/database/step"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/server/database/worker"
)

// Interface represents the interface for Vela integrating
// with the different supported Database backends.
type Interface interface {
	// Database Interface Functions

	// Driver defines a function that outputs
	// the configured database driver.
	Driver() string

	// BuildInterface defines the interface for builds stored in the database.
	build.BuildInterface

	// HookInterface defines the interface for hooks stored in the database.
	hook.HookInterface

	// LogInterface defines the interface for logs stored in the database.
	log.LogInterface

	// PipelineInterface defines the interface for pipelines stored in the database.
	pipeline.PipelineInterface

	// RepoInterface defines the interface for repos stored in the database.
	repo.RepoInterface

	// ScheduleInterface defines the interface for schedules stored in the database.
	schedule.ScheduleInterface

	// SecretInterface defines the interface for secrets stored in the database.
	secret.SecretInterface

	// ServiceInterface defines the interface for services stored in the database.
	service.ServiceInterface

	// StepInterface defines the interface for steps stored in the database.
	step.StepInterface

	// UserInterface defines the interface for users stored in the database.
	user.UserInterface

	// WorkerInterface defines the interface for workers stored in the database.
	worker.WorkerInterface
}
