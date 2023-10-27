// SPDX-License-Identifier: Apache-2.0

package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types"
)

// syncRepo has a param :repo returns mock JSON for a http PATCH.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func syncRepo(c *gin.Context) {
	r := c.Param("repo")
	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Repo %s has been synced", r))
}

// syncRepos has a param :org returns mock JSON for a http PATCH.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func syncRepos(c *gin.Context) {
	o := c.Param("org")
	if strings.Contains(o, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", o)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Org %s repos have been synced", o))
}
