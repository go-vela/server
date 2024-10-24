// SPDX-License-Identifier: Apache-2.0

package types

// QueueInfo is the API representation of a QueueInfo.
//
// swagger:model QueueInfo
type QueueInfo struct {
	QueuePublicKey *string `json:"queue_public_key,omitempty"`
	QueueAddress   *string `json:"queue_address,omitempty"`
}

// GetPublicKey returns the QueuePublicKey field.
//
// When the provided QueueInfo type is nil, or the field within
// the type is nil, it returns an empty string for the field.
func (w *QueueInfo) GetPublicKey() string {
	// return zero value if QueueInfo type or QueuePublicKey field is nil
	if w == nil || w.QueuePublicKey == nil {
		return ""
	}

	return *w.QueuePublicKey
}

// GetQueueAddress returns the QueueAddress field.
//
// When the provided QueueInfo type is nil, or the field within
// the type is nil, it returns an empty string for the field.
func (w *QueueInfo) GetQueueAddress() string {
	// return zero value if QueueInfo type or QueueAddress field is nil
	if w == nil || w.QueueAddress == nil {
		return ""
	}

	return *w.QueueAddress
}

// SetPublicKey sets the QueuePublicKey field.
//
// When the provided QueueInfo type is nil, it
// will set nothing and immediately return.
func (w *QueueInfo) SetPublicKey(v string) {
	// return if QueueInfo type is nil
	if w == nil {
		return
	}

	w.QueuePublicKey = &v
}

// SetQueueAddress sets the QueueAddress field.
//
// When the provided QueueInfo type is nil, it
// will set nothing and immediately return.
func (w *QueueInfo) SetQueueAddress(v string) {
	// return if QueueInfo type is nil
	if w == nil {
		return
	}

	w.QueueAddress = &v
}
