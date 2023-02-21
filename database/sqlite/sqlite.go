// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"fmt"
	"time"

	"github.com/go-vela/server/database/pipeline"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/server/database/sqlite/ddl"
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
		// https://pkg.go.dev/github.com/go-vela/server/database/pipeline#PipelineService
		pipeline.PipelineService
		// https://pkg.go.dev/github.com/go-vela/server/database/repo#RepoService
		repo.RepoService
		// https://pkg.go.dev/github.com/go-vela/server/database/user#UserService
		user.UserService
		// https://pkg.go.dev/github.com/go-vela/server/database/worker#WorkerService
		worker.WorkerService
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

	// create the hooks table
	err = c.Sqlite.Exec(ddl.CreateHookTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %w", constants.TableHook, err)
	}

	// create the logs table
	err = c.Sqlite.Exec(ddl.CreateLogTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %w", constants.TableLog, err)
	}

	// create the secrets table
	err = c.Sqlite.Exec(ddl.CreateSecretTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %w", constants.TableSecret, err)
	}

	// create the services table
	err = c.Sqlite.Exec(ddl.CreateServiceTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %w", constants.TableService, err)
	}

	// create the steps table
	err = c.Sqlite.Exec(ddl.CreateStepTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %w", constants.TableStep, err)
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

	// create the hooks_repo_id index for the hooks table
	err = c.Sqlite.Exec(ddl.CreateHookRepoIDIndex).Error
	if err != nil {
		return fmt.Errorf("unable to create hooks_repo_id index for the %s table: %w", constants.TableHook, err)
	}

	// create the logs_build_id index for the logs table
	err = c.Sqlite.Exec(ddl.CreateLogBuildIDIndex).Error
	if err != nil {
		return fmt.Errorf("unable to create logs_build_id index for the %s table: %w", constants.TableLog, err)
	}

	// create the secrets_type_org_repo index for the secrets table
	err = c.Sqlite.Exec(ddl.CreateSecretTypeOrgRepo).Error
	if err != nil {
		return fmt.Errorf("unable to create secrets_type_org_repo index for the %s table: %w", constants.TableSecret, err)
	}

	// create the secrets_type_org_team index for the secrets table
	err = c.Sqlite.Exec(ddl.CreateSecretTypeOrgTeam).Error
	if err != nil {
		return fmt.Errorf("unable to create secrets_type_org_team index for the %s table: %w", constants.TableSecret, err)
	}

	// create the secrets_type_org index for the secrets table
	err = c.Sqlite.Exec(ddl.CreateSecretTypeOrg).Error
	if err != nil {
		return fmt.Errorf("unable to create secrets_type_org index for the %s table: %w", constants.TableSecret, err)
	}

	return nil
}

// createServices is a helper function to create the database services.
func createServices(c *client) error {
	var err error

	// create the database agnostic pipeline service
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/pipeline#New
	c.PipelineService, err = pipeline.New(
		pipeline.WithClient(c.Sqlite),
		pipeline.WithCompressionLevel(c.config.CompressionLevel),
		pipeline.WithLogger(c.Logger),
		pipeline.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic repo service
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/repo#New
	c.RepoService, err = repo.New(
		repo.WithClient(c.Sqlite),
		repo.WithEncryptionKey(c.config.EncryptionKey),
		repo.WithLogger(c.Logger),
		repo.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic user service
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/user#New
	c.UserService, err = user.New(
		user.WithClient(c.Sqlite),
		user.WithEncryptionKey(c.config.EncryptionKey),
		user.WithLogger(c.Logger),
		user.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	// create the database agnostic worker service
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/worker#New
	c.WorkerService, err = worker.New(
		worker.WithClient(c.Sqlite),
		worker.WithLogger(c.Logger),
		worker.WithSkipCreation(c.config.SkipCreation),
	)
	if err != nil {
		return err
	}

	return nil
}
