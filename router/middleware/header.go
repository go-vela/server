// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/version"
)

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Options is a middleware function that appends headers
// for OPTIONS preflight requests and aborts then exits
// the middleware chain and ends the request.
func Options(c *gin.Context) {
	m := c.MustGet("metadata").(*internal.Metadata)

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")

		if len(m.Vela.WebAddress) > 0 {
			c.Header("Access-Control-Allow-Origin", m.Vela.WebAddress)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

// Secure is a middleware function that appends security
// and resource access headers.
func Secure(c *gin.Context) {
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	// TODO: consider adding Content-Security-Policy headers
	// c.Header("Content-Security-Policy", "script-src 'self' https://cdnjs.cloudflare.com")
	c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
}

// Cors is a middleware function that appends headers for
// CORS related requests. These are attached to actual requests
// unlike the OPTIONS preflight requests.
func Cors(c *gin.Context) {
	m := c.MustGet("metadata").(*internal.Metadata)

	c.Header("Access-Control-Allow-Origin", "*")

	if len(m.Vela.WebAddress) > 0 {
		c.Header("Access-Control-Allow-Origin", m.Vela.WebAddress)
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	c.Header("Access-Control-Expose-Headers", "link, x-total-count")
}

// RequestID is a middleware function that generates a unique request ID and
// sets it in both the request and response headers. This is intended for
// debugging and troubleshooting purposes, to help trace requests through
// the system.
func RequestID(c *gin.Context) {
	// on error "id" will be a zero value UUID (UUID.Nil), not nil
	id, _ := uuid.NewV7()

	idString := id.String()

	// set request ID in the request for other middleware to consume
	c.Request.Header.Set("X-Request-ID", idString)

	// set request ID in the response for clients to consume
	c.Header("X-Request-ID", idString)
}

// RequestVersion is a middleware function that injects the Vela API version
// information into the request so it will be logged. This is
// intended for debugging and troubleshooting.
func RequestVersion(c *gin.Context) {
	v := version.New()

	if gin.Mode() == "debug" {
		c.Request.Header.Set("X-Vela-Version", v.Semantic())
	} else { // in prod we don't want the build number metadata
		c.Request.Header.Set("X-Vela-Version", v.Semantic())
	}
}

// ResponseVersion is a middleware function that injects the Vela API version
// information into the response so it will be logged. This is
// intended for debugging and troubleshooting.
func ResponseVersion(c *gin.Context) {
	v := version.New()

	if gin.Mode() == "debug" {
		c.Header("X-Vela-Version", v.Semantic())
	} else { // in prod we don't want the build number metadata
		c.Header("X-Vela-Version", v.Semantic())
	}
}
