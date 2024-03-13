package middleware

import (
	"github.com/gin-gonic/gin"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/service"
	"oplian/utils"
	"strconv"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

// interceptor
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		waitUse, _ := utils.GetClaims(c)
		obj := c.Request.URL.Path
		act := c.Request.Method
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		e := casbinService.Casbin()
		success, _ := e.Enforce(sub, obj, act)
		if global.ZC_CONFIG.System.Env == "develop" || success {
			c.Next()
		} else {
			response.FailWithDetailed(gin.H{}, "Not enough permissions", c)
			c.Abort()
			return
		}
	}
}
