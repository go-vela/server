// SPDX-License-Identifier: Apache-2.0

package constants

// TextFileExtensions and MediaFileExtensions package-level vars (Go doesn't allow const slices).
var TextFileExtensions = []string{".xml", ".json", ".html", ".txt"}
var MediaFileExtensions = []string{".png", ".jpg", ".jpeg", ".gif", ".mp4", ".mov"}

// AllAllowedExtensions is the union of test + media.
var AllAllowedExtensions = append(TextFileExtensions, MediaFileExtensions...)
