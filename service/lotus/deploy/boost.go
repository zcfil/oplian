package deploy

import (
	"oplian/global"
	model "oplian/model/lotus"
	"oplian/model/lotus/response"
)

// AddBoost
// @author: nathan
// @function: Add Boost
// @description:  Add Boost
// @Param     model.LotusBoostInfo
// @return:  error
func (deploy *DeployService) AddBoost(boost *model.LotusBoostInfo) error {
	if boost.ID == 0 {
		return global.ZC_DB.Save(boost).Error
	}
	return global.ZC_DB.Updates(boost).Error
}

// GetBoost
//@author: nathan
//@function: GetBoost
//@description: Get boost
//@param: id uint64
//@return: model.LotusBoostInfo, error

func (deploy *DeployService) GetBoost(id uint64) (model.LotusBoostInfo, error) {
	var boost model.LotusBoostInfo
	return boost, global.ZC_DB.Model(model.LotusBoostInfo{}).Where("id = ?", id).First(&boost).Error
}

// GetBoost
//@author: nathan
//@function: GetBoostByOpId
//@description: boost Details
//@param: opId string
//@return: model.LotusBoostInfo, error

func (deploy *DeployService) GetBoostByOpId(OpId string) (response.BoostInfo, error) {

	var boost response.BoostInfo
	return boost, global.ZC_DB.Model(model.LotusBoostInfo{}).Where("op_id = ?", OpId).First(&boost).Error
}
