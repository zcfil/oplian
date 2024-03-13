package request

import (
	"oplian/model/common/request"
	"oplian/model/system"
)

type SysDictionarySearch struct {
	system.SysDictionary
	request.PageInfo
}
