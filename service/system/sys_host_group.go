package system

import (
	"gorm.io/gorm"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/system"
	systemReq "oplian/model/system/request"
	"oplian/model/system/response"
)


//@function: CreateSysHostGroup
//@description: Example Create host packet data
//@param: sysHostGroup model.SysHostGroup
//@return: err error

type HostGroupService struct{}

func (hostGroupService *HostGroupService) CreateSysHostGroup(sysHostGroup system.SysHostGroup) (err error) {
	err = global.ZC_DB.Create(&sysHostGroup).Error
	return err
}


// @function: CreateSysHostGroupList
// @description: Create a grouping data list
// @param: id int
// @return: sysHostGroup system.SysHostGroup, err error

func (hostGroupService *HostGroupService) CreateSysHostGroupList(groupList []string) (err error) {
	tx := global.ZC_DB.Begin()
	defer func() {
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	for _, val := range groupList {
		// 查询分组信息
		var hostGroupInfo system.SysHostGroup
		err = tx.Where("group_name = ?", val).First(&hostGroupInfo).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				var hostGroup system.SysHostGroup
				hostGroup.GroupName = val
				err = tx.Create(&hostGroup).Error
				if err != nil {
					return
				}
			} else {
				return
			}
		} else {
			err = global.GroupNameRepeatError
			return
		}
	}
	return
}


//@function: GetSysHostGroup
//@description: Obtain a single piece of host information based on the id
//@param: id int
//@return: sysHostGroup system.SysHostGroup, err error

func (hostGroupService *HostGroupService) GetSysHostGroup(groupId int) (sysHostGroup system.SysHostGroup, err error) {
	err = global.ZC_DB.Where("id = ?", groupId).First(&sysHostGroup).Error
	return
}


//@function: UpdateSysHostGroup
//@description: Update host grouping data
//@param: sysHostGroup *model.SysHostGroup
//@return: err error

func (hostGroupService *HostGroupService) UpdateSysHostGroup(sysHostGroup *system.SysHostGroup) (err error) {
	var dict system.SysHostGroup
	sysHostGroupMap := map[string]interface{}{
		"GroupName": sysHostGroup.GroupName,
	}
	db := global.ZC_DB.Where("id = ?", sysHostGroup.ID).First(&dict)
	err = db.Updates(sysHostGroupMap).Error
	return err
}

//@function: DeleteSysHostGroupByRoomId
//@description: Batch delete record
//@param: ids request.IdsReq
//@return: err error

func (hostGroupService *HostGroupService) DeleteSysHostGroupByIds(ids request.IdsReq) (err error) {
	err = global.ZC_DB.Delete(&[]system.SysHostGroup{}, "id in (?)", ids.Ids).Unscoped().Error
	return err
}


// @function: GetList
// @description: Paginate to get a list of host groups
// @param: info request.PageInfo
// @return: list interface{}, total int64, err error

func (hostGroupService *HostGroupService) GetSysHostGroupInfoList(info request.PageInfo) (list []response.SysHostGroup, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&system.SysHostGroup{})
	var entities []response.SysHostGroup
	err = db.Count(&total).Error
	if err != nil {
		return nil, total, err
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&entities).Error
	return entities, total, err
}


//@function: DealSysHostGroupList
//@description: Batch processing host packet data
//@param: sysHostGroup *model.SysHostGroup
//@return: err error

func (hostGroupService *HostGroupService) DealSysHostGroupList(updateRecord []systemReq.UpdateSysHostGroupReq, createRecord []string, keepIds []int) (questionId int, err error) {
	tx := global.ZC_DB.Begin()
	defer func() {
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	for _, val := range updateRecord {

		var hostGroupInfo system.SysHostGroup
		err = tx.Where("id = ?", val.ID).First(&hostGroupInfo).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				questionId = val.ID
			}
			return
		}
		hostGroupInfo.GroupName = val.GroupName
		err = tx.Where("id = ?", val.ID).First(&system.SysHostGroup{}).Updates(&hostGroupInfo).Error
		if err != nil {
			return
		}
	}

	var delIds []int
	err = tx.Model(&system.SysHostGroup{}).Where("id NOT IN (?)", keepIds).Pluck("id", &delIds).Error
	if err != nil {
		return
	}

	for _, val := range delIds {
		var total int64
		db := tx.Model(&system.SysHostRecord{})
		err = db.Where("host_group_id = ?", val).Count(&total).Error
		if err != nil {
			return
		}
		if total > 0 {
			questionId = val
			err = global.GroupIDDeletedBoundError
			return
		}
	}

	err = tx.Delete(&[]system.SysHostGroup{}, "id in (?)", delIds).Error
	if err != nil {
		return
	}

	for _, val := range createRecord {
		dict := system.SysHostGroup{GroupName: val}
		err = tx.Create(&dict).Error
		if err != nil {
			return
		}
	}
	return
}

// @function: GetHostGroupList
// @description: Paginate to get a list of host groups
// @param: info request.PageInfo
// @return: list interface{}, total int64, err error

func (hostGroupService *HostGroupService) GetHostGroupListMap() (map[int]string, error) {
	list := map[int]string{}
	var hostGroups []response.SysHostGroup
	err := global.ZC_DB.Model(&system.SysHostGroup{}).Order("id desc").Find(&hostGroups).Error
	if err != nil {
		return nil, err
	}

	for _, val := range hostGroups {
		list[int(val.ID)] = val.GroupName
	}
	return list, err
}
