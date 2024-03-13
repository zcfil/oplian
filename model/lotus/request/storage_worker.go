package request

//type BatchWorker struct {
//	GateId  string
//	Id      []uint64
//	OpId    []string
//	Ip      []string
//	MinerId uint64
//}

type BatchStorage struct {
	GateId      string
	OpId        []string
	Ip          []string
	StorageType int
	Node        string
}
