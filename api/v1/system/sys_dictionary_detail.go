package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system"
	"oplian/model/system/request"
	"oplian/utils"
)

type DictionaryDetailApi struct{}

// CreateSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   create SysDictionaryDetail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysDictionaryDetail/createSysDictionaryDetail [post]
func (s *DictionaryDetailApi) CreateSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindJSON(&detail)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryDetailService.CreateSysDictionaryDetail(detail)
	if err != nil {
		global.ZC_LOG.Error("Creation failed!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
		return
	}
	response.OkWithMessage("Created successfully", c)
}

// DeleteSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   delete SysDictionaryDetail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysDictionaryDetail/deleteSysDictionaryDetail [delete]
func (s *DictionaryDetailApi) DeleteSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindJSON(&detail)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryDetailService.DeleteSysDictionaryDetail(detail)
	if err != nil {
		global.ZC_LOG.Error("Delete failed!", zap.Error(err))
		response.FailWithMessage("Delete failed", c)
		return
	}
	response.OkWithMessage("Delete successful", c)
}

// UpdateSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   update SysDictionaryDetail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysDictionaryDetail/updateSysDictionaryDetail [put]
func (s *DictionaryDetailApi) UpdateSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindJSON(&detail)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryDetailService.UpdateSysDictionaryDetail(&detail)
	if err != nil {
		global.ZC_LOG.Error("Update failed!", zap.Error(err))
		response.FailWithMessage("Update failed", c)
		return
	}
	response.OkWithMessage("Update success", c)
}

// FindSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   Query SysDictionaryDetail with ID
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /sysDictionaryDetail/findSysDictionaryDetail [get]
func (s *DictionaryDetailApi) FindSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindQuery(&detail)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(detail, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	reSysDictionaryDetail, err := dictionaryDetailService.GetSysDictionaryDetail(detail.ID)
	if err != nil {
		global.ZC_LOG.Error("Query failed!", zap.Error(err))
		response.FailWithMessage("Query failed", c)
		return
	}
	response.OkWithDetailed(gin.H{"reSysDictionaryDetail": reSysDictionaryDetail}, "query was successful", c)
}

// GetSysDictionaryDetailList
// @Tags      SysDictionaryDetail
// @Summary   Paging to retrieve the SysDictionaryDetail list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysDictionaryDetail/getSysDictionaryDetailList [get]
func (s *DictionaryDetailApi) GetSysDictionaryDetailList(c *gin.Context) {
	var pageInfo request.SysDictionaryDetailSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := dictionaryDetailService.GetSysDictionaryDetailInfoList(pageInfo)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}
