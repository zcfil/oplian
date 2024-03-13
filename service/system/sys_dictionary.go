package system

import (
	"errors"

	"gorm.io/gorm"
	"oplian/global"
	"oplian/model/system"
	"oplian/model/system/request"
)


//@function: DeleteSysDictionary
//@description: Create dictionary data
//@param: sysDictionary model.SysDictionary
//@return: err error

type DictionaryService struct{}

func (dictionaryService *DictionaryService) CreateSysDictionary(sysDictionary system.SysDictionary) (err error) {
	if (!errors.Is(global.ZC_DB.First(&system.SysDictionary{}, "type = ?", sysDictionary.Type).Error, gorm.ErrRecordNotFound)) {
		return errors.New("the same type exists and cannot be created")
	}
	err = global.ZC_DB.Create(&sysDictionary).Error
	return err
}


//@function: DeleteSysDictionary
//@description: Delete dictionary data
//@param: sysDictionary model.SysDictionary
//@return: err error

func (dictionaryService *DictionaryService) DeleteSysDictionary(sysDictionary system.SysDictionary) (err error) {
	err = global.ZC_DB.Where("id = ?", sysDictionary.ID).Preload("SysDictionaryDetails").First(&sysDictionary).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("请不要搞事")
	}
	if err != nil {
		return err
	}
	err = global.ZC_DB.Delete(&sysDictionary).Error
	if err != nil {
		return err
	}

	if sysDictionary.SysDictionaryDetails != nil {
		return global.ZC_DB.Where("sys_dictionary_id=?", sysDictionary.ID).Delete(sysDictionary.SysDictionaryDetails).Error
	}
	return
}


//@function: UpdateSysDictionary
//@description: Update dictionary data
//@param: sysDictionary *model.SysDictionary
//@return: err error

func (dictionaryService *DictionaryService) UpdateSysDictionary(sysDictionary *system.SysDictionary) (err error) {
	var dict system.SysDictionary
	sysDictionaryMap := map[string]interface{}{
		"Name":   sysDictionary.Name,
		"Type":   sysDictionary.Type,
		"Status": sysDictionary.Status,
		"Desc":   sysDictionary.Desc,
	}
	db := global.ZC_DB.Where("id = ?", sysDictionary.ID).First(&dict)
	if dict.Type != sysDictionary.Type {
		if !errors.Is(global.ZC_DB.First(&system.SysDictionary{}, "type = ?", sysDictionary.Type).Error, gorm.ErrRecordNotFound) {
			return errors.New("The same type exists and cannot be created")
		}
	}
	err = db.Updates(sysDictionaryMap).Error
	return err
}


//@function: GetSysDictionary
//@description: Obtain the dictionary data by id or type
//@param: Type string, Id uint
//@return: err error, sysDictionary model.SysDictionary

func (dictionaryService *DictionaryService) GetSysDictionary(Type string, Id uint, status *bool) (sysDictionary system.SysDictionary, err error) {
	var flag = false
	if status == nil {
		flag = true
	} else {
		flag = *status
	}
	err = global.ZC_DB.Where("(type = ? OR id = ?) and status = ?", Type, Id, flag).Preload("SysDictionaryDetails", "status = ?", true).First(&sysDictionary).Error
	return
}


//@function: GetSysDictionaryInfoList
//@description: Page for a list of dictionaries
//@param: info request.SysDictionarySearch
//@return: err error, list interface{}, total int64

func (dictionaryService *DictionaryService) GetSysDictionaryInfoList(info request.SysDictionarySearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.ZC_DB.Model(&system.SysDictionary{})
	var sysDictionarys []system.SysDictionary

	if info.Name != "" {
		db = db.Where("`name` LIKE ?", "%"+info.Name+"%")
	}
	if info.Type != "" {
		db = db.Where("`type` LIKE ?", "%"+info.Type+"%")
	}
	if info.Status != nil {
		db = db.Where("`status` = ?", info.Status)
	}
	if info.Desc != "" {
		db = db.Where("`desc` LIKE ?", "%"+info.Desc+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&sysDictionarys).Error
	return sysDictionarys, total, err
}
