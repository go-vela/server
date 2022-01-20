// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package registry

// Source represents a registry object
// for retrieving templates.
type Source struct {
	Host string
	Org  string
	Repo string
	Name string
	Ref  string
}
