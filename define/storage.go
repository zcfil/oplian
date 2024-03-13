package define

const (
	StorageConfig     = "storage.json"
	SectorStoreConfig = "sectorstore.json"
)

const (
	StorageTypeNFS = iota + 1
	StorageTypeWorker
)
