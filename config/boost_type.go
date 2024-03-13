package config

type StorageConfig struct {
	// The maximum number of concurrent fetch operations to the storage subsystem
	ParallelFetchLimit int
	// How frequently Boost should refresh the state of sectors with Lotus. (default: 1hour)
	// When run, Boost will trigger a storage redeclare on the miner in addition to a storage list.
	// This ensures that index metadata for sectors reflects their status (removed, unsealed, etc).
	StorageListRefreshDuration string
	// Whether or not Boost should have lotus redeclare its storage list (default: true).
	// Disable this if you wish to manually handle the refresh. If manually managing the redeclare
	// and it is not triggered, retrieval quality for users will be impacted.
	RedeclareOnStorageListRefresh bool
}
type WalletsConfig struct {
	// The "owner" address of the miner
	Miner string
	// The wallet used to send PublishStorageDeals messages.
	// Must be a control or worker address of the miner.
	PublishStorageDeals string
	// The wallet used as the source for storage deal collateral
	DealCollateral string
	// Deprecated: Renamed to DealCollateral
	PledgeCollateral string
}
type GraphqlConfig struct {
	// The port that the graphql server listens on
	Port uint64
}
type MonitoringConfig struct {
	// The number of epochs after which alert is generated for a local pending
	// message in lotus mpool
	MpoolAlertEpochs int64
}
type TracingConfig struct {
	Enabled     bool
	ServiceName string
	Endpoint    string
}
type LocalIndexDirectoryYugabyteConfig struct {
	Enabled bool
	// The yugabyte postgres connect string eg "postgresql://postgres:postgres@localhost"
	ConnectString string
	// The yugabyte cassandra hosts eg ["127.0.0.1"]
	Hosts []string
}
type LocalIndexDirectoryLeveldbConfig struct {
	Enabled bool
}
type LocalIndexDirectoryConfig struct {
	Yugabyte LocalIndexDirectoryYugabyteConfig
	Leveldb  LocalIndexDirectoryLeveldbConfig
	// The maximum number of add index operations allowed to execute in parallel.
	// The add index operation is executed when a new deal is created - it fetches
	// the piece from the sealing subsystem, creates an index of where each block
	// is in the piece, and adds the index to the local index directory.
	ParallelAddIndexLimit int
	// The port that the embedded local index directory data service runs on.
	// Set this value to zero to disable the embedded local index directory data service
	// (in that case the local index directory data service must be running externally)
	EmbeddedServicePort uint64
	// The connect string for the local index directory data service RPC API eg "ws://localhost:8042"
	// Set this value to "" if the local index directory data service is embedded.
	ServiceApiInfo string
	// The RPC timeout when making requests to the boostd-data service
	ServiceRPCTimeout Duration
}
type ContractDealsConfig struct {
	// Whether to enable chain monitoring in order to accept contract deals
	Enabled bool

	// Allowlist for contracts that this SP should accept deals from
	AllowlistContracts []string

	// From address for eth_ state call
	From string
}
type HttpDownloadConfig struct {
	// NChunks is a number of chunks to split HTTP downloads into. Each chunk is downloaded in the goroutine of its own
	// which improves the overall download speed. NChunks is always equal to 1 for libp2p transport because libp2p server
	// doesn't support range requests yet. NChunks must be greater than 0 and less than 16, with the default of 5.
	NChunks int
	// AllowPrivateIPs defines whether boost should allow HTTP downloads from private IPs as per https://en.wikipedia.org/wiki/Private_network.
	// The default is false.
	AllowPrivateIPs bool
}
type Duration string
type ResourceFilteringStrategy string
