package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// config holds the configuration for the MinIO client.
type config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Secure    bool
}

// MinioClient implements the Storage interface using MinIO.
type MinioClient struct {
	config  *config
	client  *minio.Client
	Options *minio.Options
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
}

// New creates a new MinIO client.
func New(endpoint string, opts ...ClientOpt) (*MinioClient, error) {
	// create new Minio client
	c := new(MinioClient)

	// default to secure connection
	urlEndpoint := "s3.amazonaws.com"
	useSSL := true

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
	c.Options.Creds = credentials.NewStaticV4(c.config.AccessKey, c.config.SecretKey, "")
	c.Options.Secure = c.config.Secure
	logrus.Debugf("secure: %v", c.config.Secure)

	if len(endpoint) > 0 {
		useSSL = strings.HasPrefix(endpoint, "https://")

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
	//minioClient, err := minio.New(endpoint, &minio.Options{
	//	Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
	//	Secure: useSSL,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//return &MinioClient{client: minioClient}, nil
}

// pingBucket checks if the specified bucket exists.
func pingBucket(c *MinioClient, bucket string) error {
	for i := 0; i < 10; i++ {
		_, err := c.client.BucketExists(context.Background(), bucket)
		if err != nil {
			c.Logger.Debugf("unable to ping %s. Retrying in %v", bucket, time.Duration(i)*time.Second)
			time.Sleep(1 * time.Second)

			continue
		}
	}

	return nil
}

// NewTest returns a Storage implementation that
// integrates with a local MinIO instance.
//
// This function is intended for running tests only.
//
//nolint:revive // ignore returning unexported client
func NewTest(endpoint, accessKey, secretKey string, secure bool) (*MinioClient, error) {
	//var cleanup func() error
	//var err error
	//endpoint, cleanup, err = miniotest.StartEmbedded()
	//
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "while starting embedded server: %s", err)
	//	os.Exit(1)
	//}
	//
	//err = cleanup()
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "while stopping embedded server: %s", err)
	//}

	// create a local fake MinIO instance
	//
	// https://pkg.go.dev/github.com/minio/minio-go/v7#New
	//minioClient, err := minio.New(endpoint, &minio.Options{
	//	Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
	//	Secure: false,
	//})
	//if err != nil {
	//	return nil, err
	//}

	return New(endpoint, WithAccessKey(accessKey), WithSecretKey(secretKey), WithSecure(secure))
}

//// UploadArtifact uploads an artifact to storage.
//func (c *MinioClient) UploadArtifact(ctx context.Context, workflowID, artifactName string, data []byte) error {
//	key := path.Join("artifacts", workflowID, artifactName)
//	bucket := "vela-artifacts"
//	return c.upload(ctx, bucket, key, data)
//}
//
//// DownloadArtifact downloads an artifact from storage.
//func (c *MinioClient) DownloadArtifact(ctx context.Context, workflowID, artifactName string) ([]byte, error) {
//	key := path.Join("artifacts", workflowID, artifactName)
//	bucket := "vela-artifacts"
//	return c.download(ctx, bucket, key)
//}
//
//// UploadCache uploads cache data to storage.
//func (c *MinioClient) UploadCache(ctx context.Context, key string, data []byte) error {
//	cacheKey := path.Join("cache", key)
//	bucket := "vela-cache"
//	return c.upload(ctx, bucket, cacheKey, data)
//}
//
//// DownloadCache downloads cache data from storage.
//func (c *MinioClient) DownloadCache(ctx context.Context, key string) ([]byte, error) {
//	cacheKey := path.Join("cache", key)
//	bucket := "vela-cache"
//	return c.download(ctx, bucket, cacheKey)
//}
//
//// DeleteCache deletes cache data from storage.
//func (c *MinioClient) DeleteCache(ctx context.Context, key string) error {
//	cacheKey := path.Join("cache", key)
//	bucket := "vela-cache"
//	return c.client.RemoveObject(ctx, bucket, cacheKey, minio.RemoveObjectOptions{})
//}
