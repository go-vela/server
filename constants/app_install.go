// SPDX-License-Identifier: Apache-2.0

// App Install vars.
package constants

// see: https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28
const (
	// The string value for GitHub App install read permissions.
	AppInstallPermissionRead = "read"
	// The string value for GitHub App install write permissions.
	AppInstallPermissionWrite = "write"
)

const (
	// The string value for GitHub App install contents resource.
	AppInstallResourceContents = "contents"
	// The string value for GitHub App install checks resource.
	AppInstallResourceChecks = "checks"
)
