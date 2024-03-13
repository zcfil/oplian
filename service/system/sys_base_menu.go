package system

import (
	"errors"

	"gorm.io/gorm"
	"oplian/global"
	"oplian/model/system"
)

type BaseMenuService struct{}


//@function: DeleteBaseMenu
//@description: Delete a basic Route
//@param: id float64
//@return: err error

func (baseMenuService *BaseMenuService) DeleteBaseMenu(id int) (err error) {
	err = global.ZC_DB.Preload("MenuBtn").Preload("Parameters").Where("parent_id = ?", id).First(&system.SysBaseMenu{}).Error
	if err != nil {
		var menu system.SysBaseMenu
		db := global.ZC_DB.Preload("SysAuthoritys").Where("id = ?", id).First(&menu).Delete(&menu)
		err = global.ZC_DB.Delete(&system.SysBaseMenuParameter{}, "sys_base_menu_id = ?", id).Error
		err = global.ZC_DB.Delete(&system.SysBaseMenuBtn{}, "sys_base_menu_id = ?", id).Error
		err = global.ZC_DB.Delete(&system.SysAuthorityBtn{}, "sys_menu_id = ?", id).Error
		if err != nil {
			return err
		}
		if len(menu.SysAuthoritys) > 0 {
			err = global.ZC_DB.Model(&menu).Association("SysAuthoritys").Delete(&menu.SysAuthoritys)
		} else {
			err = db.Error
			if err != nil {
				return
			}
		}
	} else {
		return errors.New("this menu has submenus and cannot be deleted")
	}
	return err
}


//@function: UpdateBaseMenu
//@description: Update route
//@param: menu model.SysBaseMenu
//@return: err error

func (baseMenuService *BaseMenuService) UpdateBaseMenu(menu system.SysBaseMenu) (err error) {
	var oldMenu system.SysBaseMenu
	upDateMap := make(map[string]interface{})
	upDateMap["keep_alive"] = menu.KeepAlive
	upDateMap["close_tab"] = menu.CloseTab
	upDateMap["default_menu"] = menu.DefaultMenu
	upDateMap["parent_id"] = menu.ParentId
	upDateMap["path"] = menu.Path
	upDateMap["name"] = menu.Name
	upDateMap["hidden"] = menu.Hidden
	upDateMap["component"] = menu.Component
	upDateMap["title"] = menu.Title
	upDateMap["active_name"] = menu.ActiveName
	upDateMap["icon"] = menu.Icon
	upDateMap["sort"] = menu.Sort

	err = global.ZC_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", menu.ID).Find(&oldMenu)
		if oldMenu.Name != menu.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", menu.ID, menu.Name).First(&system.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
				global.ZC_LOG.Debug("description Failed to modify the same name")
				return errors.New("description Failed to modify the same name")
			}
		}
		txErr := tx.Unscoped().Delete(&system.SysBaseMenuParameter{}, "sys_base_menu_id = ?", menu.ID).Error
		if txErr != nil {
			global.ZC_LOG.Debug(txErr.Error())
			return txErr
		}
		txErr = tx.Unscoped().Delete(&system.SysBaseMenuBtn{}, "sys_base_menu_id = ?", menu.ID).Error
		if txErr != nil {
			global.ZC_LOG.Debug(txErr.Error())
			return txErr
		}
		if len(menu.Parameters) > 0 {
			for k := range menu.Parameters {
				menu.Parameters[k].SysBaseMenuID = menu.ID
			}
			txErr = tx.Create(&menu.Parameters).Error
			if txErr != nil {
				global.ZC_LOG.Debug(txErr.Error())
				return txErr
			}
		}

		if len(menu.MenuBtn) > 0 {
			for k := range menu.MenuBtn {
				menu.MenuBtn[k].SysBaseMenuID = menu.ID
			}
			txErr = tx.Create(&menu.MenuBtn).Error
			if txErr != nil {
				global.ZC_LOG.Debug(txErr.Error())
				return txErr
			}
		}

		txErr = db.Updates(upDateMap).Error
		if txErr != nil {
			global.ZC_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}


//@function: GetBaseMenuById
//@description: Returns the currently selected menu
//@param: id float64
//@return: menu system.SysBaseMenu, err error

func (baseMenuService *BaseMenuService) GetBaseMenuById(id int) (menu system.SysBaseMenu, err error) {
	err = global.ZC_DB.Preload("MenuBtn").Preload("Parameters").Where("id = ?", id).First(&menu).Error
	return
}
