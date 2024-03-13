package initialize

import (
	"log"
	"oplian/global"
	"oplian/model/example"
	"oplian/model/gateway"
	"oplian/model/lotus"
	"oplian/model/slot"
	"oplian/model/system"
	"oplian/service"
	"os"

	adapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Gorm Initializes the database and generates database global variables
// Author SliverHorn
func Gorm() *gorm.DB {
	switch global.ZC_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	default:
		return GormMysql()
	}
}

// RegisterTables Register database tables dedicated
// Author SliverHorn
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(

		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysAutoCode{},
		adapter.CasbinRule{},
		system.SysHostRecord{},
		system.SysMachineRoomRecord{},
		system.SysHostGroup{},
		system.SysHostMonitorRecord{},
		system.SysHostTestRecord{},
		system.SysHostPatrolRecord{},
		system.SysPatrolConfig{},
		system.SysFileManage{},
		system.SysJobExecuteRecords{},
		system.SysOpRelations{},
		system.SysWarnManage{},
		system.SysWarnStrategies{},
		system.SysDictionaryConfig{},
		system.SysColony{},
		system.SysFileManageUpload{},

		example.ExaFileUploadAndDownload{},

		lotus.LotusPledgeConfig{},
		lotus.LotusEnv{},
		lotus.LoutsWorkerConfig{},
		lotus.LotusWorkerInfo{},
		lotus.LotusStorageInfo{},
		lotus.LotusMinerInfo{},
		lotus.LotusInfo{},
		lotus.LotusSectorInfo{},
		lotus.LotusSectorPiece{},
		lotus.LotusSectorLog{},
		lotus.LotusSectorQueue{},
		lotus.LotusSectorQueueDetail{},
		lotus.LotusSectorRecover{},
		lotus.LotusSectorTask{},
		lotus.LotusSectorTaskDetail{},
		lotus.LotusWorkerCluster{},
		lotus.LotusWorkerRelations{},
		lotus.LotusBoostInfo{},
		lotus.LotusSectorAbnormal{},

		gateway.DownloadRecord{},
		//slot
		slot.WorkerCarTask{},
		slot.WorkerCarTaskDetail{},
		slot.WorkerCarTaskNo{},
		slot.WorkerCarRand{},
		slot.LotusWorkerTask{},
		slot.WorkerCarFiles{},
	)
	if err != nil {
		global.ZC_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	log.Println("register table success")
}

// InitMysqlData Initialize data
func InitMysqlData() error {
	return service.ServiceGroupApp.SystemServiceGroup.InitDBService.InitData("mysql")
}
