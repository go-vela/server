// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"
	"strings"
)

// ClientOpt represents a configuration option to initialize the scm client for GitHub.
type ClientOpt func(*client) error

// WithAddress sets the GitHub address in the scm client for GitHub.
func WithAddress(address string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring address in github scm client")

		// set a default address for the client
		c.config.Address = defaultURL
		// set a default API address for the client
		c.config.API = defaultAPI

		// check if an address was provided
		if len(address) > 0 {
			// check if the address is equal to the defaults
			if !strings.EqualFold(c.config.Address, address) {
				c.config.Address = strings.TrimSuffix(address, "/")
				if !strings.Contains(c.config.Address, "https://github.com") {
					c.config.API = fmt.Sprintf("%s/%s", c.config.Address, "api/v3/")
				}
			}
		}

		return nil
	}
}

// WithClientID sets the OAuth client ID in the scm client for GitHub.
func WithClientID(id string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring OAuth client ID in github scm client")

		// check if the OAuth client ID provided is empty
		if len(id) == 0 {
			return fmt.Errorf("no GitHub OAuth client ID provided")
		}

		// set the OAuth client ID in the github client
		c.config.ClientID = id

		return nil
	}
}

// WithClientSecret sets the OAuth client secret in the scm client for GitHub.
func WithClientSecret(secret string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring OAuth client secret in github scm client")

		// check if the OAuth client secret provided is empty
		if len(secret) == 0 {
			return fmt.Errorf("no GitHub OAuth client secret provided")
		}

		// set the OAuth client secret in the github client
		c.config.ClientSecret = secret

		return nil
	}
}

// WithServerAddress sets the Vela server address in the scm client for GitHub.
func WithServerAddress(address string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring Vela server address in github scm client")

		// check if the Vela server address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("no Vela server address provided")
		}

		// set the Vela server address in the github client
		c.config.ServerAddress = address

		return nil
	}
}

// WithServerWebhookAddress sets the Vela server webhook address in the scm client for GitHub.
func WithServerWebhookAddress(address string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring Vela server webhook address in github scm client")

		// fallback to Vela server address if the provided Vela server webhook address is empty
		if len(address) == 0 {
			c.config.ServerWebhookAddress = c.config.ServerAddress
			return nil
		}

		// set the Vela server webhook address in the github client
		c.config.ServerWebhookAddress = address

		return nil
	}
}

// WithStatusContext sets the context for commit statuses in the scm client for GitHub.
func WithStatusContext(context string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring context for commit statuses in github scm client")

		// check if the context for the commit statuses provided is empty
		if len(context) == 0 {
			return fmt.Errorf("no GitHub context for commit statuses provided")
		}

		// set the context for the commit status in the github client
		c.config.StatusContext = context

		return nil
	}
}

// WithWebUIAddress sets the Vela web UI address in the scm client for GitHub.
func WithWebUIAddress(address string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring Vela web UI address in github scm client")

		// set the Vela web UI address in the github client
		c.config.WebUIAddress = address

		return nil
	}
}

// WithScopes sets the OAuth scopes in the scm client for GitHub.
func WithScopes(scopes []string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring oauth scopes in github scm client")

		// check if the scopes provided is empty
		if len(scopes) == 0 {
			return fmt.Errorf("no GitHub OAuth scopes provided")
		}

		// set the scopes in the github client
		c.config.Scopes = scopes

		return nil
	}
}
