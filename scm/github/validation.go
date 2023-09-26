// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v54/github"
)

func (c *client) ValidateOAuthToken(ctx context.Context, token string) error {
	// create http client to connect to GitHub API
	transport := github.BasicAuthTransport{
		Username: c.config.ClientID,
		Password: c.config.ClientSecret,
	}
	// create client to connect to GitHub API
	client := github.NewClient(transport.Client())
	// check if github url was set
	if c.config.Address != "" && c.config.Address != "https://github.com" {
		// check if address has trailing slash
		if !strings.HasSuffix(c.config.Address, "/") {
			// add trailing slash
			c.config.Address = c.config.Address + "/api/v3/"
		}
		// parse the provided url into url type
		enterpriseURL, err := url.Parse(c.config.Address)
		if err != nil {
			return err
		}
		// set the base and upload url
		client.BaseURL = enterpriseURL
		client.UploadURL = enterpriseURL
	}
	// check if the provided token was created by Vela
	_, resp, err := client.Authorizations.Check(context.Background(), c.config.ClientID, token)
	// check if the error is of type ErrorResponse
	var gerr *github.ErrorResponse
	if errors.As(err, &gerr) {
		// check the status code
		switch gerr.Response.StatusCode {
		// 404 is expected when non vela token is used
		case http.StatusNotFound:
			break
		default:
			return err
		}
	} else if err != nil {
		return err
	}

	// return error if the token was not created by Vela
	if resp.StatusCode != http.StatusOK {
		return errors.New("token was not created by vela")
	}

	return nil
}
