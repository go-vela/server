// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router"
	"github.com/go-vela/server/router/middleware"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"golang.org/x/sync/errgroup"

	"k8s.io/apimachinery/pkg/util/wait"
)

func server(c *cli.Context) error {
	// validate all input
	err := validate(c)
	if err != nil {
		return err
	}

	// set log level for logrus
	switch c.String("log-level") {
	case "t", "trace", "Trace", "TRACE":
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	case "i", "info", "Info", "INFO":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.InfoLevel)
	case "w", "warn", "Warn", "WARN":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.PanicLevel)
	}

	compiler, err := setupCompiler(c)
	if err != nil {
		return err
	}

	tp, err := tracerProvider()
	if err != nil {
		logrus.Fatal(err)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logrus.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	database, err := database.FromCLIContext(c, tp)
	if err != nil {
		return err
	}

	queue, err := setupQueue(c)
	if err != nil {
		return err
	}

	secrets, err := setupSecrets(c, database)
	if err != nil {
		return err
	}

	scm, err := setupSCM(c)
	if err != nil {
		return err
	}

	metadata, err := setupMetadata(c)
	if err != nil {
		return err
	}

	router := router.Load(
		middleware.Compiler(compiler),
		middleware.Database(database),
		middleware.Logger(logrus.StandardLogger(), time.RFC3339),
		middleware.Metadata(metadata),
		middleware.TokenManager(setupTokenManager(c)),
		middleware.Queue(queue),
		middleware.RequestVersion,
		middleware.Secret(c.String("vela-secret")),
		middleware.Secrets(secrets),
		middleware.Scm(scm),
		middleware.Allowlist(c.StringSlice("vela-repo-allowlist")),
		middleware.DefaultBuildLimit(c.Int64("default-build-limit")),
		middleware.DefaultTimeout(c.Int64("default-build-timeout")),
		middleware.MaxBuildLimit(c.Int64("max-build-limit")),
		middleware.WebhookValidation(!c.Bool("vela-disable-webhook-validation")),
		middleware.SecureCookie(c.Bool("vela-enable-secure-cookie")),
		middleware.Worker(c.Duration("worker-active-interval")),
		middleware.DefaultRepoEvents(c.StringSlice("default-repo-events")),
		middleware.AllowlistSchedule(c.StringSlice("vela-schedule-allowlist")),
		middleware.ScheduleFrequency(c.Duration("schedule-minimum-frequency")),

		// inject service middleware
		otelgin.Middleware("vela-server", otelgin.WithTracerProvider(tp)),
	)

	addr, err := url.Parse(c.String("server-addr"))
	if err != nil {
		return err
	}

	port := addr.Port()
	// check if a port is part of the address
	if len(port) == 0 {
		port = c.String("server-port")
	}

	// gin expects the address to be ":<port>" ie ":8080"
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadHeaderTimeout: 60 * time.Second,
	}

	// create the context for controlling the worker subprocesses
	ctx, done := context.WithCancel(context.Background())
	// create the errgroup for managing worker subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group
	g, gctx := errgroup.WithContext(ctx)

	// spawn goroutine to check for signals to gracefully shutdown
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			logrus.Infof("received signal: %s", sig)
			err := srv.Shutdown(ctx)
			if err != nil {
				logrus.Error(err)
			}
			done()
		case <-gctx.Done():
			logrus.Info("closing signal goroutine")
			err := srv.Shutdown(ctx)
			if err != nil {
				logrus.Error(err)
			}
			return gctx.Err()
		}

		return nil
	})

	// spawn goroutine for starting the server
	g.Go(func() error {
		logrus.Infof("starting server on %s", addr.Host)
		err = srv.ListenAndServe()
		if err != nil {
			// log a message indicating the failure of the server
			logrus.Errorf("failing server: %v", err)
		}

		return err
	})

	// spawn goroutine for starting the scheduler
	g.Go(func() error {
		logrus.Info("starting scheduler")
		for {
			// cut the configured minimum frequency duration for schedules in half
			//
			// We need to sleep for some amount of time before we attempt to process schedules
			// setup in the database. Since the minimum frequency is configurable, we cut it in
			// half and use that as the base duration to determine how long to sleep for.
			base := c.Duration("schedule-minimum-frequency") / 2
			logrus.Infof("sleeping for %v before scheduling builds", base)

			// sleep for a duration of time before processing schedules
			//
			// This should prevent multiple servers from processing schedules at the same time by
			// leveraging a base duration along with a standard deviation of randomness a.k.a.
			// "jitter". To create the jitter, we use the configured minimum frequency duration
			// along with a scale factor of 0.1.
			time.Sleep(wait.Jitter(base, 0.1))

			err = processSchedules(compiler, database, metadata, queue, scm)
			if err != nil {
				logrus.WithError(err).Warn("unable to process schedules")
			} else {
				logrus.Trace("successfully processed schedules")
			}
		}
	})

	// wait for errors from server subprocesses
	return g.Wait()
}
