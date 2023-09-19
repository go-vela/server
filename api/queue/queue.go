// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

import (
	"github.com/go-vela/server/router/middleware/claims"
	"net/http"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/queue/queue-registration admin QueueRegistration
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
//     description: Successfully generated registration token
//     schema:
//       "$ref": "#/definitions/Token"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// QueueRegistration represents the API handler to
// provide queue credentials after worker onboarded with registration token.
func QueueRegistration(c *gin.Context) {
	cl := claims.Retrieve(c)
	// retrieve user from context

	logrus.WithFields(logrus.Fields{
		"user": cl.Subject,
	}).Info("requesting queue credentials with registration token")

	// extract the public key that was packed into gin context
	k := c.MustGet("public-key").(*string)

	// extract the queue-address that was packed into gin context
	a := c.MustGet("queue-address").(*string)

	wr := library.QueueRegistration{
		QueuePublicKey: k,
		QueueAddress:   a,
	}

	c.JSON(http.StatusOK, wr)
}
