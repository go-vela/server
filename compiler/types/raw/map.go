// SPDX-License-Identifier: Apache-2.0

package raw

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

// StringSliceMap represents an array of strings or a map of strings.
type StringSliceMap map[string]string

// Value returns the map in JSON format.
func (s StringSliceMap) Value() (driver.Value, error) {
	value, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return string(value), nil
}

// Scan decodes the JSON string into map[string]string.
func (s *StringSliceMap) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}

	return json.Unmarshal([]byte(b), &s)
}

// UnmarshalJSON implements the Unmarshaler interface for the StringSlice type.
func (s *StringSliceMap) UnmarshalJSON(b []byte) error {
	// return nil if no input is provided
	if len(b) == 0 {
		return nil
	}

	// target map we want to return
	targetMap := map[string]string{}

	// json slice we try unmarshalling to
	jsonSlice := StringSlice{}

	// attempt to unmarshal as a string slice type
	err := json.Unmarshal(b, &jsonSlice)
	if err == nil {
		// iterate through each element in the json slice
		for _, v := range jsonSlice {
			// split each slice element into key/value pairs
			kvPair := strings.SplitN(v, "=", 2)

			if len(kvPair) != 2 {
				return errors.New("unable to unmarshal into StringSliceMap")
			}

			// append each key/value pair to our target map
			targetMap[kvPair[0]] = kvPair[1]
		}

		// overwrite existing StringSliceMap
		*s = targetMap

		return nil
	}

	// json map we try unmarshalling to
	jsonMap := map[string]string{}

	// attempt to unmarshal as map of strings
	err = json.Unmarshal(b, &jsonMap)
	if err == nil {
		// iterate through each item in the json map
		for k, v := range jsonMap {
			// append each key/value pair to our target map
			targetMap[k] = v
		}

		// overwrite existing StringSliceMap
		*s = targetMap

		return nil
	}

	return errors.New("unable to unmarshal into StringSliceMap")
}

// UnmarshalYAML implements the Unmarshaler interface for the StringSliceMap type.
func (s *StringSliceMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// target map we want to return
	targetMap := map[string]string{}

	// yaml slice we try unmarshalling to
	yamlSlice := StringSlice{}

	// attempt to unmarshal as a string slice type
	err := unmarshal(&yamlSlice)
	if err == nil {
		// iterate through each element in the yaml slice
		for _, v := range yamlSlice {
			// split each slice element into key/value pairs
			kvPair := strings.SplitN(v, "=", 2)

			if len(kvPair) != 2 {
				return errors.New("unable to unmarshal into StringSliceMap")
			}

			// append each key/value pair to our target map
			targetMap[kvPair[0]] = kvPair[1]
		}

		// overwrite existing StringSliceMap
		*s = targetMap

		return nil
	}

	// yaml map we try unmarshalling to
	yamlMap := map[string]string{}

	// attempt to unmarshal as map of strings
	err = unmarshal(&yamlMap)
	if err == nil {
		// iterate through each item in the yaml map
		for k, v := range yamlMap {
			// append each key/value pair to our target map
			targetMap[k] = v
		}

		// overwrite existing StringSliceMap
		*s = targetMap

		return nil
	}

	return errors.New("unable to unmarshal into StringSliceMap")
}
