// 自动生成模板SysHostTTestRecord
package system

import (
	"oplian/global"
)

// 如果含有time.Time 请自行import time包
type SysHostTestRecord struct {
	global.ZC_MODEL
	HostUUID        string `json:"hostUUID" form:"hostUUID" gorm:"column:host_uuid;comment:资产编号UUID"`
	TestType        int64  `json:"testType" form:"testType" gorm:"column:test_type;comment:测试类型, 1 节点机测试, 2 worker机测试, 3 存储机测试, 4 C2 Worker机测试"`
	TestBeginAt     int64  `json:"testBeginAt" form:"testBeginAt" gorm:"column:test_begin_at;comment:测试开始时间 (时间戳)"`
	TestEndAt       int64  `json:"testEndAt" form:"testEndAt" gorm:"column:test_end_at;comment:测试结束时间 (时间戳)"`
	TestResult      int64  `json:"testResult" form:"testResult" gorm:"column:test_result;comment:测试结果, 1 测试中, 2 达标, 3 不达标, 4 测试失败"`
	TestMode        int64  `json:"testMode" form:"testMode" gorm:"column:test_mode;comment:测试方式, 1 自动触发, 2 手动触发"`
	CPUHardInfo     string `json:"cpuHardInfo" form:"cpuHardInfo" gorm:"column:cpu_hard_info;comment:CPU硬件基本参数"`
	CPUHardScore    int64  `json:"cpuHardScore" form:"cpuHardScore" gorm:"column:cpu_hard_score;comment:CPU硬件评分, 1 合格, 2 不满足"`
	GPUHardInfo     string `json:"gpuHardInfo" form:"gpuHardInfo" gorm:"type:varchar(512);column:gpu_hard_info;comment:GPU硬件基本参数"`
	GPUHardScore    int64  `json:"gpuHardScore" form:"gpuHardScore" gorm:"column:gpu_hard_score;comment:GPU硬件评分, 1 合格, 2 不满足"`
	MemoryHardInfo  string `json:"memoryHardInfo" form:"memoryHardInfo" gorm:"column:memory_hard_info;comment:内存硬件基本参数"`
	MemoryHardScore int64  `json:"memoryHardScore" form:"memoryHardScore" gorm:"column:memory_hard_score;comment:内存硬件评分, 1 合格, 2 不满足"`
	DiskHardInfo    string `json:"diskHardInfo" form:"memoryHardInfo" gorm:"column:disk_hard_info;comment:磁盘硬件基本参数"`
	DiskHardScore   int64  `json:"diskHardScore" form:"diskHardScore" gorm:"column:disk_hard_score;comment:磁盘硬件评分, 1 合格, 2 不满足"`
	NetTestInfo     string `json:"netTestInfo" form:"netTestInfo" gorm:"column:net_test_info;comment:网络测试结果"`
	NetTestScore    int64  `json:"netTestScore" form:"netTestScore" gorm:"column:net_test_score;comment:网络测试评分, 1 合格, 2 不满足"`
	GPUTestInfo     string `json:"gpuTestInfo" form:"gpuTestInfo" gorm:"column:gpu_test_info;comment:GPU压力测试结果,时长(单位Min)"`
	GPUTestScore    int64  `json:"gpuTestScore" form:"gpuTestScore" gorm:"column:gpu_test_score;comment:GPU压力测试评分, 1 合格, 2 不满足"`
	//MemoryTestInfo  string `json:"memoryTestInfo" form:"memoryTestInfo" gorm:"column:memory_test_info;comment:内存测试结果"`
	//MemoryTestScore int64  `json:"memoryTestScore" form:"memoryTestScore" gorm:"column:memory_test_score;comment:内存测试评分, 1 合格, 2 不满足"`
	//DiskTestInfo    string `json:"diskTestInfo" form:"diskTestInfo" gorm:"column:disk_test_info;comment:磁盘容量测试结果"`
	//DiskTestScore   int64  `json:"diskTestScore" form:"diskTestScore" gorm:"column:disk_test_score;comment:磁盘容量测试评分, 1 合格, 2 不满足"`
	DiskIO           int64  `json:"diskIO" form:"diskIO" gorm:"column:disk_io;comment:磁盘有无IO, 0 无, 1 有"`
	DiskAllRate      string `json:"diskAllRate" form:"diskAllRate" gorm:"column:disk_all_rate;comment:全盘读写速率"`
	DiskAllRateScore int64  `json:"diskAllRateScore" form:"diskAllRateScore" gorm:"column:disk_all_rate_score;comment:全盘读写速率评分, 1 合格, 2 不满足"`
	DiskNFSRate      string `json:"diskNFSRate" form:"diskNFSRate" gorm:"column:disk_nfs_rate;comment:NFS全盘读写速率"`
	DiskNFSRateScore int64  `json:"diskNFSRateScore" form:"diskNFSRateScore" gorm:"column:disk_nfs_rate_score;comment:NFS全盘读写速率评分, 1 合格, 2 不满足"`
	DiskSSDRate      string `json:"diskSSDRate" form:"diskSSDRate" gorm:"column:disk_ssd_rate;comment:SSD全盘读写速率"`
	DiskSSDRateScore int64  `json:"diskSSDRateScore" form:"diskSSDRateScore" gorm:"column:disk_ssd_rate_score;comment:SSD全盘读写速率评分, 1 合格, 2 不满足"`
	IsAddPower       bool   `json:"isAddPower" form:"isAddPower" gorm:"column:is_add_power;comment:是否新增算力, 1 新增, 2 不新增"`
	SelectHostUUIDs  string `json:"selectHostUUIDs" form:"selectHostUUIDs" gorm:"type:varchar(512);column:select_host_uuids;comment:网络主机uuid拼接的字符串,逗号隔开"`
	SelectHostIPs    string `json:"selectHostIPs" form:"selectHostIPs" gorm:"column:select_host_ips;comment:网络主机ip拼接的字符串,逗号隔开"`
}

func (SysHostTestRecord) TableName() string {
	return "sys_host_test_records"
}
