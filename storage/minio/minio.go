// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

// config holds the configuration for the MinIO client.
type config struct {
	Enable    bool
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Secure    bool
	Token     string
}

// Client implements the Storage interface using MinIO.
type Client struct {
	config  *config
	client  *minio.Client
	Options *minio.Options
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
}

// New creates a new MinIO client.
func New(endpoint string, opts ...ClientOpt) (*Client, error) {
	// create new Minio client
	c := new(Client)

	// default to secure connection
	var urlEndpoint string

	// create new fields
	c.config = new(config)
	c.Options = new(minio.Options)

	// create new logger for the client
	logger := logrus.StandardLogger()
	c.Logger = logrus.NewEntry(logger).WithField("minio", "minio")

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	c.Options.Creds = credentials.NewStaticV4(c.config.AccessKey, c.config.SecretKey, c.config.Token)
	c.Options.Secure = c.config.Secure
	logrus.Debugf("secure: %v", c.config.Secure)

	if len(endpoint) > 0 {
		useSSL := strings.HasPrefix(endpoint, "https://")

		if !useSSL {
			if !strings.HasPrefix(endpoint, "http://") {
				return nil, fmt.Errorf("invalid server %s: must to be a HTTP URI", endpoint)
			}

			urlEndpoint = endpoint[7:]
		} else {
			urlEndpoint = endpoint[8:]
		}
	}

	// create the Minio client from the provided endpoint and options
	minioClient, err := minio.New(urlEndpoint, c.Options)
	if err != nil {
		return nil, err
	}

	c.client = minioClient

	return c, nil
}

// NewTest returns a Storage implementation that
// integrates with a local MinIO instance.
//
// This function is intended for running tests only.
//

func NewTest(endpoint, accessKey, secretKey, bucket string, secure bool) (*Client, error) {
	return New(endpoint,
		WithOptions(true, secure,
			endpoint, accessKey, secretKey, bucket, ""))
}
