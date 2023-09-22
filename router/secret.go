// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/go-vela/server/api/secret"
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
// DELETE /api/v1/secrets/:engine/:type/:org/:name/:secret .
func SecretHandlers(base *gin.RouterGroup) {
	// Secrets endpoints
	secrets := base.Group("/secrets/:engine/:type/:org/:name", perm.MustSecretAdmin())
	{
		secrets.POST("", secret.CreateSecret)
		secrets.GET("", secret.ListSecrets)
		secrets.GET("/*secret", secret.GetSecret)
		secrets.PUT("/*secret", secret.UpdateSecret)
		secrets.DELETE("/*secret", secret.DeleteSecret)
	} // end of secrets endpoints
}
