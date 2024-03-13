package request

import (
	"oplian/model/common/request"
	"oplian/model/system"
)

// api分页条件查询及排序结构体
type SearchApiParams struct {
	system.SysApi
	request.PageInfo
	OrderKey string `json:"orderKey"`
	Desc     bool   `json:"desc"`
}
