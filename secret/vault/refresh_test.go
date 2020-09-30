package vault

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func Test_client_getAwsToken(t *testing.T) {
	tests := []struct {
		name         string
		responseFile string
		responseCode int
		stsClient    stsiface.STSAPI
		vaultRole    string
		want         string
		want1        time.Duration
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
			want:      "s.5RGnjF5aUbhz2XWWM9nHxO57",
			want1:     1 * time.Minute,
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

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := ioutil.ReadFile(fmt.Sprintf("testdata/refresh/%s", tt.responseFile))
				if err != nil {
					t.Error(err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseCode)
				_, _ = w.Write(data)
			}))
			defer ts.Close()
			c, err := New(ts.URL, "", "2", "", "", tt.vaultRole, 5*time.Minute)
			if err != nil {
				t.Error(err)
			}
			c.Aws.StsClient = tt.stsClient

			got, got1, err := c.getAwsToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("getAwsToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("getAwsToken() got = %v, want %v", got, tt.want)
			}

			if got1 != tt.want1 {
				t.Errorf("getAwsToken() got1 = %v, want %v", got1, tt.want1)
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
