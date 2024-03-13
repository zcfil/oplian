package request

type GetMinerID struct {
	Actor string `json:"actor" form:"actor"`
}

type IDActor struct {
	ID    uint   `json:"id" form:"id"`
	Actor string `json:"actor" form:"actor"`
}
