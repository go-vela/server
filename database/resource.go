// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"

	"github.com/go-vela/server/database/build"
	"github.com/go-vela/server/database/dashboard"
	"github.com/go-vela/server/database/deployment"
	"github.com/go-vela/server/database/executable"
	"github.com/go-vela/server/database/hook"
	"github.com/go-vela/server/database/jwk"
	"github.com/go-vela/server/database/log"
	"github.com/go-vela/server/database/pipeline"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/server/database/schedule"
	"github.com/go-vela/server/database/secret"
	"github.com/go-vela/server/database/service"
	"github.com/go-vela/server/database/settings"
	"github.com/go-vela/server/database/step"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/server/database/worker"
)

// NewResources creates and returns the database agnostic engines for resources.
//
//nolint:funlen // ignore function length
func (e *engine) NewResources(ctx context.Context) error {
	var err error

	// create the database agnostic engine for settings
	e.SettingsInterface, err = settings.New(
		settings.WithContext(e.ctx),
		settings.WithClient(e.client),
		settings.WithLogger(e.logger),
		settings.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for builds
	e.BuildInterface, err = build.New(
		build.WithContext(e.ctx),
		build.WithClient(e.client),
		build.WithLogger(e.logger),
		build.WithSkipCreation(e.config.SkipCreation),
		build.WithEncryptionKey(e.config.EncryptionKey),
	)
	if err != nil {
		return err
	}

	e.DashboardInterface, err = dashboard.New(
		dashboard.WithContext(e.ctx),
		dashboard.WithClient(e.client),
		dashboard.WithLogger(e.logger),
		dashboard.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for build_executables
	e.BuildExecutableInterface, err = executable.New(
		executable.WithContext(e.ctx),
		executable.WithClient(e.client),
		executable.WithLogger(e.logger),
		executable.WithSkipCreation(e.config.SkipCreation),
		executable.WithEncryptionKey(e.config.EncryptionKey),
		executable.WithDriver(e.config.Driver),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for deployments
	e.DeploymentInterface, err = deployment.New(
		deployment.WithContext(e.ctx),
		deployment.WithClient(e.client),
		deployment.WithLogger(e.logger),
		deployment.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for hooks
	e.HookInterface, err = hook.New(
		hook.WithContext(e.ctx),
		hook.WithClient(e.client),
		hook.WithLogger(e.logger),
		hook.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for JWKs
	e.JWKInterface, err = jwk.New(
		jwk.WithContext(e.ctx),
		jwk.WithClient(e.client),
		jwk.WithLogger(e.logger),
		jwk.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for logs
	e.LogInterface, err = log.New(
		log.WithContext(e.ctx),
		log.WithClient(e.client),
		log.WithCompressionLevel(e.config.CompressionLevel),
		log.WithLogger(e.logger),
		log.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for pipelines
	e.PipelineInterface, err = pipeline.New(
		pipeline.WithContext(e.ctx),
		pipeline.WithClient(e.client),
		pipeline.WithCompressionLevel(e.config.CompressionLevel),
		pipeline.WithLogger(e.logger),
		pipeline.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for repos
	e.RepoInterface, err = repo.New(
		repo.WithContext(e.ctx),
		repo.WithClient(e.client),
		repo.WithEncryptionKey(e.config.EncryptionKey),
		repo.WithLogger(e.logger),
		repo.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for schedules
	e.ScheduleInterface, err = schedule.New(
		schedule.WithContext(e.ctx),
		schedule.WithClient(e.client),
		schedule.WithEncryptionKey(e.config.EncryptionKey),
		schedule.WithLogger(e.logger),
		schedule.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for secrets
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/secret#New
	e.SecretInterface, err = secret.New(
		secret.WithContext(e.ctx),
		secret.WithClient(e.client),
		secret.WithEncryptionKey(e.config.EncryptionKey),
		secret.WithLogger(e.logger),
		secret.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for services
	e.ServiceInterface, err = service.New(
		service.WithClient(e.client),
		service.WithLogger(e.logger),
		service.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for steps
	e.StepInterface, err = step.New(
		step.WithContext(e.ctx),
		step.WithClient(e.client),
		step.WithLogger(e.logger),
		step.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for users
	e.UserInterface, err = user.New(
		user.WithContext(e.ctx),
		user.WithClient(e.client),
		user.WithEncryptionKey(e.config.EncryptionKey),
		user.WithLogger(e.logger),
		user.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for workers
	e.WorkerInterface, err = worker.New(
		worker.WithContext(e.ctx),
		worker.WithClient(e.client),
		worker.WithLogger(e.logger),
		worker.WithSkipCreation(e.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	return nil
}
