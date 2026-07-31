// SPDX-License-Identifier: Apache-2.0

package models

import (
	"net/http"
	"time"
)

type TemplateEntry struct {
	ETag      string
	Status    int
	Header    http.Header
	Body      []byte
	UpdatedAt time.Time
}
