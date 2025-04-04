// SPDX-License-Identifier: Apache-2.0

package util

import (
	"context"
	"html"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"

	api "github.com/go-vela/server/api/types"
)

// HandleError appends the error to the handler chain for logging and outputs it.
func HandleError(c context.Context, status int, err error) {
	msg := err.Error()

	switch ctx := c.(type) {
	case *gin.Context:
		//nolint:errcheck // ignore checking error
		ctx.Error(err)
		ctx.AbortWithStatusJSON(status, api.Error{Message: &msg})

		return
	default:
		return
	}
}

// FormParameter safely captures a form parameter from the context
// by removing any new lines and HTML escaping the value.
func FormParameter(c *gin.Context, parameter string) string {
	return EscapeValue(c.Request.FormValue(parameter))
}

// QueryParameter safely captures a query parameter from the context
// by removing any new lines and HTML escaping the value.
func QueryParameter(c *gin.Context, parameter, value string) string {
	return EscapeValue(c.DefaultQuery(parameter, value))
}

// PathParameter safely captures a path parameter from the context
// by removing any new lines and HTML escaping the value.
func PathParameter(c *gin.Context, parameter string) string {
	return EscapeValue(c.Param(parameter))
}

// SplitFullName safely splits the repo.FullName field into an org and name.
func SplitFullName(value string) (string, string) {
	// split repo full name into org and repo
	repoSlice := strings.Split(value, "/")
	if len(repoSlice) != 2 {
		return "", ""
	}

	org := repoSlice[0]
	repo := repoSlice[1]

	return org, repo
}

// EscapeValue safely escapes any string by removing any new lines and HTML escaping it.
func EscapeValue(value string) string {
	// replace all new lines in the value
	escaped := strings.Replace(strings.Replace(value, "\n", "", -1), "\r", "", -1)

	// HTML escape the new line escaped value
	return html.EscapeString(escaped)
}

// Unique is a helper function that takes a slice and
// validates that there are no duplicate entries.
func Unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true

			list = append(list, entry)
		}
	}

	return list
}

// CheckAllowlist is a helper function to ensure only repos in the
// allowlist are specified.
//
// a single entry of '*' allows any repo to be enabled.
func CheckAllowlist(r *api.Repo, allowlist []string) bool {
	// check if all repos are allowed to be enabled
	if len(allowlist) == 1 && allowlist[0] == "*" {
		return true
	}

	for _, repo := range allowlist {
		// allow all repos in org
		if strings.Contains(repo, "/*") {
			if strings.HasPrefix(repo, r.GetOrg()) {
				return true
			}
		}

		// allow specific repo within org
		if repo == r.GetFullName() {
			return true
		}
	}

	return false
}

// Sanitize is a helper function to verify the provided input
// field does not contain HTML content. If the input field
// does contain HTML, then the function will sanitize and
// potentially remove the HTML if deemed malicious.
func Sanitize(field string) string {
	// create new HTML input microcosm-cc/bluemonday policy
	p := bluemonday.StrictPolicy()

	// create a URL query unescaped string from the field
	queryUnescaped, err := url.QueryUnescape(field)
	if err != nil {
		// overwrite URL query unescaped string with field
		queryUnescaped = field
	}

	// create an HTML escaped string from the field
	htmlEscaped := html.EscapeString(queryUnescaped)

	// create a microcosm-cc/bluemonday escaped string from the field
	bluemondayEscaped := p.Sanitize(queryUnescaped)

	// check if the field contains html
	if !strings.EqualFold(htmlEscaped, bluemondayEscaped) {
		// create new HTML input microcosm-cc/bluemonday policy
		return bluemondayEscaped
	}

	// return the unmodified field
	return field
}

// Ptr is a helper routine that allocates a new T value
// to store v and returns a pointer to it.
func Ptr[T any](v T) *T {
	return &v
}
