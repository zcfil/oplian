package system

import (
	"errors"
	"fmt"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/system"

	"gorm.io/gorm"
)

type ApiService struct{}

var ApiServiceApp = new(ApiService)

func (apiService *ApiService) CreateApi(api system.SysApi) (err error) {
	if !errors.Is(global.ZC_DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&system.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("same API exists")
	}
	return global.ZC_DB.Create(&api).Error
}

func (apiService *ApiService) DeleteApi(api system.SysApi) (err error) {
	var entity system.SysApi
	err = global.ZC_DB.Where("id = ?", api.ID).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	err = global.ZC_DB.Delete(&entity).Error
	if err != nil {
		return err
	}
	success := CasbinServiceApp.ClearCasbin(1, entity.Path, entity.Method)
	if !success {
		return errors.New(entity.Path + ":" + entity.Method + "casbin synchronization cleanup failed")
	}
	e := CasbinServiceApp.Casbin()
	err = e.InvalidateCache()
	if err != nil {
		return err
	}
	return nil
}

func (apiService *ApiService) GetAPIInfoList(api system.SysApi, info request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&system.SysApi{})
	var apiList []system.SysApi

	if api.Path != "" {
		db = db.Where("path LIKE ?", "%"+api.Path+"%")
	}

	if api.Description != "" {
		db = db.Where("description LIKE ?", "%"+api.Description+"%")
	}

	if api.Method != "" {
		db = db.Where("method = ?", api.Method)
	}

	if api.ApiGroup != "" {
		db = db.Where("api_group = ?", api.ApiGroup)
	}

	err = db.Count(&total).Error

	if err != nil {
		return apiList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			orderMap := make(map[string]bool, 5)
			orderMap["id"] = true
			orderMap["path"] = true
			orderMap["api_group"] = true
			orderMap["description"] = true
			orderMap["method"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("illegal sorting field: %v", order)
				return apiList, total, err
			}

			err = db.Order(OrderStr).Find(&apiList).Error
		} else {
			err = db.Order("api_group").Find(&apiList).Error
		}
	}
	return apiList, total, err
}

func (apiService *ApiService) GetAllApis() (apis []system.SysApi, err error) {
	err = global.ZC_DB.Find(&apis).Error
	return
}

func (apiService *ApiService) GetApiById(id int) (api system.SysApi, err error) {
	err = global.ZC_DB.Where("id = ?", id).First(&api).Error
	return
}

func (apiService *ApiService) UpdateApi(api system.SysApi) (err error) {
	var oldA system.SysApi
	err = global.ZC_DB.Where("id = ?", api.ID).First(&oldA).Error
	if oldA.Path != api.Path || oldA.Method != api.Method {
		if !errors.Is(global.ZC_DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&system.SysApi{}).Error, gorm.ErrRecordNotFound) {
			return errors.New("same API path exists")
		}
	}
	if err != nil {
		return err
	} else {
		err = CasbinServiceApp.UpdateCasbinApi(oldA.Path, api.Path, oldA.Method, api.Method)
		if err != nil {
			return err
		} else {
			err = global.ZC_DB.Save(&api).Error
		}
	}
	return err
}

func (apiService *ApiService) DeleteApisByIds(ids request.IdsReq) (err error) {
	var apis []system.SysApi
	err = global.ZC_DB.Find(&apis, "id in ?", ids.Ids).Delete(&apis).Error
	if err != nil {
		return err
	} else {
		for _, sysApi := range apis {
			success := CasbinServiceApp.ClearCasbin(1, sysApi.Path, sysApi.Method)
			if !success {
				return errors.New(sysApi.Path + ":" + sysApi.Method + "casbin synchronization cleanup failed")
			}
		}
		e := CasbinServiceApp.Casbin()
		err = e.InvalidateCache()
		if err != nil {
			return err
		}
	}
	return err
}
