// SPDX-License-Identifier: Apache-2.0

package types

import (
	"github.com/minio/minio-go/v7"
)

// Bucket is the API types representation of an object storage.
//
// swagger:model CreateBucket
type Bucket struct {
	BucketName         string                   `json:"bucket_name,omitempty"`
	MakeBucketOptions  minio.MakeBucketOptions  `json:"make_bucket_options"`
	ListObjectsOptions minio.ListObjectsOptions `json:"list_objects_options"`
	Recursive          bool                     `json:"recursive"`
}

type Object struct {
	ObjectName string `json:"object_name,omitempty"`
	Bucket     Bucket `json:"bucket"`
	FilePath   string `json:"file_path,omitempty"`
}
