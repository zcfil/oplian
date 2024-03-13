package test

import (
	"oplian/core"
	"oplian/define"
	"oplian/global"
	"oplian/initialize"
	"oplian/service"
	"oplian/service/pb"
	"testing"
)

func TestGetFileName(t *testing.T) {
	global.ZC_VP = core.Viper("../config/config.yaml")
	global.ZC_DB = initialize.Gorm()
	service.ServiceGroupApp.SystemServiceGroup.GetFileName(&pb.FileNameInfo{GateWayId: "44928f41-3b8e-43e8-9e47-46f7504f4562", FileType: define.ProveFile.Int64()})
}
