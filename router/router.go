// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
)

const (
	base = "/api/v1"
)

// Load is a server function that returns the engine for processing web requests
// on the host it's running on
func Load(options ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(middleware.RequestVersion)
	r.Use(middleware.NoCache)
	r.Use(middleware.Options)
	r.Use(middleware.Cors)
	r.Use(middleware.Secure)

	r.Use(options...)

	r.GET("/health", api.Health)

	r.GET("/logout", api.Login)
	r.GET("/login", api.Login)
	r.POST("/login", api.Login)
	r.GET("/metrics", api.CustomMetrics, gin.WrapH(api.BaseMetrics()))
	r.GET("/badge/:org/:repo/status.svg", repo.Establish(), api.Badge)

	r.POST("/webhook", api.PostWebhook)

	// Authentication endpoints
	authenticate := r.Group("/authenticate")
	{
		authenticate.GET("", api.Authenticate)
		authenticate.POST("", api.Authenticate)
	}

	// API endpoints
	baseAPI := r.Group(base, user.Establish())
	{
		// Admin endpoints
		AdminHandlers(baseAPI)

		// Hook endpoints
		HookHandlers(baseAPI)

		// Repo endpoints
		// * Build endpoints
		//   * Service endpoints
		//     * Log endpoints
		//   * Step endpoints
		//     * Log endpoints
		RepoHandlers(baseAPI)

		// Secret endpoints
		SecretHandlers(baseAPI)

		// User endpoints
		UserHandlers(baseAPI)
	} // end of api

	return r
}
