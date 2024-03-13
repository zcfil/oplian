package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/common/response"
	systemReq "oplian/model/system/request"
	"oplian/utils"
	"strconv"
)

type HostGroupApi struct{}

// CreateSysHostGroup
// @Tags      SysHostGroup
// @Summary   create SysHostGroup
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysHostGroup/createSysHostGroup [post]
func (s *HostGroupApi) CreateSysHostGroup(c *gin.Context) {
	var sysHostGroup systemReq.CreateSysHostGroupReq
	err := c.ShouldBindJSON(&sysHostGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if !utils.IsStringsUnique(sysHostGroup.GroupNames) {
		response.FailWithMessage("The array name passed in is duplicated", c)
		return
	}

	// 创建主机分组列表
	err = hostGroupService.CreateSysHostGroupList(sysHostGroup.GroupNames)
	if err != nil {
		if err == global.GroupNameRepeatError {
			response.FailWithMessage("The array name is duplicated with the original name", c)
			return
		}
		global.ZC_LOG.Error("Creation failed!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
		return
	}

	response.OkWithMessage("Created successfully", c)
}

// DeleteSysHostGroupByIds
// @Tags      SysHostGroup
// @Summary   Batch deletion of SysHostGroup
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysHostGroup/deleteSysHostGroupByIds [delete]
func (s *HostGroupApi) DeleteSysHostGroupByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	for _, val := range ids.Ids {
		hostNum, err := hostRecordService.GetSysHostRecordCountByHostGroupId(val)
		if err != nil {
			global.ZC_LOG.Error("Failed to obtain the number of group associated hosts!", zap.Error(err))
			response.FailWithMessage("Failed to obtain the number of group associated hosts!", c)
			return
		}
		if hostNum > 0 {
			response.FailWithMessage("The deleted group has an associated host, deletion failed!", c)
			return
		}
	}
	err = hostGroupService.DeleteSysHostGroupByIds(ids)
	if err != nil {
		global.ZC_LOG.Error("Batch deletion failed!", zap.Error(err))
		response.FailWithMessage("Batch deletion failed", c)
		return
	}
	response.OkWithMessage("Batch deletion successful", c)
}

// GetSysHostGroupList
// @Tags      SysHostGroup
// @Summary   Paging to obtain the SysHostGroup list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostGroup/getSysHostGroupList [get]
func (s *HostGroupApi) GetSysHostGroupList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := hostGroupService.GetSysHostGroupInfoList(pageInfo)
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

// DealSysHostGroupInfo
// @Tags      SysHostGroup
// @Summary   Batch processing of host grouping
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  json     request.UpdateSysHostGroupReq
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostGroup/dealSysHostGroupInfo [post]
func (s *HostGroupApi) DealSysHostGroupInfo(c *gin.Context) {
	var sysHostRecord []systemReq.UpdateSysHostGroupReq
	err := c.ShouldBindJSON(&sysHostRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var updateRecord []systemReq.UpdateSysHostGroupReq
	var createRecord []string
	var keepIds []int
	var groupNames []string

	for _, val := range sysHostRecord {
		if val.ID != 0 {
			var info systemReq.UpdateSysHostGroupReq
			info.ID = val.ID
			info.GroupName = val.GroupName
			updateRecord = append(updateRecord, info)
			keepIds = append(keepIds, val.ID)
		} else {
			createRecord = append(createRecord, val.GroupName)
		}
		groupNames = append(groupNames, val.GroupName)
	}

	// Determine if the inserted data name is duplicate
	if !utils.IsStringsUnique(groupNames) {
		response.FailWithMessage("The name of the incoming data group is duplicated", c)
		return
	}

	// Processing front-end data
	questionId, err := hostGroupService.DealSysHostGroupList(updateRecord, createRecord, keepIds)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.FailWithMessage("Incoming data group ID:"+strconv.Itoa(questionId)+" error", c)
			return
		} else if err == global.GroupIDDeletedBoundError {
			response.FailWithMessage("The ID to be deleted:"+strconv.Itoa(questionId)+"bound with host information", c)
			return
		}
		global.ZC_LOG.Error("Modification failed!", zap.Error(err))
		response.FailWithMessage("Modification failed", c)
		return
	}
	response.OkWithMessage("Modified successfully", c)
}
