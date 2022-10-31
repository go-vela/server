// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
	"gopkg.in/tomb.v2"
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

	db, err := setupDatabase(c)
	if err != nil {
		return err
	}

	queue, err := setupQueue(c)
	if err != nil {
		return err
	}

	secrets, err := setupSecrets(c, db)
	if err != nil {
		return err
	}

	scm, err := setupSCM(c)
	if err != nil {
		return err
	}

	metadata, err := setupMetadata(c)
	if err != nil {
		return err
	}

	router := router.Load(
		middleware.Compiler(compiler),
		middleware.Database(db),
		middleware.Logger(logrus.StandardLogger(), time.RFC3339),
		middleware.Metadata(metadata),
		middleware.Queue(queue),
		middleware.RequestVersion,
		middleware.Secret(c.String("vela-secret")),
		middleware.Secrets(secrets),
		middleware.Scm(scm),
		middleware.Allowlist(c.StringSlice("vela-repo-allowlist")),
		middleware.DefaultBuildLimit(c.Int64("default-build-limit")),
		middleware.DefaultTimeout(c.Int64("default-build-timeout")),
		middleware.MaxBuildLimit(c.Int64("max-build-limit")),
		middleware.WebhookValidation(!c.Bool("vela-disable-webhook-validation")),
		middleware.SecureCookie(c.Bool("vela-enable-secure-cookie")),
		middleware.Worker(c.Duration("worker-active-interval")),
	)

	addr, err := url.Parse(c.String("server-addr"))
	if err != nil {
		return err
	}

	// vader: thread for listening to the queue
	go func() {
		for {
			workers, err := db.GetWorkerList()
			if err != nil {
				logrus.Infof("unable to get worker list: %w", err)
				continue
			}

			w := func() *library.Worker {
				// vader: determine if worker is active and has availability
				for _, worker := range workers {
					if worker.GetActive() {
						return worker
					}
				}
				return nil
			}()

			if w == nil {
				logrus.Info("unable to find active worker")
				continue
			}

			logrus.Info("Popping item from queue...")

			// capture an item from the queue
			item, err := queue.Pop(context.Background())
			if err != nil {
				logrus.Infof("Error popping item from queue: %w", err)
				continue
			}

			if item == nil {
				logrus.Info("Popped nil item from queue")
				continue
			}
			logrus.Info("Popped item from queue")

			//send challenge, listen on /send for an actual build request

			pkg, err := packageBuild(db, item)
			if err != nil {
				logrus.Infof("Error packaging item: %w", err)
				continue
			}

			logrus.Infof("Sending packaged build to worker: %s", w.GetHostname())

			// TODO: remove hardcoded internal reference debug
			err = sendPackagedBuild("http://host.docker.internal:8081", c.String("vela-secret"), pkg)
			if err != nil {
				logrus.Infof("unable to send package to worker %s: %w", w.GetHostname(), err)
				continue
			}
		}
	}()

	var tomb tomb.Tomb
	// start http server
	tomb.Go(func() error {
		port := addr.Port()

		// check if a port is part of the address
		if len(port) == 0 {
			port = c.String("server-port")
		}

		// gin expects the address to be ":<port>" ie ":8080"
		srv := &http.Server{
			Addr:              fmt.Sprintf(":%s", port),
			Handler:           router,
			ReadHeaderTimeout: 60 * time.Second,
		}

		logrus.Infof("running server on %s", addr.Host)
		go func() {
			logrus.Info("Starting HTTP server...")
			err := srv.ListenAndServe()
			if err != nil {
				tomb.Kill(err)
			}
		}()

		//nolint:gosimple // ignore this for now
		for {
			select {
			case <-tomb.Dying():
				logrus.Info("Stopping HTTP server...")
				return srv.Shutdown(context.Background())
			}
		}
	})

	// Wait for stuff and watch for errors
	err = tomb.Wait()
	if err != nil {
		return err
	}

	return tomb.Err()
}

func packageBuild(db database.Service, item *types.Item) (*types.BuildPackage, error) {
	secrets := item.Pipeline.Secrets
	buildPackage := new(types.BuildPackage).
		WithBuild(item.Build).
		WithPipeline(item.Pipeline).
		WithRepo(item.Repo).
		WithUser(item.User).
		WithToken("123abc")

	for _, s := range secrets {
		switch s.Type {
		// handle org secrets
		case constants.SecretOrg:
			org, key, err := s.ParseOrg(item.Repo.GetOrg())
			if err != nil {
				return nil, err
			}

			// send API call to capture the org secret
			//
			// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#SecretService.Get
			_secret, err := db.GetSecret(s.Type, org, "*", key)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve secret: %w", err)
			}

			buildPackage.Secrets = append(buildPackage.Secrets, _secret)

		// handle repo secrets
		case constants.SecretRepo:
			org, repo, key, err := s.ParseRepo(item.Repo.GetOrg(), item.Repo.GetName())
			if err != nil {
				return nil, err
			}

			// send API call to capture the repo secret
			//
			// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#SecretService.Get
			_secret, err := db.GetSecret(s.Type, org, repo, key)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve secret: %w", err)
			}

			buildPackage.Secrets = append(buildPackage.Secrets, _secret)

		// handle shared secrets
		case constants.SecretShared:
			org, team, key, err := s.ParseShared()
			if err != nil {
				return nil, err
			}

			// send API call to capture the repo secret
			//
			// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#SecretService.Get
			_secret, err := db.GetSecret(s.Type, org, team, key)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve secret: %w", err)
			}

			buildPackage.Secrets = append(buildPackage.Secrets, _secret)

		default:
			return nil, fmt.Errorf("unrecognized secret type: %s", s.Type)
		}
	}

	return buildPackage, nil
}

func sendPackagedBuild(workerAddress, secret string, data interface{}) error {
	// prepare the request to the worker
	client := http.DefaultClient
	client.Timeout = 30 * time.Second

	// set the API endpoint path we send the request to
	u := fmt.Sprintf("%s/api/v1/exec", workerAddress)
	it, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", u, bytes.NewBuffer(it))
	if err != nil {
		return err
	}

	// add the token to authenticate to the worker
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))

	// perform the request to the worker
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO: some kind of logging?
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
