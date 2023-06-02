// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"fmt"
	"time"

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
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type (
	// Config represents the settings required to create the engine that implements the Interface.
	Config struct {
		// specifies the address to use for the database engine
		Address string
		// specifies the level of compression to use for the database engine
		CompressionLevel int
		// specifies the maximum idle connections for the database engine
		ConnectionIdle int
		// specifies the connection duration to use for the database engine
		ConnectionLife time.Duration
		// specifies the maximum open connections for the database engine
		ConnectionOpen int
		// specifies the driver to use for the database engine
		Driver string
		// specifies the encryption key to use for the database engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the database engine
		SkipCreation bool
	}

	// engine represents the functionality that implements the Interface.
	engine struct {
		Config   *Config
		Database *gorm.DB
		Logger   *logrus.Entry

		build.BuildInterface
		hook.HookInterface
		log.LogInterface
		pipeline.PipelineInterface
		repo.RepoInterface
		schedule.ScheduleInterface
		secret.SecretInterface
		service.ServiceInterface
		step.StepInterface
		user.UserInterface
		worker.WorkerInterface
	}
)

// New creates and returns a Vela service capable of
// integrating with the configured database provider.
//
// Currently the following database providers are supported:
//
// * Postgres
// * Sqlite
// .
func New(s *Setup) (Interface, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/server/database?tab=doc#Setup.Validate
	err := s.Validate()
	if err != nil {
		return nil, err
	}

	logrus.Debug("creating database service from setup")
	// process the database driver being provided
	switch s.Driver {
	case constants.DriverPostgres:
		// handle the Postgres database driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/database?tab=doc#Setup.Postgres
		return s.Postgres()
	case constants.DriverSqlite:
		// handle the Sqlite database driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/database?tab=doc#Setup.Sqlite
		return s.Sqlite()
	default:
		// handle an invalid database driver being provided
		return nil, fmt.Errorf("invalid database driver provided: %s", s.Driver)
	}
}
