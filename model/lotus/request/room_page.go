package request

// PageInfo Paging common input parameter structure
type RoomPageInfo struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Keyword  string `json:"keyword" form:"keyword"`
	GateId   string `json:"gateId" form:"gateId"`
	OnOff    int    `json:"onOff" form:"onOff"`
	Actor    string `json:"actor" form:"actor"`
}
