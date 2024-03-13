package slot

import "oplian/global"

type WorkerCarTaskDetail struct {
	global.ZC_MODEL
	MinerId        string `gorm:"comment:miner节点" json:"minerId"`                                                         // miner节点
	WorkerIp       string `gorm:"index;comment:worker机IP" json:"workerIp"`                                                // worker机IP
	TaskId         int    `gorm:"index;comment:任务ID" json:"taskId"`                                                       // 任务ID
	TaskName       string `gorm:"comment:任务名称" json:"taskName"`                                                           // 任务名称
	DealId         string `gorm:"comment:订单ID" json:"dealId"`                                                             // 订单ID
	DealExpireDate string `gorm:"comment:订单到期时间" json:"dealExpireDate"`                                                   // 订单到期时间
	WalletAddr     string `gorm:"comment:钱包地址" json:"walletAddr"`                                                         // 钱包地址
	CarName        string `gorm:"comment:原值car名称" json:"carName"`                                                         // 原值car名称
	PieceCid       string `gorm:"comment:car名称ID" json:"pieceCid"`                                                        // car名称ID
	PieceSize      int    `gorm:"comment:Piece大小" json:"pieceSize"`                                                       // Piece大小
	CarSize        int    `gorm:"comment:car文件大小" json:"carSize"`                                                         // car文件大小
	DataCid        string `gorm:"comment:数据CID" json:"dataCid"`                                                           // 数据CID
	SectorSize     string `gorm:"comment:扇区大小" json:"sectorSize"`                                                         // 扇区大小
	TaskNo         int    `gorm:"comment:任务编号" json:"taskNo"`                                                             // 任务编号
	TaskStatus     int    `gorm:"default:0;comment:任务状态 0待创建,1成功,2创建中,3创建失败,4匹配失败,5car复制完成,6发单,7导入car" json:"taskStatus"` //任务状态 0待创建,1成功,2创建中,3创建失败,4匹配失败,5car复制完成,6发单,7导入car
	TaskRun        int    `gorm:"default:0;comment:状态 0等待,1运行" json:"taskRun"`                                            // 状态 0等待,1运行
	SectorId       string `gorm:"index;comment:扇区ID" json:"sectorId"`                                                     // 扇区ID
	WaitTime       string `gorm:"comment:等待时间" json:"waitTime"`                                                           // 等待时间
	ErrMsg         string `gorm:"remark;comment:备注" json:"errMsg"`                                                        // 错误信息
}

func (WorkerCarTaskDetail) TableName() string {
	return "worker_car_task_detail"
}
