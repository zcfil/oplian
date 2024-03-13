package lotus

type Duty int

const (
	WnPost   Duty = 0 << 1
	WdPost   Duty = 1 << 1
	Dispatch Duty = 2 << 1
)

type EnvConfigType int

const (
	Default EnvConfigType = iota
	WorkerC2
	Worker
	Miner
	Lotus
)
