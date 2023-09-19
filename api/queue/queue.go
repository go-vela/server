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
	k := c.MustGet("public-key")

	// extract the queue-address that was packed into gin context
	a := c.MustGet("queue-address")

	pk, ok := k.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, "public key in the context is the wrong type")
		return
	}

	qa, ok := a.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, "queue address in the context is the wrong type")
		return
	}

	wr := library.QueueRegistration{
		QueuePublicKey: &pk,
		QueueAddress:   &qa,
	}

	c.JSON(http.StatusOK, wr)
}
