// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

import (
	"net/http"

	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/admin/queue-credentials admin QueueRegistration
//
// Get queue credentials
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: worker
//   description: Hostname of the worker
//   required: true
//   type: string
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
	// retrieve user from context
	u := user.Retrieve(c)

	logrus.Infof("Platform admin %s: fetching queue creds", u.GetName())

	// extract the public key that was packed into gin context
	k, ok := c.Get("public-key")
	if !ok {
		c.JSON(http.StatusInternalServerError, "no public-key in the context")
		return
	}
	// extract the queue-address that was packed into gin context
	a, ok := c.Get("queue-address")
	if !ok {
		c.JSON(http.StatusInternalServerError, "no queue-address in the context")
		return
	}

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
	logrus.Infof("current creds %s: ", wr)
	c.JSON(http.StatusOK, wr)
}
