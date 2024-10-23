// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore duplicate with user code
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

const (
	// HookResp represents a JSON return for a single hook.
	HookResp = `{
  "id": 1,
  "repo": {
	"id": 1,
    "owner": {
	  	"id": 1,
	  	"name": "octocat",
	  	"favorites": [],
		"active": true,
    	"admin": false
  	},
  	"org": "github",
	"counter": 10,
	"name": "octocat",
	"full_name": "github/octocat",
	"link": "https://github.com/github/octocat",
	"clone": "https://github.com/github/octocat",
	"branch": "main",
	"build_limit": 10,
	"timeout": 60,
	"visibility": "public",
	"private": false,
	"trusted": true,
	"pipeline_type": "yaml",
	"topics": [],
	"active": true,
	"allow_events": {
		"push": {
			"branch": true,
			"tag": true
		},
		"pull_request": {
			"opened": true,
			"synchronize": true,
			"reopened": true,
			"edited": false
		},
		"deployment": {
			"created": true
		},
		"comment": {
			"created": false,
			"edited": false
		}
  	},
  "approve_build": "fork-always",
  "previous_name": ""
 },
  "build": {
  "id": 1,
  "repo": {
	"id": 1
 },
  "pipeline_id": 1,
  "number": 1,
  "parent": 1,
  "event": "push",
  "event_action": "",
  "status": "created",
  "error": "",
  "enqueued": 1563474077,
  "created": 1563474076,
  "started": 1563474077,
  "finished": 0,
  "deploy": "",
  "deploy_number": 1,
  "deploy_payload": {},
  "clone": "https://github.com/github/octocat.git",
  "source": "https://github.com/github/octocat/commit/48afb5bdc41ad69bf22588491333f7cf71135163",
  "title": "push received from https://github.com/github/octocat",
  "message": "First commit...",
  "commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
  "sender": "OctoKitty",
  "sender_scm_id": "0",
  "author": "OctoKitty",
  "email": "octokitty@github.com",
  "link": "https://vela.example.company.com/github/octocat/1",
  "branch": "main",
  "ref": "refs/heads/main",
  "head_ref": "",
  "base_ref": "",
  "host": "example.company.com",
  "runtime": "docker",
  "distribution": "linux",
  "approved_at": 0,
  "approved_by": ""
},
  "number": 1,
  "source_id": "c8da1302-07d6-11ea-882f-4893bca275b8",
  "created": 1563475419,
  "host": "github.com",
  "event": "push",
  "event_action": "",
  "webhook_id": 1234,
  "branch": "main",
  "error": "",
  "status": "success",
  "link": "https://github.com/github/octocat/settings/hooks/1"
}`

	// HooksResp represents a JSON return for one to many hooks.
	HooksResp = `[
{
  "id": 1,
  "repo": {
	"id": 1,
    "owner": {
	  	"id": 1,
	  	"name": "octocat",
	  	"favorites": [],
		"active": true,
    	"admin": false
  	},
  	"org": "github",
	"counter": 10,
	"name": "octocat",
	"full_name": "github/octocat",
	"link": "https://github.com/github/octocat",
	"clone": "https://github.com/github/octocat",
	"branch": "main",
	"build_limit": 10,
	"timeout": 60,
	"visibility": "public",
	"private": false,
	"trusted": true,
	"pipeline_type": "yaml",
	"topics": [],
	"active": true,
	"allow_events": {
		"push": {
			"branch": true,
			"tag": true
		},
		"pull_request": {
			"opened": true,
			"synchronize": true,
			"reopened": true,
			"edited": false
		},
		"deployment": {
			"created": true
		},
		"comment": {
			"created": false,
			"edited": false
		}
  	},
  "approve_build": "fork-always",
  "previous_name": ""
 },
  "build": {
  "id": 1,
  "repo": {
	"id": 1
 },
  "pipeline_id": 1,
  "number": 1,
  "parent": 1,
  "event": "push",
  "event_action": "",
  "status": "created",
  "error": "",
  "enqueued": 1563474077,
  "created": 1563474076,
  "started": 1563474077,
  "finished": 0,
  "deploy": "",
  "deploy_number": 1,
  "deploy_payload": {},
  "clone": "https://github.com/github/octocat.git",
  "source": "https://github.com/github/octocat/commit/48afb5bdc41ad69bf22588491333f7cf71135163",
  "title": "push received from https://github.com/github/octocat",
  "message": "First commit...",
  "commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
  "sender": "OctoKitty",
  "sender_scm_id": "0",
  "author": "OctoKitty",
  "email": "octokitty@github.com",
  "link": "https://vela.example.company.com/github/octocat/1",
  "branch": "main",
  "ref": "refs/heads/main",
  "head_ref": "",
  "base_ref": "",
  "host": "example.company.com",
  "runtime": "docker",
  "distribution": "linux",
  "approved_at": 0,
  "approved_by": ""
},
  "number": 1,
  "source_id": "c8da1302-07d6-11ea-882f-4893bca275b8",
  "created": 1563475419,
  "host": "github.com",
  "event": "push",
  "event_action": "",
  "webhook_id": 1234,
  "branch": "main",
  "error": "",
  "status": "success",
  "link": "https://github.com/github/octocat/settings/hooks/1"
},
  {
    "id": 1,
    "repo": {
	"id": 1,
    "owner": {
	  	"id": 1,
	  	"name": "octocat",
	  	"favorites": [],
		"active": true,
    	"admin": false
  	},
  	"org": "github",
	"counter": 10,
	"name": "octocat",
	"full_name": "github/octocat",
	"link": "https://github.com/github/octocat",
	"clone": "https://github.com/github/octocat",
	"branch": "main",
	"build_limit": 10,
	"timeout": 60,
	"visibility": "public",
	"private": false,
	"trusted": true,
	"pipeline_type": "yaml",
	"topics": [],
	"active": true,
	"allow_events": {
		"push": {
			"branch": true,
			"tag": true
		},
		"pull_request": {
			"opened": true,
			"synchronize": true,
			"reopened": true,
			"edited": false
		},
		"deployment": {
			"created": true
		},
		"comment": {
			"created": false,
			"edited": false
		}
  	},
  "approve_build": "fork-always",
  "previous_name": ""
 },
    "number": 1,
    "source_id": "c8da1302-07d6-11ea-882f-4893bca275b8",
    "created": 1563475419,
    "host": "github.com",
    "event": "push",
    "branch": "main",
    "error": "",
    "status": "success",
    "link": "https://github.com/github/octocat/settings/hooks/1"
  }
]`
)

// getHooks returns mock JSON for a http GET.
func getHooks(c *gin.Context) {
	data := []byte(HooksResp)

	var body []api.Hook
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getHook has a param :hook returns mock JSON for a http GET.
//
// Pass "0" to :hook to test receiving a http 404 response.
func getHook(c *gin.Context) {
	s := c.Param("hook")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Hook %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	data := []byte(HookResp)

	var body api.Hook
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addHook returns mock JSON for a http POST.
func addHook(c *gin.Context) {
	data := []byte(HookResp)

	var body api.Hook
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateHook has a param :hook returns mock JSON for a http PUT.
//
// Pass "0" to :hook to test receiving a http 404 response.
func updateHook(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		s := c.Param("hook")

		if strings.EqualFold(s, "0") {
			msg := fmt.Sprintf("Hook %s does not exist", s)

			c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

			return
		}
	}

	data := []byte(HookResp)

	var body api.Hook
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeHook has a param :hook returns mock JSON for a http DELETE.
//
// Pass "0" to :hook to test receiving a http 404 response.
func removeHook(c *gin.Context) {
	s := c.Param("hook")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Hook %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Hook %s removed", s))
}
