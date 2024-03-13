package request

import "oplian/model/common/request"

type SysAutoHistory struct {
	request.PageInfo
}

// GetById Find by id structure
type RollBack struct {
	ID          int  `json:"id" form:"id"`
	DeleteTable bool `json:"deleteTable" form:"deleteTable"`
}
