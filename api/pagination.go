// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Pagination holds basic information pertaining to paginating results
type Pagination struct {
	PerPage int
	Page    int
	Total   int64
}

// HeaderLink will hold the information needed to form a link element in the header
type HeaderLink map[string]int

// SetHeaderLink sets the Link HTTP header element to provide clients with paging information
// refs:
// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Link
// - https://tools.ietf.org/html/rfc5988
func (p *Pagination) SetHeaderLink(c *gin.Context) {
	l := []string{}
	r := c.Request

	hl := HeaderLink{
		"first": 1,
		"last":  p.TotalPages(),
		"next":  p.NextPage(),
		"prev":  p.PrevPage(),
	}

	// don't return link info if there is no pagination
	if !p.HasPages() {
		return
	}

	// drop first, prev on the first page
	if p.Page == 1 {
		delete(hl, "first")
		delete(hl, "prev")
	}

	// drop last, next on the last page
	if p.Page == p.TotalPages() {
		delete(hl, "last")
		delete(hl, "next")
	}

	for rel, page := range hl {
		ls := fmt.Sprintf(`<%s://%s%s?per_page=%d&page=%d>; rel="%s"`, resolveScheme(r), r.Host, r.URL.Path, p.PerPage, page, rel)
		l = append(l, ls)
	}

	c.Header("X-Total-Count", strconv.FormatInt(p.Total, 10))
	c.Header("Link", strings.Join(l, ", "))
}

// NextPage returns the next page number
func (p *Pagination) NextPage() int {
	if !p.HasNext() {
		return p.Page
	}

	return p.Page + 1
}

// PrevPage returns the previous page number
func (p *Pagination) PrevPage() int {
	if !p.HasPrev() {
		return 1
	}

	return p.Page - 1
}

// HasPrev will return true if there is a previous page
func (p *Pagination) HasPrev() bool {
	return p.Page > 1
}

// HasNext will return true if there is a next page
func (p *Pagination) HasNext() bool {
	return p.Page < p.TotalPages()
}

// HasPages returns true if there is need to deal with pagination
func (p *Pagination) HasPages() bool {
	return p.Total > int64(p.PerPage)
}

// TotalPages will return the total number of pages
func (p *Pagination) TotalPages() int {
	n := int(math.Ceil(float64(p.Total) / float64(p.PerPage)))
	if n == 0 {
		n = 1
	}

	return n
}

// resolveScheme is a helper to determine the protocol scheme
// c.Request.URL.Scheme does not seem to reliably provide this
func resolveScheme(r *http.Request) string {
	switch {
	case r.Header.Get("X-Forwarded-Proto") == "https":
		return "https"
	case r.URL.Scheme == "https":
		return "https"
	case r.TLS != nil:
		return "https"
	case strings.HasPrefix(strings.ToLower(r.Proto), "https"):
		return "https"
	default:
		return "http"
	}
}
