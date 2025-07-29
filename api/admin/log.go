// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
)

// swagger:operation DELETE /api/v1/admin/log/cleanup admin AdminCleanLogs
//
// Delete old log records in batches
//
// ---
// produces:
// - application/json
// parameters:
// - in: query
//   name: before
//   description: Unix timestamp - delete logs created before this time
//   required: true
//   type: integer
// - in: query
//   name: batch_size
//   description: Number of records to delete per batch (default 1000, max 10000)
//   required: false
//   type: integer
// - in: query
//   name: vacuum
//   description: Whether to run VACUUM after deletion (default false)
//   required: false
//   type: boolean
// - in: body
//   name: body
//   description: Optional message for logging purposes
//   required: false
//   schema:
//     "$ref": "#/definitions/Error"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully cleaned log records
//     schema:
//       "$ref": "#/definitions/LogCleanupResponse"
//   '400':
//     description: Invalid request parameters
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// CleanLogs represents the API handler to delete old log records in batches.
func CleanLogs(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Info("platform admin: cleaning log records")

	startTime := time.Now()

	// capture and validate before query parameter (required)
	beforeStr := c.Query("before")
	if beforeStr == "" {
		retErr := fmt.Errorf("before query parameter is required")
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	before, err := strconv.ParseInt(beforeStr, 10, 64)
	if err != nil {
		retErr := fmt.Errorf("unable to convert before query parameter %s to int64: %w", beforeStr, err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// validate that before timestamp is not too recent (must be at least 24 hours ago)
	oneDayAgo := time.Now().Add(-24 * time.Hour).Unix()
	if before > oneDayAgo {
		retErr := fmt.Errorf("before timestamp must be at least 24 hours ago (provided: %d, minimum: %d)", before, oneDayAgo)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture and validate batch_size query parameter (optional, default 1000, max 10000)
	batchSize := 1000
	if batchSizeStr := c.Query("batch_size"); batchSizeStr != "" {
		batchSize, err = strconv.Atoi(batchSizeStr)
		if err != nil {
			retErr := fmt.Errorf("unable to convert batch_size query parameter %s to int: %w", batchSizeStr, err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		if batchSize < 1 || batchSize > 10000 {
			retErr := fmt.Errorf("batch_size must be between 1 and 10000 (provided: %d)", batchSize)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	}

	// capture vacuum query parameter (optional, default false)
	withVacuum := false
	if vacuumStr := c.Query("vacuum"); vacuumStr != "" {
		withVacuum, err = strconv.ParseBool(vacuumStr)
		if err != nil {
			retErr := fmt.Errorf("unable to convert vacuum query parameter %s to bool: %w", vacuumStr, err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	}

	// capture optional message from request body
	message := "Log cleanup completed by platform admin"
	input := new(types.Error)

	// only attempt to bind if there's actually a request body
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBind(input); err == nil && input.Message != nil {
			message = util.EscapeValue(*input.Message)
		}
	}

	l.WithFields(logrus.Fields{
		"before":      before,
		"batch_size":  batchSize,
		"with_vacuum": withVacuum,
	}).Info("starting log cleanup operation")

	// get database driver for vacuum operations
	driver := database.FromContext(c).Driver()

	// perform the cleanup operation
	result, err := database.FromContext(c).CleanLogs(ctx, before, batchSize, withVacuum, driver)
	if err != nil {
		retErr := fmt.Errorf("unable to clean logs: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	duration := time.Since(startTime)
	batchesProcessed := int64(0)

	if batchSize > 0 {
		batchesProcessed = (result.DeletedCount + int64(batchSize) - 1) / int64(batchSize) // ceiling division
	}

	l.WithFields(logrus.Fields{
		"deleted_count":     result.DeletedCount,
		"batches_processed": batchesProcessed,
		"duration_seconds":  duration.Seconds(),
		"vacuum_performed":  withVacuum,
	}).Info("log cleanup operation completed")

	// return cleanup statistics
	// check if partitioned mode is being used by calling the database engine's partition check method
	partitionedMode := database.FromContext(c).IsLogPartitioned()

	response := types.LogCleanupResponse{
		DeletedCount:       result.DeletedCount,
		BatchesProcessed:   batchesProcessed,
		DurationSeconds:    duration.Seconds(),
		VacuumPerformed:    withVacuum && result.DeletedCount > 0,
		PartitionedMode:    partitionedMode,
		AffectedPartitions: result.AffectedPartitions,
		Message:            fmt.Sprintf("%s. Deleted %d log records in %d batches", message, result.DeletedCount, batchesProcessed),
	}

	c.JSON(http.StatusOK, response)
}
