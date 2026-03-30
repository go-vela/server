// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/random"
)

// Authorize uses the given access token to authorize the user.
func (c *Client) Authorize(ctx context.Context, token string) (string, error) {
	c.Logger.Trace("authorizing user with token")

	// create GitHub OAuth client with user's token
	client := c.newUserOAuthTokenClient(ctx, &api.User{Token: &token})

	// send API call to capture the current user making the call
	u, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return "", err
	}

	return u.GetLogin(), nil
}

// Login begins the authentication workflow for the session.
func (c *Client) Login(_ context.Context, w http.ResponseWriter, r *http.Request) (string, error) {
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
func (c *Client) Authenticate(ctx context.Context, _ http.ResponseWriter, r *http.Request, oAuthState string) (*api.User, error) {
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
	token, err := c.OAuth.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// authorize the user for the token
	u, err := c.Authorize(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}

	exp := token.Expiry.Unix()

	return &api.User{
		Name:              &u,
		Token:             &token.AccessToken,
		OAuthRefreshToken: &token.RefreshToken,
		TokenExp:          &exp,
	}, nil
}

// AuthenticateToken completes the authentication workflow
// for the session and returns the remote user details.
func (c *Client) AuthenticateToken(ctx context.Context, r *http.Request) (*api.User, error) {
	c.Logger.Trace("authenticating user via token")

	token := r.Header.Get("Token")
	if len(token) == 0 {
		return nil, errors.New("no token provided")
	}

	u, err := c.Authorize(ctx, token)
	if err != nil {
		return nil, err
	}

	return &api.User{
		Name:  &u,
		Token: &token,
	}, nil
}
