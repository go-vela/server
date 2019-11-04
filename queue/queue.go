// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

// Service represents the interface for Vela integrating
// with the different supported Queue backends.
type Service interface {

	// Publish defines a function that inserts an
	// item to the specified channel in the queue.
	Publish(string, []byte) error
}
