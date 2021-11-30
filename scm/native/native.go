// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/fake"
	"github.com/jenkins-x/go-scm/scm/factory"
)

const (
	// events for repo webhooks.
	eventPush         = "push"
	eventPullRequest  = "pull_request"
	eventDeployment   = "deployment"
	eventIssueComment = "issue_comment"
)

var ctx = context.Background()

type (
	config struct {
		// specifies the address to use for the SCM client
		Address string
		// specifies the API endpoint to use for the SCM client
		API string
		// specifies the OAuth client ID from SCM to use for the SCM client
		ClientID string
		// specifies the OAuth client secret from SCM to use for the SCM client
		ClientSecret string
		// specifies which driver to use in the scm package
		Kind string
		// specifies the Vela server address to use for the SCM client
		ServerAddress string
		// specifies the Vela server address that the scm provider should use to send Vela webhooks
		ServerWebhookAddress string
		// specifies the context for the commit status to use for the SCM client
		StatusContext string
		// specifies the Vela web UI address to use for the SCM client
		WebUIAddress string
		// specifies the OAuth scopes to use for the SCM client
		Scopes []string
	}

	client struct {
		config *config
		OAuth  *oauth2.Config
	}
)

// New returns a SCM implementation that integrates with
// a SCM or a SCM Enterprise instance.
//
// nolint: revive // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// initialize new SCM client
	c := new(client)
	c.config = new(config)
	c.OAuth = new(oauth2.Config)
	// c.AuthReq = new(github.AuthorizationRequest)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// create the SCM OAuth config object
	c.OAuth = &oauth2.Config{
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		Scopes:       c.config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/login/oauth/authorize", c.config.Address),
			TokenURL: fmt.Sprintf("%s/login/oauth/access_token", c.config.Address),
		},
	}

	return c, nil
}

// NewTest returns a SCM implementation that integrates with the provided
// mock server. Only the url from the mock server is required.
//
// This function is intended for running tests only.
//
// nolint: revive // ignore returning unexported client
func NewTest(urls ...string) (*client, error) {
	address := urls[0]
	server := address

	// check if multiple URLs were provided
	if len(urls) > 1 {
		server = urls[1]
	}

	return New(
		WithAddress(address),
		WithClientID("foo"),
		WithClientSecret("bar"),
		WithKind("fake"),
		WithServerAddress(server),
		WithServerWebhookAddress(""),
		WithStatusContext("continuous-integration/vela"),
		WithWebUIAddress(address),
	)
}

// helper function to return the SCM OAuth client.
func (c *client) newClientToken(token string) (*scm.Client, error) {
	// return a fake client when testing
	if strings.EqualFold(c.config.Kind, "fake") {
		fakeSCM, data := fake.NewDefault()

		// load the fake data for testing
		load(data)

		return fakeSCM, nil
	}

	return factory.NewClient(c.config.Kind, c.config.Address, token)
}

// helper function to load fake data into the test SCM client for Go tests.
func load(d *fake.Data) {
	d.UserPermissions["github/octocat"] = map[string]string{}
	d.UserPermissions["github/octocat"]["foo"] = "admin"
	d.UserPermissions["github/octocat"]["notfound"] = ""
}
