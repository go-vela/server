// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"testing"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

func Test_validateAllowlist(t *testing.T) {
	// setup types
	repoSecret := new(types.Secret)
	repoSecret.SetType(constants.SecretRepo)

	repoSecretWithAllowlist := new(types.Secret)
	repoSecretWithAllowlist.SetType(constants.SecretRepo)
	repoSecretWithAllowlist.SetRepoAllowlist([]string{"alpha/beta"})

	orgSecret := new(types.Secret)
	orgSecret.SetType(constants.SecretOrg)
	orgSecret.SetOrg("alpha")

	orgSecretWithAllowlist := new(types.Secret)
	orgSecretWithAllowlist.SetType(constants.SecretOrg)
	orgSecretWithAllowlist.SetOrg("alpha")
	orgSecretWithAllowlist.SetRepoAllowlist([]string{"alpha/beta", "alpha/gamma"})

	orgSecretWithBadFormat := new(types.Secret)
	orgSecretWithBadFormat.SetType(constants.SecretOrg)
	orgSecretWithBadFormat.SetOrg("alpha")
	orgSecretWithBadFormat.SetRepoAllowlist([]string{"alpha.beta", "alpha/gamma"})

	orgSecretWithInvalidRepo := new(types.Secret)
	orgSecretWithInvalidRepo.SetType(constants.SecretOrg)
	orgSecretWithInvalidRepo.SetOrg("alpha")
	orgSecretWithInvalidRepo.SetRepoAllowlist([]string{"alpha/beta", "gamma/delta"})

	sharedSecret := new(types.Secret)
	sharedSecret.SetType(constants.SecretShared)

	sharedSecretWithAllowlist := new(types.Secret)
	sharedSecretWithAllowlist.SetType(constants.SecretShared)
	sharedSecretWithAllowlist.SetRepoAllowlist([]string{"alpha/beta", "gamma/delta"})

	tests := []struct {
		name    string
		secret  *types.Secret
		wantErr bool
	}{
		{
			name:   "repo no allowlist",
			secret: repoSecret,
		},
		{
			name:    "repo with allowlist",
			secret:  repoSecretWithAllowlist,
			wantErr: true,
		},
		{
			name:   "org secret",
			secret: orgSecret,
		},
		{
			name:   "org secret with allowlist",
			secret: orgSecretWithAllowlist,
		},
		{
			name:    "org secret bad format",
			secret:  orgSecretWithBadFormat,
			wantErr: true,
		},
		{
			name:    "org secret invalid repo",
			secret:  orgSecretWithInvalidRepo,
			wantErr: true,
		},
		{
			name:   "shared secret",
			secret: sharedSecret,
		},
		{
			name:   "shared secret with allowlist",
			secret: sharedSecretWithAllowlist,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateAllowlist(test.secret)
			if test.wantErr != (err != nil) {
				t.Errorf("want %t, got %v", test.wantErr, err)
			}
		})
	}
}
