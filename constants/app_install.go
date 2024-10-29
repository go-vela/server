// SPDX-License-Identifier: Apache-2.0

// App Install vars.
package constants

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
)

const (
	// GitHub App install repositories selection when "all" repositories are selected.
	AppInstallRepositoriesSelectionAll = "all"
	// GitHub App install repositories selection when a subset of repositories are selected.
	AppInstallRepositoriesSelectionSelected = "selected"
)

const (
	// GitHub App install setup_action type 'install'.
	AppInstallSetupActionInstall = "install"
	// GitHub App install event type 'created'.
	AppInstallCreated = "created"
	// GitHub App install event type 'deleted'.
	AppInstallDeleted = "deleted"
)
