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

const (
	defaultPage       = 1
	defaultPerPage    = 10
	defaultMinPerPage = 1
	defaultMaxPerPage = 100
)

// PaginationOpt represents a configuration option for parsing pagination.
type PaginationOpt func(*paginationConfig) error

type paginationConfig struct {
	DefaultPage    int
	DefaultPerPage int
	MinPerPage     int
	MaxPerPage     int
	Errorf         func(string, error) error
}

func defaultPaginationConfig() paginationConfig {
	return paginationConfig{
		DefaultPage:    defaultPage,
		DefaultPerPage: defaultPerPage,
		MinPerPage:     defaultMinPerPage,
		MaxPerPage:     defaultMaxPerPage,
		Errorf:         defaultPaginationErrorf,
	}
}

// WithDefaultPage sets the default page value.
func WithDefaultPage(page int) PaginationOpt {
	return func(cfg *paginationConfig) error {
		if page <= 0 {
			return fmt.Errorf("default page must be positive")
		}

		cfg.DefaultPage = page
		return nil
	}
}

// WithDefaultPerPage sets the default per-page value.
func WithDefaultPerPage(perPage int) PaginationOpt {
	return func(cfg *paginationConfig) error {
		if perPage <= 0 {
			return fmt.Errorf("default per_page must be positive")
		}

		cfg.DefaultPerPage = perPage
		return nil
	}
}

// WithPerPageMin sets the minimum per-page value.
func WithPerPageMin(minPerPage int) PaginationOpt {
	return func(cfg *paginationConfig) error {
		if minPerPage <= 0 {
			return fmt.Errorf("min per_page must be positive")
		}

		cfg.MinPerPage = minPerPage
		return nil
	}
}

// WithPerPageMax sets the maximum per-page value.
func WithPerPageMax(maxPerPage int) PaginationOpt {
	return func(cfg *paginationConfig) error {
		if maxPerPage <= 0 {
			return fmt.Errorf("max per_page must be positive")
		}

		cfg.MaxPerPage = maxPerPage
		return nil
	}
}

// WithPaginationErrorf sets the error formatter used when parsing fails.
func WithPaginationErrorf(errorf func(string, error) error) PaginationOpt {
	return func(cfg *paginationConfig) error {
		if errorf == nil {
			return fmt.Errorf("pagination error formatter cannot be nil")
		}

		cfg.Errorf = errorf
		return nil
	}
}

// ParsePagination parses the page and per_page query parameters from a request.
func ParsePagination(c *gin.Context, opts ...PaginationOpt) (Pagination, error) {
	cfg := defaultPaginationConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if err := opt(&cfg); err != nil {
			return Pagination{}, err
		}
	}

	if cfg.MinPerPage <= 0 {
		return Pagination{}, fmt.Errorf("min per_page must be positive")
	}

	if cfg.MaxPerPage < cfg.MinPerPage {
		return Pagination{}, fmt.Errorf("max per_page must be greater than or equal to min per_page")
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(cfg.DefaultPage)))
	if err != nil {
		return Pagination{}, cfg.Errorf("page", err)
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", strconv.Itoa(cfg.DefaultPerPage)))
	if err != nil {
		return Pagination{}, cfg.Errorf("per_page", err)
	}

	perPage = max(cfg.MinPerPage, min(cfg.MaxPerPage, perPage))

	return Pagination{Page: page, PerPage: perPage}, nil
}

func defaultPaginationErrorf(param string, err error) error {
	return fmt.Errorf("unable to convert %s query parameter: %w", param, err)
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
