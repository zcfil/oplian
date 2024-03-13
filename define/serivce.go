package define

type ServiceType int

const (
	ServiceLotus ServiceType = 1 + iota
	ServiceMiner
	ServiceWorkerTask
	ServiceWorkerStorage
)

func (s ServiceType) Int32() int32 {
	return int32(s)
}
