package plugin

import (
	"github.com/gin-gonic/gin"
)

const (
	OnlyFuncName = "Plugin"
)

type Plugin interface {
	Register(group *gin.RouterGroup)

	RouterPath() string
}
