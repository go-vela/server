// SPDX-License-Identifier: Apache-2.0

package executors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// Retrieve gets the executors in the given context.
func Retrieve(c *gin.Context) []library.Executor {
	return FromContext(c)
}

// Establish sets the executors in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		e := new([]library.Executor)
		b := build.Retrieve(c)
		ctx := c.Request.Context()

		// if build is pending or pending approval, there is no host to establish executors
		if strings.EqualFold(b.GetStatus(), constants.StatusPending) ||
			strings.EqualFold(b.GetStatus(), constants.StatusPendingApproval) ||
			len(b.GetHost()) == 0 {
			ToContext(c, *e)
			c.Next()

			return
		}

		// retrieve the worker
		w, err := database.FromContext(c).GetWorkerForHostname(ctx, b.GetHost())
		if err != nil {
			retErr := fmt.Errorf("unable to get worker: %w", err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// prepare the request to the worker to retrieve executors
		client := http.DefaultClient
		client.Timeout = 30 * time.Second
		endpoint := fmt.Sprintf("%s/api/v1/executors", w.GetAddress())

		req, err := http.NewRequestWithContext(context.Background(), "GET", endpoint, nil)
		if err != nil {
			retErr := fmt.Errorf("unable to form request to %s: %w", endpoint, err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		tm := c.MustGet("token-manager").(*token.Manager)

		// set mint token options
		mto := &token.MintTokenOpts{
			Hostname:      "vela-server",
			TokenType:     constants.WorkerAuthTokenType,
			TokenDuration: time.Minute * 1,
		}

		// mint token
		tkn, err := tm.MintToken(mto)
		if err != nil {
			retErr := fmt.Errorf("unable to generate auth token: %w", err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// add the token to authenticate to the worker
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

		// make the request to the worker and check the response
		resp, err := client.Do(req)
		if err != nil || resp == nil {
			// abandoned builds might have ran on a worker that no longer exists
			// if the worker is unavailable write an empty slice ToContext
			ToContext(c, *e)
			c.Next()
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			retErr := fmt.Errorf("unable to read response from %s: %w", endpoint, err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// parse response and validate at least one item was returned
		err = json.Unmarshal(respBody, e)
		if err != nil {
			retErr := fmt.Errorf("unable to parse response from %s: %w", endpoint, err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		ToContext(c, *e)
		c.Next()
	}
}
