// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/token"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// RefreshAccessToken will return a new access token if the provided
// refresh token via cookie is valid.
func RefreshAccessToken(c *gin.Context) {
	// capture the refresh token
	// TODO: move this into token package and do it internally
	// since we are already passsing context
	rt, err := token.RetrieveRefreshToken(c.Request)
	if err != nil {
		retErr := fmt.Errorf("refresh token error: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// validate the refresh token and return a new access token
	newAccessToken, err := token.Refresh(c, rt)
	if err != nil {
		retErr := fmt.Errorf("unable to refresh token: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Login{Token: &newAccessToken})
}
