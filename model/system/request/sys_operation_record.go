package request

import (
	"oplian/model/common/request"
	"oplian/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
