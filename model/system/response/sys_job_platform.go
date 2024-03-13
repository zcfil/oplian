package response

import "time"

type DeviceResponse struct {
}

type DistributionResponse struct {
	Data interface{} `json:"data"`
}

type FileInfo struct {
	FileName  string
	FileSize  int
	CreatedAt time.Time
}

type MinerInfoRes struct {
	Ip   string `json:"ip"`
	OpId string `json:"op_id"`
}
