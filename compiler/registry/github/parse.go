// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"
	"strings"

	"github.com/go-vela/server/compiler/registry"

	"github.com/goware/urlx"
)

// Parse creates the registry source object from
// a template path and default branch.
func (c *client) Parse(path string) (*registry.Source, error) {
	// ref will hold the reference identifier,
	// eg. <org>/<repo>/<filename>@<reference>
	ref := ""

	// parse the path provided
	//
	// goware/urlx is used over net/url because it is more consistent for parsing
	// the template paths we use (similar to go imports)
	u, err := urlx.Parse(path)
	if err != nil {
		return nil, err
	}

	u.Path = strings.TrimPrefix(u.Path, "/")

	// this will handle multiple cases for the path:
	// * <org>/<repo>/<filename>
	// * <org>/<repo>/<path>/<to>/<filename>
	parts := strings.SplitN(u.Path, "/", 3)

	// ensure org, repo and filename parts exist
	if len(parts) < 3 {
		return &registry.Source{}, fmt.Errorf("invalid template source %s, must contain org/repo/path_to_template", path)
	}

	// check for reference provided in filename:
	// * <org>/<repo>/<filename>@<reference>
	// * <org>/<repo>/<path>/<to>/<filename>@<reference>
	if strings.Contains(parts[2], "@") {
		// capture the filename and reference
		refParts := strings.Split(parts[2], "@")
		// set the filename
		parts[2] = refParts[0]
		// set the reference
		ref = refParts[1]
	}

	return &registry.Source{
		Host: u.Host,
		Org:  parts[0],
		Repo: parts[1],
		Name: parts[2],
		Ref:  ref,
	}, nil
}
