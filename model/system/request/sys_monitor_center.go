package request

type MonitorCenterReq struct {
	GateWayId string `json:"gateWayId" form:"gateWayId"`
	MinerId   string `json:"minerId" form:"minerId"`
}
