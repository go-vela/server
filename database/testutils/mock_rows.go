// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"database/sql/driver"
	"reflect"
	"slices"

	"github.com/DATA-DOG/go-sqlmock"
)

func CreateMockRows(data []any) *sqlmock.Rows {
	t := reflect.TypeOf(data[0])

	headers := []string{}

	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("sql")

		if tag == "" || tag == "-" {
			continue
		}

		headers = append(headers, tag)
	}

	rows := sqlmock.NewRows(headers)

	for _, d := range data {
		v := reflect.ValueOf(d)

		values := []driver.Value{}

		for i := range v.NumField() {
			if !slices.Contains(headers, t.Field(i).Tag.Get("sql")) {
				continue
			}

			values = append(values, v.Field(i).Interface())
		}

		rows.AddRow(values...)
	}

	return rows
}
