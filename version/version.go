// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package version

import "github.com/coreos/go-semver/semver"

var (
	// VersionMajor is for an API incompatible changes
	VersionMajor int64
	// VersionMinor is for functionality in a backwards-compatible manner
	VersionMinor int64 = 3
	// VersionPatch is for backwards-compatible bug fixes
	VersionPatch int64
	// VersionDev indicates drone build number. Releases will be empty string.
	VersionDev string
)

// Version is the specification version that the package types support.
var Version = semver.Version{
	Major:    VersionMajor,
	Minor:    VersionMinor,
	Patch:    VersionPatch,
	Metadata: VersionDev,
}
