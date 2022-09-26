// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// Package router Vela server
//
// API for the Vela server
//
//	Version: 0.0.0-dev
//	Schemes: http, https
//	Host: localhost
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	SecurityDefinitions:
//	  ApiKeyAuth:
//	    description: Bearer token
//	    type: apiKey
//	    in: header
//	    name: Authorization
//	  CookieAuth:
//	    description: Refresh token sent as cookie (swagger 2.0 doesn't support cookie auth)
//	    type: apiKey
//	    in: header
//	    name: vela_refresh_token
//
// swagger:meta
package router

import (
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"

	"github.com/gin-gonic/gin"
)

const (
	base = "/api/v1"
)

// Load is a server function that returns the engine for processing web requests
// on the host it's running on.
func Load(options ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(options...)
	r.Use(middleware.RequestVersion)
	r.Use(middleware.NoCache)
	r.Use(middleware.Options)
	r.Use(middleware.Cors)
	r.Use(middleware.Secure)

	// Badge endpoint
	r.GET("/badge/:org/:repo/status.svg", org.Establish(), repo.Establish(), api.GetBadge)

	// Health endpoint
	r.GET("/health", api.Health)

	// Login endpoint
	r.GET("/login", api.Login)

	// Logout endpoint
	r.GET("/logout", user.Establish(), api.Logout)

	// Refresh Access Token endpoint
	r.GET("/token-refresh", api.RefreshAccessToken)

	// Metric endpoint
	r.GET("/metrics", api.CustomMetrics, gin.WrapH(api.BaseMetrics()))

	// Version endpoint
	r.GET("/version", api.Version)

	// Webhook endpoint
	r.POST("/webhook", api.PostWebhook)

	// Authentication endpoints
	authenticate := r.Group("/authenticate")
	{
		authenticate.GET("", api.Authenticate)
		authenticate.GET("/:type", api.AuthenticateType)
		authenticate.GET("/:type/:port", api.AuthenticateType)
		authenticate.POST("/token", api.AuthenticateToken)
	}

	// API endpoints
	baseAPI := r.Group(base, user.Establish())
	{
		// Admin endpoints
		AdminHandlers(baseAPI)

		// Deployment endpoints
		DeploymentHandlers(baseAPI)

		// Hook endpoints
		HookHandlers(baseAPI)

		// Repo endpoints
		// * Build endpoints
		//   * Service endpoints
		//     * Log endpoints
		//   * Step endpoints
		//     * Log endpoints
		RepoHandlers(baseAPI)

		// Source code management endpoints
		ScmHandlers(baseAPI)

		// Search endpoints
		SearchHandlers(baseAPI)

		// Secret endpoints
		SecretHandlers(baseAPI)

		// User endpoints
		UserHandlers(baseAPI)

		// Worker endpoints
		WorkerHandlers(baseAPI)

		// Pipeline endpoints
		PipelineHandlers(baseAPI)
	} // end of api

	return r
}
