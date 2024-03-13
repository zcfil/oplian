package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system/request"
)

type AuthorityBtnApi struct{}

// GetAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   Get permission button
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=response.SysAuthorityBtnRes,msg=string}
// @Router    /authorityBtn/getAuthorityBtn [post]
func (a *AuthorityBtnApi) GetAuthorityBtn(c *gin.Context) {
	var req request.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	res, err := authorityBtnService.GetAuthorityBtn(req)
	if err != nil {
		global.ZC_LOG.Error("Query failed!", zap.Error(err))
		response.FailWithMessage("Query failed", c)
		return
	}
	response.OkWithDetailed(res, "query was successful", c)
}

// SetAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   Set permission button
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /authorityBtn/setAuthorityBtn [post]
func (a *AuthorityBtnApi) SetAuthorityBtn(c *gin.Context) {
	var req request.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = authorityBtnService.SetAuthorityBtn(req)
	if err != nil {
		global.ZC_LOG.Error("allocation failure!", zap.Error(err))
		response.FailWithMessage("allocation failure", c)
		return
	}
	response.OkWithMessage("Allocation successful", c)
}

// CanRemoveAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   Set permission button
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}  "Delete successful"
// @Router    /authorityBtn/canRemoveAuthorityBtn [post]
func (a *AuthorityBtnApi) CanRemoveAuthorityBtn(c *gin.Context) {
	id := c.Query("id")
	err := authorityBtnService.CanRemoveAuthorityBtn(id)
	if err != nil {
		global.ZC_LOG.Error("Delete failed!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Delete successful", c)
}
