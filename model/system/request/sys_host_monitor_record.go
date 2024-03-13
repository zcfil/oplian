package request

import (
	"oplian/model/common/request"
)

type SysHostMonitorRecordSearch struct {
	HostUUID   string  `json:"hostUUID" form:"hostUUID"`
	UseRate    float32 `json:"useRate" form:"useRate"`
	IsLessThan bool    `json:"isLessThan" form:"isLessThan"`
	request.PageInfo
}

type HostUUIDsReq struct {
	HostUUIDs []HostChartReq `json:"hostUUIDs" form:"hostUUIDs"`
	Keyword   string         `json:"keyword" form:"keyword"`
}

type HostChartReq struct {
	HostUUID string `json:"hostUUID" form:"hostUUID"`
	GPUID    string `json:"gpuId" form:"gpuId"`
}

type GetHostRunListReq struct {
	RoomId string `json:"roomId" form:"roomId"`
	request.PageInfo
}
