// SPDX-License-Identifier: Apache-2.0

package types

type OpenIDConfig struct {
	Issuer          string   `json:"issuer"`
	JWKSAddress     string   `json:"jwks_uri"`
	SupportedClaims []string `json:"supported_claims"`
	Algorithms      []string `json:"id_token_signing_alg_values_supported"`
}

// JWKS is a slice of JWKs
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key parsed with fields as the correct Go types.
type JWK struct {
	Algorithm string   `json:"alg"`
	Use       string   `json:"use"`
	X5t       string   `json:"x5t"`
	Kid       string   `json:"kid"`
	Kty       string   `json:"kty"`
	X5c       []string `json:"x5c"`

	N string `json:"n"` // modulus
	E string `json:"e"` // public exponent
}
