package request

type ProductFileReq struct {
	ProductId  int64    `json:"productId"  form:"productId"`
	FileId     []string `json:"fileId"  form:"fileId"`
	Name       string   `json:"productName" form:"productName"`
	SuitSystem string   `json:"suitSystem" form:"suitSystem"`
	Version    string   `json:"version" form:"version"`
}
