package types

import "time"

type STSCreds struct {
	AccessKey    string `json:"access_key"`
	SecretKey    string `json:"secret_key"`
	SessionToken string `json:"session_token"`

	Endpoint string `json:"endpoint"`
	Bucket   string `json:"bucket"`
	Region   string `json:"region,omitempty"`
	Prefix   string `json:"prefix,omitempty"`

	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
