package request

import (
	"oplian/model/common/request"
)

type SysMachineRoomRecordSearch struct {
	request.PageInfo
}

type BindSysHostRecordsReq struct {
	RoomId    string   `json:"roomId" form:"roomId"`
	HostUUIDs []string `json:"hostUUIDs" form:"hostUUIDs"`
}

type BindSysHostRecordListReq struct {
	Bind    bool   `json:"bind" form:"bind"`
	RoomId  string `json:"roomId" form:"roomId"`
	Keyword string `json:"keyword" form:"keyword"`
}

type UnbindSysHostRecordsReq struct {
	RoomId    string   `json:"roomId" form:"roomId"`
	HostUUIDs []string `json:"hostUUIDs" form:"hostUUIDs"`
}

type SysHostRecordListByClassifyReq struct {
	GateWayId    string `json:"gateWayId" form:"gateWayId"`
	RoomId       string `json:"roomId" form:"roomId"`
	UUID         string `json:"uuid" form:"uuid"`
	HostClassify int    `json:"hostClassify" form:"hostClassify"`
	KeyWord      string `json:"keyWord" form:"keyWord"`
}

type GetOpHardwareInfoReq struct {
	UUID string `json:"uuid" form:"uuid"`
}

type GetNetHostListReq struct {
	UUID    string `json:"uuid" form:"uuid"`
	Keyword string `json:"keyword" form:"keyword"`
}

type GetPatrolHostListReq struct {
	PatrolType int64  `json:"patrolType" form:"patrolType"`
	Keyword    string `json:"keyword" form:"keyword"`
}
