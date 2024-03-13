package response

import (
	"oplian/global"
	"time"
)

type SectorQueueDetail struct {
	global.ZC_MODEL
	SectorId       uint64    `json:"sectorId" `
	Actor          string    `json:"actor" `
	JobStatus      int       `json:"jobStatus"`
	SectorType     int       `json:"sectorType"`
	SectorSize     uint64    `json:"sectorSize"`
	TaskCreateAt   time.Time `json:"taskCreateAt"`
	TaskName       string    `json:"taskName"`
	RunIndex       int       `json:"runIndex"`
	JobTotal       int       `json:"jobTotal"`
	DealUuid       string    `json:"dealUuid"`
	PieceCid       string    `json:"pieceCid"`
	ExpirationTime time.Time `json:"expirationTime"`
	CarPath        string    `json:"carPath"`
}
