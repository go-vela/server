// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"fmt"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
)

// config holds the configuration for the MinIO client.
//
// but it is necessary for the MinIO client to function properly.
type config struct {
	Enable       bool
	Endpoint     string
	AccessKey    string
	SecretKey    string
	Bucket       string
	Secure       bool
	Token        string
	Driver       string
	PublicPolicy bool
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

	urlEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	// create the Minio client from the provided endpoint and options
	minioClient, err := minio.New(urlEndpoint.Host, c.Options)
	if err != nil {
		return nil, err
	}

	c.client = minioClient

	if c.config.PublicPolicy {
		if err := c.applyPublicPolicy(context.Background()); err != nil {
			logrus.Warnf("storage: failed to apply public bucket policy: %v", err)
		}
	}

	return c, nil
}

// applyPublicPolicy sets an anonymous read-only policy on the public/* prefix of the bucket,
// allowing unauthenticated GET requests for objects stored under that prefix.
func (c *Client) applyPublicPolicy(ctx context.Context) error {
	policy := fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {"AWS": ["*"]},
    "Action": ["s3:GetObject"],
    "Resource": ["arn:aws:s3:::%s/public/*"]
  }]
}`, c.config.Bucket)

	c.Logger.Infof("storage: applying public-read policy for public/* prefix in bucket %s", c.config.Bucket)

	return c.client.SetBucketPolicy(ctx, c.config.Bucket, policy)
}

// NewTest returns a Storage implementation that
// integrates with a local MinIO instance.
//
// This function is intended for running tests only.
func NewTest(endpoint, accessKey, secretKey, bucket string, secure bool) (*Client, error) {
	return New(endpoint,
		WithOptions(true, secure,
			endpoint, accessKey, secretKey, bucket, "", constants.DriverMinio, false))
}
