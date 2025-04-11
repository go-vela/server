// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/dashboard"
	"github.com/go-vela/server/router/middleware/hook"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/schedule"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/server/util"
)

// This file, in part, reproduces portions of
// https://github.com/elastic/ecs-logging-go-logrus/blob/v1.0.0/formatter.go
// to handle our custom fields in Format().

// ECSFormatter holds ECS parameter information for logging.
type ECSFormatter struct {
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
func Logger(logger *logrus.Logger, _ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := util.EscapeValue(c.Request.URL.Path)

		fields := logrus.Fields{
			"ip":   util.EscapeValue(c.ClientIP()),
			"path": path,
		}

		entry := logger.WithFields(fields)

		// set the logger in the context so
		// downstream handlers can use it
		c.Set("logger", entry)

		c.Next()

		latency := time.Since(start)

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
				fields["build_id"] = build.ID
			}

			org := org.Retrieve(c)
			if org != "" {
				fields["org"] = org
			}

			pipeline := pipeline.Retrieve(c)
			if pipeline != nil {
				fields["pipeline_id"] = pipeline.ID
			}

			repo := repo.Retrieve(c)
			if repo != nil {
				fields["repo"] = repo.Name
				fields["repo_id"] = repo.ID
			}

			service := service.Retrieve(c)
			if service != nil {
				fields["service"] = service.Number
				fields["service_id"] = service.ID
			}

			hook := hook.Retrieve(c)
			if hook != nil {
				fields["hook"] = hook.Number
				fields["hook_id"] = hook.ID
			}

			step := step.Retrieve(c)
			if step != nil {
				fields["step"] = step.Number
				fields["step_id"] = step.ID
			}

			schedule := schedule.Retrieve(c)
			if schedule != nil {
				fields["schedule"] = schedule.Name
				fields["schedule_id"] = schedule.ID
			}

			dashboard := dashboard.Retrieve(c)
			if dashboard != nil {
				fields["dashboard"] = dashboard.Name
				fields["dashboard_id"] = dashboard.ID
			}

			user := user.Retrieve(c)
			// we check to make sure user name is populated
			// because when it's not a user token, we still
			// inject an empty user object into the context
			// which results in log entries with 'user: null'
			if user != nil && user.GetName() != "" {
				fields["user"] = user.Name
				fields["user_id"] = user.ID
			}

			worker := worker.Retrieve(c)
			if worker != nil {
				fields["worker"] = worker.Hostname
				fields["worker_id"] = worker.ID
			}

			// if there's no user or worker in the context
			// of this request, we log claims subject
			_, hasUser := fields["user"]
			_, hasWorker := fields["worker"]

			if !hasUser && !hasWorker {
				claims := claims.Retrieve(c)
				if claims != nil {
					fields["claims_subject"] = claims.Subject
				}
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

func sanitize(body any) any {
	if m, ok := body.(map[string]any); ok {
		if _, ok = m["email"]; ok {
			m["email"] = constants.SecretMask
			body = m
		}
	}

	return body
}

// Format formats logrus.Entry as ECS-compliant JSON,
// mapping our custom fields to ECS fields.
func (f *ECSFormatter) Format(e *logrus.Entry) ([]byte, error) {
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
			// map fields attached to requests
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

			// map other fields
			case "user":
				data["user.name"] = v

			default:
				extraData[k] = v
			}
		}

		if f.DataKey != "" && len(extraData) > 0 {
			data[f.DataKey] = extraData
		}
	}

	// ecsVersion holds the version of ECS with which the formatter is compatible.
	data["ecs.version"] = "1.6.0"
	ecopy := *e
	ecopy.Data = data
	e = &ecopy

	ecsFieldMap := logrus.FieldMap{
		logrus.FieldKeyTime:  "@timestamp",
		logrus.FieldKeyMsg:   "message",
		logrus.FieldKeyLevel: "log.level",
	}

	jf := logrus.JSONFormatter{
		TimestampFormat: time.RFC3339, // same as default in logrus
		FieldMap:        ecsFieldMap,
	}

	return jf.Format(e)
}
