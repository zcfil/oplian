package response

import "oplian/utils"

type SysHostTestRecord struct {
	ID          uint   `json:"ID" gorm:"column:id"`
	HostUUID    string `json:"hostUUID" gorm:"column:host_uuid"`
	TestType    int    `json:"testType" gorm:"column:test_type"`
	TestBeginAt int64  `json:"-" gorm:"column:test_begin_at"`
	TestEndAt   int64  `json:"-" gorm:"column:test_end_at"`
	BeginTime   string `json:"beginTime" gorm:"-"`
	EndTime     string `json:"endTime" gorm:"-"`
	TestResult  int    `json:"testResult" gorm:"column:test_result"`
	TestMode    int    `json:"testMode" gorm:"column:test_mode"`
	IsAddPower  bool   `json:"isAddPower" gorm:"column:is_add_power"`

	HostName     string `json:"hostName" gorm:"-"`
	IntranetIP   string `json:"intranetIP" gorm:"-"`
	DeviceSN     string `json:"deviceSN" gorm:"-"`
	AssetNumber  string `json:"assetNumber" gorm:"-"`
	RoomId       string `json:"roomId" gorm:"-"`
	RoomName     string `json:"roomName" gorm:"-"`
	TestTakeTime string `json:"testTakeTime" gorm:"-"`
}

type SysHostTestReport struct {
	ID               uint              `json:"ID" gorm:"column:id"`
	HostUUID         string            `json:"hostUUID" gorm:"column:host_uuid"`
	AssetNumber      string            `json:"assetNumber" form:"assetNumber"`
	DeviceSN         string            `json:"deviceSN" form:"deviceSN"`
	IntranetIP       string            `json:"intranetIp" form:"intranetIp"`
	HostName         string            `json:"hostName" gorm:"-"`
	RoomName         string            `json:"roomName" form:"roomName"`
	RoomId           string            `json:"roomId" form:"roomId"`
	TestBeginAt      string            `json:"testBeginAt" form:"testBeginAt"`
	TestEndAt        string            `json:"testEndAt" form:"testEndAt"`
	TestResult       int64             `json:"testResult" form:"testResult"`
	TestType         int64             `json:"testType" form:"testType"`
	TestTakeTime     string            `json:"testTakeTime" form:"testTakeTime"`
	CPUHardInfo      utils.OpCPUInfo   `json:"cpuHardInfo" form:"cpuHardInfo"`
	CPUHardScore     int64             `json:"cpuHardScore" form:"cpuHardScore"`
	GPUHardInfo      []utils.OpGPUInfo `json:"gpuHardInfo" form:"gpuHardInfo"`
	GPUHardScore     int64             `json:"gpuHardScore" form:"gpuHardScore"`
	MemoryHardInfo   utils.OpRamInfo   `json:"memoryHardInfo" form:"memoryHardInfo"`
	MemoryHardScore  int64             `json:"memoryHardScore" form:"memoryHardScore"`
	DiskHardInfo     utils.HostDisk    `json:"diskHardInfo" form:"memoryHardInfo"`
	DiskHardScore    int64             `json:"diskHardScore" form:"diskHardScore"`
	NetTestInfo      string            `json:"netTestInfo" form:"netTestInfo"`
	NetTestScore     int64             `json:"netTestScore" form:"netTestScore"`
	GPUTestInfo      string            `json:"gpuTestInfo" form:"gpuTestInfo"`
	GPUTestScore     int64             `json:"gpuTestScore" form:"gpuTestScore"`
	DiskIO           int64             `json:"diskIO" form:"diskIO"`
	DiskAllRate      string            `json:"diskAllRate" form:"diskAllRate"`
	DiskAllRateScore int64             `json:"diskAllRateScore" form:"diskAllRateScore"`
	DiskNFSRate      string            `json:"diskNFSRate,omitempty" form:"diskNFSRate"`
	DiskNFSRateScore int64             `json:"diskNFSRateScore,omitempty" form:"diskNFSRateScore"`
	DiskSSDRate      string            `json:"diskSSDRate" form:"diskSSDRate"`
	DiskSSDRateScore int64             `json:"diskSSDRateScore" form:"diskSSDRateScore"`
	IsAddPower       bool              `json:"isAddPower" form:"isAddPower"`
	SelectHostUUIDs  string            `json:"selectHostUUIDs" form:"selectHostUUIDs"`
}

type DefaultHostTestReport struct {
	CPUThreads    int     `json:"cpuThreads"`
	GPUModel      int     `json:"gpuModel"`
	RamTotalMB    string  `json:"ramTotalMb"`
	RamTotalMBAdd string  `json:"ramTotalMbAdd"`
	DiskSize      string  `json:"diskSize"`
	NetSpeed      string  `json:"netSpeed"`
	NetSpeedAdd   string  `json:"netSpeedAdd"`
	GPURunTime    float64 `json:"gpuRunTime"`
	DiskIO        bool    `json:"diskIO"`
	SSDDiskRate   string  `json:"ssdDiskRate"`
	AllDiskRate   string  `json:"allDiskRate"`
}
