// SPDX-License-Identifier: Apache-2.0

package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/router/middleware/auth"
	"github.com/go-vela/types"
)

// rotateKeys returns success message. Pass `invalid` to auth header to test 401 error.
func rotateKeys(c *gin.Context) {
	tkn, _ := auth.RetrieveAccessToken(c.Request)

	if strings.EqualFold(tkn, "invalid") {
		data := "unauthorized"
		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Error{Message: &data})

		return
	}

	c.JSON(http.StatusOK, "keys rotated successfully")
}
