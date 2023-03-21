// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-vela/server/router"
	"github.com/go-vela/server/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
	"gopkg.in/tomb.v2"
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

	database, err := setupDatabase(c)
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
	)

	addr, err := url.Parse(c.String("server-addr"))
	if err != nil {
		return err
	}

	var tomb tomb.Tomb
	// start http server
	tomb.Go(func() error {
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

		logrus.Infof("running server on %s", addr.Host)
		go func() {
			logrus.Info("Starting HTTP server...")
			err := srv.ListenAndServe()
			if err != nil {
				tomb.Kill(err)
			}
		}()

		//nolint:gosimple // ignore this for now
		for {
			select {
			case <-tomb.Dying():
				logrus.Info("Stopping HTTP server...")
				return srv.Shutdown(context.Background())
			}
		}
	})

	// Wait for stuff and watch for errors
	err = tomb.Wait()
	if err != nil {
		return err
	}

	return tomb.Err()
}
