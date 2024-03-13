package request

// PageInfo Paging common input parameter structure
type LotusInfoPage struct {
	Page         int    `json:"page" form:"page"`
	PageSize     int    `json:"pageSize" form:"pageSize"`
	Keyword      string `json:"keyword" form:"keyword"`
	SyncStatus   int    `json:"syncStatus" form:"syncStatus"`
	DeployStatus int    `json:"deployStatus" form:"deployStatus"`
	GateId       string `json:"gateId" form:"gateId"`
}

type LotusInfoApi struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
}

type AddLotusInfo struct {
	Id            uint64
	GateId        string
	OpId          string
	Ip            string
	SecpCount     int32
	BlsCount      int32
	ImportMode    int32
	FileName      string
	WalletNewMode int32
	Wallets       []OpWallet
}

type OpWallet struct {
	OpId    string  `json:"op_id,omitempty"`
	Address string  `json:"address"`
	OpIp    string  `json:"op_ip,omitempty"`
	Balance float64 `json:"balance"`
}
