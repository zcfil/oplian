package lotus

import (
	"oplian/global"
	"time"
)

type LotusBoostInfo struct {
	global.ZC_MODEL
	OpId          string    `json:"opId"  gorm:"index;unique;comment:设备ID"`
	GateId        string    `json:"gateId" gorm:"index;comment:网关ID"`
	LotusId       uint64    `json:"lotusId"  gorm:"comment:连接lotus"`
	MinerId       uint64    `json:"minerId"  gorm:"comment:连接miner"`
	Token         string    `json:"token"  gorm:"comment:节点token"`
	LanIp         string    `json:"lanIp" gorm:"comment:局域网IP"`
	LanPort       string    `json:"lanPort" gorm:"comment:局域网端口号"`
	InternetIp    string    `json:"internetIp" gorm:"comment:公网IP"`
	InternetPort  string    `json:"internetPort" gorm:"comment:公网端口号"`
	DeployStatus  int       `json:"deployStatus" gorm:"comment:部署状态：1部署中，2部署成功，3部署失败"`
	NetworkType   int       `json:"networkType" gorm:"comment:网络模式：0局域网映射，1独立公网"`
	AskStatus     int       `json:"askStatus" gorm:"comment:询价状态：0未询价，1正常，2询价失败"`
	RunStatus     int       `json:"runStatus" gorm:"comment:运行状态"`
	ErrMsg        string    `json:"errMsg" gorm:"type:text;comment:部署错误信息"`
	StartAt       time.Time `json:"startAt" gorm:"default:CURRENT_TIMESTAMP(3);comment:部署开始时间"`
	FinishAt      time.Time `json:"finishAt" gorm:"default:null;comment:部署结束时间"`
	DcQuotaWallet string    `json:"dcQuotaWallet" gorm:"type:text;comment:dc额度钱包"`
}

func (LotusBoostInfo) TableName() string {
	return "lotus_boost_info"
}
