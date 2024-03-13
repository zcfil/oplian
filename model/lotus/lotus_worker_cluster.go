package lotus

import "oplian/global"

type LotusWorkerCluster struct {
	global.ZC_MODEL
	OpId       string `gorm:"index;op_id;comment:主机OpId" json:"opId"`
	Ip         string `gorm:"ip;comment:IP" json:"ip"`
	ServerName string `gorm:"server_name;comment:主机名称" json:"serverName"`
	AssetsNum  string `gorm:"assets_num;comment:资产编号" json:"assetsNum"`
	DeviceSn   string `gorm:"device_sn;comment:设备SN" json:"deviceSn"`
	RoomNum    string `gorm:"index;room_num;comment:机房编号" json:"roomNum"`
	RoomName   string `gorm:"room_name;comment:机房名称" json:"roomName"`
	SectorSize string `gorm:"sector_size;comment:扇区大小" json:"sectorSize"`
	NodeNum    int    `gorm:"node_num;comment:节点机数量" json:"nodeNum"`
	GpuNum     int    `gorm:"gpu_num;comment:显卡数量" json:"gpuNum"`
	TaskNum    int    `gorm:"task_num;comment:c2任务数量(已完成)" json:"taskNum"`
	Remark     string `gorm:"remark;comment:备注" json:"remark"`
}

func (LotusWorkerCluster) TableName() string {
	return "lotus_worker_cluster"
}
