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
//     description: Successful retrieval of the Vela metrics
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
	// return the total number of users in the application
	u, err := database.FromContext(c).GetUserCount()
	if err != nil {
		logrus.Errorf("Error while reading all users: %v", err)
	}

	// return the total number of users in the application
	r, err := database.FromContext(c).GetRepoCount()
	if err != nil {
		logrus.Errorf("Error while reading all repos: %v", err)
	}

	b, err := database.FromContext(c).GetBuildCount()
	if err != nil {
		logrus.Errorf("Error while reading all builds: %v", err)
	}

	bRun, err := database.FromContext(c).GetBuildCountByStatus("running")
	if err != nil {
		logrus.Errorf("Error while reading all running builds: %v", err)
	}

	bPen, err := database.FromContext(c).GetBuildCountByStatus("pending")
	if err != nil {
		logrus.Errorf("Error while reading all pending builds: %v", err)
	}

	bFail, err := database.FromContext(c).GetBuildCountByStatus("failure")
	if err != nil {
		logrus.Errorf("Error while reading all failed builds: %v", err)
	}

	bKill, err := database.FromContext(c).GetBuildCountByStatus("killed")
	if err != nil {
		logrus.Errorf("Error while reading all killed builds: %v", err)
	}

	bSucc, err := database.FromContext(c).GetBuildCountByStatus("success")
	if err != nil {
		logrus.Errorf("Error while reading all success builds: %v", err)
	}

	bErr, err := database.FromContext(c).GetBuildCountByStatus("error")
	if err != nil {
		logrus.Errorf("Error while reading all error builds: %v", err)
	}

	stepStatusMap, err := database.FromContext(c).GetStepStatusCount()
	if err != nil {
		logrus.Errorf("Error while reading all error builds: %v", err)
	}

	stepImageMap, err := database.FromContext(c).GetStepImageCount()
	if err != nil {
		logrus.Errorf("Error while reading all images: %v", err)
	}

	serviceImageMap, err := database.FromContext(c).GetServiceImageCount()
	if err != nil {
		logrus.Errorf("Error while reading all images: %v", err)
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

	// Add step metrics
	for status, count := range stepStatusMap {
		totals.WithLabelValues("steps", "status", status).Set(count)
	}

	// Add image metrics
	for image, count := range stepImageMap {
		stepImages.WithLabelValues(image).Set(count)
	}

	for image, count := range serviceImageMap {
		serviceImages.WithLabelValues(image).Set(count)
	}
}
