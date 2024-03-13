package gateway

import "oplian/global"

type DownloadRecord struct {
	global.ZC_MODEL
	Url        string `json:"url" gorm:"comment:下载地址"`
	FilePath   string `json:"file_path"  gorm:"comment:文件路径"`
	FileSize   uint64 `json:"file_size"  gorm:"comment:文件大小"`
	LoadSize   uint64 `json:"load_size"  gorm:"comment:已经下载"`
	FileStatus int    `json:"file_status" gorm:"comment:文件状态 0进行中，1已完成，2失败，3已移除"`
	ErrorMsg   string `json:"ErrorMsg" gorm:"comment:下载错误信息"`
}

func (DownloadRecord) TableName() string {
	return "download_record"
}
