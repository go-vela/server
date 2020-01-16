// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// Badge represents the API handler to
// return a build status badge.
func Badge(c *gin.Context) {
	// TODO: allow getting lastbuild by branch and then allow query via `?branch=...`
	// capture middleware values
	r := repo.Retrieve(c)

	logrus.Infof("Creating badge for latest build on %s", r.GetFullName())

	// set default badge
	badge := buildBadge("unknown", "#9f9f9f")

	// send API call to capture the last build for the repo
	b, err := database.FromContext(c).GetLastBuild(r)
	if err != nil {
		c.String(http.StatusOK, badge)
		return
	}

	// set headers to prevent caching
	c.Header("Content-Type", "image/svg+xml")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0") // passing invalid date sets resource as expired

	switch b.GetStatus() {
	case constants.StatusRunning, constants.StatusPending:
		badge = buildBadge("running", "#dfb317")
	case constants.StatusFailure, constants.StatusKilled:
		badge = buildBadge("failed", "#e05d44")
	case constants.StatusSuccess:
		badge = buildBadge("success", "#44cc11")
	case constants.StatusError:
		badge = buildBadge("error", "#fe7d37")
	default:
		c.String(http.StatusOK, badge)
		return
	}

	c.String(http.StatusOK, badge)
}

// buildBadge is a helper that actually builds creates the SVG for the badge
func buildBadge(title, color string) string {
	const (
		t         = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="88" height="20"> <linearGradient id="b" x2="0" y2="100%"> <stop offset="0" stop-color="#bbb" stop-opacity=".1"/> <stop offset="1" stop-opacity=".1"/> </linearGradient> <clipPath id="a"> <rect width="88" height="20" rx="3" fill="#fff"/> </clipPath> <g clip-path="url(#a)"> <path fill="#555" d="M0 0h37v20H0z"/> <path fill="{{ .Color }}" d="M37 0h51v20H37z"/> <path fill="url(#b)" d="M0 0h88v20H0z"/> </g> <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110"> <text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text> <text x="195" y="140" transform="scale(.1)" textLength="270">build</text> <text x="615" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="410">{{ .Title }}</text> <text x="615" y="140" transform="scale(.1)" textLength="410">{{ .Title }}</text> </g> </svg>`
		tFallback = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="88" height="20"> <linearGradient id="b" x2="0" y2="100%"> <stop offset="0" stop-color="#bbb" stop-opacity=".1"/> <stop offset="1" stop-opacity=".1"/> </linearGradient> <clipPath id="a"> <rect width="88" height="20" rx="3" fill="#fff"/> </clipPath> <g clip-path="url(#a)"> <path fill="#555" d="M0 0h37v20H0z"/> <path fill="#9f9f9f" d="M37 0h51v20H37z"/> <path fill="url(#b)" d="M0 0h88v20H0z"/> </g> <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110"> <text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text> <text x="195" y="140" transform="scale(.1)" textLength="270">build</text> <text x="615" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="410">unknown</text> <text x="615" y="140" transform="scale(.1)" textLength="410">unknown</text> </g> </svg>`
	)

	tmpl, err := template.New("StatusBadge").Parse(t)
	if err != nil {
		return tFallback
	}

	buffer := &bytes.Buffer{}

	err = tmpl.Execute(buffer, struct {
		Title string
		Color string
	}{
		title,
		color,
	})
	if err != nil {
		return tFallback
	}

	return buffer.String()
}
