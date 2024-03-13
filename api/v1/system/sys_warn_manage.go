package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system/request"
	"oplian/utils"
)

type WarnManageApi struct {
}

// WarnList
// @Tags      WarnList
// @Summary   Alarm Center List
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /jobPlatform/fileDistribution [post]
func (w *WarnManageApi) WarnList(c *gin.Context) {

	var wr request.WarnReq
	err := c.ShouldBindJSON(&wr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wr, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := warnManageService.GetWarnList(wr)
	if err != nil {
		global.ZC_LOG.Error("warnManageService.GetWarnList:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     wr.Page,
		PageSize: wr.PageSize,
	}, "Successfully obtained", c)
}

// ModifyWarnStatus
// @Tags      ModifyWarnStatus
// @Summary   Alarm Center List
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /jobPlatform/fileDistribution [post]
func (w *WarnManageApi) ModifyWarnStatus(c *gin.Context) {

	var wr request.WarnReq
	err := c.ShouldBindJSON(&wr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wr, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	b, err := warnManageService.ModifyWarnStatus(wr)
	if err != nil {
		global.ZC_LOG.Error("warnManageService.ModifyWarnStatus:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(b, "Successfully obtained", c)
}

// WarnTotal
// @Tags      WarnTotal
// @Summary   Alarm Center Data Statistics
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /jobPlatform/fileDistribution [post]
func (w *WarnManageApi) WarnTotal(c *gin.Context) {

	var wr request.WarnReq
	err := c.ShouldBindJSON(&wr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wr, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fb, err := warnManageService.GetWarnTotal(wr)
	if err != nil {
		global.ZC_LOG.Error("warnManageService.GetWarnTrend:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(fb, "Successfully obtained", c)
}

// WarnTrend
// @Tags      WarnTrend
// @Summary   Trend chart of alarm center
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /jobPlatform/fileDistribution [post]
func (w *WarnManageApi) WarnTrend(c *gin.Context) {

	var wr request.WarnReq
	err := c.ShouldBindJSON(&wr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wr, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fb, err := warnManageService.GetWarnTrend(wr)
	if err != nil {
		global.ZC_LOG.Error("warnManageService.GetWarnTrend:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(fb, "Successfully obtained", c)
}

// SaveStrategy
// @Tags      SaveStrategy
// @Summary   Save/Modify Alarm Policies
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /warnManage/saveStrategy [post]
func (w *WarnManageApi) SaveStrategy(c *gin.Context) {

	var wr request.WarnStrategiesReq
	err := c.ShouldBindJSON(&wr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wr, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = warnManageService.SaveStrategy(wr)
	if err != nil {
		global.ZC_LOG.Error("warnManageService.GetWarnTrend:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "Successfully obtained", c)
}

// StrategyType
// @Tags      StrategyType
// @Summary   Policy type
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /warnManage/strategyList [post]
func (w *WarnManageApi) StrategyType(c *gin.Context) {

	var wr request.StrategyReq
	err := c.ShouldBindJSON(&wr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wr, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := warnManageService.StrategyType(wr.StrategyType)
	if err != nil {
		global.ZC_LOG.Error("warnManageService.GetStrategyList:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// StrategyId
// @Tags      StrategyId
// @Summary   Policy ID
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /warnManage/strategyList [post]
func (w *WarnManageApi) StrategyId(c *gin.Context) {

	res := warnManageService.GetStrategyId()
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// StrategyList
// @Tags      StrategyList
// @Summary   Alarm strategy list
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /warnManage/strategyList [post]
func (w *WarnManageApi) StrategyList(c *gin.Context) {

	var wr request.StrategyReq
	err := c.ShouldBindJSON(&wr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wr, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := warnManageService.GetStrategyList(wr)
	if err != nil {
		global.ZC_LOG.Error("warnManageService.GetStrategyList:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     wr.Page,
		PageSize: wr.PageSize,
	}, "Successfully obtained", c)
}

// StrategyDetail
// @Tags      StrategyDetail
// @Summary   Alarm strategy details
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: data}
// @Router    /warnManage/strategyDetail [post]
func (w *WarnManageApi) StrategyDetail(c *gin.Context) {

	var wr request.StrategyReq
	err := c.ShouldBindJSON(&wr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wr, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := warnManageService.GetStrategyDetail(wr)
	if err != nil {
		global.ZC_LOG.Error("warnManageService.GetStrategyDetail:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "Successfully obtained", c)
}

// DelStrategy
// @Tags      DelStrategy
// @Summary   Delete alarm strategy
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /warnManage/delStrategy [post]
func (w *WarnManageApi) DelStrategy(c *gin.Context) {

	var wr request.StrategyReq
	err := c.ShouldBindJSON(&wr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wr, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = warnManageService.DelStrategy(wr)
	if err != nil {
		global.ZC_LOG.Error("warnManageService.delStrategy:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "Successfully obtained", c)
}

// ModifyStrategyStatus
// @Tags      ModifyStrategyStatus
// @Summary   Change alarm policy status
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /warnManage/modifyStrategyStatus [post]
func (w *WarnManageApi) ModifyStrategyStatus(c *gin.Context) {

	var wr request.StrategyReq
	err := c.ShouldBindJSON(&wr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wr, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = warnManageService.ModifyStrategyStatus(wr)
	if err != nil {
		global.ZC_LOG.Error("warnManageService.ModifyStrategyStatus:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "Successfully obtained", c)
}
