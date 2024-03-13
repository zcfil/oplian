package lotusrpc

type JsonRpc struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

func NewJsonRpc(Method string, params interface{}) JsonRpc {
	return JsonRpc{Jsonrpc: "2.0", Method: Method, Params: params, Id: 0}
}

type JsonRpcResult struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Id      int         `json:"id"`
	Error   interface{} `json:"error"`
}

type ActorControl struct {
	Owner            string
	Worker           string
	ControlAddresses []string
	PeerId           string
	SectorSize       uint64
}

type Sector struct {
	Miner  uint64 `json:"Miner"`
	Number uint64 `json:"Number"`
}

type SectorNext struct {
	Next uint64 `json:"Next"`
}

type MinerControl struct {
	Owner            string
	Worker           string
	ControlAddresses string
}

type AskInfo struct {
	Price         string `json:"Price"`
	VerifiedPrice string `json:"VerifiedPrice"`
	MinPieceSize  uint64 `json:"MinPieceSize"`
	MaxPieceSize  uint64 `json:"MaxPieceSize"`
	Miner         string `json:"Miner"`
}

type SectorID struct {
	Miner  uint64
	Number uint64
}
type SchedDiagRequestInfo struct {
	Sector   SectorID
	TaskType string
	Priority int
	SchedId  string
}

type SchedDiagInfo struct {
	Requests    []SchedDiagRequestInfo
	OpenWindows []string
}

type MinerSectors struct {
	// Live sectors that should be proven.
	Live uint64
	// Sectors actively contributing to power.
	Active uint64
	// Sectors with failed proofs.
	Faulty uint64
}

type PowerMap struct {
	MinerPower map[string]string
	TotalPower map[string]string
}

type Power struct {
	MinerPower string
	TotalPower string
}
