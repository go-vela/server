// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/types/settings"
)

const (
	// SettingsResp represents a JSON return for a single settings.
	SettingsResp = `
		{
			"id": 1,
			"compiler": {
				"clone_image": "target/vela-git",
				"template_depth": 3,
				"starlark_exec_limit": 100
			},
			"queue": {
				"routes": [
					"vela"
				]
			},
			"repo_allowlist": [
				"*"
			],
			"schedule_allowlist": [
				"octocat/hello-world"
			]
		}`

	// UpdateSettingsResp represents a JSON return for modifying a settings field.
	UpdateSettingsResp = `
		{
			"id": 1,
			"compiler": {
				"clone_image": "target/vela-git:latest",
				"template_depth": 5,
				"starlark_exec_limit": 123
			},
			"queue": {
				"routes": [
					"vela",
					"large"
				]
			},
			"repo_allowlist": [],
			"schedule_allowlist": [
				"octocat/hello-world",
				"octocat/*"
			]
		}`
)

// getSettings has a param :settings returns mock JSON for a http GET.
func getSettings(c *gin.Context) {
	data := []byte(SettingsResp)

	var body settings.Platform
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// updateSettings returns mock JSON for a http PUT.
func updateSettings(c *gin.Context) {
	data := []byte(UpdateSettingsResp)

	var body settings.Platform
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
