// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	gauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vela_totals",
			Help: "The Vela totals collect the total number for a resource type.",
		},
		[]string{"type"},
	)
)

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

	gauge.WithLabelValues("users").Set(float64(u))
	gauge.WithLabelValues("repos").Set(float64(r))
	gauge.WithLabelValues("builds").Set(float64(b))
	gauge.WithLabelValues("running_builds").Set(float64(bRun))
	gauge.WithLabelValues("pending_builds").Set(float64(bPen))
}
