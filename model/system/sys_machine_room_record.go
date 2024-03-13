// 自动生成模板SysMachineRoomRecords
package system

import (
	"oplian/global"
)

// 如果含有time.Time 请自行import time包
type SysMachineRoomRecord struct {
	global.ZC_MODEL
	RoomId          string `json:"roomId" form:"roomId" gorm:"column:room_id;comment:机房编号"`
	RoomName        string `json:"roomName" form:"roomName" gorm:"column:room_name;comment:机房名称"`
	RoomType        int    `json:"roomType" form:"roomType" gorm:"column:room_type;comment:机房类型"`
	CabinetsNum     int    `json:"cabinetsNum" form:"cabinetsNum" gorm:"column:cabinets_num;comment:机柜数量"`
	RoomTemp        int    `json:"roomTemp" form:"roomTemp" gorm:"column:room_temp;comment:机房温度"`
	RoomLeader      string `json:"roomLeader" form:"roomLeader" gorm:"column:room_leader;comment:机房负责人(uuid)"`
	RoomLeaderPhone string `json:"roomLeaderPhone" form:"roomLeaderPhone" gorm:"column:room_leader_phone;comment:机房负责人电话"`
	RoomSupplier    string `json:"roomSupplier" form:"roomSupplier" gorm:"column:room_supplier;comment:机房供应商"`
	SupplierContact string `json:"supplierContact" form:"supplierContact" gorm:"column:supplier_contact;comment:供应商联系人"`
	SupplierPhone   string `json:"supplierPhone" form:"supplierPhone" gorm:"column:supplier_phone;comment:供应商联系电话"`
	RoomAdmin       string `json:"roomAdmin" form:"roomAdmin" gorm:"column:room_admin;comment:机房管理员(uuid)"`
	RoomOwner       string `json:"roomOwner" form:"roomOwner" gorm:"column:room_owner;comment:机房owner(uuid)"`
	PhysicalAddress string `json:"physicalAddress" form:"physicalAddress" gorm:"column:physical_address;comment:物理地址"`
	RoomArea        int    `json:"roomArea" form:"roomArea" gorm:"column:room_area;comment:机房面积"`
	GatewayId       string `json:"gatewayId" form:"gatewayId" gorm:"column:gateway_id;comment:oplian-gateway对应的网关ID"`
	IntranetIP      string `json:"intranetIp" form:"intranetIp" gorm:"column:intranet_ip;comment:内网IP"`
}

func (SysMachineRoomRecord) TableName() string {
	return "sys_machine_room_records"
}
