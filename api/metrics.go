// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

// swagger:operation GET /metrics router BaseMetrics
//
// Retrieve metrics from the  Vela api
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
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
// nolint: funlen // ignore function length due to comments
func recordGauges(c *gin.Context) {
	// send API call to capture the total number of users
	u, err := database.FromContext(c).GetUserCount()
	if err != nil {
		logrus.Errorf("unable to get count of all users: %v", err)
	}

	// send API call to capture the total number of repos
	r, err := database.FromContext(c).GetRepoCount()
	if err != nil {
		logrus.Errorf("unable to get count of all repos: %v", err)
	}

	// send API call to capture the total number of builds
	b, err := database.FromContext(c).GetBuildCount()
	if err != nil {
		logrus.Errorf("unable to get count of all builds: %v", err)
	}

	// send API call to capture the total number of running builds
	bRun, err := database.FromContext(c).GetBuildCountByStatus("running")
	if err != nil {
		logrus.Errorf("unable to get count of all running builds: %v", err)
	}

	// send API call to capture the total number of pending builds
	bPen, err := database.FromContext(c).GetBuildCountByStatus("pending")
	if err != nil {
		logrus.Errorf("unable to get count of all pending builds: %v", err)
	}

	// send API call to capture the total number of failure builds
	bFail, err := database.FromContext(c).GetBuildCountByStatus("failure")
	if err != nil {
		logrus.Errorf("unable to get count of all failure builds: %v", err)
	}

	// send API call to capture the total number of killed builds
	bKill, err := database.FromContext(c).GetBuildCountByStatus("killed")
	if err != nil {
		logrus.Errorf("unable to get count of all killed builds: %v", err)
	}

	// send API call to capture the total number of success builds
	bSucc, err := database.FromContext(c).GetBuildCountByStatus("success")
	if err != nil {
		logrus.Errorf("unable to get count of all success builds: %v", err)
	}

	// send API call to capture the total number of error builds
	bErr, err := database.FromContext(c).GetBuildCountByStatus("error")
	if err != nil {
		logrus.Errorf("unable to get count of all error builds: %v", err)
	}

	// send API call to capture the total number of step images
	stepImageMap, err := database.FromContext(c).GetStepImageCount()
	if err != nil {
		logrus.Errorf("unable to get count of all step images: %v", err)
	}

	// send API call to capture the total number of step statuses
	stepStatusMap, err := database.FromContext(c).GetStepStatusCount()
	if err != nil {
		logrus.Errorf("unable to get count of all step statuses: %v", err)
	}

	// send API call to capture the total number of service images
	serviceImageMap, err := database.FromContext(c).GetServiceImageCount()
	if err != nil {
		logrus.Errorf("unable to get count of all service images: %v", err)
	}

	// send API call to capture the total number of service statuses
	serviceStatusMap, err := database.FromContext(c).GetServiceStatusCount()
	if err != nil {
		logrus.Errorf("unable to get count of all service statuses: %v", err)
	}

	// send API call to capture the workers
	workers, err := database.FromContext(c).GetWorkerList()
	if err != nil {
		logrus.Errorf("unable to get workers: %v", err)
	}

	// Add platform metrics
	totals.WithLabelValues("platform", "count", "users").Set(float64(u))
	totals.WithLabelValues("platform", "count", "repos").Set(float64(r))
	totals.WithLabelValues("platform", "count", "builds").Set(float64(b))

	// Add build metrics
	totals.WithLabelValues("build", "status", "running").Set(float64(bRun))
	totals.WithLabelValues("build", "status", "pending").Set(float64(bPen))
	totals.WithLabelValues("build", "status", "failed").Set(float64(bFail))
	totals.WithLabelValues("build", "status", "killed").Set(float64(bKill))
	totals.WithLabelValues("build", "status", "success").Set(float64(bSucc))
	totals.WithLabelValues("build", "status", "error").Set(float64(bErr))

	// Add worker metrics
	var buildLimit int64
	var activeWorkers int64
	var inactiveWorkers int64
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

	totals.WithLabelValues("worker", "sum", "build_limit").Set(float64(buildLimit))
	totals.WithLabelValues("worker", "count", "active").Set(float64(activeWorkers))
	totals.WithLabelValues("worker", "count", "inactive").Set(float64(inactiveWorkers))

	// Add step status metrics
	for status, count := range stepStatusMap {
		totals.WithLabelValues("steps", "status", status).Set(count)
	}

	// Add service status metrics
	for status, count := range serviceStatusMap {
		totals.WithLabelValues("services", "status", status).Set(count)
	}

	// Add step image metrics
	for image, count := range stepImageMap {
		stepImages.WithLabelValues(image).Set(count)
	}

	// Add service image metrics
	for image, count := range serviceImageMap {
		serviceImages.WithLabelValues(image).Set(count)
	}
}
