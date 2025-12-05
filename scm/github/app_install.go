// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
)

// ProcessInstallation takes a GitHub installation and processes the changes.
func (c *Client) ProcessInstallation(ctx context.Context, _ *http.Request, webhook *internal.Webhook, db database.Interface) error {
	c.Logger.Tracef("processing GitHub App installation")

	errs := []string{}

	// set install_id for repos added to the installation
	for _, repo := range webhook.Installation.RepositoriesAdded {
		r, err := db.GetRepoForOrg(ctx, webhook.Installation.Org, repo)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				errs = append(errs, fmt.Sprintf("%s:%s", repo, err.Error()))
			}

			// skip repos that dont exist in vela
			continue
		}

		err = updateRepoInstallationID(ctx, webhook, r, db, webhook.Installation.ID)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s:%s", repo, err.Error()))
		}
	}

	// set install_id for repos removed from the installation
	for _, repo := range webhook.Installation.RepositoriesRemoved {
		r, err := db.GetRepoForOrg(ctx, webhook.Installation.Org, repo)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				errs = append(errs, fmt.Sprintf("%s:%s", repo, err.Error()))
			}

			// skip repos that dont exist in vela
			continue
		}

		err = updateRepoInstallationID(ctx, webhook, r, db, 0)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s:%s", repo, err.Error()))
		}
	}

	// combine all errors
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
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
	h.SetEvent(constants.EventInstallation)
	h.SetStatus(webhook.Hook.GetStatus())

	r, err := db.UpdateRepo(ctx, r)
	if err != nil {
		h.SetStatus(constants.StatusFailure)
		h.SetError(err.Error())
	}

	h.Repo = r

	// send hook update to db
	_, err = db.CreateHook(ctx, h)
	if err != nil {
		return err
	}

	return nil
}

// FinishInstallation completes the web flow for a GitHub App installation, returning a redirect to the app installation page.
func (c *Client) FinishInstallation(ctx context.Context, _ *http.Request, installID int64) (string, error) {
	c.Logger.Tracef("finishing GitHub App installation for ID %d", installID)

	client, err := c.newGithubAppClient()
	if err != nil {
		return "", err
	}

	install, _, err := client.Apps.GetInstallation(ctx, installID)
	if err != nil {
		return "", err
	}

	return install.GetHTMLURL(), nil
}
