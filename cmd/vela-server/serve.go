// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"time"

	"github.com/go-vela/server/router"
	"github.com/go-vela/server/router/middleware"

	"github.com/sirupsen/logrus"
)

// serve is a helper function to listen and serve
// traffic for web and API requests for the Server.
func (s *Server) serve() error {
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
		middleware.Metadata(s.Metadata),
		middleware.Queue(s.Queue),
		middleware.RequestVersion,
		middleware.Secret(s.Config.API.Secret),
		middleware.Secrets(s.Secrets),
		middleware.SecureCookie(s.Config.Security.SecureCookie),
		middleware.Source(s.Source),
		middleware.WebhookValidation(s.Config.Security.WebhookValidation),
		middleware.Worker(s.Config.Metrics.WorkerActive),
	)

	// set the port from the provided server address
	port := s.Config.API.Url.Port()
	// check if a port is part of the server address
	if len(port) == 0 {
		port = s.Config.API.Port
	}

	// log a message indicating the start of serving traffic
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Tracef
	logrus.Tracef("serving traffic on %s", port)

	// serve over http
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Engine.Run
	return _server.Run(fmt.Sprintf(":%s", port))
}
