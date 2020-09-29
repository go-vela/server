// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/go-vela/types/library"
	"github.com/hashicorp/vault/api"
	"time"
)

type awsCfg struct {
	Role      string
	StsClient stsiface.STSAPI
}

type client struct {
	Vault      *api.Client
	Prefix     string
	AuthMethod string
	Aws        awsCfg
	Renewal    time.Duration
	TTL        time.Duration
}

const PrefixVaultV1 = "secret"
const PrefixVaultV2 = "secret/data"

// New returns a Secret implementation that integrates with a Vault secrets engine.
func New(addr, token, version, pathPrefix, authMethod, awsRole string, renewal time.Duration) (*client, error) {
	var prefix string
	switch version {
	case "1":
		prefix = PrefixVaultV1
	case "2":
		prefix = PrefixVaultV2
	default:
		return nil, fmt.Errorf("unrecognized vault version of %s", version)
	}

	// append admin defined prefix if not empty
	if pathPrefix != "" {
		prefix = fmt.Sprintf("%s/%s", prefix, pathPrefix)
	}

	conf := api.Config{Address: addr}

	// create Vault client
	c, err := api.NewClient(&conf)
	if err != nil {
		return nil, err
	}
	if token != "" {
		c.SetToken(token)
	}

	client := &client{
		Vault:      c,
		Prefix:     prefix,
		AuthMethod: authMethod,
		Renewal:    renewal,
		Aws: awsCfg{
			Role: awsRole,
		},
	}

	if authMethod != "" {
		err = client.initialize()
		if err != nil {
			return nil, err
		}

		// start the routine to refresh the token
		go client.refreshToken()
	}

	return client, nil
}

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

	return s
}

func vaultFromSecret(s *library.Secret) *api.Secret {
	data := make(map[string]interface{})
	vault := new(api.Secret)
	vault.Data = data

	// set events if found in Database secret
	if len(s.GetEvents()) > 0 {
		vault.Data["events"] = s.GetEvents()
	}

	// set images if found in Database secret
	if len(s.GetImages()) > 0 {
		vault.Data["images"] = s.GetImages()
	}

	// set name if found in Database secret
	if len(s.GetName()) > 0 {
		vault.Data["name"] = s.GetName()
	}

	// set org if found in Database secret
	if len(s.GetOrg()) > 0 {
		vault.Data["org"] = s.GetOrg()
	}

	// set repo if found in Database secret
	if len(s.GetRepo()) > 0 {
		vault.Data["repo"] = s.GetRepo()
	}

	// set team if found in Database secret
	if len(s.GetTeam()) > 0 {
		vault.Data["team"] = s.GetTeam()
	}

	// set type if found in Database secret
	if len(s.GetType()) > 0 {
		vault.Data["type"] = s.GetType()
	}

	// set value if found in Database secret
	if len(s.GetValue()) > 0 {
		vault.Data["value"] = s.GetValue()
	}

	// set allow_command if found in Database secret
	if s.AllowCommand != nil {
		vault.Data["allow_command"] = s.GetAllowCommand()
	}

	return vault
}
