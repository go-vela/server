// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"fmt"
	"time"

	"github.com/go-vela/server/database/hook"
	"github.com/go-vela/server/database/log"
	"github.com/go-vela/server/database/pipeline"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/server/database/schedule"
	"github.com/go-vela/server/database/secret"
	"github.com/go-vela/server/database/service"
	"github.com/go-vela/server/database/sqlite/ddl"
	"github.com/go-vela/server/database/step"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/server/database/worker"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	config struct {
		// specifies the address to use for the Sqlite client
		Address string
		// specifies the level of compression to use for the Sqlite client
		CompressionLevel int
		// specifies the connection duration to use for the Sqlite client
		ConnectionLife time.Duration
		// specifies the maximum idle connections for the Sqlite client
		ConnectionIdle int
		// specifies the maximum open connections for the Sqlite client
		ConnectionOpen int
		// specifies the encryption key to use for the Sqlite client
		EncryptionKey string
		// specifies to skip creating tables and indexes for the Sqlite client
		SkipCreation bool
	}

	client struct {
		config *config
		// https://pkg.go.dev/gorm.io/gorm#DB
		Sqlite *gorm.DB
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		Logger *logrus.Entry
		// https://pkg.go.dev/github.com/go-vela/server/database/hook#HookInterface
		hook.HookInterface
		// https://pkg.go.dev/github.com/go-vela/server/database/log#LogInterface
		log.LogInterface
		// https://pkg.go.dev/github.com/go-vela/server/database/pipeline#PipelineInterface
		pipeline.PipelineInterface
		// https://pkg.go.dev/github.com/go-vela/server/database/repo#RepoInterface
		repo.RepoInterface
		// https://pkg.go.dev/github.com/go-vela/server/database/schedule#ScheduleInterface
		schedule.ScheduleInterface
		// https://pkg.go.dev/github.com/go-vela/server/database/secret#SecretInterface
		secret.SecretInterface
		// https://pkg.go.dev/github.com/go-vela/server/database/service#ServiceInterface
		service.ServiceInterface
		// https://pkg.go.dev/github.com/go-vela/server/database/step#StepInterface
		step.StepInterface
		// https://pkg.go.dev/github.com/go-vela/server/database/user#UserInterface
		user.UserInterface
		// https://pkg.go.dev/github.com/go-vela/server/database/worker#WorkerInterface
		worker.WorkerInterface
	}
)

// New returns a Database implementation that integrates with a Sqlite instance.
//
//nolint:revive // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new Sqlite client
	c := new(client)

	// create new fields
	c.config = new(config)
	c.Sqlite = new(gorm.DB)

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
	c.Logger = logrus.NewEntry(logger).WithField("database", c.Driver())

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// create the new Sqlite database client
	//
	// https://pkg.go.dev/gorm.io/gorm#Open
	_sqlite, err := gorm.Open(sqlite.Open(c.config.Address), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// set the Sqlite database client in the Sqlite client
	c.Sqlite = _sqlite

	// setup database with proper configuration
	err = setupDatabase(c)
	if err != nil {
		return nil, err
	}

	// create the services for the database
	err = createServices(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewTest returns a Database implementation that integrates with a fake Sqlite instance.
//
// This function is intended for running tests only.
//
//nolint:revive // ignore returning unexported client
func NewTest() (*client, error) {
	// create new Sqlite client
	c := new(client)

	// create new fields
	c.config = &config{
		CompressionLevel: 3,
		ConnectionLife:   30 * time.Minute,
		ConnectionIdle:   2,
		ConnectionOpen:   0,
		EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
		SkipCreation:     false,
	}
	c.Sqlite = new(gorm.DB)

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
	c.Logger = logrus.NewEntry(logger).WithField("database", c.Driver())

	// create the new Sqlite database client
	//
	// https://pkg.go.dev/gorm.io/gorm#Open
	_sqlite, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		return nil, err
	}

	c.Sqlite = _sqlite

	// setup database with proper configuration
	err = createServices(c)
	if err != nil {
		return nil, err
	}

	// create the tables in the database
	err = createTables(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// setupDatabase is a helper function to setup
// the database with the proper configuration.
func setupDatabase(c *client) error {
	// capture database/sql database from gorm database
	//
	// https://pkg.go.dev/gorm.io/gorm#DB.DB
	_sql, err := c.Sqlite.DB()
	if err != nil {
		return err
	}

	// set the maximum amount of time a connection may be reused
	//
	// https://golang.org/pkg/database/sql/#DB.SetConnMaxLifetime
	_sql.SetConnMaxLifetime(c.config.ConnectionLife)

	// set the maximum number of connections in the idle connection pool
	//
	// https://golang.org/pkg/database/sql/#DB.SetMaxIdleConns
	_sql.SetMaxIdleConns(c.config.ConnectionIdle)

	// set the maximum number of open connections to the database
	//
	// https://golang.org/pkg/database/sql/#DB.SetMaxOpenConns
	_sql.SetMaxOpenConns(c.config.ConnectionOpen)

	// verify connection to the database
	err = c.Ping()
	if err != nil {
		return err
	}

	// check if we should skip creating database objects
	if c.config.SkipCreation {
		c.Logger.Warning("skipping creation of data tables and indexes in the sqlite database")

		return nil
	}

	// create the tables in the database
	err = createTables(c)
	if err != nil {
		return err
	}

	// create the indexes in the database
	err = createIndexes(c)
	if err != nil {
		return err
	}

	return nil
}

// createTables is a helper function to setup
// the database with the necessary tables.
func createTables(c *client) error {
	c.Logger.Trace("creating data tables in the sqlite database")

	// create the builds table
	err := c.Sqlite.Exec(ddl.CreateBuildTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %w", constants.TableBuild, err)
	}

	return nil
}

// createIndexes is a helper function to setup
// the database with the necessary indexes.
func createIndexes(c *client) error {
	c.Logger.Trace("creating data indexes in the sqlite database")

	// create the builds_repo_id index for the builds table
	err := c.Sqlite.Exec(ddl.CreateBuildRepoIDIndex).Error
	if err != nil {
		return fmt.Errorf("unable to create builds_repo_id index for the %s table: %w", constants.TableBuild, err)
	}

	// create the builds_status index for the builds table
	err = c.Sqlite.Exec(ddl.CreateBuildStatusIndex).Error
	if err != nil {
		return fmt.Errorf("unable to create builds_status index for the %s table: %w", constants.TableBuild, err)
	}

	// create the builds_created index for the builds table
	err = c.Sqlite.Exec(ddl.CreateBuildCreatedIndex).Error
	if err != nil {
		return fmt.Errorf("unable to create builds_created index for the %s table: %w", constants.TableBuild, err)
	}

	// create the builds_source index for the builds table
	err = c.Sqlite.Exec(ddl.CreateBuildSourceIndex).Error
	if err != nil {
		return fmt.Errorf("unable to create builds_source index for the %s table: %w", constants.TableBuild, err)
	}

	return nil
}

// createServices is a helper function to create the database services.
func createServices(c *client) error {
	var err error

	// create the database agnostic engine for hooks
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/hook#New
	c.HookInterface, err = hook.New(
		hook.WithClient(c.Sqlite),
		hook.WithLogger(c.Logger),
		hook.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for logs
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/log#New
	c.LogInterface, err = log.New(
		log.WithClient(c.Sqlite),
		log.WithCompressionLevel(c.config.CompressionLevel),
		log.WithLogger(c.Logger),
		log.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for pipelines
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/pipeline#New
	c.PipelineInterface, err = pipeline.New(
		pipeline.WithClient(c.Sqlite),
		pipeline.WithCompressionLevel(c.config.CompressionLevel),
		pipeline.WithLogger(c.Logger),
		pipeline.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for repos
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/repo#New
	c.RepoInterface, err = repo.New(
		repo.WithClient(c.Sqlite),
		repo.WithEncryptionKey(c.config.EncryptionKey),
		repo.WithLogger(c.Logger),
		repo.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for schedules
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/schedule#New
	c.ScheduleInterface, err = schedule.New(
		schedule.WithClient(c.Sqlite),
		schedule.WithLogger(c.Logger),
		schedule.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for secrets
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/secret#New
	c.SecretInterface, err = secret.New(
		secret.WithClient(c.Sqlite),
		secret.WithEncryptionKey(c.config.EncryptionKey),
		secret.WithLogger(c.Logger),
		secret.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for services
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/service#New
	c.ServiceInterface, err = service.New(
		service.WithClient(c.Sqlite),
		service.WithLogger(c.Logger),
		service.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for steps
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/step#New
	c.StepInterface, err = step.New(
		step.WithClient(c.Sqlite),
		step.WithLogger(c.Logger),
		step.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for users
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/user#New
	c.UserInterface, err = user.New(
		user.WithClient(c.Sqlite),
		user.WithEncryptionKey(c.config.EncryptionKey),
		user.WithLogger(c.Logger),
		user.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic engine for workers
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/worker#New
	c.WorkerInterface, err = worker.New(
		worker.WithClient(c.Sqlite),
		worker.WithLogger(c.Logger),
		worker.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	return nil
}
