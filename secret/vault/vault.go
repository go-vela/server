// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"github.com/go-vela/types/library"

	"github.com/hashicorp/vault/api"
)

type client struct {
	Vault *api.Client
}

// New returns a Secret implementation that integrates with a Vault secrets engine.
func New(addr, token string) (*client, error) {
	conf := api.Config{Address: addr}

	// create Vault client
	c, err := api.NewClient(&conf)
	if err != nil {
		return nil, err
	}

	// set Vault API token in client
	c.SetToken(token)

	client := &client{
		Vault: c,
	}

	return client, nil
}

func secretFromVault(vault *api.Secret) *library.Secret {
	s := new(library.Secret)

	// set events if found in Vault secret
	v, ok := vault.Data["events"]
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
	v, ok = vault.Data["images"]
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
	v, ok = vault.Data["name"]
	if ok {
		name, ok := v.(string)
		if ok {
			s.SetName(name)
		}
	}

	// set org if found in Vault secret
	v, ok = vault.Data["org"]
	if ok {
		org, ok := v.(string)
		if ok {
			s.SetOrg(org)
		}
	}

	// set repo if found in Vault secret
	v, ok = vault.Data["repo"]
	if ok {
		repo, ok := v.(string)
		if ok {
			s.SetRepo(repo)
		}
	}

	// set team if found in Vault secret
	v, ok = vault.Data["team"]
	if ok {
		team, ok := v.(string)
		if ok {
			s.SetTeam(team)
		}
	}

	// set type if found in Vault secret
	v, ok = vault.Data["type"]
	if ok {
		secretType, ok := v.(string)
		if ok {
			s.SetType(secretType)
		}
	}

	// set value if found in Vault secret
	v, ok = vault.Data["value"]
	if ok {
		value, ok := v.(string)
		if ok {
			s.SetValue(value)
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

	return vault
}
