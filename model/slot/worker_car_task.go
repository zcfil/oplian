package slot

import (
	"oplian/global"
	"time"
)

type WorkerCarTask struct {
	global.ZC_MODEL
	TaskName       string    `gorm:"comment:任务名称" json:"taskName"`                                           // 任务名称
	TaskType       int       `gorm:"comment:任务类型 1自动,2手动" json:"taskType"`                                   // 任务类型 1自动,2手动
	SectorType     int       `gorm:"comment:扇区类型 1 cc,2 dc" json:"sectorType"`                               // 扇区类型 1 cc,2 dc
	SectorSize     string    `gorm:"comment:扇区大小" json:"sectorSize"`                                         // 扇区大小
	MinerId        string    `gorm:"comment:miner节点" json:"minerId"`                                         // miner节点
	QuotaWallet    string    `gorm:"comment:DC份额钱包" json:"quotaWallet"`                                      // DC份额钱包
	DataSourcePath string    `gorm:"comment:数据源地址" json:"dataSourcePath"`                                    // 数据源地址
	OriginalOpId   string    `gorm:"comment:源值主机opId" json:"originalOpId"`                                   // 源值主机opId
	OriginalDir    string    `gorm:"comment:源值主机目录" json:"originalDir"`                                      // 源值主机目录
	OrderRange     string    `gorm:"comment:订单数量范围" json:"orderRange"`                                       // 订单数量范围
	WorkerTaskNum  int       `gorm:"comment:worker任务数" json:"workerTaskNum"`                                 // worker任务数
	CurrentNo      int       `gorm:"comment:当前任务数" json:"currentNo"`                                         // 当前任务数
	FinishNum      int       `gorm:"comment:完成任务数" json:"finishNum"`                                         // 完成任务数
	ValidityDays   int       `gorm:"comment:有效期天数" json:"validityDays"`                                      // 有效期天数
	FinishTime     time.Time `gorm:"default:null;comment:完成时间" json:"finishTime"`                            // 完成时间
	TaskStatus     int       `gorm:"default:0;comment:任务状态 0进行中,1已完成,2暂停,3已终止,4已过期,5异常失败" json:"taskStatus"` // 任务状态
	Remark         string    `gorm:"remark;comment:备注" json:"remark"`                                        // 备注
}

func (WorkerCarTask) TableName() string {
	return "worker_car_task"
}
