// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(200)
	}
}

// Secure is a middleware function that appends security
// and resource access headers.
func Secure(c *gin.Context) {
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}

	// Also consider adding Content-Security-Policy headers
	// c.Header("Content-Security-Policy", "script-src 'self' https://cdnjs.cloudflare.com")
}

// Cors is a middleware function that appends headers for
// CORS related requests. These are attached to actual requests
// unlike the OPTIONS preflight requests.
func Cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Expose-Headers", "link, x-total-count")
}

// RequestVersion is a middleware function that injects the Vela API version
// information into the request so it will be logged. This is
// intended for debugging and troubleshooting.
func RequestVersion(c *gin.Context) {
	apiVersion := version.Version

	if gin.Mode() == "debug" {
		c.Request.Header.Set("X-Vela-Version", apiVersion.String())
	} else { // in prod we don't want the build number metadata
		apiVersion.Metadata = ""
		c.Request.Header.Set("X-Vela-Version", apiVersion.String())
	}
}

// ResponseVersion is a middleware function that injects the Vela API version
// information into the response so it will be logged. This is
// intended for debugging and troubleshooting.
func ResponseVersion(c *gin.Context) {
	apiVersion := version.Version

	if gin.Mode() == "debug" {
		c.Header("X-Vela-Version", apiVersion.String())
	} else { // in prod we don't want the build number metadata
		apiVersion.Metadata = ""
		c.Header("X-Vela-Version", apiVersion.String())
	}
}
