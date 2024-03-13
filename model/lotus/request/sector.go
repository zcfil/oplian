package request

type AutoSector struct {
	Enabled bool
	Balance int
	Id      uint64
}

type SectorParam struct {
	OpId        string
	GateId      string
	SectorType  int //CC：1，DC：2
	Actor       string
	Count       int
	StartNumber uint64
}

type SectorPage struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
	Actor      string `json:"actor" form:"actor"`
	Number     uint64 `json:"sectorId" form:"sectorId"`
	SectorType int    `json:"sectorType" form:"sectorType"`
	Status     string `json:"sectorStatus" form:"sectorStatus"`
}
