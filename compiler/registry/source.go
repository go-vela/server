// SPDX-License-Identifier: Apache-2.0

package registry

// Source represents a registry object
// for retrieving templates.
type Source struct {
	Host string
	Org  string
	Repo string
	Name string
	Ref  string
}
