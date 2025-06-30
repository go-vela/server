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

	"github.com/aws/aws-sdk-go-v2/aws"
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
			c.AWS.StsClient = newSuccessfulMockSTSClient()

			err = c.initialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_client_generateAwsAuthHeaders(t *testing.T) {
	// setup tests
	tests := []struct {
		name         string
		responseFile string
		responseCode int
		stsClient    STSClient
		vaultRole    string
		wantToken    string
		wantTTL      time.Duration
		wantErr      bool
	}{
		{
			name:         "get token success - vela aws auth flow",
			responseFile: "auth-response-success.json",
			responseCode: 200,
			stsClient:    newSuccessfulMockSTSClient(),
			vaultRole:    "local-testing",
			wantToken:    "s.5RGnjF5aUbhz2XWWM9nHxO57",
			wantTTL:      1 * time.Minute,
			wantErr:      false,
		},
		{
			name:         "get token error role not found - vela vault config issue",
			responseFile: "auth-response-error-role-not-found.json",
			responseCode: 400,
			stsClient:    newSuccessfulMockSTSClient(),
			vaultRole:    "local-testing",
			wantErr:      true,
		},
		{
			name:         "get token error no auth values - vela vault setup issue",
			responseFile: "auth-response-error-no-auth-values.json",
			responseCode: 400,
			stsClient:    newSuccessfulMockSTSClient(),
			vaultRole:    "local-testing",
			wantErr:      true,
		},
		{
			name:         "get token error nil secret - vela vault response issue",
			responseFile: "auth-response-error-nil-secret.json",
			responseCode: 200,
			stsClient:    newSuccessfulMockSTSClient(),
			vaultRole:    "local-testing",
			wantErr:      true,
		},
		{
			name:         "get token aws sts error - vela credential failure",
			responseFile: "testdata/auth-response-error-no-auth-values.json",
			responseCode: 400,
			stsClient:    newFailingMockSTSClient("AWS credentials error for Vela"),
			vaultRole:    "local-testing",
			wantErr:      true,
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

			c.AWS.StsClient = tt.stsClient

			// create mock AWS config
			mockConfig := createMockAWSConfig()

			// test generateAwsAuthHeaders with mock config
			headers, err := c.generateAwsAuthHeaders(context.Background(), *mockConfig)
			if err != nil && !tt.wantErr {
				t.Errorf("generateAwsAuthHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// test AWS auth header generation
			if !tt.wantErr && headers != nil {
				// verify that all required Vault AWS auth headers are present
				if _, ok := headers["role"]; !ok {
					t.Error("Expected 'role' in auth headers for Vela AWS authentication")
				}
				if _, ok := headers["iam_http_request_method"]; !ok {
					t.Error("Expected 'iam_http_request_method' in auth headers for Vela AWS authentication")
				}
				if _, ok := headers["iam_request_url"]; !ok {
					t.Error("Expected 'iam_request_url' in auth headers for Vela AWS authentication")
				}
				if _, ok := headers["iam_request_headers"]; !ok {
					t.Error("Expected 'iam_request_headers' in auth headers for Vela AWS authentication")
				}
				if _, ok := headers["iam_request_body"]; !ok {
					t.Error("Expected 'iam_request_body' in auth headers for Vela AWS authentication")
				}
			}
		})
	}
}

type mockSTSClient struct {
	mockGetCallerIdentity func(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

func (m *mockSTSClient) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return m.mockGetCallerIdentity(ctx, input, optFns...)
}

// newSuccessfulMockSTSClient creates a mock STS client that returns successful responses.
func newSuccessfulMockSTSClient() STSClient {
	return &mockSTSClient{
		mockGetCallerIdentity: func(_ context.Context, _ *sts.GetCallerIdentityInput, _ ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
			return &sts.GetCallerIdentityOutput{}, nil
		},
	}
}

// newFailingMockSTSClient creates a mock STS client that returns errors.
func newFailingMockSTSClient(errMsg string) STSClient {
	return &mockSTSClient{
		mockGetCallerIdentity: func(_ context.Context, _ *sts.GetCallerIdentityInput, _ ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
			return nil, fmt.Errorf("%s", errMsg)
		},
	}
}

// createMockAWSConfig creates a mock AWS config to avoid real credential loading.
func createMockAWSConfig() *aws.Config {
	return &aws.Config{
		Region: "us-east-1",
		Credentials: aws.NewCredentialsCache(
			aws.CredentialsProviderFunc(func(_ context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     "test-access-key",
					SecretAccessKey: "test-secret-key",
				}, nil
			}),
		),
	}
}
