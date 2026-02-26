// SPDX-License-Identifier: Apache-2.0

package types

// STSCreds defines the structure for temporary credentials used for object storage access.
//
// swagger:model STSCreds
type STSCreds struct {
	AccessKey    string `json:"access_key"`
	SecretKey    string `json:"secret_key"`
	SessionToken string `json:"session_token"`

	Endpoint string `json:"endpoint"`
	Enable   bool   `json:"enable"`
	Driver   string `json:"driver"`
	Bucket   string `json:"bucket"`
	Region   string `json:"region,omitempty"`
	Prefix   string `json:"prefix,omitempty"`
	Secure   bool   `json:"secure,omitempty"`

}
