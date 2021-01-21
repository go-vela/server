// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/perm"

	"github.com/gin-gonic/gin"
)

// SecretHandlers is a function that extends the provided base router group
// with the API handlers for secret functionality.
//
// POST   /api/v1/secrets/:engine/:type/:org/:name
// GET    /api/v1/secrets/:engine/:type/:org/:name
// GET    /api/v1/secrets/:engine/:type/:org/:name/:secret
// PUT    /api/v1/secrets/:engine/:type/:org/:name/:secret
// DELETE /api/v1/secrets/:engine/:type/:org/:name/:secret
func SecretHandlers(base *gin.RouterGroup) {
	// Secrets endpoints
	secrets := base.Group("/secrets/:engine/:type/:org/:name", perm.MustSecretAdmin())
	{
		secrets.POST("", api.CreateSecret)
		secrets.GET("", api.GetSecrets)
		secrets.GET("/*secret", api.GetSecret)
		secrets.PUT("/*secret", api.UpdateSecret)
		secrets.DELETE("/*secret", api.DeleteSecret)
	} // end of secrets endpoints
}
