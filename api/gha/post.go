// SPDX-License-Identifier: Apache-2.0

package gha

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

var baseErr = "unable to process webhook"

// TODO: Add swagger docs

// PostWebhook represents the API handler to capture
// a webhook from a source control provider and
// publish it to the configure queue.
//
//nolint:funlen,gocyclo // ignore function length and cyclomatic complexity
func PostWebhook(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("webhook received")

	// duplicate request so we can perform operations on the request body
	//
	// https://golang.org/pkg/net/http/#Request.Clone
	dupRequest := c.Request.Clone(ctx)

	// -------------------- Start of TODO: --------------------
	//
	// Remove the below code once http.Request.Clone()
	// actually performs a deep clone.
	//
	// This code is required due to a known bug:
	//
	// * https://github.com/golang/go/issues/36095

	// create buffer for reading request body
	var buf bytes.Buffer

	// read the request body for duplication
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		retErr := fmt.Errorf("unable to read webhook body: %w", err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// add the request body to the original request
	c.Request.Body = io.NopCloser(&buf)

	// add the request body to the duplicate request
	dupRequest.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	//
	// -------------------- End of TODO: --------------------

	// TODO: Implement
	// process the webhook from the source control provider
	err = scm.FromContext(c).ProcessGitHubAppWebhook(ctx, c.Request)
	if err != nil {
		retErr := fmt.Errorf("unable to parse webhook: %w", err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}
}
