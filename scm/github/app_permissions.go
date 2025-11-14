// SPDX-License-Identifier: Apache-2.0

package github

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v79/github"
)

// see: https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28
const (
	// GitHub App install permission 'none'.
	AppInstallPermissionNone = "none"
	// GitHub App install permission 'read'.
	AppInstallPermissionRead = "read"
	// GitHub App install permission 'write'.
	AppInstallPermissionWrite = "write"
)

const (
	// GitHub App install contents resource.
	AppInstallResourceContents = "contents"
	// GitHub App install checks resource.
	AppInstallResourceChecks = "checks"
	// GitHub App install packages resource.
	AppInstallResourcePackages = "packages"
	// add more supported resources as needed.
)

// GetInstallationPermission takes permissions and returns the permission level if valid.
func GetInstallationPermission(resource string, appPermissions *github.InstallationPermissions) (string, error) {
	switch resource {
	case AppInstallResourceContents:
		return appPermissions.GetContents(), nil
	case AppInstallResourceChecks:
		return appPermissions.GetChecks(), nil
	case AppInstallResourcePackages:
		return appPermissions.GetPackages(), nil
	// add more supported resources as needed.
	default:
		return "", fmt.Errorf("given permission resource not supported: %s", resource)
	}
}

// ApplyInstallationPermissions takes permissions and applies a new permission if valid.
func ApplyInstallationPermissions(resource, perm string, perms *github.InstallationPermissions) (*github.InstallationPermissions, error) {
	// convert permissions from string
	switch strings.ToLower(perm) {
	case AppInstallPermissionNone:
	case AppInstallPermissionRead:
	case AppInstallPermissionWrite:
		break
	default:
		return perms, fmt.Errorf("invalid permission level given for <resource>:<level> in %s:%s", resource, perm)
	}

	// convert resource from string
	switch strings.ToLower(resource) {
	case AppInstallResourceContents:
		perms.Contents = github.Ptr(perm)
	case AppInstallResourceChecks:
		perms.Checks = github.Ptr(perm)
	case AppInstallResourcePackages:
		perms.Packages = github.Ptr(perm)
	// add more supported resources as needed.
	default:
		return perms, fmt.Errorf("invalid permission resource given for <resource>:<level> in %s:%s", resource, perm)
	}

	return perms, nil
}

// InstallationHasPermission takes a resource:perm pair and checks if the actual permission matches the expected permission or is supersceded by a higher permission.
func InstallationHasPermission(resource, requiredPerm, actualPerm string) error {
	if len(actualPerm) == 0 {
		return fmt.Errorf("github app missing permission %s:%s", resource, requiredPerm)
	}

	permitted := false

	switch requiredPerm {
	case AppInstallPermissionNone:
		permitted = true
	case AppInstallPermissionRead:
		if actualPerm == AppInstallPermissionRead ||
			actualPerm == AppInstallPermissionWrite {
			permitted = true
		}
	case AppInstallPermissionWrite:
		if actualPerm == AppInstallPermissionWrite {
			permitted = true
		}
	default:
		return fmt.Errorf("invalid required permission type: %s", requiredPerm)
	}

	if !permitted {
		return fmt.Errorf("github app requires permission %s:%s, found: %s", AppInstallResourceContents, AppInstallPermissionRead, actualPerm)
	}

	return nil
}
