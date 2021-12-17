// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"time"

	"github.com/go-vela/server/router/middleware/org"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/sirupsen/logrus"
)

// Logger returns a gin.HandlerFunc (middleware) that logs requests using logrus.
//
// Requests with errors are logged using logrus.Error().
// Requests without errors are logged using logrus.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
func Logger(logger *logrus.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		// prevent us from logging the health endpoint
		if c.Request.URL.Path != "/health" {
			fields := logrus.Fields{
				"ip":         c.ClientIP(),
				"latency":    latency,
				"method":     c.Request.Method,
				"path":       path,
				"status":     c.Writer.Status(),
				"user-agent": c.Request.UserAgent(),
				"version":    c.GetHeader("X-Vela-Version"),
			}

			body := c.Value("payload")
			if body != nil {
				fields["body"] = body
			}

			build := build.Retrieve(c)
			if build != nil {
				fields["build"] = build.Number
			}

			org := org.Retrieve(c)
			if org != "" {
				fields["org"] = org
			}

			repo := repo.Retrieve(c)
			if repo != nil {
				fields["repo"] = repo.Name
			}

			service := service.Retrieve(c)
			if service != nil {
				fields["service"] = service.Number
			}

			step := step.Retrieve(c)
			if step != nil {
				fields["step"] = step.Number
			}

			user := user.Retrieve(c)
			if user != nil {
				fields["user"] = user.Name
			}

			worker := worker.Retrieve(c)
			if worker != nil {
				fields["worker"] = worker.Hostname
			}

			entry := logger.WithFields(fields)

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				entry.Error(c.Errors.String())
			} else {
				entry.Info()
			}
		}
	}
}
