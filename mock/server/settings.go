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
			"compiler": {},
			"queue": {},
			"repo_allowlist": [],
			"schedule_allowlist": []
		}`

	// CreateSettingsResp represents a JSON return for creating a settings record.
	CreateSettingsResp = `
		{
			"id": 1,
			"compiler": {},
			"queue": {},
			"repo_allowlist": [],
			"schedule_allowlist": []
		}`
	// UpdateSettingsResp represents a JSON return for modifying a settings field.
	UpdateSettingsResp = `
		{
			"id": 1,
			"compiler": {},
			"queue": {},
			"repo_allowlist": [],
			"schedule_allowlist": []
		}`
	// DeleteSettingsResp represents a JSON return for deleting a settings record.
	DeleteSettingsResp = `
		{
			"id": 1,
			"compiler": {},
			"queue": {},
			"repo_allowlist": [],
			"schedule_allowlist": []
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

// deleteSettings has a param :settings returns mock JSON for a http DELETE.
//
// Pass "0" to :settings to test receiving a http 404 response.
func deleteSettings(c *gin.Context) {
	data := []byte(DeleteSettingsResp)

	var body settings.Platform
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
