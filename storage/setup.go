// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"fmt"

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

	// verify an endpoint was provided
	if len(s.Endpoint) == 0 {
		return fmt.Errorf("no storage endpoint provided")
	}

	// verify an access key was provided
	if len(s.AccessKey) == 0 {
		return fmt.Errorf("no storage access key provided")
	}

	// verify a secret key was provided
	if len(s.SecretKey) == 0 {
		return fmt.Errorf("no storage secret key provided")
	}

	// verify a bucket was provided
	//if len(s.Bucket) == 0 {
	//	return fmt.Errorf("no storage bucket provided")
	//}

	// setup is valid
	return nil
}
