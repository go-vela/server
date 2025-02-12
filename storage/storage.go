package storage

import (
	"fmt"
	"github.com/go-vela/server/constants"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// FromCLIContext helper function to setup Minio Client from the CLI arguments.
func FromCLIContext(c *cli.Context) (Storage, error) {
	logrus.Debug("creating Minio client from CLI configuration")
	logrus.Debugf("STORAGE Key: %s", c.String("storage.access.key"))
	logrus.Debugf("STORAGE Secret: %s", c.String("storage.secret.key"))
	// S3 configuration
	_setup := &Setup{
		Enable:    c.Bool("storage.enable"),
		Driver:    c.String("storage.driver"),
		Endpoint:  c.String("storage.addr"),
		AccessKey: c.String("storage.access.key"),
		SecretKey: c.String("storage.secret.key"),
		Bucket:    c.String("storage.bucket.name"),
		Secure:    c.Bool("storage.use.ssl"),
	}

	return New(_setup)

}

// New creates and returns a Vela service capable of
// integrating with the configured storage environment.
// Currently, the following storages are supported:
//
// * minio
// .
func New(s *Setup) (Storage, error) {
	// validate the setup being provided
	//
	err := s.Validate()
	if err != nil {
		return nil, fmt.Errorf("unable to validate storage setup: %w", err)
	}
	logrus.Debug("creating storage client from setup")
	// process the storage driver being provided
	switch s.Driver {
	case constants.DriverMinio:
		// handle the Kafka queue driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/queue?tab=doc#Setup.Kafka
		return s.Minio()
	default:
		// handle an invalid queue driver being provided
		return nil, fmt.Errorf("invalid storage driver provided: %s", s.Driver)
	}

}
