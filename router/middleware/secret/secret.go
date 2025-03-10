// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
)

// Retrieve gets the hook in the given context.
func Retrieve(c *gin.Context) *api.Secret {
	return FromContext(c)
}
