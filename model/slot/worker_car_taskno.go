package slot

import "oplian/global"

type WorkerCarTaskNo struct {
	global.ZC_MODEL
	TaskId     int    `gorm:"index;comment:任务Id" json:"taskId"`                            // 任务Id
	MinerId    string `gorm:"comment:miner节点" json:"minerId"`                              // miner节点
	WorkerIp   string `gorm:"comment:workerIp" json:"workerIp"`                            // workerIp
	CarNo      string `gorm:"comment:文件编号" json:"orderRange"`                              // 文件编号
	InputDir   string `gorm:"comment:输入目录" json:"inputDir"`                                // 输入目录
	OutputDir  string `gorm:"comment:输出目录" json:"outputDir"`                               // 输出目录
	StartNo    int    `gorm:"comment:开始文件编号" json:"startNo"`                               // 开始文件编号
	EndNo      int    `gorm:"comment:结束文件编号" json:"endNo"`                                 // 结束文件编号
	TaskStatus int    `gorm:"default:0;comment:任务状态 0进行中,1已完成,2暂停,3已终止" json:"taskStatus"` // 任务状态 0进行中,1已完成,2暂停,3已终止
}

func (WorkerCarTaskNo) TableName() string {
	return "worker_car_task_no"
}
