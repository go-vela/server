// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/queue"
	"github.com/go-vela/server/router/middleware/perm"
)

// QueueHandlers is a function that extends the provided base router group
// with the API handlers for queue registration functionality.
//
// POST   /api/v1/queue/register.
func QueueHandlers(base *gin.RouterGroup) {
	// Queue endpoints
	_queue := base.Group("/queue")
	{
		_queue.GET("/info", perm.MustWorkerRegisterToken(), queue.Info)
	} // end of queue endpoints
}
