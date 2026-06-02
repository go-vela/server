// SPDX-License-Identifier: Apache-2.0

package github

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/cache"
	"github.com/go-vela/server/cache/models"
)

type ContentsTransport struct {
	Base  http.RoundTripper
	Store cache.Service
}

// NewContentsTransport creates a new http.RoundTripper that caches responses from GitHub API requests for template contents.
func NewContentsTransport(store cache.Service, base http.RoundTripper) *ContentsTransport {
	if base == nil {
		base = http.DefaultTransport
	}

	return &ContentsTransport{Base: base, Store: store}
}

func (t *ContentsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Store == nil {
		return t.Base.RoundTrip(req)
	}

	// generate cache key by hashing the request URL
	sum := sha256.Sum256([]byte(req.URL.String()))
	key := "github:contents:" + hex.EncodeToString(sum[:])

	entry, err := t.Store.GetTemplateContents(req.Context(), key)
	if err != nil {
		logrus.Errorf("failed to get template contents from cache for key %s: %v", key, err)

		// continue to request if cache is acting up
		return t.Base.RoundTrip(req)
	}

	clonedReq := req.Clone(req.Context())

	if entry != nil && entry.ETag != "" {
		logrus.Debug("adding If-None-Match header to request for cache validation")

		clonedReq.Header = clonedReq.Header.Clone()
		clonedReq.Header.Set("If-None-Match", entry.ETag)
	}

	resp, err := t.Base.RoundTrip(clonedReq)
	if err != nil {
		return nil, err
	}

	// true cache hit
	if resp.StatusCode == http.StatusNotModified && entry != nil {
		_ = resp.Body.Close()

		err = t.Store.ExtendTemplateExpiry(req.Context(), key)
		if err != nil {
			logrus.Errorf("failed to extend cache expiry for key %s: %v", key, err)
		}

		return &http.Response{
			Status:        fmt.Sprintf("%d %s", entry.Status, http.StatusText(entry.Status)),
			StatusCode:    entry.Status,
			Header:        entry.Header.Clone(),
			Body:          io.NopCloser(bytes.NewReader(entry.Body)),
			ContentLength: int64(len(entry.Body)),
			Request:       req,
		}, nil
	}

	// if success, cache response
	if resp.StatusCode == http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if readErr != nil {
			return nil, readErr
		}

		etag := resp.Header.Get("Etag")
		if etag != "" {
			err = t.Store.StoreTemplateContents(req.Context(), key, &models.TemplateEntry{
				ETag:      etag,
				Status:    resp.StatusCode,
				Header:    resp.Header.Clone(),
				Body:      append([]byte(nil), body...),
				UpdatedAt: time.Now().UTC(),
			})

			if err != nil {
				logrus.Errorf("failed to store template contents in cache for key %s: %v", key, err)
			}
		}

		resp.Body = io.NopCloser(bytes.NewReader(body))
		resp.ContentLength = int64(len(body))
	}

	return resp, nil
}
