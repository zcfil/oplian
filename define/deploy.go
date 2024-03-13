package define

// 部署状态
type DeployStatus int

const (
	DeployRunning DeployStatus = 1 + iota
	DeployFinish
	DeployFail
	DeployReset
)

func (d DeployStatus) Int32() int32 {
	return int32(d)
}
func (d DeployStatus) Int() int {
	return int(d)
}

const (
	MinerDepolyFile int32 = 1 + iota
	MinerDepolyNew
	MinerDepolyWorker
)
