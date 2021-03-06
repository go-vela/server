// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-vela/compiler/compiler"
	"github.com/go-vela/pkg-queue/queue"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/source"

	"github.com/sirupsen/logrus"
)

type (
	// API represents the server configuration for API information.
	API struct {
		Address *url.URL
	}

	// Build represents the server configuration for build information.
	Build struct {
		Timeout int64
	}

	// Github represents the compiler configuration for the github information.
	Github struct {
		Address string
		Token   string
	}

	// Modification represents the compiler configuration for the modification information.
	Modification struct {
		Address  string
		Secret   string
		Duration time.Duration
		Retries  int
	}

	// Compiler represents the server configuration for compiler information.
	Compiler struct {
		Github       *Github
		Modification *Modification
	}

	// Database represents the server configuration for database information.
	Database struct {
		Driver           string
		Address          string
		CompressionLevel int
		ConnectionIdle   int
		ConnectionLife   time.Duration
		ConnectionOpen   int
		EncryptionKey    string
	}

	// Logger represents the server configuration for logger information.
	Logger struct {
		Format string
		Level  string
	}

	// Metrics represents the server configuration for metrics information.
	Metrics struct {
		WorkerActive time.Duration
	}

	// Security represents the server configuration for security information.
	Security struct {
		AccessToken       time.Duration
		RefreshToken      time.Duration
		RepoAllowList     []string
		SecureCookie      bool
		WebhookValidation bool
	}

	// Source represents the server configuration for source information.
	Source struct {
		Driver       string
		Address      string
		ClientID     string
		ClientSecret string
		Context      string
	}

	// WebUI represents the server configuration for web UI information.
	WebUI struct {
		Address       string
		OAuthEndpoint string
	}

	// Config represents the server configuration.
	Config struct {
		Address  string
		Port     string
		Secret   string
		API      *API
		Build    *Build
		Compiler *Compiler
		Database *Database
		Logger   *Logger
		Metrics  *Metrics
		Queue    *queue.Setup
		Security *Security
		Source   *Source
		WebUI    *WebUI
	}

	// Server represents all configuration and
	// system processes for the server.
	Server struct {
		Config   *Config
		Compiler compiler.Engine
		Database database.Service
		Queue    queue.Service
		Source   source.Service
	}
)

// server is a helper function to listen and serve
// traffic for web and API requests for the Server.
func (s *Server) server() error {
	// log a message indicating the setup of the server handlers
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Trace
	logrus.Trace("loading router with server handlers")

	// create the server router to listen and serve traffic
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router?tab=doc#Load
	_server := router.Load(
		middleware.Allowlist(s.Config.Security.RepoAllowList),
		middleware.Compiler(s.Compiler),
		middleware.Database(s.Database),
		middleware.DefaultTimeout(s.Config.Build.Timeout),
		middleware.Logger(logrus.StandardLogger(), time.RFC3339, true),
		middleware.Metadata(metadata),
		middleware.Queue(s.Queue),
		middleware.RequestVersion,
		middleware.Secret(s.Config.Secret),
		middleware.Secrets(secrets),
		middleware.SecureCookie(s.Config.Security.SecureCookie),
		middleware.Source(s.Source),
		middleware.WebhookValidation(s.Config.Security.WebhookValidation),
		middleware.Worker(s.Config.Metrics.WorkerActive),
	)

	// set the port from the provided server address
	port := s.Config.API.Address.Port()
	// check if a port is part of the server address
	if len(port) == 0 {
		port = s.Config.Port
	}

	// log a message indicating the start of serving traffic
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Tracef
	logrus.Tracef("serving traffic on %s", port)

	// else serve over http
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Engine.Run
	return _server.Run(fmt.Sprintf(":%s", port))
}
