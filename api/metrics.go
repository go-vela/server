// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// predefine Prometheus metrics else they will be regenerated
// each function call which will throw error:
// "duplicate metrics collector registration attempted"
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

// BaseMetrics returns a Prometheus handler for serving go metrics
func BaseMetrics() http.Handler {
	return promhttp.Handler()
}

// CustomMetrics returns custom Prometheus metrics from private functions
func CustomMetrics(c *gin.Context) {
	// call helper function to return total users
	recordGauges(c)
}

// helper function to get the totals of resource types
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
