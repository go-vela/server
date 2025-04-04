// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Pagination holds basic information pertaining to paginating results.
type Pagination struct {
	PerPage int
	Page    int
	Results int
}

// HeaderLink will hold the information needed to form a link element in the header.
type HeaderLink map[string]int

// SetHeaderLink sets the Link HTTP header element to provide clients with paging information
// refs:
// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Link
// - https://tools.ietf.org/html/rfc5988
func (p *Pagination) SetHeaderLink(c *gin.Context) {
	l := []string{}
	r := c.Request

	// grab the current query params
	q := r.URL.Query()

	hl := HeaderLink{
		"first": 1,
		"next":  p.NextPage(),
		"prev":  p.PrevPage(),
	}

	// don't return link info if there is no pagination
	if !p.HasPages() {
		return
	}

	// reset per config
	q.Set("per_page", strconv.Itoa(p.PerPage))

	// drop first, prev on the first page
	if p.Page == 1 {
		delete(hl, "first")
		delete(hl, "prev")
	}

	// drop last, next on the last page
	if p.Results < p.PerPage {
		delete(hl, "last")
		delete(hl, "next")
	}

	// loop over the fields that make up the header links
	for rel, page := range hl {
		// set the page info for the current field
		q.Set("page", strconv.Itoa(page))

		ls := fmt.Sprintf(
			`<%s://%s%s?%s>; rel="%s"`,
			resolveScheme(r),
			r.Host,
			r.URL.Path,
			q.Encode(),
			rel,
		)

		l = append(l, ls)
	}

	c.Header("Link", strings.Join(l, ", "))
}

// NextPage returns the next page number.
func (p *Pagination) NextPage() int {
	if !p.HasNext() {
		return p.Page
	}

	return p.Page + 1
}

// PrevPage returns the previous page number.
func (p *Pagination) PrevPage() int {
	if !p.HasPrev() {
		return 1
	}

	return p.Page - 1
}

// HasPrev will return true if there is a previous page.
func (p *Pagination) HasPrev() bool {
	return p.Page > 1
}

// HasNext will return true if there is a next page.
func (p *Pagination) HasNext() bool {
	return p.PerPage == p.Results
}

// HasPages returns true if there is need to deal with pagination.
func (p *Pagination) HasPages() bool {
	return !(p.Page == 1 && p.Results < p.PerPage)
}

// resolveScheme is a helper to determine the protocol scheme
// c.Request.URL.Scheme does not seem to reliably provide this.
//
//nolint:goconst // ignore making constant for https
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
