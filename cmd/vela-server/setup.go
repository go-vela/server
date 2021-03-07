// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"net/url"

	"github.com/go-vela/compiler/compiler"
	"github.com/go-vela/compiler/compiler/native"
	cnative "github.com/go-vela/compiler/compiler/native"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/secret"
	snative "github.com/go-vela/server/secret/native"
	"github.com/go-vela/server/secret/vault"
	"github.com/go-vela/server/source"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// Setup prepares the API for execution.
func (a *API) Setup() error {
	// parse the server address into a url structure
	parsed, err := url.Parse(a.Address)
	if err != nil {
		return fmt.Errorf("unable to parse server address: %v", err)
	}

	// save the parsed server address
	a.Url = parsed

	return nil
}

// Setup prepares the Compiler for execution.
func (c *Compiler) Setup() (compiler.Engine, error) {
	logrus.Trace("preparing compiler for execution")

	// check if the github driver is enabled
	if c.Github.Driver {
		// parse the compiler GitHub address into a url structure
		parsed, err := url.Parse(c.Github.Address)
		if err != nil {
			return nil, fmt.Errorf("unable to parse compiler GitHub address: %v", err)
		}

		// save the parsed compiler GitHub address
		c.Github.Url = parsed
	}

	// check if an address for the modification service is provided
	if len(c.Modification.Address) > 0 {
		// parse the compiler modification address into a url structure
		parsed, err := url.Parse(c.Modification.Address)
		if err != nil {
			return nil, fmt.Errorf("unable to parse compiler modification address: %v", err)
		}

		// save the parsed compiler modification address
		c.Modification.Url = parsed
	}

	return cnative.New(s.Config.Compiler), nil
}

// Setup prepares the Database for execution.
func (d *Database) Setup() (database.Service, error) {
	logrus.Trace("preparing database for execution")

	// parse the database address into a url structure
	parsed, err := url.Parse(d.Address)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database address: %v", err)
	}

	// save the parsed database address
	d.Url = parsed

	return database.New(s.Config.Database), nil
}

// Setup prepares the Secrets for execution.
func (s *Secrets) Setup() (map[string]secret.Service, error) {
	logrus.Trace("preparing secrets for execution")

	secrets := make(map[string]secret.Service)

	_native, err := native.New(d)
	if err != nil {
		return nil, err
	}

	secrets[constants.DriverNative] = _native

	// check if the vault driver is enabled
	if s.Vault.Driver {
		// parse the vault address into a url structure
		parsed, err := url.Parse(s.Vault.Address)
		if err != nil {
			return nil, fmt.Errorf("unable to parse secrets Vault address: %v", err)
		}

		// save the parsed vault address
		s.Vault.Url = parsed

		// setup the vault from the configuration
		_vault, err := vault.New(vault.Config{
			Address:    s.Vault.Address,
			AuthMethod: s.Vault.AuthMethod,
			AwsRole:    s.Vault.AwsRole,
			Prefix:     s.Vault.Prefix,
			Renewal:    s.Vault.TokenDuration,
			Token:      s.Vault.Token,
			Version:    s.Vault.Version,
		})
		if err != nil {
			return nil, err
		}

		secrets[constants.DriverVault] = _vault
	}

	return secrets, nil
}

// Setup prepares the Source for execution.
func (s *Source) Setup() (source.Service, error) {
	logrus.Trace("preparing source for execution")

	// parse the source address into a url structure
	parsed, err := url.Parse(s.Address)
	if err != nil {
		return nil, fmt.Errorf("unable to parse source address: %v", err)
	}

	// save the parsed source address
	s.Url = parsed

	return snative.New(c), nil
}

// Setup prepares the WebUI for execution.
func (w *WebUI) Setup() error {
	// parse the web UI address into a url structure
	parsed, err := url.Parse(w.Address)
	if err != nil {
		return fmt.Errorf("unable to parse web UI address: %v", err)
	}

	// save the parsed web UI address
	w.Url = parsed

	return nil
}

// Setup prepares the Server for execution.
func (s *Server) Setup() error {
	// log a message indicating the configuration preparation
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
	logrus.Info("preparing server for execution")

	// prepare the API for execution
	err := s.Config.API.Setup()
	if err != nil {
		return err
	}

	// prepare the compiler for execution
	s.Compiler, err = s.Config.Compiler.Setup()
	if err != nil {
		return err
	}

	// prepare the database for execution
	s.Database, err = s.Config.Database.Setup()
	if err != nil {
		return err
	}

	// prepare the secrets for execution
	s.Secrets, err = s.Config.Secrets.Setup()
	if err != nil {
		return err
	}

	// prepare the source for execution
	s.Source, err = s.Config.Source.Setup()
	if err != nil {
		return err
	}

	// prepare the web UI for execution
	err = s.Config.WebUI.Setup()
	if err != nil {
		return err
	}

	// prepare the metadata for execution
	s.Metadata = &types.Metadata{
		Database: s.Config.Database.Metadata(),
		Queue: &types.Queue{
			Driver: s.Config.Queue.Driver,
			Host:   s.Config.Queue.Url.Host,
		},
		Source: s.Config.Source.Metadata(),
		Vela: &types.Vela{
			Address:              s.Config.API.Address,
			WebAddress:           s.Config.WebUI.Address,
			WebOauthCallbackPath: s.Config.WebUI.OAuthEndpoint,
			AccessTokenDuration:  s.Config.Security.AccessToken,
			RefreshTokenDuration: s.Config.Security.RefreshToken,
		},
	}

	return nil
}
