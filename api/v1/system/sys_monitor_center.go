package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system/request"
	"oplian/utils"
)

type MonitorCenterApi struct{}

// BusinessMonitor
// @Tags      BusinessMonitor
// @Summary   Monitoring Center - Business Monitoring
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{res: res}
// @Router    /monitorCenter/businessMonitor [post]
func (m *MonitorCenterApi) BusinessMonitor(c *gin.Context) {

	var req request.MonitorCenterReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.MinerIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := monitorCenterServer.GetBusinessMonitor(req)
	if err != nil {
		global.ZC_LOG.Error("monitorCenterServer.GetBusinessMonitor:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}
