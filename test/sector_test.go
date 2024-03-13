package test

import (
	"fmt"
	"oplian/core"
	"oplian/define"
	"oplian/global"
	"oplian/initialize"
	model "oplian/model/lotus"
	"oplian/service"
	"oplian/service/pb"
	"testing"
)

func TestWorkerConfigList(t *testing.T) {
	global.ZC_VP = core.Viper("../config/config.yaml")
	global.ZC_DB = initialize.Gorm()
	sql := `CREATE TABLE IF NOT EXISTS lotus_sector_info_` + "test" + `(
				  id bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
				  created_at datetime(3) NULL DEFAULT NULL,
				  updated_at datetime(3) NULL DEFAULT NULL,
				  deleted_at datetime(3) NULL DEFAULT NULL,
				  sector_id bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '扇区ID',
				  actor varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '节点ID',
				  status varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '扇区状态',
				  cid_comm_d varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'Unsealed',
				  cid_comm_r varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'Sealed',
				  ticket varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '重做扇区字段',
				  ticket_h bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT 'ticket获取高度',
				  seed varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'WaitSeed',
				  seed_h bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT 'WaitSeed高度',
				  pre_cid varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'P2 消息ID',
				  commit_cid varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'C2 消息ID',
				  proof varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'C2证明结果',
				  commit1 varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'C1返回结果',
				  PRIMARY KEY (id) USING BTREE,
				  UNIQUE INDEX sector_id(sector_id) USING BTREE,
				  INDEX idx_lotus_sector_info_deleted_at(deleted_at) USING BTREE,
				  INDEX idx_lotus_sector_info_sector_id(sector_id) USING BTREE
				) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;
				`

	fmt.Println(global.ZC_DB.Exec(sql).Error)
}

func TestUpdateSectorStatus(t *testing.T) {
	global.ZC_VP = core.Viper("../config/config.yaml")
	global.ZC_DB = initialize.Gorm()
	service.ServiceGroupApp.LotusServiceGroup.UpdateSectorStatus(&pb.SectorStatus{Sector: &pb.SectorID{Miner: "f023013", Number: 0}, Status: define.SealPreCommit1})
}

func TestAdd(t *testing.T) {
	global.ZC_VP = core.Viper("../config/config.yaml")
	global.ZC_DB = initialize.Gorm()
	var date model.LotusSectorQueue
	global.ZC_DB.Where("id = ?", 1).First(&date)
	date.CompleteCount++
	fmt.Println(global.ZC_DB.Updates(date).Error)
}

func TestUpdateSectorQueue(t *testing.T) {
	global.ZC_VP = core.Viper("../config/config.yaml")
	global.ZC_DB = initialize.Gorm()
	var date model.LotusSectorQueue
	global.ZC_DB.Where("id = ?", 1).First(&date)
	date.CompleteCount++
	date.ID = 0
	fmt.Println(global.ZC_DB.Updates(date).Error)
}

func TestGetNextRunIndex(t *testing.T) {
	global.ZC_VP = core.Viper("../config/config.yaml")
	global.ZC_DB = initialize.Gorm()
	index := service.ServiceGroupApp.LotusServiceGroup.GetNextRunIndex("f018147", 7)
	fmt.Println(index)
}

func TestCreateSectorTable(t *testing.T) {
	global.ZC_VP = core.Viper("../config/config.yaml")
	global.ZC_DB = initialize.Gorm()
	index := service.ServiceGroupApp.LotusServiceGroup.CreateSectorTable("aaaa")
	fmt.Println(index)
}

func TestGetSectorQueueDetailID(t *testing.T) {
	global.ZC_VP = core.Viper("../config/config.yaml")
	global.ZC_DB = initialize.Gorm()
	index, err := service.ServiceGroupApp.LotusServiceGroup.GetSectorQueueDetail("f018147", 7)
	fmt.Println(index.QueueId, err)
}
