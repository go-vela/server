// SPDX-License-Identifier: Apache-2.0

package raw

import (
	"encoding/json"
	"errors"
)

// StringSlice represents a string or an array of strings.
type StringSlice []string

// UnmarshalJSON implements the Unmarshaler interface for the StringSlice type.
func (s *StringSlice) UnmarshalJSON(b []byte) error {
	// return nil if no input is provided
	if len(b) == 0 {
		return nil
	}

	// json string we try unmarshalling to
	jsonString := ""

	// attempt to unmarshal as a string type
	err := json.Unmarshal(b, &jsonString)
	if err == nil {
		// overwrite existing StringSlice
		*s = []string{jsonString}

		return nil
	}

	// json slice we try unmarshalling to
	jsonSlice := []string{}

	// attempt to unmarshal as a string slice type
	err = json.Unmarshal(b, &jsonSlice)
	if err == nil {
		// overwrite existing StringSlice
		*s = jsonSlice

		return nil
	}

	return errors.New("unable to unmarshal into StringSlice")
}

// UnmarshalYAML implements the Unmarshaler interface for the StringSlice type.
func (s *StringSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// yaml string we try unmarshalling to
	yamlString := ""

	// attempt to unmarshal as a string type
	err := unmarshal(&yamlString)
	if err == nil {
		// overwrite existing StringSlice
		*s = []string{yamlString}

		return nil
	}

	yamlSlice := []string{}

	// attempt to unmarshal as a string slice type
	err = unmarshal(&yamlSlice)
	if err == nil {
		// overwrite existing StringSlice
		*s = yamlSlice

		return nil
	}

	return errors.New("unable to unmarshal into StringSlice")
}
