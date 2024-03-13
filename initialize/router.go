package initialize

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"oplian/global"
	"oplian/middleware"
	"oplian/router"
)

// 初始化总路由

func routers() *gin.Engine {
	Router := gin.Default()
	systemRouter := router.RouterGroupApp.System
	exampleRouter := router.RouterGroupApp.Example

	// VUE_APP_BASE_API = /
	// VUE_APP_BASE_PATH = http://localhost
	// Router.LoadHTMLGlob("./dist/*.html") //
	// Router.Static("/favicon.ico", "./dist/favicon.ico")
	// Router.Static("/static", "./dist/assets")
	// Router.StaticFile("/", "./dist/index.html")
	//
	Router.StaticFS(global.ZC_CONFIG.Local.Path, http.Dir(global.ZC_CONFIG.Local.StorePath))
	// Router.UsePercent(middleware.LoadTls())

	//Router.UsePercent(middleware.Cors())
	//Router.UsePercent(middleware.CorsByRules())
	//log.Println("use middleware cors")

	PublicGroup := Router.Group("")
	{
		// Health monitoring
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}
	{
		systemRouter.InitBaseRouter(PublicGroup) // Basic function Routes are not authenticated
		systemRouter.InitInitRouter(PublicGroup) // Automatic initialization related
	}
	//systemRouter.InitEquipMonitorRouter(PublicGroup) // Equipment monitoring
	PrivateGroup := Router.Group("")
	//PrivateGroup.UsePercent(middleware.JWTAuth()).UsePercent(middleware.CasbinHandler())
	//PrivateGroup.UsePercent(middleware.JWTAuth()).UsePercent()
	PrivateGroup.Use(middleware.JWTAuth())
	{
		systemRouter.InitApiRouter(PrivateGroup)
		systemRouter.InitUserRouter(PrivateGroup)
		systemRouter.InitMenuRouter(PrivateGroup)
		systemRouter.InitSystemRouter(PrivateGroup)
		systemRouter.InitCasbinRouter(PrivateGroup)
		systemRouter.InitAuthorityRouter(PrivateGroup)
		systemRouter.InitSysDictionaryRouter(PrivateGroup)
		systemRouter.InitSysOperationRecordRouter(PrivateGroup)
		systemRouter.InitSysDictionaryDetailRouter(PrivateGroup)
		systemRouter.InitAuthorityBtnRouterRouter(PrivateGroup)
		systemRouter.InitSysMachineRoomRecordRouter(PrivateGroup)
		systemRouter.InitSysHostGroupRouter(PrivateGroup)
		systemRouter.InitSysHostRecordRouter(PrivateGroup)
		systemRouter.InitSysHostMonitorRecordRouter(PrivateGroup)
		systemRouter.InitSysHostTestRecordRouter(PrivateGroup)
		systemRouter.InitSysHostPatrolRecordRouter(PrivateGroup)
		systemRouter.InitSysPatrolConfigRouter(PrivateGroup)
		systemRouter.InitHomeInterfaceRouter(PrivateGroup)
		systemRouter.InitEquipMonitorRouter(PrivateGroup)
		systemRouter.InitEquipWSRouter(PublicGroup)
		systemRouter.InitAppStoreRouter(PrivateGroup)

		systemRouter.InitJobPlatformRouter(PrivateGroup)
		systemRouter.InitWarnManageRouter(PrivateGroup)
		systemRouter.InitMonitorCenterRouter(PrivateGroup)

		exampleRouter.InitFileUploadAndDownloadRouter(PrivateGroup)
	}

	log.Println("router register success")
	return Router
}

func WebRouters() *gin.Engine {
	Router := routers()
	gatewayRouter := router.RouterGroupApp.Gateway
	connGroup := Router.Group("")
	//connGroup.Use(middleware.JWTAuth())
	{
		gatewayRouter.ConnRouter.InitConnRouter(connGroup)
	}
	lotusRouter := router.RouterGroupApp.Lotus
	lotusGroup := Router.Group("")
	lotusGroup.Use(middleware.JWTAuth())
	{
		lotusRouter.DeployRouter.InitDeployRouter(lotusGroup)
		lotusRouter.DispatchRouter.InitDispatchRouter(lotusGroup)
		lotusRouter.DispatchRouter.InitWorkerClusterRouter(lotusGroup)
		lotusRouter.DispatchRouter.InitSectorsRecoverRouter(lotusGroup)
	}

	return Router
}
