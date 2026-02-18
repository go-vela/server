// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
)

func (c *Client) AssumeRole(_ context.Context, durationSeconds int, prefix, sessionName string) (*types.STSCreds, error) {
	c.Logger.WithFields(logrus.Fields{
		"sessionName": sessionName,
	}).Tracef("creating STS assume role credentials")

	if durationSeconds <= 0 {
		durationSeconds = 900
	}

	opts := credentials.STSAssumeRoleOptions{
		AccessKey:       c.config.AccessKey, // server long-lived
		SecretKey:       c.config.SecretKey, // server long-lived
		RoleARN:         "arn:minio:iam:::role/vela-uploader",
		RoleSessionName: sessionName,
		DurationSeconds: durationSeconds,
		Policy:          c.GetPolicy(prefix),
	}

	stsCreds, err := credentials.NewSTSAssumeRole(c.config.Endpoint, opts)
	if err != nil {
		return nil, fmt.Errorf("unable to assume role: %w", err)
	}

	val, err := stsCreds.GetWithContext(c.client.CredContext())
	if err != nil {
		return nil, fmt.Errorf("unable to get credentials: %w", err)
	}

	return &types.STSCreds{
		AccessKey:    val.AccessKeyID,
		SecretKey:    val.SecretAccessKey,
		SessionToken: val.SessionToken,
		ExpiresAt:    time.Now().Add(time.Duration(durationSeconds) * time.Second),
		Endpoint:     c.GetAddress(),
		Bucket:       c.config.Bucket,
		Secure:       c.config.Secure,
	}, nil
}
