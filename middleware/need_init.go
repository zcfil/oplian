package middleware

import (
	"github.com/gin-gonic/gin"
	"oplian/global"
	"oplian/model/common/response"
)

// Handles cross-domain requests and supports options access
func NeedInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.ZC_DB == nil {
			response.OkWithDetailed(gin.H{
				"needInit": true,
			}, "Go to initialize the database", c)
			c.Abort()
		} else {
			c.Next()
		}
		// 处理请求
	}
}
