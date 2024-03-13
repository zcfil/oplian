package lotus

import (
	"oplian/global"
	"time"
)

type LotusMinerInfo struct {
	global.ZC_MODEL
	OpId         string    `json:"opId"  gorm:"index;unique;comment:设备ID"`
	GateId       string    `json:"gateId" gorm:"index;comment:网关ID"`
	LotusId      uint64    `json:"lotusId"  gorm:"comment:连接lotus"`
	Actor        string    `json:"minerId" gorm:"comment:节点号"`
	SectorSize   uint64    `json:"sectorSize" gorm:"comment:扇区大小"`
	Token        string    `json:"token"  gorm:"comment:节点token"`
	Ip           string    `json:"ip" gorm:"comment:ip"`
	Port         string    `json:"port" gorm:"comment:端口号"`
	DeployStatus int       `json:"deployStatus" gorm:"comment:部署状态：1部署中，2部署成功，3部署失败"`
	RunStatus    int       `json:"runStatus" gorm:"comment:运行状态"`
	IsManage     bool      `json:"isManage" gorm:"comment:是否打开调度"`
	IsWdpost     bool      `json:"isWdpost" gorm:"comment:是否打开wdpost"`
	IsWnpost     bool      `json:"isWnpost" gorm:"comment:是否打开wnpost"`
	Partitions   string    `json:"partitions" gorm:"comment:wdpost partitions号"`
	ErrMsg       string    `json:"errMsg" gorm:"type:text;comment:部署错误信息"`
	StartAt      time.Time `json:"startAt" gorm:"default:CURRENT_TIMESTAMP(3);comment:部署开始时间"`
	FinishAt     time.Time `json:"finishAt" gorm:"default:null;comment:部署结束时间"`
	AddType      int       `json:"addType"  gorm:"comment:部署方式：1 系统内节点部署,2.全新节点部署,3.链上节点部署"`
}

func (LotusMinerInfo) TableName() string {
	return "lotus_miner_info"
}
