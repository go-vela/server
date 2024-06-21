// SPDX-License-Identifier: Apache-2.0

package image

import (
	"github.com/distribution/reference"
)

// ParseWithError digests the provided image into a
// fully qualified canonical reference. If an error
// occurs, it will return the last digested form of
// the image.
func ParseWithError(_image string) (string, error) {
	// parse the image provided into a
	// named, fully qualified reference
	//
	// https://pkg.go.dev/github.com/distribution/reference#ParseAnyReference
	_reference, err := reference.ParseAnyReference(_image)
	if err != nil {
		return _image, err
	}

	// ensure we have the canonical form of the named reference
	//
	// https://pkg.go.dev/github.com/distribution/reference#ParseNamed
	_canonical, err := reference.ParseNamed(_reference.String())
	if err != nil {
		return _reference.String(), err
	}

	// ensure the canonical reference has a tag
	//
	// https://pkg.go.dev/github.com/distribution/reference#TagNameOnly
	return reference.TagNameOnly(_canonical).String(), nil
}
