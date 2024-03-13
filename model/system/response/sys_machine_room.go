package response

import (
	"time"
)

type SysMachineRoomRecord struct {
	CreatedAt       time.Time    `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt       time.Time    `json:"UpdatedAt" gorm:"column:updated_at"`
	RoomId          string       `json:"roomId" gorm:"column:room_id"`
	RoomName        string       `json:"roomName" gorm:"column:room_name"`
	RoomType        int          `json:"roomType" gorm:"column:room_type"`
	RoomTemp        int          `json:"roomTemp" gorm:"column:room_temp"`
	RoomLeader      string       `json:"roomLeader" gorm:"column:room_leader"`
	RoomLeaderName  string       `json:"roomLeaderName" gorm:"-"`
	RoomLeaderPhone string       `json:"roomLeaderPhone" gorm:"column:room_leader_phone"`
	RoomSupplier    string       `json:"roomSupplier" gorm:"column:room_supplier"`
	SupplierContact string       `json:"supplierContact" gorm:"column:supplier_contact"`
	SupplierPhone   string       `json:"supplierPhone" gorm:"column:supplier_phone"`
	RoomAdmin       string       `json:"roomAdmin" gorm:"column:room_admin"`
	RoomAdminName   string       `json:"roomAdminName" gorm:"-"`
	RoomOwner       string       `json:"roomOwner" gorm:"column:room_owner"`
	RoomOwnerName   string       `json:"roomOwnerName" gorm:"-"`
	PhysicalAddress string       `json:"physicalAddress" gorm:"column:physical_address"`
	RoomArea        int          `json:"roomArea" gorm:"column:room_area"`
	GatewayId       string       `json:"gatewayId" gorm:"column:gateway_id"`
	BindNodeHostNum int64        `json:"bindNodeHostNum" gorm:"-"`
	RoomBindHost    RoomBindHost `json:"roomBindHost" gorm:"-"`
}

type RoomBindHost struct {
	BindHostNum   int64    `json:"bindHostNum" gorm:"-"`
	BindHostUUIDs []string `json:"bindHostUUIDs" gorm:"-"`
}

type RoomRecordList struct {
	RoomId    string `json:"roomId" gorm:"column:room_id"`
	RoomName  string `json:"roomName" gorm:"column:room_name"`
	GatewayId string `json:"gatewayId" gorm:"column:gateway_id"`
}
