package define

const (
	RunStatusStop RunStatus = iota
	RunStatusRunning
)

type RunStatus int

func RunStatusBool(run bool) RunStatus {
	if run {
		return RunStatusRunning
	}
	return RunStatusStop
}

func (r RunStatus) Int32() int32 {
	return int32(r)
}

func (r RunStatus) Int() int {
	return int(r)
}
