package response

import "time"

type WorkerCarDetail struct {
	Page     int   `json:"page"`      // 页码
	PageSize int   `json:"pageSize" ` // 每页大小
	Total    int64 `json:"total"`
	Detail   []WorkerCarTaskDetail
}

type TaskNoInfo struct {
	TaskNo int
}

type WorkerCarMiner struct {
	Remark      string `json:"remark"`
	Wallet      string `json:"wallet"`
	MaxSectorId string `json:"maxSectorId"`
	SectorSize  string `json:"sectorSize"`
}

type WorkerCarTaskDetail struct {
	TaskName       string ` json:"taskName"`                          // 任务名称
	CreatedAt      string ` json:"createdAt"`                         // 任务名称
	DealId         string ` json:"dealId"`                            // 订单ID
	MinerId        string ` json:"minerId"`                           // 节点
	DealExpireDate string ` json:"dealExpireDate"`                    // 订单到期时间
	CarName        string ` json:"carName"`                           // 原值car名称
	SectorSize     string `json:"sectorSize"`                         // 扇区大小
	TaskNo         int    `json:"taskNo"`                             // 任务编号
	TaskStatus     int    `json:"taskStatus"`                         //任务状态 1.待创建，2.创建中,3.已完成,4.创建失败,5.匹配失败,6.已过期
	SectorId       string `gorm:"index;comment:扇区ID" json:"sectorId"` // 扇区ID
	WaitTime       string `gorm:"comment:等待时间" json:"waitTime"`       // 等待时间

}

type CarTaskDetailInfo struct {
	Id           int    `json:"id"`            //ID
	PieceCid     string `json:"pieceCid"`      //piece id
	PieceSize    uint64 `json:"pieceSize"`     // 订单大小：32,64
	CarSize      int    `json:"carSize"`       // car文件大小
	DataCid      string `json:"dataCid"`       // DataCid
	QuotaWallet  string ` json:"quotaWallet"`  // 钱包
	JobStatus    int    `json:"jobStatus"`     //状态
	CarOutputDir string ` json:"carOutputDir"` //car文件目录
	OriginalOpId string ` json:"originalOpId"` //源主机opId
	OriginalDir  string ` json:"originalDir"`  //源主机目录
	ValidityDays int    ` json:"validityDays"` //订单生命周期
}

type WorkerTask struct {
	WorkerIp string
	Total    int
}

type BoostConfig struct {
	LanIp   string `json:"lanIp"`
	LanPort string `json:"lanPort"`
	Token   string `json:"token"`
}

type ConfigInfo struct {
	Actor         string `json:"actor"`
	DcQuotaWallet string `json:"dc_quota_wallet"`
	LotusToken    string `json:"lotus_token"`
	LotusIp       string `json:"lotus_ip"`
	MinerToken    string `json:"miner_token"`
	MinerIp       string `json:"miner_ip"`
}

type C2TaskInfo struct {
	MinerId string
	Number  int
	RunTime time.Time
}
