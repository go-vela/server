// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/go-vela/types/library"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	client struct {
		config *config
		AWS    *awsCfg
		Vault  *api.Client
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		Logger *logrus.Entry
	}
)

// New returns a Secret implementation that integrates with a Vault secrets engine.
//
//nolint:revive // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new Vault client
	c := new(client)

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
func secretFromVault(vault *api.Secret) *library.Secret {
	s := new(library.Secret)

	var data map[string]interface{}
	// handle k/v v2
	if _, ok := vault.Data["data"]; ok {
		data = vault.Data["data"].(map[string]interface{})
	} else {
		data = vault.Data
	}

	// set events if found in Vault secret
	v, ok := data["events"]
	if ok {
		events, ok := v.([]interface{})
		if ok {
			for _, element := range events {
				event, ok := element.(string)
				if ok {
					s.SetEvents(append(s.GetEvents(), event))
				}
			}
		}
	}

	// set images if found in Vault secret
	v, ok = data["images"]
	if ok {
		images, ok := v.([]interface{})
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

	// set created_at if found in Vault secret
	v, ok = data["created_at"]
	if ok {
		createdAt, ok := v.(int64)
		if ok {
			s.SetCreatedAt(createdAt)
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
		updatedAt, ok := v.(int64)
		if ok {
			s.SetUpdatedAt(updatedAt)
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
func vaultFromSecret(s *library.Secret) *api.Secret {
	data := make(map[string]interface{})
	vault := new(api.Secret)
	vault.Data = data

	// set events if found in Vela secret
	if len(s.GetEvents()) > 0 {
		vault.Data["events"] = s.GetEvents()
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
