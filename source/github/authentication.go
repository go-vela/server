// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-vela/server/random"

	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Authorize uses the given access token to authorize the user.
func (c *client) Authorize(token string) (string, error) {
	logrus.Trace("Authorizing user with token")

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
func (c *client) Login(w http.ResponseWriter, r *http.Request) (string, error) {
	logrus.Trace("Processing login request")

	// generate a random string for creating the OAuth state
	oAuthState, err := random.GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	// temporarily redirect request to Github to begin workflow
	http.Redirect(w, r, c.OConfig.AuthCodeURL(oAuthState), http.StatusTemporaryRedirect)

	return oAuthState, nil
}

func (c *client) LoginCLI(username, password, otp string) (*library.User, error) {
	logrus.Trace("Processing CLI login request")

	// create GitHub Basic auth client with user's credentials
	client := c.newClientBasicAuth(username, password, otp)

	// create authorization for user
	authorization, _, err := client.Authorizations.Create(ctx, c.AuthReq)
	if err != nil {
		return nil, err
	}

	// authorize the user for the token
	u, err := c.Authorize(authorization.GetToken())
	if err != nil {
		return nil, err
	}

	return &library.User{
		Name:  &u,
		Token: authorization.Token,
	}, nil
}

// Authenticate completes the authentication workflow for the session and returns the remote user details.
func (c *client) Authenticate(w http.ResponseWriter, r *http.Request, oAuthState string) (*library.User, error) {
	logrus.Trace("Authenticating user")

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

	// exchange OAuth code for token
	token, err := c.OConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	// authorize the user for the token
	u, err := c.Authorize(token.AccessToken)
	if err != nil {
		return nil, err
	}

	return &library.User{
		Name:  &u,
		Token: &token.AccessToken,
	}, nil
}
