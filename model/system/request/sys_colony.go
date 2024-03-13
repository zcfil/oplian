package request

type ColonyPageInfo struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
	Keyword    string `json:"keyword" form:"keyword"`
	ColonyType int    `json:"colonyType" form:"colonyType" `
}
