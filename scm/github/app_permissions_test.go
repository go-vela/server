// SPDX-License-Identifier: Apache-2.0

package github

import (
	"testing"

	"github.com/google/go-github/v72/github"
)

func TestGetInstallationPermission(t *testing.T) {
	tests := []struct {
		name          string
		resource      string
		permissions   *github.InstallationPermissions
		expectedPerm  string
		expectedError bool
	}{
		{
			name:         "valid contents permission",
			resource:     AppInstallResourceContents,
			permissions:  &github.InstallationPermissions{Contents: github.Ptr(AppInstallPermissionRead)},
			expectedPerm: AppInstallPermissionRead,
		},
		{
			name:         "valid checks permission",
			resource:     AppInstallResourceChecks,
			permissions:  &github.InstallationPermissions{Checks: github.Ptr(AppInstallPermissionWrite)},
			expectedPerm: AppInstallPermissionWrite,
		},
		{
			name:         "valid packages permission",
			resource:     AppInstallResourcePackages,
			permissions:  &github.InstallationPermissions{Packages: github.Ptr(AppInstallPermissionNone)},
			expectedPerm: AppInstallPermissionNone,
		},
		{
			name:          "invalid resource",
			resource:      "invalid_resource",
			permissions:   &github.InstallationPermissions{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perm, err := GetInstallationPermission(tt.resource, tt.permissions)
			if (err != nil) != tt.expectedError {
				t.Errorf("GetInstallationPermission() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if perm != tt.expectedPerm {
				t.Errorf("GetInstallationPermission() = %v, expected %v", perm, tt.expectedPerm)
			}
		})
	}
}

func TestApplyInstallationPermissions(t *testing.T) {
	tests := []struct {
		name          string
		resource      string
		perm          string
		initialPerms  *github.InstallationPermissions
		expectedPerms *github.InstallationPermissions
		expectedError bool
	}{
		{
			name:         "apply read permission to contents",
			resource:     AppInstallResourceContents,
			perm:         AppInstallPermissionRead,
			initialPerms: &github.InstallationPermissions{},
			expectedPerms: &github.InstallationPermissions{
				Contents: github.Ptr(AppInstallPermissionRead),
			},
		},
		{
			name:         "apply write permission to checks",
			resource:     AppInstallResourceChecks,
			perm:         AppInstallPermissionWrite,
			initialPerms: &github.InstallationPermissions{},
			expectedPerms: &github.InstallationPermissions{
				Checks: github.Ptr(AppInstallPermissionWrite),
			},
		},
		{
			name:         "apply none permission to packages",
			resource:     AppInstallResourcePackages,
			perm:         AppInstallPermissionNone,
			initialPerms: &github.InstallationPermissions{},
			expectedPerms: &github.InstallationPermissions{
				Packages: github.Ptr(AppInstallPermissionNone),
			},
		},
		{
			name:          "invalid permission level",
			resource:      AppInstallResourceContents,
			perm:          "invalid_perm",
			initialPerms:  &github.InstallationPermissions{},
			expectedError: true,
		},
		{
			name:          "invalid resource",
			resource:      "invalid_resource",
			perm:          AppInstallPermissionRead,
			initialPerms:  &github.InstallationPermissions{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms, err := ApplyInstallationPermissions(tt.resource, tt.perm, tt.initialPerms)
			if (err != nil) != tt.expectedError {
				t.Errorf("ApplyInstallationPermissions() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if !tt.expectedError && !comparePermissions(perms, tt.expectedPerms) {
				t.Errorf("ApplyInstallationPermissions() = %v, expected %v", perms, tt.expectedPerms)
			}
		})
	}
}

func TestInstallationHasPermission(t *testing.T) {
	tests := []struct {
		name          string
		resource      string
		requiredPerm  string
		actualPerm    string
		expectedError bool
	}{
		{
			name:         "valid read permission",
			resource:     AppInstallResourceContents,
			requiredPerm: AppInstallPermissionRead,
			actualPerm:   AppInstallPermissionRead,
		},
		{
			name:         "valid write permission",
			resource:     AppInstallResourceChecks,
			requiredPerm: AppInstallPermissionWrite,
			actualPerm:   AppInstallPermissionWrite,
		},
		{
			name:         "valid none permission",
			resource:     AppInstallResourcePackages,
			requiredPerm: AppInstallPermissionNone,
			actualPerm:   AppInstallPermissionNone,
		},
		{
			name:         "read permission superseded by write",
			resource:     AppInstallResourceContents,
			requiredPerm: AppInstallPermissionRead,
			actualPerm:   AppInstallPermissionWrite,
		},
		{
			name:          "missing permission",
			resource:      AppInstallResourceChecks,
			requiredPerm:  AppInstallPermissionWrite,
			actualPerm:    "",
			expectedError: true,
		},
		{
			name:          "invalid required permission",
			resource:      AppInstallResourcePackages,
			requiredPerm:  "invalid_perm",
			actualPerm:    AppInstallPermissionRead,
			expectedError: true,
		},
		{
			name:          "insufficient permission",
			resource:      AppInstallResourceContents,
			requiredPerm:  AppInstallPermissionWrite,
			actualPerm:    AppInstallPermissionRead,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InstallationHasPermission(tt.resource, tt.requiredPerm, tt.actualPerm)
			if (err != nil) != tt.expectedError {
				t.Errorf("InstallationHasPermission() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func comparePermissions(a, b *github.InstallationPermissions) bool {
	if a == nil || b == nil {
		return a == b
	}
	return github.Stringify(a) == github.Stringify(b)
}
