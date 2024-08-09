// SPDX-License-Identifier: Apache-2.0

package internal

import "time"

type (
	// Database is the extra set of database data passed to the compiler.
	Database struct {
		Driver string `json:"driver"`
		Host   string `json:"host"`
	}

	// Queue is the extra set of queue data passed to the compiler.
	Queue struct {
		Channel string `json:"channel"`
		Driver  string `json:"driver"`
		Host    string `json:"host"`
	}

	// Source is the extra set of source data passed to the compiler.
	Source struct {
		Driver string `json:"driver"`
		Host   string `json:"host"`
	}

	// Vela is the extra set of Vela data passed to the compiler.
	Vela struct {
		Address              string        `json:"address"`
		WebAddress           string        `json:"web_address"`
		WebOauthCallbackPath string        `json:"web_oauth_callback_path"`
		AccessTokenDuration  time.Duration `json:"access_token_duration"`
		RefreshTokenDuration time.Duration `json:"refresh_token_duration"`
		OpenIDIssuer         string        `json:"oidc_issuer"`
	}

	// Metadata is the extra set of data passed to the compiler for
	// converting a yaml configuration to an executable pipeline.
	Metadata struct {
		Database *Database `json:"database"`
		Queue    *Queue    `json:"queue"`
		Source   *Source   `json:"source"`
		Vela     *Vela     `json:"vela"`
	}
)
