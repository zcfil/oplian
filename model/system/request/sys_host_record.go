package request

import (
	"oplian/model/common/request"
)

type SysHostRecordSearch struct {
	HostKeyword  string `json:"hostKeyword" form:"hostKeyword"`
	RoomKeyword  string `json:"roomKeyword" form:"roomKeyword"`
	UUID         string `json:"uuid" form:"uuid"`
	RoomId       string `json:"roomId" form:"roomId"`
	DeviceSN     string `json:"deviceSN" form:"deviceSN"`
	HostType     *int   `json:"hostType" form:"hostType"`
	HostClassify *int   `json:"hostClassify" form:"hostClassify"`
	HostGroupId  *int   `json:"hostGroupId" form:"hostGroupId"`
	request.PageInfo
}
