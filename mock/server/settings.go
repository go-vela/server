// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
)

const (
	// SettingsResp represents a JSON return for a single settings.
	SettingsResp = `
		{
			"id": 1
		}`

	// CreateSettingsResp represents a JSON return for creating a settings record.
	CreateSettingsResp = `
		{
			"id": 1
		}`
	// UpdateSettingsResp represents a JSON return for modifying a settings field.
	UpdateSettingsResp = `
		{
			"id": 1
		}`
	// RemoveSettingsResp represents a JSON return for deleting a settings record.
	RemoveSettingsResp = `
		{
			"id": 1
		}`
)

// getSettings has a param :settings returns mock JSON for a http GET.
func getSettings(c *gin.Context) {
	data := []byte(SettingsResp)

	var body api.Settings
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// createSettings returns mock JSON for a http POST.
func createSettings(c *gin.Context) {
	data := []byte(CreateSettingsResp)

	var body api.Settings
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateSettings returns mock JSON for a http PUT.
func updateSettings(c *gin.Context) {
	data := []byte(UpdateSettingsResp)

	var body api.Settings
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeSettings has a param :settings returns mock JSON for a http DELETE.
//
// Pass "0" to :settings to test receiving a http 404 response.
func removeSettings(c *gin.Context) {
	data := []byte(RemoveSettingsResp)

	var body api.Settings
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
