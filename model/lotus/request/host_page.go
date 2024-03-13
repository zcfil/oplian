package request

// PageInfo Paging common input parameter structure
type HostPage struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Keyword  string `json:"keyword" form:"keyword"`
	Classify int    `json:"classify" form:"classify"`
	GateId   string `json:"gateId" form:"gateId"`
}
