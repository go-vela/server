// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// postServiceStream returns a nock response for an http POST.
func postServiceStream(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// postStepStream returns a nock response for an http POST.
func postStepStream(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}
