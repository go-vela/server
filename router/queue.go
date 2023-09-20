// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/queue"
	"github.com/go-vela/server/router/middleware/perm"
)

// QueueHandlers is a function that extends the provided base router group
// with the API handlers for queue registration functionality.
//
// POST   /api/v1/queue/queue-registration.
func QueueHandlers(base *gin.RouterGroup) {
	// Workers endpoints
	_queue := base.Group("/queue")
	{
		_queue.POST("/register", perm.MustWorkerRegisterToken(), queue.Registration)
	} // end of queue endpoints
}
