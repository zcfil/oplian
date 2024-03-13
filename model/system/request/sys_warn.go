package request

type WarnLog struct {
	WarnId         string `json:"warn_id"`
	WarnName       string `json:"warn_name"`
	WarnType       int    `json:"warn_type"`
	Ip             string `json:"ip"`
	Sn             string `json:"sn"`
	WarnInfo       string `json:"warn_info"`
	NotifyPerson   string `json:"notify_person"`
	ComputerType   string `json:"computer_type"`
	ComputerRoomId string `json:"computer_room_id"`
}
