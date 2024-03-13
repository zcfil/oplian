package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system/request"
	systemRes "oplian/model/system/response"
	"oplian/utils"
)

type CasbinApi struct{}

// UpdateCasbin
// @Tags      Casbin
// @Summary   Update role API permissions
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /casbin/UpdateCasbin [post]
func (cas *CasbinApi) UpdateCasbin(c *gin.Context) {
	var cmr request.CasbinInReceive
	err := c.ShouldBindJSON(&cmr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(cmr, utils.AuthorityIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = casbinService.UpdateCasbin(cmr.AuthorityId, cmr.CasbinInfos)
	if err != nil {
		global.ZC_LOG.Error("Update failed!", zap.Error(err))
		response.FailWithMessage("Update failed", c)
		return
	}
	response.OkWithMessage("Update success", c)
}

// GetPolicyPathByAuthorityId
// @Tags      Casbin
// @Summary   Get permission list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=systemRes.PolicyPathResponse,msg=string}
// @Router    /casbin/getPolicyPathByAuthorityId [post]
func (cas *CasbinApi) GetPolicyPathByAuthorityId(c *gin.Context) {
	var casbin request.CasbinInReceive
	err := c.ShouldBindJSON(&casbin)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(casbin, utils.AuthorityIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	paths := casbinService.GetPolicyPathByAuthorityId(casbin.AuthorityId)
	response.OkWithDetailed(systemRes.PolicyPathResponse{Paths: paths}, "Successfully obtained", c)
}
