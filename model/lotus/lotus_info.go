package lotus

import (
	"oplian/global"
	"time"
)

type LotusInfo struct {
	global.ZC_MODEL
	OpId         string    `json:"opId"  gorm:"index;unique;comment:设备ID"`
	GateId       string    `json:"gateId" gorm:"index;comment:网关ID"`
	Token        string    `json:"token"  gorm:"comment:节点token"`
	Ip           string    `json:"ip" gorm:"comment:ip"`
	Port         string    `json:"port" gorm:"comment:端口号"`
	Actor        string    `json:"actor" gorm:"comment:节点号"`
	DeployStatus int       `json:"deployStatus" gorm:"comment:部署状态:部署中，2部署成功，3部署失败"`
	SyncStatus   int       `json:"syncStatus" gorm:"comment:同步状态"`
	RunStatus    int       `json:"runStatus" gorm:"comment:运行状态"`
	StartAt      time.Time `json:"startAt" gorm:"default:CURRENT_TIMESTAMP(3);comment:部署开始时间"`
	FinishAt     time.Time `json:"finishAt" gorm:"default:null;comment:部署完成时间"`
	SnapshotAt   time.Time `json:"snapshotAt" gorm:"default:CURRENT_TIMESTAMP(3);comment:更换快照时间"`
	ErrMsg       string    `json:"errMsg" gorm:"type:text;comment:部署失败信息"`
}

func (LotusInfo) TableName() string {
	return "lotus_info"
}
