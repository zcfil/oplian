package lotus

import (
	"oplian/global"
)

type LoutsWorkerConfig struct {
	global.ZC_MODEL
	OpId      string `json:"opId" gorm:"index;unique;comment:服务器唯一标识"`
	GateId    string `json:"gateId" gorm:"index;comment:机房ID"`
	WorkerId  uint   `json:"workerId" gorm:"index;comment:worker表ID"`
	Actor     string `json:"actor"  gorm:"index;comment:节点ID"`
	PreCount1 int    `json:"preCount1" gorm:"comment:配置P1任务数量"`
	PreCount2 int    `json:"preCount2" gorm:"comment:配置P2任务数量"`
	IP        string `json:"ip" gorm:"comment:ip"`
	Port      string `json:"port" gorm:"comment:端口"`
	OnOff1    bool   `json:"onOff1" gorm:"default:true;comment:任务开关"`
}

func (LoutsWorkerConfig) TableName() string {
	return "lotus_worker_config"
}
