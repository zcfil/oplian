package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
	"oplian/middleware"
)

type UserRouter struct{}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user").Use(middleware.OperationRecord())
	userRouterWithoutRecord := Router.Group("user")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.POST("admin_register", baseApi.Register)               // Administrator Register Account
		userRouter.POST("changePassword", baseApi.ChangePassword)         // User change password
		userRouter.POST("setUserAuthority", baseApi.SetUserAuthority)     // Set user rights
		userRouter.DELETE("deleteUser", baseApi.DeleteUser)               // Delete a user
		userRouter.PUT("setUserInfo", baseApi.SetUserInfo)                // Set user information
		userRouter.PUT("setSelfInfo", baseApi.SetSelfInfo)                // Set self information
		userRouter.POST("setUserAuthorities", baseApi.SetUserAuthorities) // Set the user permission group
		userRouter.POST("resetPassword", baseApi.ResetPassword)           // Reset password
	}
	{
		userRouterWithoutRecord.POST("getUserList", baseApi.GetUserList)        // Page for a list of users
		userRouterWithoutRecord.GET("getUserInfo", baseApi.GetUserInfo)         // Access to self information
		userRouterWithoutRecord.GET("getUserPullList", baseApi.GetUserPullList) // Get user list drop-down
	}
}
