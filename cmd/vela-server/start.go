// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"

	tomb "gopkg.in/tomb.v2"
)

// Start initiates all subprocesses for the Server
// from the provided configuration. The server
// subprocess enables the Server to listen and
// serve traffic for web and API requests.
func (s *Server) Start() error {
	// create the tomb for managing server subprocesses
	//
	// https://pkg.go.dev/gopkg.in/tomb.v2?tab=doc#Tomb
	tomb := new(tomb.Tomb)

	// spawn a tomb goroutine to manage the server subprocesses
	//
	// https://pkg.go.dev/gopkg.in/tomb.v2?tab=doc#Tomb.Go
	tomb.Go(func() error {
		// spawn goroutine for starting the server
		go func() {
			// log a message indicating the start of the server
			//
			// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
			logrus.Info("starting server")

			// start serving traffic for the server
			err := s.serve()
			if err != nil {
				// log the error received from the server
				//
				// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Errorf
				logrus.Errorf("failing server: %v", err)

				// kill the server subprocesses
				//
				// https://pkg.go.dev/gopkg.in/tomb.v2?tab=doc#Tomb.Kill
				tomb.Kill(err)
			}
		}()

		// create an infinite loop to poll for errors
		//
		// nolint: gosimple // ignore this for now
		for {
			// create a select statement to check for errors
			select {
			// check if one of the server subprocesses died
			case <-tomb.Dying():
				// fatally log that we're shutting down the server
				//
				// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Fatal
				logrus.Fatal("shutting down server")

				return tomb.Err()
			}
		}
	})

	// wait for errors from server subprocesses
	//
	// https://pkg.go.dev/gopkg.in/tomb.v2?tab=doc#Tomb.Wait
	return tomb.Wait()
}
