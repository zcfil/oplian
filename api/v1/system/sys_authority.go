package system

import (
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/common/response"
	"oplian/model/system"
	systemReq "oplian/model/system/request"
	systemRes "oplian/model/system/response"
	"oplian/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorityApi struct{}

// CreateAuthority
// @Tags      Authority
// @Summary   create role
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=systemRes.SysAuthorityResponse,msg=string}
// @Router    /authority/createAuthority [post]
func (a *AuthorityApi) CreateAuthority(c *gin.Context) {
	var authority system.SysAuthority
	err := c.ShouldBindJSON(&authority)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(authority, utils.AuthorityVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if authBack, err := authorityService.CreateAuthority(authority); err != nil {
		global.ZC_LOG.Error("Creation failed!", zap.Error(err))
		response.FailWithMessage("Creation failed"+err.Error(), c)
	} else {
		_ = menuService.AddMenuAuthority(systemReq.DefaultMenu(), authority.AuthorityId)
		_ = casbinService.UpdateCasbin(authority.AuthorityId, systemReq.DefaultCasbin())
		response.OkWithDetailed(systemRes.SysAuthorityResponse{Authority: authBack}, "Created successfully", c)
	}
}

// CopyAuthority
// @Tags      Authority
// @Summary   Copy Role
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=systemRes.SysAuthorityResponse,msg=string}
// @Router    /authority/copyAuthority [post]
func (a *AuthorityApi) CopyAuthority(c *gin.Context) {
	var copyInfo systemRes.SysAuthorityCopyResponse
	err := c.ShouldBindJSON(&copyInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(copyInfo, utils.OldAuthorityVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(copyInfo.Authority, utils.AuthorityVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authBack, err := authorityService.CopyAuthority(copyInfo)
	if err != nil {
		global.ZC_LOG.Error("Copy failed!", zap.Error(err))
		response.FailWithMessage("Copy failed"+err.Error(), c)
		return
	}
	response.OkWithDetailed(systemRes.SysAuthorityResponse{Authority: authBack}, "Success copy", c)
}

// DeleteAuthority
// @Tags      Authority
// @Summary   remove role
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /authority/deleteAuthority [post]
func (a *AuthorityApi) DeleteAuthority(c *gin.Context) {
	var authority system.SysAuthority
	err := c.ShouldBindJSON(&authority)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(authority, utils.AuthorityIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = authorityService.DeleteAuthority(&authority)
	if err != nil { // Before deleting a role, it is necessary to determine if any users are using it
		global.ZC_LOG.Error("Delete failed!", zap.Error(err))
		response.FailWithMessage("Delete failed"+err.Error(), c)
		return
	}
	response.OkWithMessage("Delete successful", c)
}

// UpdateAuthority
// @Tags      Authority
// @Summary   Update role information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=systemRes.SysAuthorityResponse,msg=string}
// @Router    /authority/updateAuthority [post]
func (a *AuthorityApi) UpdateAuthority(c *gin.Context) {
	var auth system.SysAuthority
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(auth, utils.AuthorityVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authority, err := authorityService.UpdateAuthority(auth)
	if err != nil {
		global.ZC_LOG.Error("Update failed!", zap.Error(err))
		response.FailWithMessage("Update failed"+err.Error(), c)
		return
	}
	response.OkWithDetailed(systemRes.SysAuthorityResponse{Authority: authority}, "Update success", c)
}

// GetAuthorityList
// @Tags      Authority
// @Summary   Paging to obtain the role list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /authority/getAuthorityList [post]
func (a *AuthorityApi) GetAuthorityList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := authorityService.GetAuthorityInfoList(pageInfo)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// SetDataAuthority
// @Tags      Authority
// @Summary   Set role resource permissions
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /authority/setDataAuthority [post]
func (a *AuthorityApi) SetDataAuthority(c *gin.Context) {
	var auth system.SysAuthority
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(auth, utils.AuthorityIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = authorityService.SetDataAuthority(auth)
	if err != nil {
		global.ZC_LOG.Error("Setting failed!", zap.Error(err))
		response.FailWithMessage("Setting failed"+err.Error(), c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}
