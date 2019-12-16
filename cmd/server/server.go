// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"net/http"
	"time"

	"github.com/go-vela/server/router"
	"github.com/go-vela/server/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
	tomb "gopkg.in/tomb.v2"
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

	source, err := setupSource(c)
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
		middleware.Logger(logrus.StandardLogger(), time.RFC3339, true),
		middleware.Metadata(metadata),
		middleware.Queue(queue),
		middleware.RequestVersion,
		middleware.Secret(c.String("vela-secret")),
		middleware.Secrets(secrets),
		middleware.Source(source),
		middleware.Whitelist(c.StringSlice("vela-repo-whitelist")),
	)

	var tomb tomb.Tomb
	// start http server
	tomb.Go(func() error {
		srv := &http.Server{Addr: c.String("server-port"), Handler: router}

		go func() {
			logrus.Info("Starting HTTP server...")
			err := srv.ListenAndServe()
			if err != nil {
				tomb.Kill(err)
			}
		}()

		for {
			select {
			case <-tomb.Dying():
				logrus.Info("Stopping HTTP server...")
				return srv.Shutdown(nil)
			}
		}
	})

	// Wait for stuff and watch for errors
	tomb.Wait()
	return tomb.Err()
}
