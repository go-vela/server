// SPDX-License-Identifier: Apache-2.0

package scm

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/scm/github"
	"github.com/go-vela/server/tracing"
)

// Setup represents the configuration necessary for
// creating a Vela service capable of integrating
// with a configured scm system.
type Setup struct {
	// scm Configuration

	// specifies the driver to use for the scm client
	Driver string
	// specifies the address to use for the scm client
	Address string
	// specifies the OAuth client ID from the scm system to use for the scm client
	ClientID string
	// specifies the OAuth client secret from the scm system to use for the scm client
	ClientSecret string
	// specifies App integration id
	AppID int64
	// specifies App integration private key
	AppPrivateKey string
	// specifies App integration path to private key
	AppPrivateKeyPath string
	// specifies App integration permissions set
	AppPermissions []string
	// specifies the Vela server address to use for the scm client
	ServerAddress string
	// specifies the Vela server address that the scm provider should use to send Vela webhooks
	ServerWebhookAddress string
	// specifies the context for the commit status to use for the scm client
	StatusContext string
	// specifies the Vela web UI address to use for the scm client
	WebUIAddress string
	// specifies the OAuth scopes to use for the scm client
	OAuthScopes []string
	// specifies the repo role map to use for the scm client
	RepoRoleMap map[string]string
	// specifies the org role map to use for the scm client
	OrgRoleMap map[string]string
	// specifies the team role map to use for the scm client
	TeamRoleMap map[string]string
	// specifies OTel tracing configurations
	Tracing *tracing.Client
}

// Github creates and returns a Vela service capable of
// integrating with a Github scm system.
func (s *Setup) Github(ctx context.Context) (Service, error) {
	logrus.Trace("creating github scm client from setup")

	// create new Github scm service
	//
	// https://pkg.go.dev/github.com/go-vela/server/scm/github?tab=doc#New
	return github.New(
		ctx,
		github.WithAddress(s.Address),
		github.WithClientID(s.ClientID),
		github.WithClientSecret(s.ClientSecret),
		github.WithServerAddress(s.ServerAddress),
		github.WithServerWebhookAddress(s.ServerWebhookAddress),
		github.WithStatusContext(s.StatusContext),
		github.WithWebUIAddress(s.WebUIAddress),
		github.WithOAuthScopes(s.OAuthScopes),
		github.WithTracing(s.Tracing),
		github.WithGithubAppID(s.AppID),
		github.WithGithubPrivateKey(s.AppPrivateKey),
		github.WithGithubPrivateKeyPath(s.AppPrivateKeyPath),
		github.WithGitHubAppPermissions(s.AppPermissions),
		github.WithRepoRoleMap(s.RepoRoleMap),
		github.WithOrgRoleMap(s.OrgRoleMap),
		github.WithTeamRoleMap(s.TeamRoleMap),
	)
}

// Gitlab creates and returns a Vela service capable of
// integrating with a Gitlab scm system.
func (s *Setup) Gitlab(_ context.Context) (Service, error) {
	logrus.Trace("creating gitlab scm client from setup")

	return nil, fmt.Errorf("unsupported scm driver: %s", constants.DriverGitlab)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating scm setup for client")

	// verify a scm driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no scm driver provided")
	}

	// verify a scm address was provided
	if len(s.Address) == 0 {
		return fmt.Errorf("no scm address provided")
	}

	// check if the scm address has a scheme
	if !strings.Contains(s.Address, "://") {
		return fmt.Errorf("scm address must be fully qualified (<scheme>://<host>)")
	}

	// check if the scm address has a trailing slash
	if strings.HasSuffix(s.Address, "/") {
		return fmt.Errorf("scm address must not have trailing slash")
	}

	// verify a scm OAuth client ID was provided
	if len(s.ClientID) == 0 {
		return fmt.Errorf("no scm client id provided")
	}

	// verify a scm OAuth client secret was provided
	if len(s.ClientSecret) == 0 {
		return fmt.Errorf("no scm client secret provided")
	}

	// verify a scm status context secret was provided
	if len(s.StatusContext) == 0 {
		return fmt.Errorf("no scm status context provided")
	}

	if len(s.OAuthScopes) == 0 {
		return fmt.Errorf("no scm scopes provided")
	}

	// setup is valid
	return nil
}
