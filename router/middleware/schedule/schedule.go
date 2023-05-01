// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"fmt"
	"github.com/go-vela/server/api/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/util"
)

// Retrieve gets the schedule in the given context.
func Retrieve(c *gin.Context) *types.Schedule {
	return FromContext(c)
}

// Establish used to check if schedule param is used only.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		sParam := util.PathParameter(c, "schedule")
		if len(sParam) == 0 {
			retErr := fmt.Errorf("no schedule parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		ToContext(c, sParam)

		c.Next()
	}
}
