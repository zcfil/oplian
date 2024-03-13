package response

type DataScreeningReport struct {
	RemainingLicensesNum int   `json:"remainingLicensesNum"`
	WarnNum              int64 `json:"warnNum"`
	MachineRoomNum       int64 `json:"machineRoomNum"`
	FreeHostNum          int64 `json:"freeHostNum"`
	AllHostNum           int64 `json:"allHostNum"`
}

type HostUseDataResp struct {
	CPUUseRate    float32 `json:"cpuUseRate"`
	CPUAllNum     int64   `json:"cpuAllNum"`
	CPUUseNum     int     `json:"cpuUseNum"`
	MemoryUseRate float64 `json:"memoryUseRate"`
	MemoryAllSize string  `json:"memoryAllSize"`
	MemoryUseSize string  `json:"memoryUseSize"`
	DiskUseRate   float32 `json:"diskUseRate"`
	DiskAllSize   string  `json:"diskAllSize"`
	DiskUseSize   string  `json:"diskUseSize"`
}

type HostRunDataResp struct {
	ID             uint    `json:"ID"`
	HostUUID       string  `json:"hostUUID" `
	CPUUseRate     float32 `json:"cpuUseRate"`
	CPUTemperature string  `json:"cpuTemperature"`
	MemoryUseRate  float32 `json:"memoryUseRate"`
	DiskUseRate    float32 `json:"diskUseRate"`
	CreatedAt      string  `json:"createdAt"`
	HostName       string  `json:"hostName"`
	IntranetIP     string  `json:"intranetIp"`
	InternetIP     string  `json:"internetIp"`
	GroupName      string  `json:"groupName"`
	RoomId         string  `json:"roomId"`
	RoomName       string  `json:"roomName"`
}
