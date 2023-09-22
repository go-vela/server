// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// This file, in part, reproduces portions of
// https://github.com/elastic/ecs-logging-go-logrus/blob/v1.0.0/formatter.go
// to handle our custom fields in Format().

const (
	// ecsVersion holds the version of ECS with which the formatter is compatible.
	ecsVersion = "1.6.0"
)

var (
	ecsFieldMap = logrus.FieldMap{
		logrus.FieldKeyTime:  "@timestamp",
		logrus.FieldKeyMsg:   "message",
		logrus.FieldKeyLevel: "log.level",
	}
)

// Formatter holds ECS parameter information for logging.
type Formatter struct {
	// DataKey allows users to put all the log entry parameters into a
	// nested dictionary at a given key.
	//
	// DataKey is ignored for well-defined fields, such as "error",
	// which will instead be stored under the appropriate ECS fields.
	DataKey string
}

// Logger returns a gin.HandlerFunc (middleware) that logs requests using logrus.
//
// Requests with errors are logged using logrus.Error().
// Requests without errors are logged using logrus.Info().
//
// It receives:
//  1. A time package format string (e.g. time.RFC3339).
func Logger(logger *logrus.Logger, timeFormat string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := util.EscapeValue(c.Request.URL.Path)

		c.Next()

		end := time.Now()

		latency := end.Sub(start)

		// prevent us from logging the health endpoint
		if c.Request.URL.Path != "/health" {
			fields := logrus.Fields{
				"ip":         util.EscapeValue(c.ClientIP()),
				"latency":    latency,
				"method":     c.Request.Method,
				"path":       path,
				"status":     c.Writer.Status(),
				"user-agent": util.EscapeValue(c.Request.UserAgent()),
				"version":    util.EscapeValue(c.GetHeader("X-Vela-Version")),
			}

			body := c.Value("payload")
			if body != nil {
				body = sanitize(body)
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
				entry.Infof("%v %v %v %s %s", fields["status"], fields["latency"], fields["ip"], fields["method"], fields["path"])
			}
		}
	}
}

func sanitize(body interface{}) interface{} {
	if m, ok := body.(map[string]interface{}); ok {
		if _, ok = m["email"]; ok {
			m["email"] = constants.SecretMask
			body = m
		}
	}

	return body
}

// Format formats e as ECS-compliant JSON,
// mapping our custom fields to ECS fields.
func (f *Formatter) Format(e *logrus.Entry) ([]byte, error) {
	datahint := len(e.Data)
	if f.DataKey != "" {
		datahint = 2
	}
	data := make(logrus.Fields, datahint)
	if len(e.Data) > 0 {
		extraData := data
		if f.DataKey != "" {
			extraData = make(logrus.Fields, len(e.Data))
		}
		for k, v := range e.Data {
			switch k {
			case "ip":
				data["client.ip"] = v
			case "latency":
				data["event.duration"] = v
			case "method":
				data["http.request.method"] = v
			case "path":
				data["url.path"] = v
			case "status":
				data["http.response.status_code"] = v
			case "user-agent":
				data["user_agent.name"] = v
			case "version":
				data["user_agent.version"] = v
			default:
				extraData[k] = v
			}
		}
		if f.DataKey != "" && len(extraData) > 0 {
			data[f.DataKey] = extraData
		}
	}
	if e.HasCaller() {
		// Logrus has a single configurable field (logrus.FieldKeyFile)
		// for storing a combined filename and line number, but we want
		// to split them apart into two fields. Remove the event's Caller
		// field, and encode the ECS fields explicitly.
		var funcVal, fileVal string
		var lineVal int

		funcVal = e.Caller.Function
		fileVal = e.Caller.File
		lineVal = e.Caller.Line

		e.Caller = nil
		if funcVal != "" {
			data["log.origin.function"] = funcVal
		}
		if fileVal != "" {
			data["log.origin.file.name"] = fileVal
		}
		if lineVal > 0 {
			data["log.origin.file.line"] = lineVal
		}
	}
	data["ecs.version"] = ecsVersion
	ecopy := *e
	ecopy.Data = data
	e = &ecopy

	jf := logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z0700",
		FieldMap:        ecsFieldMap,
	}
	return jf.Format(e)
}
