package lotus

import (
	"oplian/define"
	"oplian/global"
	"time"
)

type LotusWorkerInfo struct {
	global.ZC_MODEL
	OpId         string            `json:"opId"  gorm:"index;unique;comment:设备ID"`
	GateId       string            `json:"gateId" gorm:"index;comment:网关ID"`
	MinerId      uint64            `json:"minerId"  gorm:"comment:连接miner"`
	Ip           string            `json:"ip" gorm:"comment:ip"`
	WorkerType   define.WorkerType `json:"workerType" gorm:"comment:worker类型 0任务worker，1存储worker"`
	DeployStatus int               `json:"deployStatus" gorm:"comment:部署状态 1部署中，2部署成功，3部署失败"`
	RunStatus    int               `json:"runStatus" gorm:"comment:运行状态"`
	ErrMsg       string            `json:"errMsg" gorm:"type:text;comment:部署错误信息"`
	StartAt      time.Time         `json:"startAt" gorm:"default:CURRENT_TIMESTAMP(3);comment:部署开始时间"`
	FinishAt     time.Time         `json:"finishAt" gorm:"default:null;comment:部署结束时间"`
}

func (LotusWorkerInfo) TableName() string {
	return "lotus_worker_info"
}
