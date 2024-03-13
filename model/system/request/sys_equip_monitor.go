package request

type HostMonitorReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Keyword  string `json:"keyword" form:"keyword"`
	GateId   string `json:"gateId" form:"gateId"`
}

type HostScriptReq struct {
	UUID         string `json:"opId" form:"opId"`
	GatewayId    string `json:"gateId" form:"gateId"`
	HostClassify int64  `json:"hostClassify" form:"hostClassify"`
	ScriptInfo   string `json:"scriptInfo" form:"scriptInfo"`
}

type DiskReMountReq struct {
	UUID         string `json:"opId" form:"opId"`
	GatewayId    string `json:"gateId" form:"gateId"`
	HostClassify int64  `json:"hostClassify" form:"hostClassify"`
	Actor        string `json:"actor" form:"actor"`
	NodeIP       string `json:"nodeIP" form:"nodeIP"`
	MountOpId    string `json:"mountOpId" form:"mountOpId"`
}

type DiskLetterReq struct {
	UUID       string `json:"opId" form:"opId"`
	GatewayId  string `json:"gateId" form:"gateId"`
	DiskLetter string `json:"diskLetter" form:"diskLetter"`
}

type GetNodeStorageInfoReq struct {
	UUID      string `json:"opId" form:"opId"`
	GatewayId string `json:"gateId" form:"gateId"`
	Actor     string `json:"actor" form:"actor"`
}

type GetNodeLogInfoReq struct {
	UUID         string `json:"opId" form:"opId"`
	GatewayId    string `json:"gateId" form:"gateId"`
	HostClassify int64  `json:"hostClassify" form:"hostClassify"`
	LogType      string `json:"logType" form:"logType"`
}

type GetNodeLogInfoInfoReq struct {
	UUID         string `json:"opId" form:"opId"`                 
	GatewayId    string `json:"gateId" form:"gateId"`             
	HostClassify int64  `json:"hostClassify" form:"hostClassify"` 
	LogType      string `json:"logType" form:"logType"`           
	BeginNum     int64  `json:"beginNum" form:"beginNum"`         
}
