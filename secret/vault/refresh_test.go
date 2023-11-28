// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

func Test_client_initialize(t *testing.T) {
	// setup tests
	tests := []struct {
		name            string
		responseFile    string
		responseCode    int
		vaultAuthMethod string
		vaultRole       string
		wantErr         bool
	}{
		{
			name:            "initialize success",
			responseFile:    "auth-response-success.json",
			responseCode:    200,
			vaultAuthMethod: "",
			vaultRole:       "local-testing",
			wantErr:         false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := os.ReadFile(fmt.Sprintf("testdata/refresh/%s", tt.responseFile))
				if err != nil {
					t.Error(err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseCode)
				_, _ = w.Write(data)
			}))
			defer ts.Close()

			c, err := New(
				WithAddress(ts.URL),
				WithAuthMethod(""),
				WithAWSRole(tt.vaultRole),
				WithPrefix(""),
				WithToken(""),
				WithTokenDuration(5*time.Minute),
				WithVersion("2"),
			)
			if err != nil {
				t.Error(err)
			}

			c.config.AuthMethod = tt.vaultAuthMethod
			c.AWS.StsClient = &mockSTSClient{
				mockGetCallerIdentityRequest: func(in *sts.GetCallerIdentityInput) (*request.Request, *sts.GetCallerIdentityOutput) {
					return &request.Request{
						HTTPRequest: &http.Request{
							Host: "sts.amazonaws.com",
							URL:  &url.URL{Host: "sts.amazonaws.com"},
						},
						Body: aws.ReadSeekCloser(strings.NewReader("the body")),
					}, nil
				},
			}

			err = c.initialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_client_getAwsToken(t *testing.T) {
	// setup tests
	tests := []struct {
		name         string
		responseFile string
		responseCode int
		stsClient    stsiface.STSAPI
		vaultRole    string
		wantToken    string
		wantTTL      time.Duration
		wantErr      bool
	}{
		{
			name:         "get token success",
			responseFile: "auth-response-success.json",
			responseCode: 200,
			stsClient: &mockSTSClient{
				mockGetCallerIdentityRequest: func(in *sts.GetCallerIdentityInput) (*request.Request, *sts.GetCallerIdentityOutput) {
					return &request.Request{
						HTTPRequest: &http.Request{
							Host: "sts.amazonaws.com",
							URL:  &url.URL{Host: "sts.amazonaws.com"},
						},
						Body: aws.ReadSeekCloser(strings.NewReader("the body")),
					}, nil
				},
			},
			vaultRole: "local-testing",
			wantToken: "s.5RGnjF5aUbhz2XWWM9nHxO57",
			wantTTL:   1 * time.Minute,
			wantErr:   false,
		},
		{
			name:         "get token error role not found",
			responseFile: "auth-response-error-role-not-found.json",
			responseCode: 400,
			stsClient: &mockSTSClient{
				mockGetCallerIdentityRequest: func(in *sts.GetCallerIdentityInput) (*request.Request, *sts.GetCallerIdentityOutput) {
					return &request.Request{
						HTTPRequest: &http.Request{
							Host: "sts.amazonaws.com",
							URL:  &url.URL{Host: "sts.amazonaws.com"},
						},
						Body: aws.ReadSeekCloser(strings.NewReader("the body")),
					}, nil
				},
			},
			vaultRole: "local-testing",
			wantErr:   true,
		},
		{
			name:         "get token error no auth values",
			responseFile: "auth-response-error-no-auth-values.json",
			responseCode: 400,
			stsClient: &mockSTSClient{
				mockGetCallerIdentityRequest: func(in *sts.GetCallerIdentityInput) (*request.Request, *sts.GetCallerIdentityOutput) {
					return &request.Request{
						HTTPRequest: &http.Request{
							Host: "sts.amazonaws.com",
							URL:  &url.URL{Host: "sts.amazonaws.com"},
						},
						Body: aws.ReadSeekCloser(strings.NewReader("the body")),
					}, nil
				},
			},
			vaultRole: "local-testing",
			wantErr:   true,
		},
		{
			name:         "get token error nil secret",
			responseFile: "auth-response-error-nil-secret.json",
			responseCode: 200,
			stsClient: &mockSTSClient{
				mockGetCallerIdentityRequest: func(in *sts.GetCallerIdentityInput) (*request.Request, *sts.GetCallerIdentityOutput) {
					return &request.Request{
						HTTPRequest: &http.Request{
							Host: "sts.amazonaws.com",
							URL:  &url.URL{Host: "sts.amazonaws.com"},
						},
						Body: aws.ReadSeekCloser(strings.NewReader("the body")),
					}, nil
				},
			},
			vaultRole: "local-testing",
			wantErr:   true,
		},
		{
			name:         "get token aws sts error",
			responseFile: "testdata/auth-response-error-no-auth-values.json",
			responseCode: 400,
			stsClient: &mockSTSClient{
				mockGetCallerIdentityRequest: func(in *sts.GetCallerIdentityInput) (*request.Request, *sts.GetCallerIdentityOutput) {
					return &request.Request{
						HTTPRequest: &http.Request{
							Host: "sts.amazonaws.com",
							URL:  &url.URL{Host: "sts.amazonaws.com"},
						},
						Body:  aws.ReadSeekCloser(strings.NewReader("the body")),
						Error: awserr.New(sts.ErrCodeExpiredTokenException, "token error", fmt.Errorf("token error")),
					}, nil
				},
			},
			vaultRole: "local-testing",
			wantErr:   true,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := os.ReadFile(fmt.Sprintf("testdata/refresh/%s", tt.responseFile))
				if err != nil {
					t.Error(err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseCode)
				_, _ = w.Write(data)
			}))
			defer ts.Close()
			c, err := New(
				WithAddress(ts.URL),
				WithAuthMethod(""),
				WithAWSRole(tt.vaultRole),
				WithPrefix(""),
				WithToken(""),
				WithTokenDuration(5*time.Minute),
				WithVersion("2"),
			)
			if err != nil {
				t.Error(err)
			}
			c.AWS.StsClient = tt.stsClient

			gotToken, gotTTL, err := c.getAwsToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("getAwsToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotToken != tt.wantToken {
				t.Errorf("getAwsToken() gotToken = %v, wantToken %v", gotToken, tt.wantToken)
			}

			if gotTTL != tt.wantTTL {
				t.Errorf("getAwsToken() gotTTL = %v, wantTTL %v", gotTTL, tt.wantTTL)
			}
		})
	}
}

type mockSTSClient struct {
	stsiface.STSAPI
	mockGetCallerIdentityRequest func(in *sts.GetCallerIdentityInput) (*request.Request, *sts.GetCallerIdentityOutput)
}

func (m *mockSTSClient) GetCallerIdentityRequest(in *sts.GetCallerIdentityInput) (*request.Request, *sts.GetCallerIdentityOutput) {
	return m.mockGetCallerIdentityRequest(in)
}
