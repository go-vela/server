// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	sigV4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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

			err = c.initialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_client_getAwsToken(t *testing.T) {
	defaultPresigned := func() *sigV4.PresignedHTTPRequest {
		return &sigV4.PresignedHTTPRequest{
			Method: http.MethodGet,
			URL:    "https://sts.amazonaws.com/?Action=GetCallerIdentity&Version=2011-06-15",
			SignedHeader: http.Header{
				"Authorization": []string{"AWS4-HMAC-SHA256 Credential=test/20250101/us-east-1/sts/aws4_request, SignedHeaders=host;x-amz-date, Signature=deadbeef"},
				"Host":          []string{"sts.amazonaws.com"},
				"X-Amz-Date":    []string{"20250101T000000Z"},
			},
		}
	}

	tests := []struct {
		name         string
		responseFile string
		responseCode int
		presigner    AWSPresigner
		vaultRole    string
		wantToken    string
		wantTTL      time.Duration
		wantErr      bool
	}{
		{
			name:         "get token success",
			responseFile: "auth-response-success.json",
			responseCode: 200,
			presigner: &mockAWSPresigner{
				mockPresignGetCallerIdentity: func(context.Context, *sts.GetCallerIdentityInput, ...func(*sts.PresignOptions)) (*sigV4.PresignedHTTPRequest, error) {
					return defaultPresigned(), nil
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
			presigner: &mockAWSPresigner{
				mockPresignGetCallerIdentity: func(context.Context, *sts.GetCallerIdentityInput, ...func(*sts.PresignOptions)) (*sigV4.PresignedHTTPRequest, error) {
					return defaultPresigned(), nil
				},
			},
			vaultRole: "local-testing",
			wantErr:   true,
		},
		{
			name:         "get token error no auth values",
			responseFile: "auth-response-error-no-auth-values.json",
			responseCode: 400,
			presigner: &mockAWSPresigner{
				mockPresignGetCallerIdentity: func(context.Context, *sts.GetCallerIdentityInput, ...func(*sts.PresignOptions)) (*sigV4.PresignedHTTPRequest, error) {
					return defaultPresigned(), nil
				},
			},
			vaultRole: "local-testing",
			wantErr:   true,
		},
		{
			name:         "get token error nil secret",
			responseFile: "auth-response-error-nil-secret.json",
			responseCode: 200,
			presigner: &mockAWSPresigner{
				mockPresignGetCallerIdentity: func(context.Context, *sts.GetCallerIdentityInput, ...func(*sts.PresignOptions)) (*sigV4.PresignedHTTPRequest, error) {
					return defaultPresigned(), nil
				},
			},
			vaultRole: "local-testing",
			wantErr:   true,
		},
		{
			name:         "get token aws sts error",
			responseFile: "testdata/auth-response-error-no-auth-values.json",
			responseCode: 400,
			presigner: &mockAWSPresigner{
				mockPresignGetCallerIdentity: func(context.Context, *sts.GetCallerIdentityInput, ...func(*sts.PresignOptions)) (*sigV4.PresignedHTTPRequest, error) {
					return nil, fmt.Errorf("token error")
				},
			},
			vaultRole: "local-testing",
			wantErr:   true,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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

			c.AWS.AWSPresigner = tt.presigner

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

type mockAWSPresigner struct {
	mockPresignGetCallerIdentity func(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.PresignOptions)) (*sigV4.PresignedHTTPRequest, error)
}

func (m *mockAWSPresigner) PresignGetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.PresignOptions)) (*sigV4.PresignedHTTPRequest, error) {
	if m.mockPresignGetCallerIdentity != nil {
		return m.mockPresignGetCallerIdentity(ctx, params, optFns...)
	}

	return nil, fmt.Errorf("mock presign not implemented")
}
