package system

import (
	"oplian/global"
	"oplian/model/system"
)

type ColonyService struct{}

func (colonyervice *ColonyService) AddColony(colony *system.SysColony) (err error) {
	global.ZC_DB.Delete(&system.SysColony{}, "colony_name = ?", colony.ColonyName)
	return global.ZC_DB.Create(colony).Error
}

func (colonyervice *ColonyService) GetColony(minerId string) (colony system.SysColony, err error) {
	return colony, global.ZC_DB.Where("colony_name = ?", minerId).Find(&colony).Error
}
