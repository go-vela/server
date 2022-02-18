// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	// ErrUnableToConvertStarlark defines the error type when the
	// toStarlark cannot convert the provided value.
	ErrUnableToConvertStarlark = errors.New("unable to convert to starlark type")

	// ErrUnableToConvertJSON defines the error type when the
	// writeJSON cannot convert the provided value.
	ErrUnableToConvertJSON = errors.New("unable to convert to json")
)

// toStarlark takes an value as an interface an
// will return the comparable Starlark type.
//
// This code is under copyright (full attribution in NOTICE) and is from:
//
// https://github.com/wonderix/shalm/blob/899b8f7787883d40619eefcc39bd12f42a09b5e7/pkg/shalm/convert.go#L14-L85
//
// nolint: gocyclo // ignore complexity
func toStarlark(value interface{}) (starlark.Value, error) {
	logrus.Tracef("converting %v to starlark type", value)

	if value == nil {
		return starlark.None, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.String:
		return starlark.String(v.String()), nil
	case reflect.Bool:
		return starlark.Bool(v.Bool()), nil
	case reflect.Int:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int16:
		return starlark.MakeInt64(v.Int()), nil
	case reflect.Uint:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uint16:
		return starlark.MakeUint64(v.Uint()), nil
	case reflect.Float32:
		return starlark.Float(v.Float()), nil
	case reflect.Float64:
		return starlark.Float(v.Float()), nil
	case reflect.Slice:
		if b, ok := value.([]byte); ok {
			return starlark.String(string(b)), nil
		}

		a := make([]starlark.Value, 0)

		for i := 0; i < v.Len(); i++ {
			val, err := toStarlark(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}

			a = append(a, val)
		}

		return starlark.Tuple(a), nil
	case reflect.Ptr:
		val, err := toStarlark(v.Elem().Interface())
		if err != nil {
			return nil, err
		}

		return val, nil
	case reflect.Map:
		d := starlark.NewDict(16)

		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)

			keyValue, err := toStarlark(key.Interface())
			if err != nil {
				return nil, err
			}

			kv, err := toStarlark(strct.Interface())
			if err != nil {
				return nil, err
			}

			err = d.SetKey(keyValue, kv)
			if err != nil {
				return nil, err
			}
		}

		return d, nil
	case reflect.Struct:
		ios, ok := value.(intstr.IntOrString)
		if ok {
			switch ios.Type {
			case intstr.String:
				return starlark.String(ios.StrVal), nil
			case intstr.Int:
				return starlark.MakeInt(int(ios.IntVal)), nil
			}
		} else {
			data, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}

			var m map[string]interface{}
			err = json.Unmarshal(data, &m)
			if err != nil {
				return nil, err
			}

			return toStarlark(m)
		}
	}

	return nil, fmt.Errorf("%s: %v", ErrUnableToConvertStarlark, value)
}

// writeJSON takes an starlark input and return the valid JSON
// for the specific type.
//
// This code is under copyright (full attribution in NOTICE) and is from:
//
// https://github.com/drone/drone-cli/blob/master/drone/starlark/starlark.go#L214-L274
//
// Note: we are using logrus log unchecked errors that the original implementation ignored.
// if/when we try to return values it breaks the recursion. Panics were swapped to error
// returns from implementation.
//
// nolint: gocyclo // ignore cyclomatic complexity
func writeJSON(out *bytes.Buffer, v starlark.Value) error {
	logrus.Tracef("converting %v to JSON", v)

	if marshaler, ok := v.(json.Marshaler); ok {
		jsonData, err := marshaler.MarshalJSON()
		if err != nil {
			return err
		}

		_, err = out.Write(jsonData)
		if err != nil {
			logrus.Error(err)
		}

		return nil
	}

	switch v := v.(type) {
	case starlark.NoneType:
		_, err := out.WriteString("null")
		if err != nil {
			logrus.Error(err)
		}
	case starlark.Bool:
		_, err := fmt.Fprintf(out, "%t", v)
		if err != nil {
			logrus.Error(err)
		}
	case starlark.Int:
		_, err := out.WriteString(v.String())
		if err != nil {
			logrus.Error(err)
		}
	case starlark.Float:
		_, err := fmt.Fprintf(out, "%g", v)
		if err != nil {
			logrus.Error(err)
		}
	case starlark.String:
		s := string(v)

		if goQuoteIsSafe(s) {
			fmt.Fprintf(out, "%q", s)
		} else {
			// vanishingly rare for text strings
			data, err := json.Marshal(s)
			if err != nil {
				logrus.Error(err)
			}

			_, err = out.Write(data)
			if err != nil {
				logrus.Error(err)
			}
		}
	case starlark.Indexable: // Tuple, List
		err := out.WriteByte('[')
		if err != nil {
			logrus.Error(err)
		}

		for i, n := 0, starlark.Len(v); i < n; i++ {
			if i > 0 {
				_, err := out.WriteString(", ")
				if err != nil {
					logrus.Error(err)
				}
			}

			err := writeJSON(out, v.Index(i))
			if err != nil {
				return err
			}
		}

		err = out.WriteByte(']')
		if err != nil {
			logrus.Error(err)
		}
	case *starlark.Dict:
		err := out.WriteByte('{')
		if err != nil {
			logrus.Error(err)
		}

		for i, itemPair := range v.Items() {
			key, value := itemPair[0], itemPair[1]

			if i > 0 {
				_, err := out.WriteString(", ")
				if err != nil {
					logrus.Error(err)
				}
			}

			err := writeJSON(out, key)
			if err != nil {
				return err
			}

			_, err = out.WriteString(": ")
			if err != nil {
				logrus.Error(err)
			}

			err = writeJSON(out, value)
			if err != nil {
				return err
			}
		}

		err = out.WriteByte('}')
		if err != nil {
			logrus.Error(err)
		}
	default:
		return fmt.Errorf("%s: %v", ErrUnableToConvertJSON, v)
	}

	return nil
}

// goQuoteIsSafe takes a string and checks if is safely quoted
//
// This code is under copyright (full attribution in NOTICE) and is from:
// https://github.com/drone/drone-cli/blob/master/drone/starlark/starlark.go#L276-L285
func goQuoteIsSafe(s string) bool {
	for _, r := range s {
		// JSON doesn't like Go's \xHH escapes for ASCII control codes,
		// nor its \UHHHHHHHH escapes for runes >16 bits.
		if r < 0x20 || r >= 0x10000 {
			return false
		}
	}

	return true
}
