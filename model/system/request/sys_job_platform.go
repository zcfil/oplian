package request

import (
	"oplian/model/system"
)

type JobReq struct {
	ID             int                     `json:"id"`
	TaskName       string                  `json:"task_name"`
	TaskType       int                     `json:"task_type"`
	GateWayId      string                  `json:"gate_way_id"`
	SysOpRelations []system.SysOpRelations `json:"op_relations"`
	ScriptContent  string                  `json:"script_content"`
	ScriptParam    string                  `json:"script_param"`
	TimeLength     int                     `json:"time_length"`
	UserName       string                  `json:"user_name"`
}

type DistributeReq struct {
	Id             int                          `json:"id"`
	TaskName       string                       `json:"task_name"`
	TaskType       int                          `json:"task_type"`
	SysOpRelations []system.SysOpRelations      `json:"op_relations"`
	GateWayId      string                       `json:"gate_way_id"`
	TimeLength     int                          `json:"time_length"`
	Enable         int                          `json:"enable"`
	LimitSpeed     int                          `json:"limit_speed"`
	FileType       int                          `json:"file_type"`
	FilePath       string                       `json:"file_path"`
	SendType       int                          `json:"send_type"`
	OpList         []system.SysFileManageUpload `json:"op_list"`
	UpLoadFile     []UpLoadFile                 `json:"up_load_file"`
	UserName       string                       `json:"user_name"`
	ScriptName     string                       `json:"script_name"`
}

type File struct {
	Path       string       `json:"path"`
	UpLoadFile []UpLoadFile `json:"up_load_file"`
	Enable     int          `json:"enable"`
	LimitSpeed int          `json:"limit_speed"`
}

type UpLoadFile struct {
	FileType int    `json:"file_type"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

type CommonParam struct {
	ID      int    `json:"id"`
	TaskNum string `json:"task_num"`
}

type ExecuteRecordsReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Keyword  string `json:"keyword" form:"keyword"`
	TaskType string `json:"task_type"`
	Status   string `json:"status"`
}

type FileManageReq struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
	Id         int    `json:"id"`
	GateWayId  string `json:"gate_way_id"`
	KeyWord    string `json:"keyword"`
	FileName   string `json:"file_name"`
	FileStatus int    `json:"file_status"`
	FileType   int    `json:"file_type"`
	FileUrl    string `json:"file_url"`
}

type AddFileReq struct {
	FileType     int64        `json:"file_type"`
	AddType      int64        `json:"add_type"`
	GateWayId    string       `json:"gate_way_id"`
	OpId         string       `json:"op_id"`
	Ip           string       `json:"ip"`
	Port         string       `json:"port"`
	FileUrl      string       `json:"file_url"`
	UpLoadFile   []UpLoadFile `json:"up_load_file"`
	OpList       []OpListReq  `json:"op_list"`
	DownloadPath string       `json:"download_path"`
	ZipFileName  string       `json:"zip_file_name"`
	RoomId       string       `json:"room_id"`
	RoomName     string       `json:"room_name"`
}

type SynFileReq struct {
	FromPath   string       `json:"from_path"`
	ToPath     string       `json:"to_path"`
	GateWayId  string       `json:"gate_way_id"`
	OpId       string       `json:"op_id"`
	ToOpId     string       `json:"to_op_id"`
	Ip         string       `json:"ip"`
	Port       string       `json:"port"`
	Files      []UpLoadFile `json:"files"`
	FileName   string       `json:"file_name"`
	LimitSpeed int          `json:"limit_speed"`
	TimeOut    string       `json:"time_out"`
	TimeLength int64        `json:"time_length"`
	FileType   int          `json:"file_type"`
}

type FileTypeReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	FileType int    `json:"fileType"`
	GateId   string `json:"gateId"`
}

type DownLoadReq struct {
	GateWayId    string       `json:"gate_way_id"`
	Files        []UpLoadFile `json:"files"`
	OpList       []OpListReq  `json:"op_list"`
	DownloadPath string       `json:"download_path"`
	ObZuFileName string       `json:"ob_zu_file_name"`
	ObZuDek      int64        `json:"ob_zu_dek"`
}

type OpListReq struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	OpId     string `json:"op_id"`
	Path     string `json:"path"`
	FileName string `json:"file_name"`
}

type GateWayReq struct {
	Id string `json:"id"`
}

type OpFileSync struct {
	GateWayId string      `json:"gate_way_id"`
	OpListReq []OpListReq `json:"op_local"`
}
