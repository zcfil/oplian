package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
	"oplian/middleware"
)

type AuthorityRouter struct{}

func (s *AuthorityRouter) InitAuthorityRouter(Router *gin.RouterGroup) {
	authorityRouter := Router.Group("authority").Use(middleware.OperationRecord())
	authorityRouterWithoutRecord := Router.Group("authority")
	authorityApi := v1.ApiGroupApp.SystemApiGroup.AuthorityApi
	{
		authorityRouter.POST("createAuthority", authorityApi.CreateAuthority)   // Create a role
		authorityRouter.POST("deleteAuthority", authorityApi.DeleteAuthority)   // Deleting a role
		authorityRouter.PUT("updateAuthority", authorityApi.UpdateAuthority)    // Update role
		authorityRouter.POST("copyAuthority", authorityApi.CopyAuthority)       // Copy role
		authorityRouter.POST("setDataAuthority", authorityApi.SetDataAuthority) // Set role resource permissions
	}
	{
		authorityRouterWithoutRecord.POST("getAuthorityList", authorityApi.GetAuthorityList) // Get a list of roles
	}
}
