package request

type SysRelationHostSearch struct {
	RoomId string `json:"roomId" form:"roomId"`
}

type DelSysRelationHostSearch struct {
	RoomId   string   `json:"roomId" form:"roomId"`
	HostUUID []string `json:"hostUUIDs" form:"hostUUIDs"`
}
