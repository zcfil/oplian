package request

import (
	"oplian/model/common/request"
)

type SysHostTestRecordSearch struct {
	TestType           int64  `json:"testType" form:"testType"`
	TestResult         int64  `json:"testResult" form:"testResult"`
	RoomId             string `json:"roomId" form:"roomId"`
	HostNameIPKeyword  string `json:"hostNameIPKeyword" form:"hostNameIPKeyword"`
	HostAssetSNKeyword string `json:"hostAssetSNKeyword" form:"hostAssetSNKeyword"`
	request.PageInfo
}

type GetHostTestReportReq struct {
	ID       int    `json:"id" form:"id"`
	HostUUID string `json:"hostUUID" form:"hostUUID"`
}

type AddHostTestByHandReq struct {
	HostUUID   string      `json:"hostUUID" form:"hostUUID"`
	TestType   int64       `json:"testType" form:"testType"`
	IsAddPower bool        `json:"isAddPower" form:"isAddPower"`
	HostUUIDs  []UUIDAndIP `json:"hostUUIDs" form:"hostUUIDs"`
}

type UUIDAndIP struct {
	IntranetIP string `json:"intranetIP" form:"intranetIP"`
	HostUUID   string `json:"hostUUID" form:"hostUUID"`
}

type CloseHostTestReq struct {
	ID       int    `json:"id" form:"id"`
	HostUUID string `json:"hostUUID" form:"hostUUID"`
}

type AddHostTestRepeatReq struct {
	ID       int    `json:"id" form:"id"`
	HostUUID string `json:"hostUUID" form:"hostUUID"`
}

type GetDefaultHostTestInfoReq struct {
	TestType int64 `json:"testType" form:"testType"`
}
