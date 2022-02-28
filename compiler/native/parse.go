// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/go-vela/server/compiler/template/native"
	"github.com/go-vela/server/compiler/template/starlark"
	"github.com/go-vela/types/constants"
	types "github.com/go-vela/types/yaml"

	"github.com/buildkite/yaml"
)

// ParseRaw converts an object to a string.
func (c *client) ParseRaw(v interface{}) (string, error) {
	switch v := v.(type) {
	case []byte:
		return string(v), nil
	case *os.File:
		return ParseFileRaw(v)
	case io.Reader:
		return ParseReaderRaw(v)
	case string:
		// check if string is path to file
		_, err := os.Stat(v)
		if err == nil {
			// parse string as path to yaml configuration
			return ParsePathRaw(v)
		}

		// parse string as yaml configuration
		return v, nil
	default:
		return "", fmt.Errorf("unable to parse yaml: unrecognized type %T", v)
	}
}

// Parse converts an object to a yaml configuration.
func (c *client) Parse(v interface{}) (*types.Build, error) {
	var p *types.Build

	switch c.repo.GetPipelineType() {
	case constants.PipelineTypeGo:
		// expand the base configuration
		parsedRaw, err := c.ParseRaw(v)
		if err != nil {
			return nil, err
		}

		p, err = native.RenderBuild(parsedRaw, c.EnvironmentBuild())
		if err != nil {
			return nil, err
		}
	case constants.PipelineTypeStarlark:
		// expand the base configuration
		parsedRaw, err := c.ParseRaw(v)
		if err != nil {
			return nil, err
		}

		p, err = starlark.RenderBuild(parsedRaw, c.EnvironmentBuild())
		if err != nil {
			return nil, err
		}
	case constants.PipelineTypeYAML, "":
		switch v := v.(type) {
		case []byte:
			return ParseBytes(v)
		case *os.File:
			return ParseFile(v)
		case io.Reader:
			return ParseReader(v)
		case string:
			// check if string is path to file
			_, err := os.Stat(v)
			if err == nil {
				// parse string as path to yaml configuration
				return ParsePath(v)
			}

			// parse string as yaml configuration
			return ParseString(v)
		default:
			return nil, fmt.Errorf("unable to parse yaml: unrecognized type %T", v)
		}
	default:
		return nil, fmt.Errorf("unable to parse config: unrecognized pipeline_type of %s", c.repo.GetPipelineType())
	}

	return p, nil
}

// ParseBytes converts a byte slice to a yaml configuration.
func ParseBytes(b []byte) (*types.Build, error) {
	config := new(types.Build)

	// unmarshal the bytes into the yaml configuration
	err := yaml.Unmarshal(b, config)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
	}

	return config, nil
}

// ParseFile converts an os.File into a yaml configuration.
func ParseFile(f *os.File) (*types.Build, error) {
	return ParseReader(f)
}

// ParseFileRaw converts an os.File into a string.
func ParseFileRaw(f *os.File) (string, error) {
	return ParseReaderRaw(f)
}

// ParsePath converts a file path into a yaml configuration.
func ParsePath(p string) (*types.Build, error) {
	// open the file for reading
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("unable to open yaml file %s: %w", p, err)
	}

	defer f.Close()

	return ParseReader(f)
}

// ParsePathRaw converts a file path into a yaml configuration.
func ParsePathRaw(p string) (string, error) {
	// open the file for reading
	f, err := os.Open(p)
	if err != nil {
		return "", fmt.Errorf("unable to open yaml file %s: %w", p, err)
	}

	defer f.Close()

	return ParseReaderRaw(f)
}

// ParseReader converts an io.Reader into a yaml configuration.
func ParseReader(r io.Reader) (*types.Build, error) {
	// read all the bytes from the reader
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read bytes for yaml: %w", err)
	}

	return ParseBytes(b)
}

// ParseReaderRaw converts an io.Reader into a yaml configuration.
func ParseReaderRaw(r io.Reader) (string, error) {
	// read all the bytes from the reader
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("unable to read bytes for yaml: %w", err)
	}

	return string(b), nil
}

// ParseString converts a string into a yaml configuration.
func ParseString(s string) (*types.Build, error) {
	return ParseBytes([]byte(s))
}
