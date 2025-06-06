// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
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
	"github.com/go-vela/server/tracing"
)

type (
	// config represents the settings required to create the engine that implements the Interface.
	config struct {
		// specifies the address to use for the database engine
		Address string
		// specifies the level of compression to use for the database engine
		CompressionLevel int
		// specifies the connection duration to use for the database engine
		ConnectionLife time.Duration
		// specifies the maximum idle connections for the database engine
		ConnectionIdle int
		// specifies the maximum open connections for the database engine
		ConnectionOpen int
		// specifies the driver to use for the database engine
		Driver string
		// specifies the encryption key to use for the database engine
		EncryptionKey string
		// specifies the database engine specific log level
		LogLevel string
		// specifies to skip logging when a record is not found
		LogSkipNotFound bool
		// specifies the threshold for slow queries in the database engine
		LogSlowThreshold time.Duration
		// specifies whether to log SQL queries in the database engine
		LogShowSQL bool
		// specifies to skip creating tables and indexes for the database engine
		SkipCreation bool
	}

	// engine represents the functionality that implements the Interface.
	engine struct {
		// gorm.io/gorm database client used in database functions
		client *gorm.DB
		// engine configuration settings used in database functions
		config *config
		// engine context used in database functions
		ctx context.Context
		// sirupsen/logrus logger used in database functions
		logger *logrus.Entry
		// configurations related to telemetry/tracing
		tracing *tracing.Client

		settings.SettingsInterface
		build.BuildInterface
		dashboard.DashboardInterface
		executable.BuildExecutableInterface
		deployment.DeploymentInterface
		hook.HookInterface
		jwk.JWKInterface
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

// New creates and returns an engine capable of integrating with the configured database provider.
//
// Currently, the following database providers are supported:
//
// * postgres
// * sqlite3
// .
func New(opts ...EngineOpt) (Interface, error) {
	// create new database engine
	e := new(engine)

	// create new fields
	e.client = new(gorm.DB)
	e.config = new(config)
	e.logger = new(logrus.Entry)
	e.ctx = context.Background()

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	// validate the configuration being provided
	err := e.config.Validate()
	if err != nil {
		return nil, err
	}

	// by default use the global logger with additional metadata
	e.logger = logrus.NewEntry(logrus.StandardLogger()).WithField("database", e.Driver())

	// translate the log level to logrus level for the database engine
	var dbLogLevel logrus.Level

	switch e.config.LogLevel {
	case "t", "trace", "Trace", "TRACE":
		dbLogLevel = logrus.TraceLevel
	case "d", "debug", "Debug", "DEBUG":
		dbLogLevel = logrus.DebugLevel
	case "i", "info", "Info", "INFO":
		dbLogLevel = logrus.InfoLevel
	case "w", "warn", "Warn", "WARN":
		dbLogLevel = logrus.WarnLevel
	case "e", "error", "Error", "ERROR":
		dbLogLevel = logrus.ErrorLevel
	case "f", "fatal", "Fatal", "FATAL":
		dbLogLevel = logrus.FatalLevel
	case "p", "panic", "Panic", "PANIC":
		dbLogLevel = logrus.PanicLevel
	}

	// if the log level for the database engine is different than
	// the global log level, create a new logrus instance
	if dbLogLevel != logrus.GetLevel() {
		log := logrus.New()

		// set the custom log level
		log.Level = dbLogLevel

		// copy the formatter from the global logger to
		// retain the same format for the database engine
		log.Formatter = logrus.StandardLogger().Formatter

		// update the logger with additional metadata
		e.logger = logrus.NewEntry(log).WithField("database", e.Driver())
	}

	e.logger.Trace("creating database engine from configuration")

	// configure gorm to use logrus as internal logger
	gormConfig := &gorm.Config{
		Logger: NewGormLogger(e.logger, e.config.LogSlowThreshold, e.config.LogSkipNotFound, e.config.LogShowSQL),
	}

	switch e.config.Driver {
	case constants.DriverPostgres:
		// create the new Postgres database client
		e.client, err = gorm.Open(postgres.Open(e.config.Address), gormConfig)
		if err != nil {
			return nil, err
		}
	case constants.DriverSqlite:
		// create the new Sqlite database client
		e.client, err = gorm.Open(sqlite.Open(e.config.Address), gormConfig)
		if err != nil {
			return nil, err
		}
	default:
		// handle an invalid database driver being provided
		return nil, fmt.Errorf("invalid database driver provided: %s", e.Driver())
	}

	// capture database/sql database from gorm.io/gorm database
	db, err := e.client.DB()
	if err != nil {
		return nil, err
	}

	// initialize otel tracing if enabled
	if e.tracing.EnableTracing {
		otelPlugin := otelgorm.NewPlugin(
			otelgorm.WithTracerProvider(e.tracing.TracerProvider),
			otelgorm.WithoutQueryVariables(),
		)

		err := e.client.Use(otelPlugin)
		if err != nil {
			return nil, err
		}
	}

	// set the maximum amount of time a connection may be reused
	db.SetConnMaxLifetime(e.config.ConnectionLife)
	// set the maximum number of connections in the idle connection pool
	db.SetMaxIdleConns(e.config.ConnectionIdle)
	// set the maximum number of open connections to the database
	db.SetMaxOpenConns(e.config.ConnectionOpen)

	// verify connection to the database
	err = e.Ping()
	if err != nil {
		return nil, err
	}

	// create database agnostic engines for resources
	err = e.NewResources(e.ctx)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// NewTest creates and returns an engine that integrates with an in-memory database provider.
//
// This function is ONLY intended to be used for testing purposes.
func NewTest() (Interface, error) {
	return New(
		WithAddress("file::memory:?cache=shared"),
		WithCompressionLevel(3),
		WithConnectionLife(30*time.Minute),
		WithConnectionIdle(2),
		WithConnectionOpen(0),
		WithDriver("sqlite3"),
		WithEncryptionKey("A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW"),
		WithSkipCreation(false),
		WithLogLevel("warn"),
		WithLogShowSQL(false),
		WithLogSkipNotFound(true),
		WithLogSlowThreshold(200*time.Millisecond),
		WithTracing(&tracing.Client{Config: tracing.Config{EnableTracing: false}}),
	)
}
