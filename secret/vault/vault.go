// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	velaAPI "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

const (
	PrefixVaultV1 = "secret"
	PrefixVaultV2 = "secret/data"
)

type (
	awsCfg struct {
		Role      string
		StsClient stsiface.STSAPI
	}

	config struct {
		// specifies the address to use for the Vault client
		Address string
		// specifies the authentication method to use for the Vault client
		AuthMethod string
		// specifies the AWS role to use for the Vault client
		AWSRole string
		// specifies the prefix to use for the Vault client
		Prefix string
		// specifies the system prefix to use for the Vault client
		SystemPrefix string
		// specifies the token to use for the Vault client
		Token string
		// specifies the token duration to use for the Vault client
		TokenDuration time.Duration
		// specifies the token time to live for the Vault client
		TokenTTL time.Duration
		// specifies the version to use for the Vault client
		Version string
	}

	Client struct {
		config *config
		AWS    *awsCfg
		Vault  *api.Client
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		Logger *logrus.Entry
	}
)

// New returns a Secret implementation that integrates with a Vault secrets engine.
func New(opts ...ClientOpt) (*Client, error) {
	// create new Vault client
	c := new(Client)

	// create new fields
	c.config = new(config)
	c.AWS = new(awsCfg)
	c.Vault = new(api.Client)

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
	c.Logger = logrus.NewEntry(logger).WithField("engine", c.Driver())

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// check if a Vault prefix was provided
	if len(c.config.Prefix) > 0 {
		// update the Vault prefix with the system prefix
		c.config.Prefix = fmt.Sprintf("%s/%s", c.config.SystemPrefix, c.config.Prefix)
	} else {
		// set the Vault prefix from the system prefix
		c.config.Prefix = c.config.SystemPrefix
	}

	// create new Vault API client
	//
	// https://pkg.go.dev/github.com/hashicorp/vault/api#NewClient
	_vault, err := api.NewClient(&api.Config{Address: c.config.Address})
	if err != nil {
		return nil, err
	}

	// check if a token was provided for the Vault client
	if len(c.config.Token) > 0 {
		// set the token in the Vault client
		_vault.SetToken(c.config.Token)
	}

	// set the AWS role in the Vault client
	c.AWS.Role = c.config.AWSRole

	// set the Vault API client in the Vault client
	c.Vault = _vault

	// check if a authentication method was provided for the Vault client
	if len(c.config.AuthMethod) > 0 {
		// initialize the Vault client
		err = c.initialize()
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize vault token")
		}

		// start the routine to refresh the token
		go c.refreshToken()
	}

	return c, nil
}

// secretFromVault is a helper function to convert a HashiCorp Vault secret to a Vela secret.
//
//nolint:gocyclo,funlen // ignore cyclomatic complexity and function length due to conditionals
func secretFromVault(vault *api.Secret) *velaAPI.Secret {
	s := new(velaAPI.Secret)

	var data map[string]any
	// handle k/v v2
	if _, ok := vault.Data["data"]; ok {
		data = vault.Data["data"].(map[string]any)
	} else {
		data = vault.Data
	}

	// set allow_events if found in Vault secret
	v, ok := data["allow_events"]
	if ok {
		maskJSON, ok := v.(json.Number)
		if ok {
			mask, err := maskJSON.Int64()
			if err == nil {
				s.SetAllowEvents(velaAPI.NewEventsFromMask(mask))
			}
		}
	} else {
		// if not found, convert events to allow_events
		// this happens when vault secret has not been updated since before v0.23
		events, ok := data["events"]
		if ok {
			allowEventsMask := int64(0)

			for _, element := range events.([]any) {
				event, ok := element.(string)
				if ok {
					switch event {
					case constants.EventPush:
						allowEventsMask |= constants.AllowPushBranch
					case constants.EventPull:
						allowEventsMask |= constants.AllowPullOpen | constants.AllowPullReopen | constants.AllowPullSync
					case constants.EventComment:
						allowEventsMask |= constants.AllowCommentCreate | constants.AllowCommentEdit
					case constants.EventDeploy:
						allowEventsMask |= constants.AllowDeployCreate
					case constants.EventTag:
						allowEventsMask |= constants.AllowPushTag
					case constants.EventSchedule:
						allowEventsMask |= constants.AllowSchedule
					}
				}
			}

			s.SetAllowEvents(velaAPI.NewEventsFromMask(allowEventsMask))
		}
	}

	// set images if found in Vault secret
	v, ok = data["images"]
	if ok {
		images, ok := v.([]any)
		if ok {
			for _, element := range images {
				image, ok := element.(string)
				if ok {
					s.SetImages(append(s.GetImages(), image))
				}
			}
		}
	}

	// set name if found in Vault secret
	v, ok = data["name"]
	if ok {
		name, ok := v.(string)
		if ok {
			s.SetName(name)
		}
	}

	// set org if found in Vault secret
	v, ok = data["org"]
	if ok {
		org, ok := v.(string)
		if ok {
			s.SetOrg(org)
		}
	}

	// set repo if found in Vault secret
	v, ok = data["repo"]
	if ok {
		repo, ok := v.(string)
		if ok {
			s.SetRepo(repo)
		}
	}

	// set team if found in Vault secret
	v, ok = data["team"]
	if ok {
		team, ok := v.(string)
		if ok {
			s.SetTeam(team)
		}
	}

	// set type if found in Vault secret
	v, ok = data["type"]
	if ok {
		secretType, ok := v.(string)
		if ok {
			s.SetType(secretType)
		}
	}

	// set value if found in Vault secret
	v, ok = data["value"]
	if ok {
		value, ok := v.(string)
		if ok {
			s.SetValue(value)
		}
	}

	// set allow_command if found in Vault secret
	v, ok = data["allow_command"]
	if ok {
		command, ok := v.(bool)
		if ok {
			s.SetAllowCommand(command)
		}
	}

	// set allow_substitution if found in Vault secret
	v, ok = data["allow_substitution"]
	if ok {
		substitution, ok := v.(bool)
		if ok {
			s.SetAllowSubstitution(substitution)
		}
	} else {
		// set allow_substitution to allow_command value if not found in Vault secret
		cmd, ok := data["allow_command"]
		if ok {
			command, ok := cmd.(bool)
			if ok {
				s.SetAllowSubstitution(command)
			}
		}
	}

	// set created_at if found in Vault secret
	v, ok = data["created_at"]
	if ok {
		createdAtJSON, ok := v.(json.Number)
		if ok {
			createdAt, err := createdAtJSON.Int64()
			if err == nil {
				s.SetCreatedAt(createdAt)
			}
		}
	}

	// set created_by if found in Vault secret
	v, ok = data["created_by"]
	if ok {
		createdBy, ok := v.(string)
		if ok {
			s.SetCreatedBy(createdBy)
		}
	}

	// set updated_at if found in Vault secret
	v, ok = data["updated_at"]
	if ok {
		updatedAtJSON, ok := v.(json.Number)
		if ok {
			updatedAt, err := updatedAtJSON.Int64()
			if err == nil {
				s.SetUpdatedAt(updatedAt)
			}
		}
	}

	// set updated_by if found in Vault secret
	v, ok = data["updated_by"]
	if ok {
		updatedBy, ok := v.(string)
		if ok {
			s.SetUpdatedBy(updatedBy)
		}
	}

	return s
}

// vaultFromSecret is a helper function to convert a Vela secret to a HashiCorp Vault secret.
func vaultFromSecret(s *velaAPI.Secret) *api.Secret {
	data := make(map[string]any)
	vault := new(api.Secret)
	vault.Data = data

	// set allow events to mask
	if s.GetAllowEvents().ToDatabase() != 0 {
		vault.Data["allow_events"] = s.GetAllowEvents().ToDatabase()
	}

	// set images if found in Vela secret
	if len(s.GetImages()) > 0 {
		vault.Data["images"] = s.GetImages()
	}

	// set name if found in Vela secret
	if len(s.GetName()) > 0 {
		vault.Data["name"] = s.GetName()
	}

	// set org if found in Vela secret
	if len(s.GetOrg()) > 0 {
		vault.Data["org"] = s.GetOrg()
	}

	// set repo if found in Vela secret
	if len(s.GetRepo()) > 0 {
		vault.Data["repo"] = s.GetRepo()
	}

	// set team if found in Vela secret
	if len(s.GetTeam()) > 0 {
		vault.Data["team"] = s.GetTeam()
	}

	// set type if found in Vela secret
	if len(s.GetType()) > 0 {
		vault.Data["type"] = s.GetType()
	}

	// set value if found in Vela secret
	if len(s.GetValue()) > 0 {
		vault.Data["value"] = s.GetValue()
	}

	// set allow_command if found in Vela secret
	if s.AllowCommand != nil {
		vault.Data["allow_command"] = s.GetAllowCommand()
	}

	// set allow_substitution if found in Vela secret
	if s.AllowSubstitution != nil {
		vault.Data["allow_substitution"] = s.GetAllowSubstitution()
	}

	// set created_at if found in Vela secret
	if s.GetCreatedAt() > 0 {
		vault.Data["created_at"] = s.GetCreatedAt()
	}

	// set created_by if found in Vela secret
	if len(s.GetCreatedBy()) > 0 {
		vault.Data["created_by"] = s.GetCreatedBy()
	}

	// set updated_at if found in Vela secret
	if s.GetUpdatedAt() > 0 {
		vault.Data["updated_at"] = s.GetUpdatedAt()
	}

	// set updated_by if found in Vela secret
	if len(s.GetUpdatedBy()) > 0 {
		vault.Data["updated_by"] = s.GetUpdatedBy()
	}

	return vault
}
