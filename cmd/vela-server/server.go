// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"net/url"
	"time"

	"github.com/go-vela/compiler/compiler"
	"github.com/go-vela/pkg-queue/queue"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/source"

	"github.com/go-vela/types"
)

type (
	// API represents the server configuration for API information.
	API struct {
		Address  string
		Hostname string
		Port     string
		Secret   string
		Url      *url.URL
	}

	// Build represents the server configuration for build information.
	Build struct {
		Timeout int64
	}

	// Github represents the compiler configuration for the github information.
	Github struct {
		Driver  bool
		Address string
		Token   string
		Url     *url.URL
	}

	// Modification represents the compiler configuration for the modification information.
	Modification struct {
		Address  string
		Secret   string
		Duration time.Duration
		Retries  int
		Url      *url.URL
	}

	// Compiler represents the server configuration for compiler information.
	Compiler struct {
		Github       *Github
		Modification *Modification
	}

	// Database represents the server configuration for database information.
	Database struct {
		Driver           string
		Address          string
		CompressionLevel int
		ConnectionIdle   int
		ConnectionLife   time.Duration
		ConnectionOpen   int
		EncryptionKey    string
		Url              *url.URL
	}

	// Logger represents the server configuration for logger information.
	Logger struct {
		Format string
		Level  string
	}

	// Metrics represents the server configuration for metrics information.
	Metrics struct {
		WorkerActive time.Duration
	}

	// Vault represents the secrets configuration for the vault information.
	Vault struct {
		Driver        bool
		Address       string
		AuthMethod    string
		AwsRole       string
		Prefix        string
		Token         string
		TokenDuration time.Duration
		Version       string
		Url           *url.URL
	}

	// Secrets represents the server configuration for secrets information.
	Secrets struct {
		Vault *Vault
	}

	// Security represents the server configuration for security information.
	Security struct {
		AccessToken       time.Duration
		RefreshToken      time.Duration
		RepoAllowList     []string
		SecureCookie      bool
		WebhookValidation bool
	}

	// Source represents the server configuration for source information.
	Source struct {
		Driver       string
		Address      string
		ClientID     string
		ClientSecret string
		Context      string
		Url          *url.URL
	}

	// WebUI represents the server configuration for web UI information.
	WebUI struct {
		Address       string
		OAuthEndpoint string
		Url           *url.URL
	}

	// Config represents the server configuration.
	Config struct {
		API      *API
		Build    *Build
		Compiler *Compiler
		Database *Database
		Logger   *Logger
		Metrics  *Metrics
		Queue    *queue.Setup
		Secrets  *Secrets
		Security *Security
		Source   *Source
		WebUI    *WebUI
	}

	// Server represents all configuration and
	// system processes for the server.
	Server struct {
		Config   *Config
		Compiler compiler.Engine
		Database database.Service
		Metadata *types.Metadata
		Queue    queue.Service
		Secrets  map[string]secret.Service
		Source   source.Service
	}
)
