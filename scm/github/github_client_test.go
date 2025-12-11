// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v80/github"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

func TestClient_installationCanReadRepo(t *testing.T) {
	// setup types
	accessibleRepo := new(api.Repo)
	accessibleRepo.SetOrg("octocat")
	accessibleRepo.SetName("Hello-World")
	accessibleRepo.SetFullName("octocat/Hello-World")
	accessibleRepo.SetInstallID(0)

	inaccessibleRepo := new(api.Repo)
	inaccessibleRepo.SetOrg("octocat")
	inaccessibleRepo.SetName("Hello-World")
	inaccessibleRepo.SetFullName("octocat/Hello-World2")
	inaccessibleRepo.SetInstallID(4)

	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/app/installations/:id/access_tokens", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/installations_access_tokens.json")
	})
	engine.GET("/api/v3/installation/repositories", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/installation_repositories.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	oauthClient, _ := NewTest(s.URL)

	appsClient, err := NewTest(s.URL)
	if err != nil {
		t.Errorf("unable to create GitHub App client: %v", err)
	}

	appsClient.AppsTransport = NewTestAppsTransport("")

	// setup tests
	tests := []struct {
		name          string
		client        *Client
		repo          *api.Repo
		installation  *github.Installation
		appsTransport bool
		want          bool
		wantErr       bool
	}{
		{
			name:   "installation can read repo",
			client: appsClient,
			repo:   accessibleRepo,
			installation: &github.Installation{
				ID: github.Ptr(int64(1)),
				Account: &github.User{
					Login: github.Ptr("github"),
				},
				RepositorySelection: github.Ptr(constants.AppInstallRepositoriesSelectionSelected),
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "installation cannot read repo",
			client: appsClient,
			repo:   inaccessibleRepo,
			installation: &github.Installation{
				ID: github.Ptr(int64(2)),
				Account: &github.User{
					Login: github.Ptr("github"),
				},
				RepositorySelection: github.Ptr(constants.AppInstallRepositoriesSelectionSelected),
			},
			want:    false,
			wantErr: false,
		},
		{
			name:   "no GitHub App client",
			client: oauthClient,
			repo:   accessibleRepo,
			installation: &github.Installation{
				ID: github.Ptr(int64(1)),
				Account: &github.User{
					Login: github.Ptr("github"),
				},
				RepositorySelection: github.Ptr(constants.AppInstallRepositoriesSelectionSelected),
			},
			want:    false,
			wantErr: true,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.installationCanReadRepo(context.Background(), tt.repo, tt.installation)
			if (err != nil) != tt.wantErr {
				t.Errorf("installationCanReadRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("installationCanReadRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
