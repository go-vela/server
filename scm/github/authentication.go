// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-vela/server/random"
	"github.com/go-vela/types/library"
	"github.com/google/go-github/v54/github"
)

// Authorize uses the given access token to authorize the user.
func (c *client) Authorize(ctx context.Context, token string) (string, error) {
	c.Logger.Trace("authorizing user with token")

	// create GitHub OAuth client with user's token
	client := c.newClientToken(token)

	// send API call to capture the current user making the call
	u, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return "", err
	}

	return u.GetLogin(), nil
}

// Login begins the authentication workflow for the session.
func (c *client) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) (string, error) {
	c.Logger.Trace("processing login request")

	// generate a random string for creating the OAuth state
	oAuthState, err := random.GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	// pass through the redirect if it exists
	redirect := r.FormValue("redirect_uri")
	if len(redirect) > 0 {
		c.OAuth.RedirectURL = redirect
	}

	// temporarily redirect request to Github to begin workflow
	http.Redirect(w, r, c.OAuth.AuthCodeURL(oAuthState), http.StatusTemporaryRedirect)

	return oAuthState, nil
}

// Authenticate completes the authentication workflow for the session
// and returns the remote user details.
func (c *client) Authenticate(ctx context.Context, w http.ResponseWriter, r *http.Request, oAuthState string) (*library.User, error) {
	c.Logger.Trace("authenticating user")

	// get the OAuth code
	code := r.FormValue("code")
	if len(code) == 0 {
		return nil, nil
	}

	// verify the OAuth state
	state := r.FormValue("state")
	if state != oAuthState {
		return nil, fmt.Errorf("unexpected oauth state: want %s but got %s", oAuthState, state)
	}

	// pass through the redirect if it exists
	redirect := r.FormValue("redirect_uri")
	if len(redirect) > 0 {
		c.OAuth.RedirectURL = redirect
	}

	// exchange OAuth code for token
	token, err := c.OAuth.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	// authorize the user for the token
	u, err := c.Authorize(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}

	return &library.User{
		Name:  &u,
		Token: &token.AccessToken,
	}, nil
}

// AuthenticateToken completes the authentication workflow
// for the session and returns the remote user details.
func (c *client) AuthenticateToken(ctx context.Context, r *http.Request) (*library.User, error) {
	c.Logger.Trace("authenticating user via token")

	token := r.Header.Get("Token")
	if len(token) == 0 {
		return nil, errors.New("no token provided")
	}

	err := c.ValidateOAuthToken(ctx, token)
	if err != nil {
		return nil, err
	}

	u, err := c.Authorize(ctx, token)
	if err != nil {
		return nil, err
	}

	return &library.User{
		Name:  &u,
		Token: &token,
	}, nil
}

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
