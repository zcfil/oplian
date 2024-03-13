package response

import "time"

type UnsealedInfo struct {
	Id           uint      `json:"id"`           //主键ID
	OpId         string    `json:"opId"`         // 设备ID
	GateId       string    `json:"gateId"`       //机房ID
	RoomId       string    `json:"roomId"`       //机房ID
	RoomName     string    `json:"roomName"`     //机房名称
	HostName     string    `json:"hostName"`     //主机名称
	DeviceSN     string    `json:"deviceSN"`     //主机编号
	Ip           string    `json:"ip"`           // IP
	DeployStatus int       `json:"deployStatus"` // 1部署中，2部署成功，3部署失败
	StartAt      time.Time `json:"startAt"`      //开始时间
	FinishAt     time.Time `json:"finishAt"`     //完成时间
	ErrMsg       string    `json:"errMsg"`       // "部署错误信息"
}

type UnsealedList struct {
	OpId         string `json:"opId"`
	GateId       string `json:"gateId"`
	Ip           string `json:"ip"`
	DeployStatus int    `json:"deployStatus"`
	NFSDisk      string `json:"nfsDisk"`
}
