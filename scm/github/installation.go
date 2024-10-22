// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"time"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ProcessInstallation takes a GitHub installation and processes the changes.
func (c *client) ProcessInstallation(ctx context.Context, request *http.Request, webhook *internal.Webhook, db database.Interface) error {
	c.Logger.Tracef("processing GitHub App installation")

	errs := []error{}

	// if action is "deleted" then the RepositoriesAdded field will indicate the repositories that
	// need to have install_id set to zero

	// set install_id for repos added to the installation
	for _, repo := range webhook.Installation.RepositoriesAdded {
		r, err := db.GetRepoForOrg(ctx, webhook.Installation.Org, repo)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				errs = append(errs, err)
			}

			// skip repos that dont exist in vela
			continue
		}

		installID := webhook.Installation.ID

		// clear install_id if the installation is deleted
		if webhook.Installation.Action == "deleted" {
			installID = 0
		}

		err = updateRepoInstallationID(ctx, webhook, r, db, installID)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// set install_id for repos removed from the installation
	for _, repo := range webhook.Installation.RepositoriesRemoved {
		r, err := db.GetRepoForOrg(ctx, webhook.Installation.Org, repo)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				errs = append(errs, err)
			}

			// skip repos that dont exist in vela
			continue
		}

		err = updateRepoInstallationID(ctx, webhook, r, db, 0)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// combine all errors
	if len(errs) > 0 {
		var combined error
		for _, e := range errs {
			if combined == nil {
				combined = e
			} else {
				combined = errors.Wrap(combined, e.Error())
			}
		}
		return combined
	}

	return nil
}

// updateRepoInstallationID updates the installation ID for a repo.
func updateRepoInstallationID(ctx context.Context, webhook *internal.Webhook, r *types.Repo, db database.Interface, installID int64) error {
	r.SetInstallID(installID)

	h := new(types.Hook)
	h.SetNumber(webhook.Hook.GetNumber())
	h.SetSourceID(webhook.Hook.GetSourceID())
	h.SetWebhookID(webhook.Hook.GetWebhookID())
	h.SetCreated(webhook.Hook.GetCreated())
	h.SetHost(webhook.Hook.GetHost())
	h.SetEvent("installation")
	h.SetStatus(webhook.Hook.GetStatus())

	r, err := db.UpdateRepo(ctx, r)
	if err != nil {
		h.SetStatus(constants.StatusFailure)
		h.SetError(err.Error())
	}

	h.Repo = r

	// number of times to retry
	retryLimit := 3
	// implement a loop to process asynchronous operations with a retry limit
	//
	// Some operations taken during the webhook workflow can lead to race conditions
	// failing to successfully process the request. This logic ensures we attempt our
	// best efforts to handle these cases gracefully.
	for i := 0; i < retryLimit; i++ {
		// check if we're on the first iteration of the loop
		if i > 0 {
			// incrementally sleep in between retries
			time.Sleep(time.Duration(i) * time.Second)
		}

		// send API call to capture the last hook for the repo
		lastHook, err := db.LastHookForRepo(ctx, r)
		if err != nil {
			// log the error for traceability
			logrus.Error(err.Error())

			// check if the retry limit has been exceeded
			if i < retryLimit {
				// continue to the next iteration of the loop
				continue
			}

			return err
		}

		// set the Number field
		if lastHook != nil {
			h.SetNumber(
				lastHook.GetNumber() + 1,
			)
		}

		// send hook update to db
		_, err = db.CreateHook(ctx, h)
		if err != nil {
			return err
		}

		break
	}

	return nil
}
