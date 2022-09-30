// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// MetricsQueryParameters holds query parameter information pertaining to requested metrics.
type MetricsQueryParameters struct {
	// UserCount represents total platform users
	UserCount bool `form:"user_count"`
	// RepoCount represents total platform repos
	RepoCount bool `form:"repo_count"`

	// BuildCount represents total number of builds
	BuildCount bool `form:"build_count"`
	// RunningBuildCount represents total number of builds with status==running
	RunningBuildCount bool `form:"running_build_count"`
	// PendingBuildCount represents total number of builds with status==pending
	PendingBuildCount bool `form:"pending_build_count"`
	// FailureBuildCount represents total number of builds with status==failure
	FailureBuildCount bool `form:"failure_build_count"`
	// KilledBuildCount represents total number of builds with status==killed
	KilledBuildCount bool `form:"killed_build_count"`
	// SuccessBuildCount represents total number of builds with status==success
	SuccessBuildCount bool `form:"success_build_count"`
	// ErrorBuildCount represents total number of builds with status==error
	ErrorBuildCount bool `form:"error_build_count"`

	// StepImageCount represents total number of step images
	StepImageCount bool `form:"step_image_count"`
	// StepStatusCount represents total number of step statuses
	StepStatusCount bool `form:"step_status_count"`
	// ServiceImageCount represents total number of service images
	ServiceImageCount bool `form:"service_image_count"`
	// ServiceStatusCount represents total number of service statuses
	ServiceStatusCount bool `form:"service_status_count"`

	// WorkerBuildLimit represents total worker build limit
	WorkerBuildLimit bool `form:"worker_build_limit"`
	// ActiveWorkerCount represents total number of active workers
	ActiveWorkerCount bool `form:"active_worker_count"`
	// InactiveWorkerCount represents total number of inactive workers
	InactiveWorkerCount bool `form:"inactive_worker_count"`
}

// predefine Prometheus metrics else they will be regenerated
// each function call which will throw error:
// "duplicate metrics collector registration attempted".
var (
	totals = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vela_totals",
			Help: "The Vela totals collect the total number for a resource type.",
		},
		[]string{"resource", "field", "value"},
	)

	stepImages = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "step_images",
			Help: "Step Images collect the number of times an image was used in a step.",
		},
		[]string{"name"},
	)

	serviceImages = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_images",
			Help: "Service Images collect the number of times an image was used in a service.",
		},
		[]string{"name"},
	)
)

// swagger:operation GET /metrics base BaseMetrics
//
// Retrieve metrics from the Vela api
//
// ---
// produces:
// - text/plain
// parameters:
// - in: query
//   name: user_count
//   description: Indicates a request for user count
//   type: boolean
//   default: false
// - in: query
//   name: repo_count
//   description: Indicates a request for repo count
//   type: boolean
//   default: false
// - in: query
//   name: build_count
//   description: Indicates a request for build count
//   type: boolean
//   default: false
// - in: query
//   name: running_build_count
//   description: Indicates a request for running build count
//   type: boolean
//   default: false
// - in: query
//   name: pending_build_count
//   description: Indicates a request for pending build count
//   type: boolean
//   default: false
// - in: query
//   name: failure_build_count
//   description: Indicates a request for failure build count
//   type: boolean
//   default: false
// - in: query
//   name: killed_build_count
//   description: Indicates a request for killed build count
//   type: boolean
//   default: false
// - in: query
//   name: success_build_count
//   description: Indicates a request for success build count
//   type: boolean
//   default: false
// - in: query
//   name: error_build_count
//   description: Indicates a request for error build count
//   type: boolean
//   default: false
// - in: query
//   name: step_image_count
//   description: Indicates a request for step image count
//   type: boolean
//   default: false
// - in: query
//   name: step_status_count
//   description: Indicates a request for step status count
//   type: boolean
//   default: false
// - in: query
//   name: service_image_count
//   description: Indicates a request for service image count
//   type: boolean
//   default: false
// - in: query
//   name: service_status_count
//   description: Indicates a request for service status count
//   type: boolean
//   default: false
// - in: query
//   name: worker_build_limit
//   description: Indicates a request for total worker build limit
//   type: boolean
//   default: false
// - in: query
//   name: active_worker_count
//   description: Indicates a request for active worker count
//   type: boolean
//   default: false
// - in: query
//   name: inactive_worker_count
//   description: Indicates a request for inactive worker count
//   type: boolean
//   default: false
// responses:
//   '200':
//     description: Successfully retrieved the Vela metrics
//     schema:
//       type: string

// BaseMetrics returns a Prometheus handler for serving go metrics.
func BaseMetrics() http.Handler {
	return promhttp.Handler()
}

// CustomMetrics returns custom Prometheus metrics from private functions.
func CustomMetrics(c *gin.Context) {
	// call helper function to return total users
	recordGauges(c)
}

// helper function to get the totals of resource types.
//
//nolint:funlen,gocyclo // ignore function length and cyclomatic complexity
func recordGauges(c *gin.Context) {
	// variable to store query parameters
	q := MetricsQueryParameters{}

	// take incoming request and bind query parameters
	err := c.ShouldBindQuery(&q)
	if err != nil {
		logrus.Errorf("unable to get bind query parameters: %v", err)
	} // continue execution with parameters defaulted to false

	// get each metric separately based on request query parameters
	// user_count
	if q.UserCount {
		// send API call to capture the total number of users
		u, err := database.FromContext(c).CountUsers()
		if err != nil {
			logrus.Errorf("unable to get count of all users: %v", err)
		}
		// add platform metrics
		totals.WithLabelValues("platform", "count", "users").Set(float64(u))
	}

	// repo_count
	if q.RepoCount {
		// send API call to capture the total number of repos
		r, err := database.FromContext(c).CountRepos()
		if err != nil {
			logrus.Errorf("unable to get count of all repos: %v", err)
		}
		// add platform metrics
		totals.WithLabelValues("platform", "count", "repos").Set(float64(r))
	}

	// build_count
	if q.BuildCount {
		// send API call to capture the total number of builds
		b, err := database.FromContext(c).GetBuildCount()
		if err != nil {
			logrus.Errorf("unable to get count of all builds: %v", err)
		}
		// add platform metrics
		totals.WithLabelValues("platform", "count", "builds").Set(float64(b))
	}

	// running_build_count
	if q.RunningBuildCount {
		// send API call to capture the total number of running builds
		bRun, err := database.FromContext(c).GetBuildCountByStatus("running")
		if err != nil {
			logrus.Errorf("unable to get count of all running builds: %v", err)
		}
		// add build metrics
		totals.WithLabelValues("build", "status", "running").Set(float64(bRun))
	}

	// pending_build_count
	if q.PendingBuildCount {
		// send API call to capture the total number of pending builds
		bPen, err := database.FromContext(c).GetBuildCountByStatus("pending")
		if err != nil {
			logrus.Errorf("unable to get count of all pending builds: %v", err)
		}
		// add build metrics
		totals.WithLabelValues("build", "status", "pending").Set(float64(bPen))
	}

	// failure_build_count
	if q.FailureBuildCount {
		// send API call to capture the total number of failure builds
		bFail, err := database.FromContext(c).GetBuildCountByStatus("failure")
		if err != nil {
			logrus.Errorf("unable to get count of all failure builds: %v", err)
		}
		// add build metrics
		totals.WithLabelValues("build", "status", "failed").Set(float64(bFail))
	}

	// killed_build_count
	if q.KilledBuildCount {
		// send API call to capture the total number of killed builds
		bKill, err := database.FromContext(c).GetBuildCountByStatus("killed")
		if err != nil {
			logrus.Errorf("unable to get count of all killed builds: %v", err)
		}
		// add build metrics
		totals.WithLabelValues("build", "status", "killed").Set(float64(bKill))
	}

	// success_build_count
	if q.SuccessBuildCount {
		// send API call to capture the total number of success builds
		bSucc, err := database.FromContext(c).GetBuildCountByStatus("success")
		if err != nil {
			logrus.Errorf("unable to get count of all success builds: %v", err)
		}
		// add build metrics
		totals.WithLabelValues("build", "status", "success").Set(float64(bSucc))
	}

	// error_build_count
	if q.ErrorBuildCount {
		// send API call to capture the total number of error builds
		bErr, err := database.FromContext(c).GetBuildCountByStatus("error")
		if err != nil {
			logrus.Errorf("unable to get count of all error builds: %v", err)
		}
		// add build metrics
		totals.WithLabelValues("build", "status", "error").Set(float64(bErr))
	}

	// step_image_count
	if q.StepImageCount {
		// send API call to capture the total number of step images
		stepImageMap, err := database.FromContext(c).GetStepImageCount()
		if err != nil {
			logrus.Errorf("unable to get count of all step images: %v", err)
		}
		// add step image metrics
		for image, count := range stepImageMap {
			stepImages.WithLabelValues(image).Set(count)
		}
	}

	// step_status_count
	if q.StepStatusCount {
		// send API call to capture the total number of step statuses
		stepStatusMap, err := database.FromContext(c).GetStepStatusCount()
		if err != nil {
			logrus.Errorf("unable to get count of all step statuses: %v", err)
		}
		// add step status metrics
		for status, count := range stepStatusMap {
			totals.WithLabelValues("steps", "status", status).Set(count)
		}
	}

	// service_image_count
	if q.ServiceImageCount {
		// send API call to capture the total number of service images
		serviceImageMap, err := database.FromContext(c).GetServiceImageCount()
		if err != nil {
			logrus.Errorf("unable to get count of all service images: %v", err)
		}
		// add service image metrics
		for image, count := range serviceImageMap {
			serviceImages.WithLabelValues(image).Set(count)
		}
	}

	// service_status_count
	if q.ServiceStatusCount {
		// send API call to capture the total number of service statuses
		serviceStatusMap, err := database.FromContext(c).GetServiceStatusCount()
		if err != nil {
			logrus.Errorf("unable to get count of all service statuses: %v", err)
		}
		// add service status metrics
		for status, count := range serviceStatusMap {
			totals.WithLabelValues("services", "status", status).Set(count)
		}
	}

	// add worker metrics
	var (
		buildLimit      int64
		activeWorkers   int64
		inactiveWorkers int64
	)

	// get worker metrics based on request query parameters
	// worker_build_limit, active_worker_count, inactive_worker_count
	if q.WorkerBuildLimit || q.ActiveWorkerCount || q.InactiveWorkerCount {
		// send API call to capture the workers
		workers, err := database.FromContext(c).ListWorkers()
		if err != nil {
			logrus.Errorf("unable to get workers: %v", err)
		}

		// get the unix time from worker_active_interval ago
		before := time.Now().UTC().Add(-c.Value("worker_active_interval").(time.Duration)).Unix()
		for _, worker := range workers {
			// check if the worker checked in within the last worker_active_interval
			if worker.GetLastCheckedIn() >= before {
				buildLimit += worker.GetBuildLimit()
				activeWorkers++
			} else {
				inactiveWorkers++
			}
		}

		// apply metrics based on request query parameters
		// worker_build_limit
		if q.WorkerBuildLimit {
			totals.WithLabelValues("worker", "sum", "build_limit").Set(float64(buildLimit))
		}

		// active_worker_count
		if q.ActiveWorkerCount {
			totals.WithLabelValues("worker", "count", "active").Set(float64(activeWorkers))
		}

		// inactive_worker_count
		if q.InactiveWorkerCount {
			totals.WithLabelValues("worker", "count", "inactive").Set(float64(inactiveWorkers))
		}
	}
}
