// Copyright (c) 2023 Target Brands, Ine. All rights reserved.
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

// NewResources creates and returns the database agnostic engines for resources.
func (e *engine) NewResources() error {
	var err error

	// create the database agnostic engine for builds
	e.BuildInterface, err = build.New(
		build.WithClient(e.Database),
		build.WithLogger(e.Logger),
		build.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for hooks
	e.HookInterface, err = hook.New(
		hook.WithClient(e.Database),
		hook.WithLogger(e.Logger),
		hook.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for logs
	e.LogInterface, err = log.New(
		log.WithClient(e.Database),
		log.WithCompressionLevel(e.Config.CompressionLevel),
		log.WithLogger(e.Logger),
		log.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for pipelines
	e.PipelineInterface, err = pipeline.New(
		pipeline.WithClient(e.Database),
		pipeline.WithCompressionLevel(e.Config.CompressionLevel),
		pipeline.WithLogger(e.Logger),
		pipeline.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for repos
	e.RepoInterface, err = repo.New(
		repo.WithClient(e.Database),
		repo.WithEncryptionKey(e.Config.EncryptionKey),
		repo.WithLogger(e.Logger),
		repo.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for schedules
	e.ScheduleInterface, err = schedule.New(
		schedule.WithClient(e.Database),
		schedule.WithLogger(e.Logger),
		schedule.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for secrets
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/secret#New
	e.SecretInterface, err = secret.New(
		secret.WithClient(e.Database),
		secret.WithEncryptionKey(e.Config.EncryptionKey),
		secret.WithLogger(e.Logger),
		secret.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for services
	e.ServiceInterface, err = service.New(
		service.WithClient(e.Database),
		service.WithLogger(e.Logger),
		service.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for steps
	e.StepInterface, err = step.New(
		step.WithClient(e.Database),
		step.WithLogger(e.Logger),
		step.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for users
	e.UserInterface, err = user.New(
		user.WithClient(e.Database),
		user.WithEncryptionKey(e.Config.EncryptionKey),
		user.WithLogger(e.Logger),
		user.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for workers
	e.WorkerInterface, err = worker.New(
		worker.WithClient(e.Database),
		worker.WithLogger(e.Logger),
		worker.WithSkipCreation(e.Config.SkipCreation),
	)
	if err != nil {
		return err
	}

	return nil
}
