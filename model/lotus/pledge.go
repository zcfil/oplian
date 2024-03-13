package lotus

import (
	"oplian/global"
)

type LotusPledgeConfig struct {
	global.ZC_MODEL
	MinerId        uint   `json:"minerId"  gorm:"index;comment:Miner ID"`
	OpId           string `json:"opId"  gorm:"index;comment:设备ID"`
	Actor          string `json:"actor"  gorm:"comment:节点ID"`
	Enabled        bool   `json:"enabled"  gorm:"comment:是否开启"`
	UndonePre      int    `json:"undonePre"  gorm:"comment:已发且未完成P1任务数"`
	DelayMinute    int    `json:"delayMinute" gorm:"comment:P1任务延迟时间"`
	WorkerBlannce  int    `json:"workerBlannce" gorm:"comment:worker钱包余额阈值"`
	FinSleepMinute int    `json:"finSleepMinute" gorm:"comment:转移延迟时间"`
}

func (LotusPledgeConfig) TableName() string {
	return "lotus_pledge_config"
}
