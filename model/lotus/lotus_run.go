package lotus

type LotusRun struct {
	OpId    string `json:"op_id" `
	GateId  string `json:"gate_id"`
	Token   string `json:"token"`
	Ip      string `json:"ip"`
	Port    string `json:"port"`
	CountF1 int    `json:"countF1"`
	CountF3 int    `json:"countF3"`
}
