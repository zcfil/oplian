package response

type RelationMiner struct {
	Id         string `json:"id"`
	OpId       string `json:"opId"`
	GateId     string `json:"gateId"`
	RoomId     string `json:"roomId"`
	RoomName   string `json:"roomName"`
	HostName   string `json:"hostName"`
	DeviceSN   string `json:"deviceSN"`
	Ip         string `json:"ip"`
	Actor      string `json:"actor"`
	SectorSize uint64 `json:"sectorSize"`
}
