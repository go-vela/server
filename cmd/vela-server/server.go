// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/tracing"
)

//nolint:funlen,gocyclo // ignore function length and cyclomatic complexity
func server(ctx context.Context, cmd *cli.Command) error {
	// set log formatter
	switch cmd.String("log-formatter") {
	case "json":
		// set logrus to log in JSON format
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "ecs":
		// set logrus to log in Elasticsearch Common Schema (ecs) format
		logrus.SetFormatter(&middleware.ECSFormatter{
			DataKey: "labels.vela",
		})
	}

	// set log level for logrus
	switch cmd.String("log-level") {
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

	compiler, err := native.FromCLICommand(ctx, cmd)
	if err != nil {
		return err
	}

	tc, err := tracing.FromCLICommand(ctx, cmd)
	if err != nil {
		return err
	}

	if tc.EnableTracing {
		defer func() {
			err := tc.TracerProvider.Shutdown(context.Background())
			if err != nil {
				logrus.Errorf("unable to shutdown tracer provider: %v", err)
			}
		}()
	}

	database, err := database.FromCLICommand(cmd, tc)
	if err != nil {
		return err
	}

	queue, err := queue.FromCLICommand(cmd)
	if err != nil {
		return err
	}

	secrets, err := setupSecrets(cmd, database)
	if err != nil {
		return err
	}

	scm, err := setupSCM(ctx, cmd, tc)
	if err != nil {
		return err
	}

	metadata, err := setupMetadata(cmd)
	if err != nil {
		return err
	}

	tm, err := setupTokenManager(ctx, cmd, database)
	if err != nil {
		return err
	}

	// determine issuer for metadata and token manager
	oidcIssuer := cmd.String("oidc-issuer")
	if len(oidcIssuer) == 0 {
		oidcIssuer = fmt.Sprintf("%s/_services/token", cmd.String("server-addr"))
	}

	metadata.Vela.OpenIDIssuer = oidcIssuer
	tm.Issuer = oidcIssuer

	jitter := wait.Jitter(5*time.Second, 2.0)

	logrus.Infof("retrieving initial platform settings after %v delay", jitter)

	time.Sleep(jitter)

	ps, err := database.GetSettings(context.Background())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// platform settings record does not exist
	if err != nil {
		logrus.Info("creating initial platform settings")

		// create initial settings record
		ps = settings.FromCLICommand(cmd)

		// singleton record ID should always be 1
		ps.SetID(1)

		ps.SetCreatedAt(time.Now().UTC().Unix())
		ps.SetUpdatedAt(time.Now().UTC().Unix())
		ps.SetUpdatedBy("vela-server")

		// read in defaults supplied from the cli runtime
		compilerSettings := compiler.GetSettings()
		ps.SetCompiler(compilerSettings)

		queueSettings := queue.GetSettings()
		ps.SetQueue(queueSettings)

		// create the settings record in the database
		_, err = database.CreateSettings(context.Background(), ps)
		if err != nil {
			return err
		}

		logrus.Info("initial platform settings created")
	}

	// update any internal settings, this occurs in middleware
	// to keep settings refreshed for each request
	queue.SetSettings(ps)
	compiler.SetSettings(ps)

	router := router.Load(
		middleware.AppWebhookSecret(cmd.String("scm.app.webhook-secret")),
		middleware.CLI(cmd),
		middleware.Settings(ps),
		middleware.Compiler(compiler),
		middleware.Database(database),
		middleware.Logger(logrus.StandardLogger(), time.RFC3339),
		middleware.Metadata(metadata),
		middleware.TokenManager(tm),
		middleware.Queue(queue),
		middleware.RequestVersion,
		middleware.Secret(cmd.String("vela-secret")),
		middleware.Secrets(secrets),
		middleware.Scm(scm),
		middleware.QueueSigningPrivateKey(cmd.String("queue.private-key")),
		middleware.QueueSigningPublicKey(cmd.String("queue.public-key")),
		middleware.QueueAddress(cmd.String("queue.addr")),
		middleware.DefaultBuildLimit(int(cmd.Int("default-build-limit"))),
		middleware.DefaultTimeout(int(cmd.Int("default-build-timeout"))),
		middleware.DefaultApprovalTimeout(int(cmd.Int("default-approval-timeout"))),
		middleware.MaxBuildLimit(int(cmd.Int("max-build-limit"))),
		middleware.WebhookValidation(!cmd.Bool("vela-disable-webhook-validation")),
		middleware.SecureCookie(cmd.Bool("vela-enable-secure-cookie")),
		middleware.Worker(cmd.Duration("worker-active-interval")),
		middleware.DefaultRepoEvents(cmd.StringSlice("default-repo-events")),
		middleware.DefaultRepoEventsMask(cmd.Int("default-repo-events-mask")),
		middleware.DefaultRepoApproveBuild(cmd.String("default-repo-approve-build")),
		middleware.ScheduleFrequency(cmd.Duration("schedule-minimum-frequency")),
		middleware.TracingClient(tc),
		middleware.TracingInstrumentation(tc),
	)

	addr, err := url.Parse(cmd.String("server-addr"))
	if err != nil {
		return err
	}

	port := addr.Port()
	// check if a port is part of the address
	if len(port) == 0 {
		port = cmd.String("server-port")
	}

	// gin expects the address to be ":<port>" ie ":8080"
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadHeaderTimeout: 60 * time.Second,
	}

	// create the context for controlling the worker subprocesses
	ctx, done := context.WithCancel(context.Background())
	// create the errgroup for managing worker subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group
	g, gctx := errgroup.WithContext(ctx)

	// spawn goroutine to check for signals to gracefully shutdown
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			logrus.Infof("received signal: %s", sig)

			err := srv.Shutdown(ctx)
			if err != nil {
				logrus.Error(err)
			}

			done()
		case <-gctx.Done():
			logrus.Info("closing signal goroutine")

			err := srv.Shutdown(ctx)
			if err != nil {
				logrus.Error(err)
			}

			return gctx.Err()
		}

		return nil
	})

	// spawn goroutine for refreshing settings
	g.Go(func() error {
		interval := cmd.Duration("settings-refresh-interval")

		logrus.Infof("refreshing platform settings every %v", interval)

		for {
			time.Sleep(interval)

			newSettings, err := database.GetSettings(context.Background())
			if err != nil {
				logrus.WithError(err).Warn("unable to refresh platform settings")

				continue
			}

			// update the internal fields for the shared settings record
			ps.FromSettings(newSettings)
		}
	})

	// spawn goroutine for starting the server
	g.Go(func() error {
		logrus.Infof("starting server on %s", addr.Host)

		err = srv.ListenAndServe()
		if err != nil {
			// log a message indicating the failure of the server
			logrus.Errorf("failing server: %v", err)
		}

		return err
	})

	// spawn go routine for cleaning up pending approval builds
	g.Go(func() error {
		logrus.Info("starting pending approval cleanup routine")

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			err := cleanupPendingApproval(ctx, database)
			if err != nil {
				logrus.WithError(err).Warn("unable to cleanup pending approval builds")
			}

			<-ticker.C
		}
	})

	// spawn goroutine for starting the scheduler
	g.Go(func() error {
		logrus.Info("starting scheduler")

		for {
			// track the starting time for when the server begins processing schedules
			//
			// This will be used to control which schedules will have a build triggered based
			// off the configured entry and last time a build was triggered for the schedule.
			start := time.Now().UTC()

			// capture the interval of time to wait before processing schedules
			//
			// We need to sleep for some amount of time before we attempt to process schedules
			// setup in the database. Since the schedule interval is configurable, we use that
			// as the base duration to determine how long to sleep for.
			interval := cmd.Duration("schedule-interval")

			// This should prevent multiple servers from processing schedules at the same time by
			// leveraging a base duration along with a standard deviation of randomness a.k.a.
			// "jitter". To create the jitter, we use the configured schedule interval duration
			// along with a scale factor of 0.5.
			jitter := wait.Jitter(interval, 0.5)

			logrus.Infof("sleeping for %v before scheduling builds", jitter)
			// sleep for a duration of time before processing schedules
			time.Sleep(jitter)

			// update internal settings updated through refresh
			compiler.SetSettings(ps)
			queue.SetSettings(ps)

			err = processSchedules(ctx, start, ps, compiler, database, metadata, queue, scm)
			if err != nil {
				logrus.WithError(err).Warn("unable to process schedules")
			} else {
				logrus.Trace("successfully processed schedules")
			}
		}
	})

	// wait for errors from server subprocesses
	return g.Wait()
}
