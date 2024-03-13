package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"oplian/global"
	"oplian/middleware"
	"oplian/plugin/email"
	"oplian/utils/plugin"
)

func PluginInit(group *gin.RouterGroup, Plugin ...plugin.Plugin) {
	for i := range Plugin {
		PluginGroup := group.Group(Plugin[i].RouterPath())
		Plugin[i].Register(PluginGroup)
	}
}

func InstallPlugin(Router *gin.Engine) {
	PublicGroup := Router.Group("")
	fmt.Println("No authentication plug-in is installed==》", PublicGroup)
	PrivateGroup := Router.Group("")
	fmt.Println("Authentication plug-in installation==》", PrivateGroup)
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	//Add plugins with role-linked permissions Example Local example mode in online warehouse mode Note that import above can be switched on its own with the same effect
	PluginInit(PrivateGroup, email.CreateEmailPlug(
		global.ZC_CONFIG.Email.To,
		global.ZC_CONFIG.Email.From,
		global.ZC_CONFIG.Email.Host,
		global.ZC_CONFIG.Email.Secret,
		global.ZC_CONFIG.Email.Nickname,
		global.ZC_CONFIG.Email.Port,
		global.ZC_CONFIG.Email.IsSSL,
	))
}
