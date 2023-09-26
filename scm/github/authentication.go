// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-vela/server/random"
	"github.com/go-vela/types/library"
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
