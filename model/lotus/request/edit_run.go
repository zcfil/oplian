package request

type EditRun struct {
	Id          uint64
	IsRun       bool
	GateId      string
	OpId        string
	ServiceType int //服务类型：1lotus,2miner,3任务worker,4存储worker
}

//type GateOp struct {
//	GateId string
//	OpId   string
//}
