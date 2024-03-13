package utils

// op测试主机返回信息的数据结构
type HostTestInfo struct {
	ID              string `json:"id"`
	HostUUID        string `json:"hostUUID"`
	TestResult      int64  `json:"testResult"`
	CPUHardInfo     int64  `json:"cpuHardInfo"`
	CPUHardScore    int64  `json:"cpuHardScore"`
	GPUHardInfo     string `json:"gpuHardInfo"`
	GPUHardScore    int64  `json:"gpuHardScore"`
	MemoryHardInfo  string `json:"memoryHardInfo"`
	MemoryHardScore int64  `json:"memoryHardScore"`
	DiskHardInfo    string `json:"diskHardInfo"`
	DiskHardScore   int64  `json:"diskHardScore"`
	NetTestInfo     string `json:"netTestInfo"`
	NetTestScore    int64  `json:"netTestScore"`
	GPUTestInfo     int64  `json:"gpuTestInfo"`
	GPUTestScore    int64  `json:"gpuTestScore"`
	MemoryTestInfo  string `json:"memoryTestInfo"`
	MemoryTestScore int64  `json:"memoryTestScore"`
	DiskTestInfo    string `json:"diskTestInfo"`
	DiskTestScore   int64  `json:"diskTestScore"`
}
