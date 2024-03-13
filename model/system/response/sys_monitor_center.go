package response

type BusinessListReq struct {
	SectorStatus   string      `json:"sectorStatus"`
	PendTotal      int         `json:"pendTotal"`
	ProcessTotal   int         `json:"processTotal"`
	TimeOutTotal   int         `json:"timeOutTotal"`
	TimeLength     string      `json:"timeLength"`
	TotalCompleted int64       `json:"totalCompleted"`
	OpListRes      []OpListRes `json:"op_list_res"`
}
type OpListRes struct {
	//GateWayId  string `json:"gateWayId"`
	MinerId  string `json:"minerId"`
	SectorID uint64 `json:"sectorID"`
	//OpId       string `json:"opId"`
	Ip         string `json:"ip"`
	Progress   int    `json:"progress"`
	TimeLength string `json:"timeLength"`
}

type TimeTotal struct {
	AvgSecond int64
	Total     int64
}
