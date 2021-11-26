// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// ClientOpt represents a configuration option to initialize the scm client.
type ClientOpt func(*client) error

// ErrMissingProvidedField defines the error type when the
// Client is missing a required option to be provided in the client.
var ErrMissingProvidedField = errors.New("missing required option")

// WithAddress sets the GitHub address in the scm client.
func WithAddress(address string) ClientOpt {
	logrus.Trace("configuring address in native scm client")
	return func(c *client) error {
		// set a default address for the client
		c.config.Address = address
		// set a default API address for the client
		c.config.API = fmt.Sprintf("%s/api/v3/", address)

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

// WithClientID sets the GitHub OAuth client ID in the scm client.
func WithClientID(id string) ClientOpt {
	logrus.Trace("configuring oauth client id in native scm client")

	return func(c *client) error {
		// check if the OAuth client ID provided is empty
		if len(id) == 0 {
			return fmt.Errorf("%s: ClientID", ErrMissingProvidedField)
		}

		// set the OAuth client ID in the github client
		c.config.ClientID = id

		return nil
	}
}

// WithClientSecret sets the GitHub OAuth client secret in the scm client.
func WithClientSecret(secret string) ClientOpt {
	logrus.Trace("configuring oauth client secret in native scm client")

	return func(c *client) error {
		// check if the OAuth client secret provided is empty
		if len(secret) == 0 {
			return fmt.Errorf("%s: ClientSecret", ErrMissingProvidedField)
		}

		// set the OAuth client secret in the github client
		c.config.ClientSecret = secret

		return nil
	}
}

// WithKind sets the driver kind in the scm client.
func WithKind(kind string) ClientOpt {
	logrus.Trace("configuring oauth client secret in native scm client")

	return func(c *client) error {
		// check if the driver kind provided is empty
		if len(kind) == 0 {
			return fmt.Errorf("%s: Kind", ErrMissingProvidedField)
		}

		// set the driver kind in the github client
		c.config.Kind = kind

		return nil
	}
}

// WithServerAddress sets the Vela server address in the scm client.
func WithServerAddress(address string) ClientOpt {
	logrus.Trace("configuring vela server address in native scm client")

	return func(c *client) error {
		// check if the Vela server address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("%s: ServerAddress", ErrMissingProvidedField)
		}

		// set the Vela server address in the github client
		c.config.ServerAddress = address

		return nil
	}
}

// WithServerWebhookAddress sets the Vela server webhook address in the scm client.
func WithServerWebhookAddress(address string) ClientOpt {
	logrus.Trace("configuring vela server webhook address in native scm client")

	return func(c *client) error {
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

// WithStatusContext sets the GitHub context for commit statuses in the scm client.
func WithStatusContext(context string) ClientOpt {
	logrus.Trace("configuring context for commit statuses in native scm client")

	return func(c *client) error {
		// check if the context for the commit statuses provided is empty
		if len(context) == 0 {
			return fmt.Errorf("%s: StatusContext", ErrMissingProvidedField)
		}

		// set the context for the commit status in the github client
		c.config.StatusContext = context

		return nil
	}
}

// WithWebUIAddress sets the Vela web UI address in the scm client.
func WithWebUIAddress(address string) ClientOpt {
	logrus.Trace("configuring Vela web ui address in native scm client")

	return func(c *client) error {
		// set the Vela web UI address in the github client
		c.config.WebUIAddress = address

		return nil
	}
}

// WithScopes sets the GitHub OAuth scopes in the scm client.
func WithScopes(scopes []string) ClientOpt {
	logrus.Trace("configuring oauth scopes in native scm client")

	return func(c *client) error {
		// check if the scopes provided is empty
		if len(scopes) == 0 {
			return fmt.Errorf("%s: Scopes", ErrMissingProvidedField)
		}

		// set the scopes in the github client
		c.config.Scopes = scopes

		return nil
	}
}
