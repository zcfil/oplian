package request

import (
	"oplian/model/common/request"
)

type SysHostPatrolRecordSearch struct {
	PatrolType         int64  `json:"patrolType" form:"patrolType"`
	PatrolResult       int64  `json:"patrolResult" form:"patrolResult"`
	RoomId             string `json:"roomId" form:"roomId"`
	HostNameIPKeyword  string `json:"hostNameIPKeyword" form:"hostNameIPKeyword"`
	HostAssetSNKeyword string `json:"hostAssetSNKeyword" form:"hostAssetSNKeyword"`
	request.PageInfo
}

type GetHostPatrolReportReq struct {
	ID       int    `json:"id" form:"id"`
	HostUUID string `json:"hostUUID" form:"hostUUID"`
}

type AddHostPatrolByHandReq struct {
	PatrolType int64  `json:"patrolType" form:"patrolType"`
	HostUUID   string `json:"hostUUID" form:"hostUUID"`
}
