// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/admin"
	"github.com/go-vela/server/router/middleware/perm"
)

// QueueHandlers is a function that extends the provided base router group
// with the API handlers for queue registration functionality.
//
// POST   /api/v1/queue-registration.
func QueueHandlers(base *gin.RouterGroup) {
	// Workers endpoints
	_queue := base.Group("/queue-registration")
	{
		_queue.POST("", perm.MustWorkerRegisterToken(), admin.QueueRegistration)
	} // end of queue endpoints
}
