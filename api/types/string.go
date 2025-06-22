// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/base64"
	"strconv"
	"strings"

	json "github.com/ghodss/yaml"
	"go.yaml.in/yaml/v3"
)

// ToString is a helper function to convert
// the provided interface value to a string.
func ToString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case []byte:
		return base64.StdEncoding.EncodeToString(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case []interface{}:
		return unmarshalSlice(v)
	default:
		return unmarshalMap(v)
	}
}

// helper function to unmarshal a parameter in map format.
func unmarshalMap(v interface{}) string {
	yml, err := yaml.Marshal(v)
	if err != nil {
		return err.Error()
	}

	out, err := json.YAMLToJSON(yml)
	if err != nil {
		return err.Error()
	}

	return string(out)
}

// helper function to unmarshal a parameter in slice format.
func unmarshalSlice(v interface{}) string {
	out, err := yaml.Marshal(v)
	if err != nil {
		return err.Error()
	}

	in := []string{}

	err = yaml.Unmarshal(out, &in)
	if err == nil {
		return strings.Join(in, ",")
	}

	out, err = json.YAMLToJSON(out)
	if err != nil {
		return err.Error()
	}

	return string(out)
}
