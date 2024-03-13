package system

import (
	"oplian/global"
	"oplian/model/system"
	"oplian/model/system/request"
)


//@function: CreateSysDictionaryDetail
//@description: Create dictionary detail data
//@param: sysDictionaryDetail model.SysDictionaryDetail
//@return: err error

type DictionaryDetailService struct{}

func (dictionaryDetailService *DictionaryDetailService) CreateSysDictionaryDetail(sysDictionaryDetail system.SysDictionaryDetail) (err error) {
	err = global.ZC_DB.Create(&sysDictionaryDetail).Error
	return err
}


//@function: DeleteSysDictionaryDetail
//@description: Delete dictionary details data
//@param: sysDictionaryDetail model.SysDictionaryDetail
//@return: err error

func (dictionaryDetailService *DictionaryDetailService) DeleteSysDictionaryDetail(sysDictionaryDetail system.SysDictionaryDetail) (err error) {
	err = global.ZC_DB.Delete(&sysDictionaryDetail).Error
	return err
}


//@function: UpdateSysDictionaryDetail
//@description: Update dictionary details data
//@param: sysDictionaryDetail *model.SysDictionaryDetail
//@return: err error

func (dictionaryDetailService *DictionaryDetailService) UpdateSysDictionaryDetail(sysDictionaryDetail *system.SysDictionaryDetail) (err error) {
	err = global.ZC_DB.Save(sysDictionaryDetail).Error
	return err
}


//@function: GetSysDictionaryDetail
//@description: Gets a single piece of dictionary detail data by id
//@param: id uint
//@return: sysDictionaryDetail system.SysDictionaryDetail, err error

func (dictionaryDetailService *DictionaryDetailService) GetSysDictionaryDetail(id uint) (sysDictionaryDetail system.SysDictionaryDetail, err error) {
	err = global.ZC_DB.Where("id = ?", id).First(&sysDictionaryDetail).Error
	return
}


//@function: GetSysDictionaryDetailInfoList
//@description: Page for a list of dictionary details
//@param: info request.SysDictionaryDetailSearch
//@return: list interface{}, total int64, err error

func (dictionaryDetailService *DictionaryDetailService) GetSysDictionaryDetailInfoList(info request.SysDictionaryDetailSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.ZC_DB.Model(&system.SysDictionaryDetail{})
	var sysDictionaryDetails []system.SysDictionaryDetail

	if info.Label != "" {
		db = db.Where("label LIKE ?", "%"+info.Label+"%")
	}
	if info.Value != 0 {
		db = db.Where("value = ?", info.Value)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}
	if info.SysDictionaryID != 0 {
		db = db.Where("sys_dictionary_id = ?", info.SysDictionaryID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("sort").Find(&sysDictionaryDetails).Error
	return sysDictionaryDetails, total, err
}
