// SPDX-License-Identifier: Apache-2.0

package types

import (
	"github.com/golang-jwt/jwt/v5"
)

// OpenIDConfig is a struct that represents the OpenID Connect configuration.
//
// swagger:model OpenIDConfig
type OpenIDConfig struct {
	Issuer                 string   `json:"issuer"`
	JWKSAddress            string   `json:"jwks_uri"`
	ClaimsSupported        []string `json:"claims_supported"`
	Algorithms             []string `json:"id_token_signing_alg_values_supported"`
	ResponseTypesSupported []string `json:"response_types_supported"`
	SubjectTypesSupported  []string `json:"subject_types_supported"`
}

// OpenIDClaims struct is an extension of the JWT standard claims. It
// includes information relevant to OIDC services.
type OpenIDClaims struct {
	BuildNumber int    `json:"build_number,omitempty"`
	BuildID     int64  `json:"build_id,omitempty"`
	Actor       string `json:"actor,omitempty"`
	ActorSCMID  string `json:"actor_scm_id,omitempty"`
	Repo        string `json:"repo,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	Image       string `json:"image,omitempty"`
	Request     string `json:"request,omitempty"`
	Commands    bool   `json:"commands,omitempty"`
	Event       string `json:"event,omitempty"`
	Ref         string `json:"ref,omitempty"`
	SHA         string `json:"sha,omitempty"`
	jwt.RegisteredClaims
}

// JWKSet exists solely to provide proper swagger documentation.
// It is not otherwise used in code.
//
// swagger:model JWKSet
type JWKSet struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	E   string `json:"e"`
	N   string `json:"n"`
}
