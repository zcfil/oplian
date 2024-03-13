package middleware

import (
	"oplian/define"
	"oplian/utils"

	"github.com/gin-gonic/gin"
	"oplian/model/common/response"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Here, jwt authentication header information x-token returns the token information upon login. Here, the front-end needs to store the token in the cookie or local localStorage, but the expiration time needs to be negotiated with the back-end to agree to refresh the token or log in again
		token := c.Request.Header.Get(define.TOKEN_NAME)
		if token == "" {
			response.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", c)
			c.Abort()
			return
		}
		j := utils.NewJWT()

		claims, err := j.ParseToken(token)
		if err != nil {
			if err == utils.TokenExpired {
				response.FailWithDetailed(gin.H{"reload": true}, "授权已过期", c)
				c.Abort()
				return
			}
			response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
