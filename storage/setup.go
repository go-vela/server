// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"fmt"
	"net/url"

	"github.com/go-vela/server/constants"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/storage/minio"
)

// Setup represents the configuration necessary for
// creating a Vela service capable of integrating
// with a configured S3 environment.
type Setup struct {
	Enable    bool
	Driver    string
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	Secure    bool
}

// Minio creates and returns a Vela service capable
// of integrating with an S3 environment.
func (s *Setup) Minio() (Storage, error) {
	//client, err := minio.New(s.Endpoint, &minio.MakeBucketOptions{
	//	Creds:  credentials.NewStaticV4(s.AccessKey, s.SecretKey, ""),
	//	Secure: s.Secure,
	//})
	//if err != nil {
	//	return nil, err
	//}
	return minio.New(
		s.Endpoint,
		minio.WithAccessKey(s.AccessKey),
		minio.WithSecretKey(s.SecretKey),
		minio.WithSecure(s.Secure),
		minio.WithBucket(s.Bucket),
	)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating Storage setup for client")

	// verify storage is enabled
	//if !s.Enable {
	//	return fmt.Errorf("Storage is not enabled")
	//}
	if s.Enable {
		if s.Bucket == "" {
			return fmt.Errorf("storage is enabled but no bucket provided")
		}
		if s.Driver != "" && s.Driver != constants.DriverMinio {
			return fmt.Errorf("invalid storage driver provided: %s", s.Driver)
		}
		if s.Endpoint == "" {
			return fmt.Errorf("storage is enabled but no endpoint provided")
		}
		if s.AccessKey == "" || s.SecretKey == "" {
			return fmt.Errorf("storage is enabled but no access key or secret key provided")
		}
		if _, err := url.ParseRequestURI(s.Endpoint); err != nil {
			return fmt.Errorf("storage is enabled but endpoint is invalid")
		}
	}
	// setup is valid
	return nil
}
