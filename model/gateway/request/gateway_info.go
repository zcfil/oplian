package request

type GateWayInfo struct {
	GateWayId string `json:"gateWayId"`
	IP        string `json:"ip"`
	Port      string `json:"port"`
	Token     string `json:"token"`
}
