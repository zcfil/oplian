package response

import (
	"time"
)

type NodeInfo struct {
	Id          int64     `json:"id"`
	ColonyName  string    `json:"colonyName"`
	ColonyType  int       `json:"colonyType"`
	Wcount      int       `json:"wcount"`
	Mcount      int       `json:"mcount"`
	Scount      int       `json:"scount"`
	CreatedAt   time.Time `json:"createdAt"`
	Power       string    `json:"power"`
	TotalPower  string    `json:"totalPower"`
	SectorSize  uint64    `json:"sectorSize"`
	Wallets     []Wallet  `json:"wallets" gorm:"-"`
	WalletCount int       `json:"walletCount"`
	Live        uint64    `json:"live"`
	Active      uint64    `json:"active"`
	Faulty      uint64    `json:"faulty"`
}
