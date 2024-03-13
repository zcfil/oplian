package request

import (
	"oplian/model/common/request"
	"oplian/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
