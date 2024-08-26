// SPDX-License-Identifier: Apache-2.0

package types

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
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
	Actor       string `json:"actor,omitempty"`
	ActorSCMID  string `json:"actor_scm_id,omitempty"`
	Branch      string `json:"branch,omitempty"`
	BuildID     int64  `json:"build_id,omitempty"`
	BuildNumber int    `json:"build_number,omitempty"`
	Commands    bool   `json:"commands,omitempty"`
	Event       string `json:"event,omitempty"`
	Fork        bool   `json:"fork,omitempty"`
	Image       string `json:"image,omitempty"`
	ImageName   string `json:"image_name,omitempty"`
	ImageTag    string `json:"image_tag,omitempty"`
	Ref         string `json:"ref,omitempty"`
	Repo        string `json:"repo,omitempty"`
	Request     string `json:"request,omitempty"`
	SHA         string `json:"sha,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	jwt.RegisteredClaims
}

// JWKSet is a wrapper of lestrrat-go/jwx/jwk.Set for API Swagger gen.
//
// swagger:model JWKSet
type JWKSet struct {
	jwk.Set
}
