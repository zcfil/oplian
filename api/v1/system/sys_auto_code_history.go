package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/common/response"
	systemReq "oplian/model/system/request"
)

type AutoCodeHistoryApi struct{}

// First
// @Tags      AutoCode
// @Summary   Obtaining Meta Information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /autoCode/getMeta [post]
func (a *AutoCodeHistoryApi) First(c *gin.Context) {
	var info request.GetById
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := autoCodeHistoryService.First(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"meta": data}, "Successfully obtained", c)
}

// Delete
// @Tags      AutoCode
// @Summary   Delete rollback record
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /autoCode/delSysHistory [post]
func (a *AutoCodeHistoryApi) Delete(c *gin.Context) {
	var info request.GetById
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = autoCodeHistoryService.Delete(&info)
	if err != nil {
		global.ZC_LOG.Error("Delete failed!", zap.Error(err))
		response.FailWithMessage("Delete failed", c)
		return
	}
	response.OkWithMessage("Delete successful", c)
}

// RollBack
// @Tags      AutoCode
// @Summary   Rollback automatic code generation
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /autoCode/rollback [post]
func (a *AutoCodeHistoryApi) RollBack(c *gin.Context) {
	var info systemReq.RollBack
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = autoCodeHistoryService.RollBack(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Rollback successful", c)
}

// GetList
// @Tags      AutoCode
// @Summary   Query rollback records
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /autoCode/getSysHistory [post]
func (a *AutoCodeHistoryApi) GetList(c *gin.Context) {
	var search systemReq.SysAutoHistory
	err := c.ShouldBindJSON(&search)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := autoCodeHistoryService.GetList(search.PageInfo)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "Successfully obtained", c)
}
