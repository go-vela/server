// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const logUpdateInterval = 1 * time.Second

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/service/{service}/stream stream PostServiceStream
//
// Stream the logs for a service
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: path
//   name: service
//   description: Service number
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Payload containing logs
//   required: true
//   schema:
//     type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '204':
//     description: Successfully received logs
//   '400':
//     description: Unable to stream the logs
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to stream the logs
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to stream the logs
//     schema:
//       "$ref": "#/definitions/Error"

// PostServiceStream represents the API handler that
// streams service logs to the database.
//
//nolint:dupl // separate service/step functions for consistency with API
func PostServiceStream(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	b := build.Retrieve(c)
	s := service.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"org":     o,
		"repo":    r.GetName(),
		"build":   b.GetNumber(),
		"service": s.GetNumber(),
		"user":    u.GetName(),
	})

	logger.Infof("streaming logs for service %s/%d", entry, s.GetNumber())

	// send API call to capture the service logs
	_log, err := database.FromContext(c).GetLogForService(s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for service %s/%d: %w", entry, s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	streamLogToDatabase(c, logger, _log, "service", entry, r.GetTimeout())

	c.JSON(http.StatusNoContent, nil)
}

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/stream stream PostStepStream
//
// Stream the logs for a step
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: path
//   name: step
//   description: Step number
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Payload containing logs
//   required: true
//   schema:
//     type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '204':
//     description: Successfully received logs
//   '400':
//     description: Unable to stream the logs
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to stream the logs
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to stream the logs
//     schema:
//       "$ref": "#/definitions/Error"

// PostStepStream represents the API handler that
// streams service logs to the database.
//
//nolint:dupl // separate service/step functions for consistency with API
func PostStepStream(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	b := build.Retrieve(c)
	s := step.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"org":   o,
		"repo":  r.GetName(),
		"build": b.GetNumber(),
		"step":  s.GetNumber(),
		"user":  u.GetName(),
	})

	logger.Infof("streaming logs for step %s/%d", entry, s.GetNumber())

	// send API call to capture the step logs
	_log, err := database.FromContext(c).GetLogForStep(s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for step %s/%d: %w", entry, s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	streamLogToDatabase(c, logger, _log, "step", entry, r.GetTimeout())

	c.JSON(http.StatusNoContent, nil)
}

// streamLogToDatabase handles streaming logs to the database.
func streamLogToDatabase(c *gin.Context, logger *logrus.Entry, log *library.Log, entryType, entry string, timeout int64) {
	// create new buffer for uploading logs
	logs := new(bytes.Buffer)
	// create new channel for processing logs
	done := make(chan bool)
	// defer closing channel to stop processing logs
	defer close(done)

	go func() {
		logger.Debugf("polling request body buffer for %s %s", entryType, entry)

		// spawn "infinite" loop that will upload logs
		// from the buffer until the channel is closed
		for {
			// sleep before attempting to upload logs
			time.Sleep(logUpdateInterval)

			// create a non-blocking select to check if the channel is closed
			select {
			// after repo timeout of idle (no response) end the stream
			//
			// this is a safety mechanism
			case <-time.After(time.Duration(timeout) * time.Minute):
				logger.Tracef("repo timeout of %d exceeded", timeout)

				return
			// channel is closed
			case <-done:
				logger.Tracef("channel closed for polling %s logs", entryType)

				// return out of the go routine
				return
				// channel is not closed
			default:
				// get the current size of log data
				currBytesSize := len(log.GetData())

				// update the existing log with the new bytes if there is new data to add
				if len(logs.Bytes()) > currBytesSize {
					// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.SetData
					log.SetData(logs.Bytes())

					// update the log in the database
					err := database.FromContext(c).UpdateLog(log)
					if err != nil {
						retErr := fmt.Errorf("unable to update logs for %s %s: %w", entryType, entry, err)

						util.HandleError(c, http.StatusInternalServerError, retErr)

						return
					}
				}
			}
		}
	}()

	logger.Debugf("scanning request body for %s %s", entryType, entry)

	scanner := bufio.NewScanner(c.Request.Body)
	for scanner.Scan() {
		// write all the logs from the scanner
		logs.Write(append(scanner.Bytes(), []byte("\n")...))
	}
}
