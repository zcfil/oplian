package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/config"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system"
	"time"
)

type PatrolConfigApi struct{}

// GetSysPatrolConfigList
// @Tags      SysPatrolConfig
// @Summary   Get the list of inspection settings
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysPatrolConfig/getSysPatrolConfigList [get]
func (s *PatrolConfigApi) GetSysPatrolConfigList(c *gin.Context) {
	list, err := patrolConfigService.GetSysPatrolConfigInfoList()
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	if len(list) == 0 {
		patrolTypes := []int64{config.HostMinerType, config.HostWorkerType, config.HostStorageType}
		for _, val := range patrolTypes {
			patrolConfig := system.SysPatrolConfig{
				ZC_MODEL:   global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
				PatrolType: val,
			}
			err = patrolConfigService.CreateSysPatrolConfig(patrolConfig)
			if err != nil {
				global.ZC_LOG.Error("Failed to initialize inspection settings", zap.Error(err))
				response.FailWithMessage("Failed to initialize inspection settings", c)
				return
			}
		}
	}
	list, err = patrolConfigService.GetSysPatrolConfigInfoList()
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain list!", zap.Error(err))
		response.FailWithMessage("Failed to obtain list", c)
		return
	}
	response.OkWithDetailed(list, "Successfully obtained", c)
}

// UpdateSysPatrolConfig
// @Tags      SysPatrolConfig
// @Summary   Update host inspection settings
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysPatrolConfig/updateSysPatrolConfig [post]
func (s *PatrolConfigApi) UpdateSysPatrolConfig(c *gin.Context) {
	var reportReq system.SysPatrolConfig
	err := c.ShouldBindJSON(&reportReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Obtain patrol setting information
	reportReq.IntervalTime = reportReq.IntervalHours*3600 + reportReq.IntervalMinutes*60 + reportReq.IntervalSeconds
	err = patrolConfigService.UpdateSysPatrolConfig(&reportReq)
	if err != nil {
		global.ZC_LOG.Error("Modification failed!", zap.Error(err))
		response.FailWithMessage("Modification failed", c)
		return
	}

	response.OkWithMessage("Modified successfully", c)
}
