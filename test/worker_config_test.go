package test

import (
	"fmt"
	"oplian/core"
	"oplian/global"
	op "oplian/initialize"
	"oplian/service"
	"testing"
)

//func TestWorkerConfigList(t *testing.T) {
//	global.ZC_VP = core.Viper("../config/config.yaml")
//	global.ZC_DB = initialize.Gorm()
//	res, err := service.ServiceGroupApp.LotusServiceGroup.GetWorkerConfigList("f023013")
//	fmt.Println(res, err)
//}
//
//func TestGetWorkerConfig(t *testing.T) {
//	global.ZC_VP = core.Viper("../config/config.yaml")
//	global.ZC_DB = initialize.Gorm()
//	res, err := service.ServiceGroupApp.LotusServiceGroup.GetWorkerConfig("1")
//	fmt.Println(res, err)
//}
//
//func TestSetWorkerConfig(t *testing.T) {
//	global.ZC_VP = core.Viper("../config/config.yaml")
//	global.ZC_DB = initialize.Gorm()
//	err := service.ServiceGroupApp.LotusServiceGroup.SetWrokerConfig(lotus.LoutsWorkerConfig{ZC_MODEL: global.ZC_MODEL{ID: 1, CreatedAt: time.Now()}, PreCount1: 14, PreCount2: 1, OpId: "1", IP: "10.0.1.197"})
//	fmt.Println(err)
//}

func TestIsExistNode(t *testing.T) {
	global.ZC_VP = core.Viper("../config/config.yaml")
	global.ZC_DB = op.Gorm()
	err := service.ServiceGroupApp.LotusServiceGroup.IsExistNode("f023013")
	fmt.Println(err)
}
