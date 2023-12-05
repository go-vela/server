// SPDX-License-Identifier: Apache-2.0

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
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/api/auth"
	"github.com/go-vela/server/api/webhook"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
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
	r.GET("/login", auth.Login)

	// Logout endpoint
	r.GET("/logout", claims.Establish(), user.Establish(), auth.Logout)

	// Refresh Access Token endpoint
	r.GET("/token-refresh", auth.RefreshAccessToken)

	// Metric endpoint
	r.GET("/metrics", api.CustomMetrics, gin.WrapH(api.BaseMetrics()))

	// Validate Server Token endpoint
	r.GET("/validate-token", claims.Establish(), auth.ValidateServerToken)

	// Validate OAuth Token endpoint
	r.GET("/validate-oauth", claims.Establish(), auth.ValidateOAuthToken)

	// Version endpoint
	r.GET("/version", api.Version)

	// Webhook endpoint
	r.POST("/webhook", webhook.PostWebhook)

	// Authentication endpoints
	authenticate := r.Group("/authenticate")
	{
		authenticate.GET("", auth.GetAuthToken)
		authenticate.GET("/:type", auth.GetAuthRedirect)
		authenticate.GET("/:type/:port", auth.GetAuthRedirect)
		authenticate.POST("/token", auth.PostAuthToken)
	}

	// API endpoints
	baseAPI := r.Group(base, claims.Establish(), user.Establish())
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

		// Schedule endpoints
		ScheduleHandler(baseAPI)

		// Source code management endpoints
		ScmHandlers(baseAPI)

		// Search endpoints
		DashboardHandlers(baseAPI)

		// Secret endpoints
		SecretHandlers(baseAPI)

		// User endpoints
		UserHandlers(baseAPI)

		// Worker endpoints
		WorkerHandlers(baseAPI)

		// Pipeline endpoints
		PipelineHandlers(baseAPI)

		// Queue endpoints
		QueueHandlers(baseAPI)
	} // end of api

	return r
}
