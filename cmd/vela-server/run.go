// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/pkg-queue/queue"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

// run executes the server based
// off the configuration provided.
//
// nolint: funlen // ignore function length due to comments
func run(c *cli.Context) error {
	// set log format for the server
	switch c.String("log.format") {
	case "t", "text", "Text", "TEXT":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "j", "json", "Json", "JSON":
		fallthrough
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// set log level for the server
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
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
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.InfoLevel)
	}

	// create a log entry with extra metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#WithFields
	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/server/",
		"docs":     "https://go-vela.github.io/docs/concepts/infrastructure/server/",
		"registry": "https://hub.docker.com/r/target/vela-server/",
	}).Info("Vela Server")

	// parse the servers address, returning any errors.
	addr, err := url.Parse(c.String("server.addr"))
	if err != nil {
		return fmt.Errorf("unable to parse server address: %w", err)
	}

	// create the server
	s := &Server{
		// server configuration
		Config: &Config{
			Address: c.String("server.addr"),
			Port:    c.String("server.port"),
			Secret:  c.String("server.secret"),
			// api configuration
			API: &API{
				Address: addr,
			},
			// build configuration
			Build: &Build{
				Timeout: c.Int64("build.timeout"),
			},
			// compiler configuration
			Compiler: &Compiler{
				Github: &Github{
					Address: c.String("compiler.github.addr"),
					Token:   c.String("compiler.github.token"),
				},
				Modification: &Modification{
					Address:  c.String("compiler.modification.addr"),
					Secret:   c.String("compiler.modification.secret"),
					Duration: c.Duration("compiler.modification.duration"),
					Retries:  c.Int("compiler.modification.retries"),
				},
			},
			// database configuration
			Database: &Database{
				Driver:           c.String("database.driver"),
				Address:          c.String("database.addr"),
				CompressionLevel: c.Int("database.compression.level"),
				ConnectionIdle:   c.Int("database.connection.idle"),
				ConnectionLife:   c.Duration("database.connection.life"),
				ConnectionOpen:   c.Int("database.connection.open"),
				EncryptionKey:    c.String("database.encryption.key"),
			},
			// logger configuration
			Logger: &Logger{
				Format: c.String("log.format"),
				Level:  c.String("log.level"),
			},
			// metrics configuration
			Metrics: &Metrics{
				WorkerActive: c.Duration("metrics.worker.active.duration"),
			},
			// queue configuration
			Queue: &queue.Setup{
				Driver:  c.String("queue.driver"),
				Config:  c.String("queue.config"),
				Cluster: c.Bool("queue.cluster"),
				Routes:  c.StringSlice("queue.worker.routes"),
			},
			Security: &Security{
				AccessToken:       c.Duration("access.token.duration"),
				RefreshToken:      c.Duration("refresh.token.duration"),
				RepoAllowList:     c.StringSlice("repo.allow.list"),
				SecureCookie:      c.Bool("secure.cookie"),
				WebhookValidation: c.Bool("webhook.validation"),
			},
			Source: &Source{
				Driver:       c.String("source.driver"),
				Address:      c.String("source.addr"),
				ClientID:     c.String("source.client"),
				ClientSecret: c.String("source.secret"),
				Context:      c.String("source.context"),
			},
			WebUI: &WebUI{
				Address:       c.String("webui.addr"),
				OAuthEndpoint: c.String("webui.oauth.endpoint"),
			},
		},
	}

	// set the server address if no flag was provided
	if len(w.Config.API.Address.String()) == 0 {
		s.Config.API.Address, _ = url.Parse(fmt.Sprintf("http://%s", hostname))
	}

	// validate the server
	err = s.Validate()
	if err != nil {
		return err
	}

	// start the server
	return s.Start()
}
