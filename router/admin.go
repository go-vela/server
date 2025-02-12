// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/admin"
	"github.com/go-vela/server/router/middleware/perm"
)

// AdminHandlers is a function that extends the provided base router group
// with the API handlers for admin functionality.
//
// GET    	 /api/v1/admin/builds/queue
// PUT    	 /api/v1/admin/build
// PUT    	 /api/v1/admin/clean
// PUT    	 /api/v1/admin/deployment
// PUT    	 /api/v1/admin/hook
// PUT    	 /api/v1/admin/repo
// PUT    	 /api/v1/admin/secret
// PUT    	 /api/v1/admin/service
// PUT    	 /api/v1/admin/step
// PUT    	 /api/v1/admin/user
// POST   	 /api/v1/admin/workers/:worker/register
// GET    	 /api/v1/admin/settings
// PUT    	 /api/v1/admin/settings
// DELETE	 /api/v1/admin/settings.
func AdminHandlers(base *gin.RouterGroup) {
	// Admin endpoints
	_admin := base.Group("/admin", perm.MustPlatformAdmin())
	{
		// Admin build queue endpoint
		_admin.GET("/builds/queue", admin.AllBuildsQueue)

		// Admin build endpoint
		_admin.PUT("/build", admin.UpdateBuild)

		// Admin clean endpoint
		_admin.PUT("/clean", admin.CleanResources)

		// Admin deployment endpoint
		_admin.PUT("/deployment", admin.UpdateDeployment)

		// Admin hook endpoint
		_admin.PUT("/hook", admin.UpdateHook)

		// Admin repo endpoint
		_admin.PUT("/repo", admin.UpdateRepo)

		// Admin rotate keys endpoint
		_admin.POST("/rotate_oidc_keys", admin.RotateOIDCKeys)

		// Admin secret endpoint
		_admin.PUT("/secret", admin.UpdateSecret)

		// Admin service endpoint
		_admin.PUT("/service", admin.UpdateService)

		// Admin step endpoint
		_admin.PUT("/step", admin.UpdateStep)

		// Admin storage bucket endpoints
		//_admin.GET("/storage/bucket", admin.ListBuckets)
		_admin.PUT("/storage/bucket", admin.CreateBucket)
		_admin.DELETE("/storage/bucket", admin.DeleteBucket)

		// Admin storage bucket lifecycle endpoint
		_admin.GET("/storage/bucket/lifecycle", admin.GetBucketLifecycle)
		_admin.PUT("/storage/bucket/lifecycle", admin.SetBucketLifecycle)

		// Admin storage object endpoints
		_admin.POST("/storage/object/download", admin.DownloadObject)
		//_admin.POST("/storage/object", admin.UploadObject)

		// Admin storage presign endpoints
		_admin.GET("/storage/presign", admin.GetPresignedURL)

		// Admin user endpoint
		_admin.PUT("/user", admin.UpdateUser)

		// Admin worker endpoint
		_admin.POST("/workers/:worker/register", admin.RegisterToken)

		// Admin settings endpoints
		_admin.GET("/settings", admin.GetSettings)
		_admin.PUT("/settings", admin.UpdateSettings)
		_admin.DELETE("/settings", admin.RestoreSettings)
	} // end of admin endpoints
}
