package lotus

import (
	"oplian/global"
	"time"
)

type LotusStorageInfo struct {
	global.ZC_MODEL
	OpId         string    `json:"opId"  gorm:"index;unique;comment:设备ID"`
	GateId       string    `json:"gateId" gorm:"index;comment:网关ID"`
	MinerId      uint64    `json:"minerId"  gorm:"comment:连接miner"`
	Ip           string    `json:"ip" gorm:"comment:ip"`
	DeployStatus int       `json:"deployStatus" gorm:"comment:部署状态:1部署中，2部署成功，3部署失败"`
	ColonyName   string    `json:"colonyName" gorm:"comment:节点名称"`
	ColonyType   int       `json:"colonyType" gorm:"comment:存储类型：1NFS,2worker"`
	ErrMsg       string    `json:"errMsg" gorm:"type:text;comment:部署错误信息"`
	StartAt      time.Time `json:"startAt" gorm:"default:CURRENT_TIMESTAMP(3);comment:部署开始时间"`
	FinishAt     time.Time `json:"finishAt" gorm:"default:null;comment:部署结束时间"`
	NFSDisk      string    `json:"nfsDisk" gorm:"type:text;comment:NFS节点挂载的磁盘信息"`
}

func (LotusStorageInfo) TableName() string {
	return "lotus_storage_info"
}
