package response

import (
	"time"
)

type LotusInfo struct {
	Id           int64     `json:"id"`
	OpId         string    `json:"opId"`
	GateId       string    `json:"gateId"`
	RoomId       string    `json:"roomId"`
	RoomName     string    `json:"roomName"`
	HostName     string    `json:"hostName"`
	DeviceSN     string    `json:"deviceSN"`
	Actor        string    `json:"actor"`
	Ip           string    `json:"ip"`
	Port         string    `json:"port"`
	WalletCount  int       `json:"walletCount"`
	MinerCount   int       `json:"minerCount"`
	DeployStatus int       `json:"deployStatus"`
	SyncStatus   int       `json:"syncStatus"`
	Online       bool      `json:"online" gorm:"-"`
	Token        string    `json:"-"`
	StartAt      time.Time `json:"startAt"`
	FinishAt     time.Time `json:"finishAt"`
	SnapshotAt   time.Time `json:"snapshotAt"`
	ErrMsg       string    `json:"errMsg"`
}

type RelationLotusInfo struct {
	Id       int64    `json:"id"`
	OpId     string   `json:"opId"`
	GateId   string   `json:"gateId"`
	RoomId   string   `json:"roomId"`
	RoomName string   `json:"roomName"`
	HostName string   `json:"hostName"`
	Actor    string   `json:"actor"`
	Ip       string   `json:"ip"`
	Token    string   `json:"-"`
	Wallets  []Wallet `json:"wallets" gorm:"-"`
}
