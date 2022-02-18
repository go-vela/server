// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"github.com/go-vela/server/database"
	"github.com/go-vela/types/library"

	"github.com/go-vela/server/util"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the worker in the given context.
func Retrieve(c *gin.Context) *library.Worker {
	return FromContext(c)
}

// Establish sets the worker in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		wParam := c.Param("worker")
		if len(wParam) == 0 {
			retErr := fmt.Errorf("no worker parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		logrus.Debugf("Reading worker %s", wParam)

		w, err := database.FromContext(c).GetWorker(wParam)
		if err != nil {
			retErr := fmt.Errorf("unable to read worker %s: %v", wParam, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		ToContext(c, w)
		c.Next()
	}
}
