// SPDX-License-Identifier: Apache-2.0

package types

// StorageInfo is the API representation of a StorageInfo.
//
// swagger:model StorageInfo
type StorageInfo struct {
	StorageAccessKey *string `json:"storage_access_key,omitempty"`
	StorageSecretKey *string `json:"storage_secret_key,omitempty"`
	StorageAddress   *string `json:"storage_address,omitempty"`
	StorageBucket    *string `json:"storage_bucket,omitempty"`
}

// GetAccessKey returns the StorageAccessKey field.
//
// When the provided StorageInfo type is nil, or the field within
// the type is nil, it returns an empty string for the field.
func (w *StorageInfo) GetAccessKey() string {
	// return zero value if StorageInfo type or StorageAccessKey field is nil
	if w == nil || w.StorageAccessKey == nil {
		return ""
	}

	return *w.StorageAccessKey
}

// GetSecretKey returns the StorageSecretKey field.
//
// When the provided StorageInfo type is nil, or the field within
// the type is nil, it returns an empty string for the field.
func (w *StorageInfo) GetSecretKey() string {
	// return zero value if StorageInfo type or StorageSecretKey field is nil
	if w == nil || w.StorageSecretKey == nil {
		return ""
	}

	return *w.StorageSecretKey
}

// GetStorageAddress returns the StorageAddress field.
//
// When the provided StorageInfo type is nil, or the field within
// the type is nil, it returns an empty string for the field.
func (w *StorageInfo) GetStorageAddress() string {
	// return zero value if StorageInfo type or StorageAddress field is nil
	if w == nil || w.StorageAddress == nil {
		return ""
	}

	return *w.StorageAddress
}

// GetStorageBucket returns the StorageBucket field.
//
// When the provided StorageInfo type is nil, or the field within
// the type is nil, it returns an empty string for the field.
func (w *StorageInfo) GetStorageBucket() string {
	// return zero value if StorageInfo type or StorageBucket field is nil
	if w == nil || w.StorageBucket == nil {
		return ""
	}

	return *w.StorageBucket
}

// SetAccessKey sets the StorageAccessKey field.
//
// When the provided StorageInfo type is nil, it
// will set nothing and immediately return.
func (w *StorageInfo) SetAccessKey(v string) {
	// return if StorageInfo type is nil
	if w == nil {
		return
	}

	w.StorageAccessKey = &v
}

// SetSecretKey sets the StorageSecretKey field.
//
// When the provided StorageInfo type is nil, it
// will set nothing and immediately return.
func (w *StorageInfo) SetSecretKey(v string) {
	// return if StorageInfo type is nil
	if w == nil {
		return
	}

	w.StorageSecretKey = &v
}

// SetStorageAddress sets the StorageAddress field.
//
// When the provided StorageInfo type is nil, it
// will set nothing and immediately return.
func (w *StorageInfo) SetStorageAddress(v string) {
	// return if StorageInfo type is nil
	if w == nil {
		return
	}

	w.StorageAddress = &v
}

// SetStorageBucket sets the StorageBucket field.
//
// When the provided StorageInfo type is nil, it
// will set nothing and immediately return.
func (w *StorageInfo) SetStorageBucket(v string) {
	// return if StorageInfo type is nil
	if w == nil {
		return
	}

	w.StorageBucket = &v
}
