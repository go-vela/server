package types

import "github.com/minio/minio-go/v7/pkg/lifecycle"

// Bucket is the API types representation of an object storage.
//
// swagger:model CreateBucket
type Bucket struct {
	BucketName      string                  `json:"bucket_name,omitempty"`
	Options         BucketOptions           `json:"options,omitempty"`
	LifecycleConfig lifecycle.Configuration `json:"life_cycle_config,omitempty"`
}

type BucketOptions struct {
	Region        string `json:"region,omitempty"`
	ObjectLocking bool   `json:"object_locking,omitempty"`
}

type Object struct {
	ObjectName string `json:"object_name,omitempty"`
	Bucket     Bucket `json:"bucket,omitempty"`
	FilePath   string `json:"file_path,omitempty"`
}
