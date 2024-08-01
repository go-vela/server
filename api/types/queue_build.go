// SPDX-License-Identifier: Apache-2.0

package types

import "fmt"

// QueueBuild is the API representation of the builds in the queue.
//
// swagger:model QueueBuild
type QueueBuild struct {
	Status   *string `json:"status,omitempty"`
	Number   *int32  `json:"number,omitempty"`
	Created  *int64  `json:"created,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

// GetStatus returns the Status field.
//
// When the provided QueueBuild type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *QueueBuild) GetStatus() string {
	// return zero value if QueueBuild type or Status field is nil
	if b == nil || b.Status == nil {
		return ""
	}

	return *b.Status
}

// GetNumber returns the Number field.
//
// When the provided QueueBuild type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *QueueBuild) GetNumber() int32 {
	// return zero value if QueueBuild type or Number field is nil
	if b == nil || b.Number == nil {
		return 0
	}

	return *b.Number
}

// GetCreated returns the Created field.
//
// When the provided QueueBuild type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *QueueBuild) GetCreated() int64 {
	// return zero value if QueueBuild type or Created field is nil
	if b == nil || b.Created == nil {
		return 0
	}

	return *b.Created
}

// GetFullName returns the FullName field.
//
// When the provided QueueBuild type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *QueueBuild) GetFullName() string {
	// return zero value if QueueBuild type or FullName field is nil
	if b == nil || b.FullName == nil {
		return ""
	}

	return *b.FullName
}

// SetStatus sets the Status field.
//
// When the provided QueueBuild type is nil, it
// will set nothing and immediately return.
func (b *QueueBuild) SetStatus(v string) {
	// return if QueueBuild type is nil
	if b == nil {
		return
	}

	b.Status = &v
}

// SetNumber sets the Number field.
//
// When the provided QueueBuild type is nil, it
// will set nothing and immediately return.
func (b *QueueBuild) SetNumber(v int32) {
	// return if QueueBuild type is nil
	if b == nil {
		return
	}

	b.Number = &v
}

// SetCreated sets the Created field.
//
// When the provided QueueBuild type is nil, it
// will set nothing and immediately return.
func (b *QueueBuild) SetCreated(v int64) {
	// return if QueueBuild type is nil
	if b == nil {
		return
	}

	b.Created = &v
}

// SetFullName sets the FullName field.
//
// When the provided QueueBuild type is nil, it
// will set nothing and immediately return.
func (b *QueueBuild) SetFullName(v string) {
	// return if QueueBuild type is nil
	if b == nil {
		return
	}

	b.FullName = &v
}

// String implements the Stringer interface for the QueueBuild type.
func (b *QueueBuild) String() string {
	return fmt.Sprintf(`{
  Created: %d,
  FullName: %s,
  Number: %d,
  Status: %s,
}`,
		b.GetCreated(),
		b.GetFullName(),
		b.GetNumber(),
		b.GetStatus(),
	)
}
