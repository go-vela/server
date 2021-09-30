// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const logUpdateInterval = 1 * time.Second

// PostServiceStream represents the API handler that
// streams service logs to the database
func PostServiceStream(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	logrus.Infof("streaming logs for service %d for build %s", s.GetNumber(), entry)

	// create new buffer for uploading logs
	logs := new(bytes.Buffer)
	// create new channel for processing logs
	done := make(chan bool)
	// defer closing channel to stop processing logs
	defer close(done)

	// send API call to capture the service logs
	_log, err := database.FromContext(c).GetServiceLog(s.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for service %s/%d: %w", entry, s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	go func() {
		logrus.Debugf("polling request body buffer for service %d for build %s", s.GetNumber(), entry)

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
			case <-time.After(time.Duration(r.GetTimeout()) * time.Minute):
				logrus.Tracef("repo timeout of %d exceeded", r.GetTimeout())

				return
			// channel is closed
			case <-done:
				logrus.Trace("channel closed for polling container logs")

				// return out of the go routine
				return
			// channel is not closed
			default:
				// get the current size of log data
				currBytesSize := len(_log.GetData())

				// update the existing log with the new bytes if there is new data to add
				if len(logs.Bytes()) > currBytesSize {
					// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.SetData
					_log.SetData(logs.Bytes())

					// update the log in the database
					err = database.FromContext(c).UpdateLog(_log)
					if err != nil {
						retErr := fmt.Errorf("unable to update logs for service %s/%d: %w", entry, s.GetNumber(), err)

						util.HandleError(c, http.StatusInternalServerError, retErr)

						return
					}
				}
			}
		}
	}()

	logrus.Debugf("scanning request body for service %d for build %s", s.GetNumber(), entry)

	scanner := bufio.NewScanner(c.Request.Body)
	for scanner.Scan() {
		// write all the logs from the scanner
		logs.Write(append(scanner.Bytes(), []byte("\n")...))
	}

	c.JSON(http.StatusNoContent, nil)
}

// PostStepStream represents the API handler that
// streams service logs to the database
func PostStepStream(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	logrus.Infof("streaming logs for step %d for build %s", s.GetNumber(), entry)

	// create new buffer for uploading logs
	logs := new(bytes.Buffer)
	// create new channel for processing logs
	done := make(chan bool)
	// defer closing channel to stop processing logs
	defer close(done)

	// send API call to capture the step logs
	_log, err := database.FromContext(c).GetStepLog(s.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for step %s/%d: %w", entry, s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	go func() {
		logrus.Debugf("polling request body buffer for step %d for build %s", s.GetNumber(), entry)

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
			case <-time.After(time.Duration(r.GetTimeout()) * time.Minute):
				logrus.Tracef("repo timeout of %d exceeded", r.GetTimeout())

				return
			// channel is closed
			case <-done:
				logrus.Trace("channel closed for polling container logs")

				// return out of the go routine
				return
				// channel is not closed
			default:
				// get the current size of log data
				currBytesSize := len(_log.GetData())

				// update the existing log with the new bytes if there is new data to add
				if len(logs.Bytes()) > currBytesSize {
					// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Log.SetData
					_log.SetData(logs.Bytes())

					// update the log in the database
					err = database.FromContext(c).UpdateLog(_log)
					if err != nil {
						retErr := fmt.Errorf("unable to update logs for step %s/%d: %w", entry, s.GetNumber(), err)

						util.HandleError(c, http.StatusInternalServerError, retErr)

						return
					}
				}
			}
		}
	}()

	logrus.Debugf("scanning request body for step %d for build %s", s.GetNumber(), entry)

	scanner := bufio.NewScanner(c.Request.Body)
	for scanner.Scan() {
		// write all the logs from the scanner
		logs.Write(append(scanner.Bytes(), []byte("\n")...))
	}

	c.JSON(http.StatusNoContent, nil)
}
