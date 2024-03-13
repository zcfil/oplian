package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
	"oplian/middleware"
)

type MenuRouter struct{}

func (s *MenuRouter) InitMenuRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	menuRouter := Router.Group("menu").Use(middleware.OperationRecord())
	menuRouterWithoutRecord := Router.Group("menu")
	authorityMenuApi := v1.ApiGroupApp.SystemApiGroup.AuthorityMenuApi
	{
		menuRouter.POST("addBaseMenu", authorityMenuApi.AddBaseMenu)           // New menu
		menuRouter.POST("addMenuAuthority", authorityMenuApi.AddMenuAuthority) // Added the relationship between menu and role
		menuRouter.POST("deleteBaseMenu", authorityMenuApi.DeleteBaseMenu)     // Delete menu
		menuRouter.POST("updateBaseMenu", authorityMenuApi.UpdateBaseMenu)     // Update menu
	}
	{
		menuRouterWithoutRecord.POST("getMenu", authorityMenuApi.GetMenu)                   // Get menu tree
		menuRouterWithoutRecord.POST("getMenuList", authorityMenuApi.GetMenuList)           // Page to get the base menu list
		menuRouterWithoutRecord.POST("getBaseMenuTree", authorityMenuApi.GetBaseMenuTree)   // Get the user dynamic route
		menuRouterWithoutRecord.POST("getMenuAuthority", authorityMenuApi.GetMenuAuthority) // Gets the specified role menu
		menuRouterWithoutRecord.POST("getBaseMenuById", authorityMenuApi.GetBaseMenuById)   // Get the menu by id
	}
	return menuRouter
}
