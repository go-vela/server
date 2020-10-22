// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// Package router Vela server
//
// API for the Vela server
//
//     Version: 0.6.1
//     Schemes: http, https
//     Host: localhost
//     BasePath: /api/v1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     SecurityDefinitions:
//       ApiKeyAuth:
//         type: apiKey
//         in: header
//         name: Authorization
//
// swagger:meta
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

	// Badge endpoint
	r.GET("/badge/:org/:repo/status.svg", repo.Establish(), api.GetBadge)

	// Health endpoint
	r.GET("/health", api.Health)

	// Login endpoints
	r.GET("/login", api.Login)
	r.POST("/login", api.Login)

	// Logout endpoint
	r.GET("/logout", api.Login)

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
		authenticate.POST("", api.Authenticate)
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

		// Secret endpoints
		SecretHandlers(baseAPI)

		// User endpoints
		UserHandlers(baseAPI)

		// Worker endpoints
		WorkerHandlers(baseAPI)

		// todo rename expand
		pipelines := baseAPI.Group("pipelines/:org/:repo", repo.Establish())
		{
			pipelines.GET("", api.GetPipeline)                // return the pipeline for the org/repo (raw YAML)
			pipelines.POST("/compile", api.Health)            // compile and return the pipeline for the org/repo (fully rendered YAML - includes env vars, extra vars etc)
			pipelines.POST("/expand", api.ExpandPipeline)     // expand the templates for the pipeline and return the YAML
			pipelines.GET("/templates", api.GetTemplates)     // return the templates from the pipeline
			pipelines.POST("/validate", api.ValidatePipeline) // validate the pipeline
		}
	} // end of api

	return r
}
