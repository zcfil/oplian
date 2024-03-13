package gateway

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type ConnRouter struct{}

// 网关连接路由
func (s *ConnRouter) InitConnRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	connRouter := Router.Group("conn")
	connApi := v1.ApiGroupApp.GateWayApiGroup.ConnApi
	{
		connRouter.POST("connectGateWay", connApi.ConnectGateWay) // Gateway connection
		connRouter.GET("ping", connApi.Ping)                      // Test connection
	}

	return connRouter
}
