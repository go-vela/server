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
	"github.com/go-vela/server/secret"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
	"gopkg.in/tomb.v2"
)

//nolint:funlen // ignore function length linter
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
			// sleep 5 seconds before querying build queue
			time.Sleep(5 * time.Second)

			fmt.Println("SERVER WILL BE LISTING THE BUILDS IN QUEUE NOW")

			builds, err := db.ListQueuedBuilds()
			if err != nil {
				logrus.Info("no builds in queue")
				fmt.Println("SERVER HAS FOUND NO BUILDS IN THE QUEUE")

				continue
			}

			for _, b := range builds {
				// define transaction
				logrus.Infof("defining transaction for build %d", b.GetBuildID())
				tx := func(_db *gorm.DB) error {
					// begin transaction
					w, err := db.GetAvailableWorker(_db, b.GetFlavor())
					if err != nil {
						logrus.Infof("no available worker for build %d", b.GetBuildID())

						return err
					}

					//send challenge, listen on /send for an actual build request
					pkg, err := packageBuild(db, secrets, b.GetBuildID(), b.GetPipeline())
					if err != nil {
						logrus.Errorf("unable to package item: %s", err)
						// update build with error
						return err
					}

					logrus.Infof("Sending packaged build to worker: %s", w.GetHostname())

					// TODO: remove hardcoded internal reference debug
					err = sendPackagedBuild("http://localhost:8081", c.String("vela-secret"), pkg)
					if err != nil {
						logrus.Infof("unable to send package to worker %s: %s", w.GetHostname(), err)
						return err
					}

					logrus.Infof("Removing build from the queue: %s", b.GetBuildID())

					err = db.PopQueuedBuild(_db, b.GetBuildID())
					if err != nil {
						logrus.Infof("unable to pop queue item for build %s", b.GetBuildID())
						return err
					}

					logrus.Infof("build popped, packaged, and sent to worker %s, completing transaction", w.GetHostname())
					return nil
				}

				logrus.Infof("performing transaction for build %d", b.GetBuildID())

				// attempt to execute and commit the transaction
				err := db.Transaction(tx)
				if err != nil {
					logrus.Errorf("unable to complete build queue transaction: %s", err)
				}
			}

			// workers, err := db.GetWorkerList()
			// if err != nil {
			// 	logrus.Errorf("unable to get worker list: %s", err)
			// 	continue
			// }

			// w := func() *library.Worker {
			// 	// vader: determine if worker is active and has availability
			// 	for _, worker := range workers {
			// 		if worker.GetActive() {
			// 			return worker
			// 		}
			// 	}

			// 	return nil
			// }()

			// if w == nil {
			// 	logrus.Error("unable to find active worker")
			// 	continue
			// }

			// logrus.Info("Popping item from queue...")

			// // capture an item from the queue
			// item, err := queue.Pop(context.Background())
			// if err != nil {
			// 	logrus.Errorf("Error popping item from queue: %s", err)
			// 	continue
			// }

			// if item == nil {
			// 	logrus.Info("Popped nil item from queue")
			// 	continue
			// }

			// logrus.Info("Popped item from queue")
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

// packageBuild is a helper function that takes in a secret service and a queue item and
// produces a BuildPackage object, which includes all build information (including sensitive content)
// that is sent to the worker.
func packageBuild(db database.Service, secretsServices map[string]secret.Service, buildID int64, data []byte) (*types.BuildPackage, error) {
	build, err := db.GetBuildByID(buildID)
	if err != nil {
		return nil, err
	}

	repo, err := db.GetRepo(build.GetRepoID())
	if err != nil {
		return nil, err
	}

	user, err := db.GetUser(repo.GetUserID())
	if err != nil {
		return nil, err
	}

	pipeline := new(pipeline.Build)

	err = json.Unmarshal(data, pipeline)
	if err != nil {
		return nil, err
	}

	// grab secret information declared in pipeline
	secrets := pipeline.Secrets

	// create BuildPackage object and populate with build information
	buildPackage := new(types.BuildPackage).
		WithBuild(build).
		WithPipeline(pipeline).
		WithRepo(repo).
		WithUser(user).
		WithToken("123abc") // TODO: insert worker-server API token for updating build status

	// iterate through pipeline secrets
	for _, s := range secrets {
		switch s.Type {
		// handle org secrets
		case constants.SecretOrg:
			org, key, err := s.ParseOrg(repo.GetOrg())
			if err != nil {
				return nil, err
			}

			// utilize secrets service to capture the org secret
			_secret, err := secretsServices[s.Engine].Get(s.Type, org, "*", key)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve secret: %w", err)
			}

			// add secret to BuildPackage
			buildPackage.Secrets = append(buildPackage.Secrets, _secret)

		// handle repo secrets
		case constants.SecretRepo:
			org, repo, key, err := s.ParseRepo(repo.GetOrg(), repo.GetName())
			if err != nil {
				return nil, err
			}

			// utilize secrets service to capture the repo secret
			_secret, err := secretsServices[s.Engine].Get(s.Type, org, repo, key)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve secret: %w", err)
			}

			// add secret to BuildPackage
			buildPackage.Secrets = append(buildPackage.Secrets, _secret)

		// handle shared secrets
		case constants.SecretShared:
			org, team, key, err := s.ParseShared()
			if err != nil {
				return nil, err
			}

			// utilize secrets service to capture the shared secret
			_secret, err := secretsServices[s.Engine].Get(s.Type, org, team, key)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve secret: %w", err)
			}

			// add secret to BuildPackage
			buildPackage.Secrets = append(buildPackage.Secrets, _secret)

		default:
			return nil, fmt.Errorf("unrecognized secret type: %s", s.Type)
		}
	}

	return buildPackage, nil
}

// sendPackageBuild is a helper function that takes a worker address and a build package
// and sends the build package to the worker via https.
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
