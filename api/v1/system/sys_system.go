package system

import (
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system"
	systemRes "oplian/model/system/response"
	"oplian/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemApi struct{}

// GetSystemConfig
// @Tags      System
// @Summary   Get configuration file content
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysConfigResponse,msg=string}
// @Router    /system/getSystemConfig [post]
func (s *SystemApi) GetSystemConfig(c *gin.Context) {
	config, err := systemConfigService.GetSystemConfig()
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(systemRes.SysConfigResponse{Config: config}, "Successfully obtained", c)
}

// SetSystemConfig
// @Tags      System
// @Summary   Set configuration file content
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=string}
// @Router    /system/setSystemConfig [post]
func (s *SystemApi) SetSystemConfig(c *gin.Context) {
	var sys system.System
	err := c.ShouldBindJSON(&sys)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemConfigService.SetSystemConfig(sys)
	if err != nil {
		global.ZC_LOG.Error("Setting failed!", zap.Error(err))
		response.FailWithMessage("Setting failed", c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}

// ReloadSystem
// @Tags      System
// @Summary   Restart the system
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}
// @Router    /system/reloadSystem [post]
func (s *SystemApi) ReloadSystem(c *gin.Context) {
	err := utils.Reload()
	if err != nil {
		global.ZC_LOG.Error("System restart failed!", zap.Error(err))
		response.FailWithMessage("System restart failed", c)
		return
	}
	response.OkWithMessage("System restart successful", c)
}

// GetServerInfo
// @Tags      System
// @Summary   Get server information
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /system/getServerInfo [post]
func (s *SystemApi) GetServerInfo(c *gin.Context) {
	server, err := systemConfigService.GetServerInfo()
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(gin.H{"server": server}, "Successfully obtained", c)
}
