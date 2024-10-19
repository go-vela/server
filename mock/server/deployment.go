// SPDX-License-Identifier: Apache-2.0

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
	// DeploymentResp represents a JSON return for a single build.
	DeploymentResp = `{
  "id": 1,
  "number": 744479,
  "repo": {
    "id": 1,
    "owner": {
      "id": 1,
      "name": "Octocat",
      "active": true
    },
    "org": "Octocat",
    "name": "myvela",
    "full_name": "Octocat/myvela",
    "link": "https://github.com/Octocat/myvela",
    "clone": "https://github.com/Octocat/myvela.git",
    "branch": "main",
    "topics": [
      "example"
    ],
    "build_limit": 10,
    "timeout": 30,
    "counter": 2,
    "visibility": "public",
    "private": false,
    "trusted": false,
    "active": true,
    "allow_events": {
      "push": {
        "branch": true,
        "tag": false,
        "delete_branch": false,
        "delete_tag": false
      },
      "pull_request": {
        "opened": false,
        "edited": false,
        "synchronize": false,
        "reopened": false,
        "labeled": false,
        "unlabeled": false
      },
      "deployment": {
        "created": true
      },
      "comment": {
        "created": false,
        "edited": false
      },
      "schedule": {
        "run": false
      }
    },
    "pipeline_type": "yaml",
    "previous_name": "",
    "approve_build": "first-time"
  },
  "url": "https://github.com/api/v3/repos/Octocat/myvela/deployments/744479",
  "commit": "14c8b131c0c5e1489811bd64c755a9b6f74c792b",
  "ref": "main",
  "task": "deploy:vela",
  "target": "production",
  "description": "Deployment request from Vela",
  "payload": {},
  "created_at": 1727710527,
  "created_by": "Octocat",
  "builds": [
    {
      "id": 1,
      "repo": {
        "id": 1
      },
      "pipeline_id": 1,
      "number": 1,
      "parent": 0,
      "event": "deployment",
      "event_action": "created",
      "status": "success",
      "error": "",
      "enqueued": 1727710528,
      "created": 1727710528,
      "started": 1727710528,
      "finished": 1727710532,
      "deploy": "production",
      "deploy_number": 744479,
      "deploy_payload": {},
      "clone": "https://github.com/Octocat/myvela.git",
      "source": "https://github.com/api/v3/repos/Octocat/myvela/deployments/744479",
      "title": "deployment received from https://github.com/Octocat/myvela",
      "message": "Deployment request from Vela",
      "commit": "14c8b131c0c5e1489811bd64c755a9b6f74c792b",
      "sender": "Octocat",
      "sender_scm_id": "17043",
      "author": "Octocat",
      "email": "",
      "link": "http://localhost:8888/Octocat/myvela/1",
      "branch": "main",
      "ref": "refs/heads/main",
      "base_ref": "",
      "head_ref": "",
      "host": "worker",
      "runtime": "docker",
      "distribution": "linux",
      "approved_at": 0,
      "approved_by": ""
    },
    {
      "id": 2,
      "repo": {
        "id": 1
      },
      "pipeline_id": 1,
      "number": 2,
      "parent": 0,
      "event": "deployment",
      "event_action": "created",
      "status": "success",
      "error": "",
      "enqueued": 1727711899,
      "created": 1727711899,
      "started": 1727711899,
      "finished": 1727711904,
      "deploy": "production",
      "deploy_number": 744479,
      "deploy_payload": {},
      "clone": "https://github.com/Octocat/myvela.git",
      "source": "https://github.com/api/v3/repos/Octocat/myvela/deployments/744479",
      "title": "deployment received from https://github.com/Octocat/myvela",
      "message": "Deployment request from Vela",
      "commit": "14c8b131c0c5e1489811bd64c755a9b6f74c792b",
      "sender": "Octocat",
      "sender_scm_id": "17043",
      "author": "Octocat",
      "email": "",
      "link": "http://localhost:8888/Octocat/myvela/2",
      "branch": "main",
      "ref": "refs/heads/main",
      "base_ref": "",
      "head_ref": "",
      "host": "worker",
      "runtime": "docker",
      "distribution": "linux",
      "approved_at": 0,
      "approved_by": ""
    }
  ]
}`

	// DeploymentsResp represents a JSON return for one to many builds.
	DeploymentsResp = `[
{
  "id": 1,
  "number": 744479,
  "repo": {
    "id": 1,
    "owner": {
      "id": 1,
      "name": "Octocat",
      "active": true
    },
    "org": "Octocat",
    "name": "myvela",
    "full_name": "Octocat/myvela",
    "link": "https://github.com/Octocat/myvela",
    "clone": "https://github.com/Octocat/myvela.git",
    "branch": "main",
    "topics": [
      "example"
    ],
    "build_limit": 10,
    "timeout": 30,
    "counter": 2,
    "visibility": "public",
    "private": false,
    "trusted": false,
    "active": true,
    "allow_events": {
      "push": {
        "branch": true,
        "tag": false,
        "delete_branch": false,
        "delete_tag": false
      },
      "pull_request": {
        "opened": false,
        "edited": false,
        "synchronize": false,
        "reopened": false,
        "labeled": false,
        "unlabeled": false
      },
      "deployment": {
        "created": true
      },
      "comment": {
        "created": false,
        "edited": false
      },
      "schedule": {
        "run": false
      }
    },
    "pipeline_type": "yaml",
    "previous_name": "",
    "approve_build": "first-time"
  },
  "url": "https://github.com/api/v3/repos/Octocat/myvela/deployments/744479",
  "commit": "14c8b131c0c5e1489811bd64c755a9b6f74c792b",
  "ref": "main",
  "task": "deploy:vela",
  "target": "production",
  "description": "Deployment request from Vela",
  "payload": {},
  "created_at": 1727710527,
  "created_by": "Octocat",
  "builds": [
    {
      "id": 1,
      "repo": {
        "id": 1
      },
      "pipeline_id": 1,
      "number": 1,
      "parent": 0,
      "event": "deployment",
      "event_action": "created",
      "status": "success",
      "error": "",
      "enqueued": 1727710528,
      "created": 1727710528,
      "started": 1727710528,
      "finished": 1727710532,
      "deploy": "production",
      "deploy_number": 744479,
      "deploy_payload": {},
      "clone": "https://github.com/Octocat/myvela.git",
      "source": "https://github.com/api/v3/repos/Octocat/myvela/deployments/744479",
      "title": "deployment received from https://github.com/Octocat/myvela",
      "message": "Deployment request from Vela",
      "commit": "14c8b131c0c5e1489811bd64c755a9b6f74c792b",
      "sender": "Octocat",
      "sender_scm_id": "17043",
      "author": "Octocat",
      "email": "",
      "link": "http://localhost:8888/Octocat/myvela/1",
      "branch": "main",
      "ref": "refs/heads/main",
      "base_ref": "",
      "head_ref": "",
      "host": "worker",
      "runtime": "docker",
      "distribution": "linux",
      "approved_at": 0,
      "approved_by": ""
    },
    {
      "id": 2,
      "repo": {
        "id": 1
      },
      "pipeline_id": 1,
      "number": 2,
      "parent": 0,
      "event": "deployment",
      "event_action": "created",
      "status": "success",
      "error": "",
      "enqueued": 1727711899,
      "created": 1727711899,
      "started": 1727711899,
      "finished": 1727711904,
      "deploy": "production",
      "deploy_number": 744479,
      "deploy_payload": {},
      "clone": "https://github.com/Octocat/myvela.git",
      "source": "https://github.com/api/v3/repos/Octocat/myvela/deployments/744479",
      "title": "deployment received from https://github.com/Octocat/myvela",
      "message": "Deployment request from Vela",
      "commit": "14c8b131c0c5e1489811bd64c755a9b6f74c792b",
      "sender": "Octocat",
      "sender_scm_id": "17043",
      "author": "Octocat",
      "email": "",
      "link": "http://localhost:8888/Octocat/myvela/2",
      "branch": "main",
      "ref": "refs/heads/main",
      "base_ref": "",
      "head_ref": "",
      "host": "worker",
      "runtime": "docker",
      "distribution": "linux",
      "approved_at": 0,
      "approved_by": ""
    }
  ]
},
{
  "id": 1,
  "number": 744479,
  "repo": {
    "id": 1,
    "owner": {
      "id": 1,
      "name": "Octocat",
      "active": true
    },
    "org": "Octocat",
    "name": "myvela",
    "full_name": "Octocat/myvela",
    "link": "https://github.com/Octocat/myvela",
    "clone": "https://github.com/Octocat/myvela.git",
    "branch": "main",
    "topics": [
      "example"
    ],
    "build_limit": 10,
    "timeout": 30,
    "counter": 2,
    "visibility": "public",
    "private": false,
    "trusted": false,
    "active": true,
    "allow_events": {
      "push": {
        "branch": true,
        "tag": false,
        "delete_branch": false,
        "delete_tag": false
      },
      "pull_request": {
        "opened": false,
        "edited": false,
        "synchronize": false,
        "reopened": false,
        "labeled": false,
        "unlabeled": false
      },
      "deployment": {
        "created": true
      },
      "comment": {
        "created": false,
        "edited": false
      },
      "schedule": {
        "run": false
      }
    },
    "pipeline_type": "yaml",
    "previous_name": "",
    "approve_build": "first-time"
  },
  "url": "https://github.com/api/v3/repos/Octocat/myvela/deployments/744479",
  "commit": "14c8b131c0c5e1489811bd64c755a9b6f74c792b",
  "ref": "main",
  "task": "deploy:vela",
  "target": "production",
  "description": "Deployment request from Vela",
  "payload": {},
  "created_at": 1727710527,
  "created_by": "Octocat",
  "builds": [
    {
      "id": 2,
      "repo": {
        "id": 1
      },
      "pipeline_id": 1,
      "number": 2,
      "parent": 0,
      "event": "deployment",
      "event_action": "created",
      "status": "success",
      "error": "",
      "enqueued": 1727711899,
      "created": 1727711899,
      "started": 1727711899,
      "finished": 1727711904,
      "deploy": "production",
      "deploy_number": 744479,
      "deploy_payload": {},
      "clone": "https://github.com/Octocat/myvela.git",
      "source": "https://github.com/api/v3/repos/Octocat/myvela/deployments/744479",
      "title": "deployment received from https://github.com/Octocat/myvela",
      "message": "Deployment request from Vela",
      "commit": "14c8b131c0c5e1489811bd64c755a9b6f74c792b",
      "sender": "Octocat",
      "sender_scm_id": "17043",
      "author": "Octocat",
      "email": "",
      "link": "http://localhost:8888/Octocat/myvela/2",
      "branch": "main",
      "ref": "refs/heads/main",
      "base_ref": "",
      "head_ref": "",
      "host": "worker",
      "runtime": "docker",
      "distribution": "linux",
      "approved_at": 0,
      "approved_by": ""
    }
  ]
}
]`
)

// getDeployments returns mock JSON for a http GET.
func getDeployments(c *gin.Context) {
	data := []byte(DeploymentsResp)

	var body []api.Deployment
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getDeployment has a param :deployment returns mock JSON for a http GET.
func getDeployment(c *gin.Context) {
	d := c.Param("deployment")

	if strings.EqualFold(d, "0") {
		msg := fmt.Sprintf("Deployment %s does not exist", d)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	data := []byte(DeploymentResp)

	var body api.Deployment
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addDeployment returns mock JSON for a http POST.
func addDeployment(c *gin.Context) {
	data := []byte(DeploymentResp)

	var body api.Deployment
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateDeployment returns mock JSON for a http PUT.
func updateDeployment(c *gin.Context) {
	data := []byte(DeploymentResp)

	var body api.Deployment
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
