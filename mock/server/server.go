// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FakeHandler returns an http.Handler that is capable of handling
// Vela API requests and returning mock responses.
// nolint:funlen // number of endpoints is causing linter warning
func FakeHandler() http.Handler {
	gin.SetMode(gin.TestMode)

	e := gin.New()

	// mock endpoints for admin calls
	e.GET("/api/v1/admin/builds", getBuilds)
	e.PUT("/api/v1/admin/build", updateBuild)
	e.GET("/api/v1/admin/builds/queue", buildQueue)
	e.GET("/api/v1/admin/deployments", getDeployments)
	e.PUT("/api/v1/admin/deployment", updateDeployment)
	e.GET("/api/v1/admin/hooks", getHooks)
	e.PUT("/api/v1/admin/hook", updateHook)
	e.GET("/api/v1/admin/repos", getRepos)
	e.PUT("/api/v1/admin/repo", updateRepo)
	e.GET("/api/v1/admin/secrets", getSecrets)
	e.PUT("/api/v1/admin/secret", updateSecret)
	e.GET("/api/v1/admin/services", getServices)
	e.PUT("/api/v1/admin/service", updateService)
	e.GET("/api/v1/admin/steps", getSteps)
	e.PUT("/api/v1/admin/step", updateStep)
	e.GET("/api/v1/admin/users", getUsers)
	e.PUT("/api/v1/admin/user", updateUser)

	// mock endpoints for build calls
	e.GET("/api/v1/repos/:org/:repo/builds/:build", getBuild)
	e.POST("/api/v1/repos/:org/:repo/builds/:build", restartBuild)
	e.DELETE("/api/v1/repos/:org/:repo/builds/:build/cancel", cancelBuild)
	e.GET("/api/v1/repos/:org/:repo/builds/:build/logs", getLogs)
	e.GET("/api/v1/repos/:org/:repo/builds", getBuilds)
	e.POST("/api/v1/repos/:org/:repo/builds", addBuild)
	e.PUT("/api/v1/repos/:org/:repo/builds/:build", updateBuild)
	e.DELETE("/api/v1/repos/:org/:repo/builds/:build", removeBuild)

	// mock endpoints for deployment calls
	e.GET("/api/v1/deployments/:org/:repo", getDeployments)
	e.POST("/api/v1/deployments/:org/:repo", addDeployment)
	e.GET("/api/v1/deployments/:org/:repo/:deployment", getDeployment)

	// mock endpoints for hook calls
	e.GET("/api/v1/hooks/:org/:repo", getHooks)
	e.GET("/api/v1/hooks/:org/:repo/:hook", getHook)
	e.POST("/api/v1/hooks/:org/:repo", addHook)
	e.PUT("/api/v1/hooks/:org/:repo/:hook", updateHook)
	e.DELETE("/api/v1/hooks/:org/:repo/:hook", removeHook)

	// mock endpoints for log calls
	e.GET("/api/v1/repos/:org/:repo/builds/:build/services/:service/logs", getServiceLog)
	e.POST("/api/v1/repos/:org/:repo/builds/:build/services/:service/logs", addServiceLog)
	e.PUT("/api/v1/repos/:org/:repo/builds/:build/services/:service/logs", updateServiceLog)
	e.DELETE("/api/v1/repos/:org/:repo/builds/:build/services/:service/logs", removeServiceLog)
	e.GET("/api/v1/repos/:org/:repo/builds/:build/steps/:step/logs", getStepLog)
	e.POST("/api/v1/repos/:org/:repo/builds/:build/steps/:step/logs", addStepLog)
	e.PUT("/api/v1/repos/:org/:repo/builds/:build/steps/:step/logs", updateStepLog)
	e.DELETE("/api/v1/repos/:org/:repo/builds/:build/steps/:step/logs", removeStepLog)

	// mock endpoints for pipeline calls
	e.GET("/api/v1/pipelines/:org/:repo", getPipeline)
	e.POST("/api/v1/pipelines/:org/:repo/compile", compilePipeline)
	e.POST("/api/v1/pipelines/:org/:repo/expand", expandPipeline)
	e.GET("/api/v1/pipelines/:org/:repo/templates", getTemplates)
	e.POST("/api/v1/pipelines/:org/:repo/validate", validatePipeline)

	// mock endpoints for repo calls
	e.GET("/api/v1/repos/:org/:repo", getRepo)
	e.GET("/api/v1/repos", getRepos)
	e.POST("/api/v1/repos", addRepo)
	e.PUT("/api/v1/repos/:org/:repo", updateRepo)
	e.DELETE("/api/v1/repos/:org/:repo", removeRepo)
	e.PATCH("/api/v1/repos/:org/:repo/repair", repairRepo)
	e.PATCH("/api/v1/repos/:org/:repo/chown", chownRepo)

	// mock endpoints for secret calls
	e.GET("/api/v1/secrets/:engine/:type/:org/:name/:secret", getSecret)
	e.GET("/api/v1/secrets/:engine/:type/:org/:name", getSecrets)
	e.POST("/api/v1/secrets/:engine/:type/:org/:name", addSecret)
	e.PUT("/api/v1/secrets/:engine/:type/:org/:name/:secret", updateSecret)
	e.DELETE("/api/v1/secrets/:engine/:type/:org/:name/:secret", removeSecret)

	// mock endpoints for step calls
	e.GET("/api/v1/repos/:org/:repo/builds/:build/steps/:step", getStep)
	e.GET("/api/v1/repos/:org/:repo/builds/:build/steps", getSteps)
	e.POST("/api/v1/repos/:org/:repo/builds/:build/steps", addStep)
	e.PUT("/api/v1/repos/:org/:repo/builds/:build/steps/:step", updateStep)
	e.DELETE("/api/v1/repos/:org/:repo/builds/:build/steps/:step", removeStep)
	e.POST("/api/v1/repos/:org/:repo/builds/:build/steps/:step/stream", postStepStream)

	// mock endpoints for service calls
	e.GET("/api/v1/repos/:org/:repo/builds/:build/services/:service", getService)
	e.GET("/api/v1/repos/:org/:repo/builds/:build/services", getServices)
	e.POST("/api/v1/repos/:org/:repo/builds/:build/services", addService)
	e.PUT("/api/v1/repos/:org/:repo/builds/:build/services/:service", updateService)
	e.DELETE("/api/v1/repos/:org/:repo/builds/:build/services/:service", removeService)
	e.POST("/api/v1/repos/:org/:repo/builds/:build/services/:service/stream", postServiceStream)

	// mock endpoints for user calls
	e.GET("/api/v1/users/:user", getUser)
	e.GET("/api/v1/users", getUsers)
	e.POST("/api/v1/users", addUser)
	e.PUT("/api/v1/users/:user", updateUser)
	e.DELETE("/api/v1/users/:user", removeUser)

	// mock endpoints for worker calls
	e.GET("/api/v1/workers", getWorkers)
	e.GET("/api/v1/workers/:worker", getWorker)
	e.POST("/api/v1/workers", addWorker)
	e.PUT("/api/v1/workers/:worker", updateWorker)
	e.DELETE("/api/v1/workers/:worker", removeWorker)

	// mock endpoints for authentication calls
	e.GET("/token-refresh", getTokenRefresh)
	e.GET("/authenticate", getAuthenticate)
	e.POST("/authenticate/token", getAuthenticateFromToken)

	return e
}
