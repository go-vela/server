// SPDX-License-Identifier: Apache-2.0

package github

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v81/github"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	acceptHeader = "application/vnd.github.v3+json"
)

// AppsTransport provides a http.RoundTripper by wrapping an existing
// http.RoundTripper and provides GitHub Apps authentication as a GitHub App.
//
// Client can also be overwritten, and is useful to change to one which
// provides retry logic if you do experience retryable errors.
//
// See https://developer.github.com/apps/building-integrations/setting-up-and-registering-github-apps/about-authentication-options-for-github-apps/
type AppsTransport struct {
	BaseURL string            // BaseURL is the scheme and host for GitHub API, defaults to https://api.github.com
	Client  HTTPClient        // Client to use to refresh tokens, defaults to http.Client with provided transport
	tr      http.RoundTripper // tr is the underlying roundtripper being wrapped
	signer  Signer            // signer signs JWT tokens.
	appID   int64             // appID is the GitHub App's ID
}

// newGitHubAppTransport creates a new GitHub App transport for authenticating as the GitHub App.
func (c *Client) newGitHubAppTransport(appID int64, baseURL string, privateKey *rsa.PrivateKey) *AppsTransport {
	transport := c.newAppsTransportFromPrivateKey(http.DefaultTransport, appID, privateKey)
	transport.BaseURL = baseURL

	// apply tracing to the transport
	if c.Tracing.Config.EnableTracing {
		transport.tr = otelhttp.NewTransport(
			transport.tr,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithoutSubSpans())
			}),
		)
	}

	return transport
}

// newAppsTransportFromPrivateKey returns an AppsTransport using a crypto/rsa.(*PrivateKey).
func (c *Client) newAppsTransportFromPrivateKey(tr http.RoundTripper, appID int64, key *rsa.PrivateKey) *AppsTransport {
	return &AppsTransport{
		BaseURL: defaultAPI,
		Client:  &http.Client{Transport: tr},
		tr:      tr,
		signer:  NewRSASigner(jwt.SigningMethodRS256, key),
		appID:   appID,
	}
}

// RoundTrip implements http.RoundTripper interface.
func (t *AppsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// GitHub rejects expiry and issue timestamps that are not an integer,
	// while the jwt-go library serializes to fractional timestamps
	// then truncate them before passing to jwt-go.
	iss := time.Now().Add(-30 * time.Second).Truncate(time.Second)
	exp := iss.Add(2 * time.Minute)
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(iss),
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    strconv.FormatInt(t.appID, 10),
	}

	ss, err := t.signer.Sign(claims)
	if err != nil {
		return nil, fmt.Errorf("could not sign jwt: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+ss)
	req.Header.Add("Accept", acceptHeader)

	return t.tr.RoundTrip(req)
}

// Transport provides a http.RoundTripper by wrapping an existing
// http.RoundTripper and provides GitHub Apps authentication as an installation.
//
// Client can also be overwritten, and is useful to change to one which
// provides retry logic if you do experience retryable errors.
//
// See https://developer.github.com/apps/building-integrations/setting-up-and-registering-github-apps/about-authentication-options-for-github-apps/
type Transport struct {
	BaseURL                  string                           // BaseURL is the scheme and host for GitHub API, defaults to https://api.github.com
	Client                   HTTPClient                       // Client to use to refresh tokens, defaults to http.Client with provided transport
	tr                       http.RoundTripper                // tr is the underlying roundtripper being wrapped
	installationID           int64                            // installationID is the GitHub App Installation ID
	InstallationTokenOptions *github.InstallationTokenOptions // parameters restrict a token's access
	appsTransport            *AppsTransport

	mu    *sync.Mutex
	token *accessToken // the installation's access token
}

// accessToken is an installation access token response from GitHub.
type accessToken struct {
	Token        string                         `json:"token"`
	ExpiresAt    time.Time                      `json:"expires_at"`
	Permissions  github.InstallationPermissions `json:"permissions,omitempty"`
	Repositories []github.Repository            `json:"repositories,omitempty"`
}

var _ http.RoundTripper = &Transport{}

// HTTPClient is a HTTP client which sends a http.Request and returns a http.Response
// or an error.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// RoundTrip implements http.RoundTripper interface.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false

	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}

	token, err := t.Token(req.Context())
	if err != nil {
		return nil, err
	}

	creq := cloneRequest(req)
	creq.Header.Set("Authorization", "token "+token)

	if creq.Header.Get("Accept") == "" {
		creq.Header.Add("Accept", acceptHeader)
	}

	reqBodyClosed = true

	return t.tr.RoundTrip(creq)
}

// getRefreshTime returns the time when the token should be refreshed.
func (at *accessToken) getRefreshTime() time.Time {
	return at.ExpiresAt.Add(-time.Minute)
}

// isExpired checks if the access token is expired.
func (at *accessToken) isExpired() bool {
	return at == nil || at.getRefreshTime().Before(time.Now())
}

// Token checks the active token expiration and renews if necessary. Token returns
// a valid access token. If renewal fails an error is returned.
func (t *Transport) Token(ctx context.Context) (string, error) {
	t.mu.Lock()

	defer t.mu.Unlock()

	if t.token.isExpired() {
		// token is not set or expired/nearly expired, so refresh
		if err := t.refreshToken(ctx); err != nil {
			return "", fmt.Errorf("could not refresh installation id %v's token: %w", t.installationID, err)
		}
	}

	return t.token.Token, nil
}

// Expiry returns a transport token's expiration time and refresh time. There is a small grace period
// built in where a token will be refreshed before it expires. expiresAt is the actual token expiry,
// and refreshAt is when a call to Token() will cause it to be refreshed.
func (t *Transport) Expiry() (expiresAt time.Time, refreshAt time.Time, err error) {
	if t.token == nil {
		return time.Time{}, time.Time{}, errors.New("Expiry() = unknown, err: nil token")
	}

	return t.token.ExpiresAt, t.token.getRefreshTime(), nil
}

func (t *Transport) refreshToken(ctx context.Context) error {
	// convert InstallationTokenOptions into a ReadWriter to pass as an argument to http.NewRequest
	body, err := GetReadWriter(t.InstallationTokenOptions)
	if err != nil {
		return fmt.Errorf("could not convert installation token parameters into json: %w", err)
	}

	requestURL := fmt.Sprintf("%s/app/installations/%v/access_tokens", strings.TrimRight(t.BaseURL, "/"), t.installationID)

	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, body)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	// set Content and Accept headers
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", acceptHeader)

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	t.appsTransport.BaseURL = t.BaseURL
	t.appsTransport.Client = t.Client

	resp, err := t.appsTransport.RoundTrip(req)
	if err != nil {
		return fmt.Errorf("could not get access_tokens from GitHub API for installation ID %v: %w", t.installationID, err)
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("received non 2xx response status %q when fetching %v", resp.Status, req.URL)
	}

	// closing body late, to provide caller a chance to inspect body in an error / non-200 response status situation
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&t.token)
}

// GetReadWriter converts a body interface into an io.ReadWriter object.
func GetReadWriter(i interface{}) (io.ReadWriter, error) {
	var buf io.ReadWriter

	if i != nil {
		buf = new(bytes.Buffer)

		enc := json.NewEncoder(buf)

		err := enc.Encode(i)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	_r := new(http.Request)

	*_r = *r

	// deep copy of the Header
	_r.Header = make(http.Header, len(r.Header))

	for k, s := range r.Header {
		_r.Header[k] = append([]string(nil), s...)
	}

	return _r
}

// Signer is a JWT token signer. This is a wrapper around [jwt.SigningMethod] with predetermined
// key material.
type Signer interface {
	// sign the given claims and returns a JWT token string, as specified
	// by [jwt.Token.SignedString]
	Sign(claims jwt.Claims) (string, error)
}

// RSASigner signs JWT tokens using RSA keys.
type RSASigner struct {
	method *jwt.SigningMethodRSA
	key    *rsa.PrivateKey
}

// NewRSASigner creates a new RSASigner with the given RSA key.
func NewRSASigner(method *jwt.SigningMethodRSA, key *rsa.PrivateKey) *RSASigner {
	return &RSASigner{
		method: method,
		key:    key,
	}
}

// Sign signs the JWT claims with the RSA key.
func (s *RSASigner) Sign(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(s.method, claims).SignedString(s.key)
}

// AppsTransportOption is a func option for configuring an AppsTransport.
type AppsTransportOption func(*AppsTransport)

// WithSigner configures the AppsTransport to use the given Signer for generating JWT tokens.
func WithSigner(signer Signer) AppsTransportOption {
	return func(at *AppsTransport) {
		at.signer = signer
	}
}

// NewTestAppClient creates a new AppsTransport for testing purposes.
func NewTestAppClient(baseURL string) *github.Client {
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	client, err := github.NewClient(
		&http.Client{
			Transport: &AppsTransport{
				BaseURL: baseURL,
				Client:  &http.Client{Transport: http.DefaultTransport},
				tr:      http.DefaultTransport,
				signer: &RSASigner{
					method: jwt.SigningMethodRS256,
					key:    pk,
				},
				appID: 1,
			},
		}).
		WithEnterpriseURLs(baseURL, baseURL)
	if err != nil {
		return nil
	}

	return client
}
