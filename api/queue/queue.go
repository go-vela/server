// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/queue/info queue Info
//
// Get queue credentials
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved queue credentials
//     schema:
//       "$ref": "#/definitions/QueueInfo"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// Info represents the API handler to
// retrieve queue credentials as part of worker onboarding.
func Info(c *gin.Context) {
	cl := claims.Retrieve(c)

	logrus.WithFields(logrus.Fields{
		"user": cl.Subject,
	}).Info("requesting queue credentials with registration token")

	// extract the public key that was packed into gin context
	k := c.MustGet("public-key").(string)

	// extract the queue-address that was packed into gin context
	a := c.MustGet("queue-address").(string)

	wr := library.QueueInfo{
		QueuePublicKey: &k,
		QueueAddress:   &a,
	}

	c.JSON(http.StatusOK, wr)
}
