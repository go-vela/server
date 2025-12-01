// SPDX-License-Identifier: Apache-2.0

package native

import (
	"fmt"
	"io"
	"os"

	"github.com/go-vela/server/compiler/template/native"
	"github.com/go-vela/server/compiler/template/starlark"
	typesRaw "github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/internal"
)

// ParseRaw converts an object to a string.
func (c *Client) ParseRaw(v interface{}) (string, error) {
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
func (c *Client) Parse(v interface{}, pipelineType string, template *yaml.Template) (*yaml.Build, []byte, []string, error) {
	var (
		p        *yaml.Build
		warnings []string
		raw      []byte
	)

	switch pipelineType {
	case constants.PipelineTypeGo, constants.PipelineTypeGoAlt:
		// expand the base configuration
		parsedRaw, err := c.ParseRaw(v)
		if err != nil {
			return nil, nil, nil, err
		}

		// capture the raw pipeline configuration
		raw = []byte(parsedRaw)

		p, warnings, err = native.RenderBuild(template.Name, parsedRaw, c.EnvironmentBuild(), template.Variables)
		if err != nil {
			return nil, raw, nil, err
		}
	case constants.PipelineTypeStarlark:
		// expand the base configuration
		parsedRaw, err := c.ParseRaw(v)
		if err != nil {
			return nil, nil, nil, err
		}

		// capture the raw pipeline configuration
		raw = []byte(parsedRaw)

		p, warnings, err = starlark.RenderBuild(template.Name, parsedRaw, c.EnvironmentBuild(), template.Variables, c.GetStarlarkExecLimit())
		if err != nil {
			return nil, raw, nil, err
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
			return nil, nil, nil, fmt.Errorf("unable to parse yaml: unrecognized type %T", v)
		}
	default:
		return nil, nil, nil, fmt.Errorf("unable to parse config: unrecognized pipeline_type of %s", c.repo.GetPipelineType())
	}

	// initializing Environment to prevent nil error
	// as it may be modified later via templates, if
	// none are defined in the base pipeline
	if p.Environment == nil {
		p.Environment = typesRaw.StringSliceMap{}
	}

	return p, raw, warnings, nil
}

// ParseBytes converts a byte slice to a yaml configuration.
func ParseBytes(data []byte) (*yaml.Build, []byte, []string, error) {
	config, warnings, err := internal.ParseYAML(data, "")
	if err != nil {
		return nil, nil, nil, err
	}

	// initializing Environment to prevent nil error
	// as it may be modified later via templates, if
	// none are defined in the base pipeline
	if config.Environment == nil {
		config.Environment = typesRaw.StringSliceMap{}
	}

	return config, data, warnings, nil
}

// ParseFile converts an os.File into a yaml configuration.
func ParseFile(f *os.File) (*yaml.Build, []byte, []string, error) {
	return ParseReader(f)
}

// ParseFileRaw converts an os.File into a string.
func ParseFileRaw(f *os.File) (string, error) {
	return ParseReaderRaw(f)
}

// ParsePath converts a file path into a yaml configuration.
func ParsePath(p string) (*yaml.Build, []byte, []string, error) {
	// open the file for reading
	f, err := os.Open(p)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to open yaml file %s: %w", p, err)
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
func ParseReader(r io.Reader) (*yaml.Build, []byte, []string, error) {
	// read all the bytes from the reader
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to read bytes for yaml: %w", err)
	}

	return ParseBytes(data)
}

// ParseReaderRaw converts an io.Reader into a yaml configuration.
func ParseReaderRaw(r io.Reader) (string, error) {
	// read all the bytes from the reader
	b, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("unable to read bytes for yaml: %w", err)
	}

	return string(b), nil
}

// ParseString converts a string into a yaml configuration.
func ParseString(s string) (*yaml.Build, []byte, []string, error) {
	return ParseBytes([]byte(s))
}
