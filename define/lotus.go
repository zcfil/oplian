package define

import (
	"sync"
	"time"
)

var Lh = LotusHeart{
	LotusMap:  make(map[string]time.Time),
	MinerMap:  make(map[string]time.Time),
	WorkerMap: make(map[string]time.Time),
}

type LotusHeart struct {
	LotusMap   map[string]time.Time
	MinerMap   map[string]time.Time
	WorkerMap  map[string]time.Time
	LotusLock  sync.RWMutex
	MinerLock  sync.RWMutex
	WorkerLock sync.RWMutex
}

type LotusInitModel int

const (
	LotusInitCopyModel LotusInitModel = iota
	LotusInitImportModel
	LotusInitRunModel
)

type SyncStatus int

const (
	Synchronizing SyncStatus = 1 + iota
	SyncFinish
	SyncFail
)

func (s SyncStatus) Int32() int32 {
	return int32(s)
}
func (s SyncStatus) Int() int {
	return int(s)
}

type WorkerType int

const (
	TaskWorker WorkerType = iota
	StorageWorker
	MinerWorker
)

func (s WorkerType) Int32() int32 {
	return int32(s)
}
func (s WorkerType) Int() int {
	return int(s)
}

type NetworkType int

const (
	NetworkLanMap NetworkType = iota
	NetworkPub
)

func (s NetworkType) Int32() int32 {
	return int32(s)
}
func (s NetworkType) Int() int {
	return int(s)
}

type WalletNewMode int

const (
	WalletNew WalletNewMode = iota
	WalletCopy
)

const (
	LocalHost           = "127.0.0.1"
	BoostPort           = "50001"
	LotusPort           = "50002"
	LotusPortOpen       = "50012"
	MinerPort           = "50003"
	WorkerPort          = "50004"
	StoragePort         = "50014"
	OplianPort          = "50005"
	GatewayPort         = "50006"
	OpPort              = "50007"
	GatewayDownLoadPort = "50009"
	OpC2Port            = "50020"
	ChainsysPort        = "50030"
	SlotUnsealedPort    = "50050"
	OpSectorC2Port      = "50060"
	OpSectorTaskPort    = "50061"
	OpP2Port1           = "50058"
	OpP2Port2           = "50059"

	ApiRouter      = "/rpc/v0"
	MinerName      = "lotusminer"
	WorkerName     = "lotusworker"
	JwtHmacSecret  = "MF2XI2BNNJ3XILLQOJUXMYLUMU"
	Libp2pHost     = "NRUWE4BSOAWWQ33TOQ"
	ConfigName     = "config.toml"
	StoragePath    = "lotusminer"
	NotStoragePath = "nostorage"
	WorkerSeal     = "workerseal"
	MsgPrefix      = "bafy2bza"
)

const (
	WalletTypeBls       = "bls"
	WalletTypeSecp256k1 = "secp256k1"
	WalletPrefix        = "wallet-"
)

// Filecoin RPC
const (
	FilecoinID = "Filecoin.ID"
	//lotus
	FilecoinStateVerifiedClientStatus = "Filecoin.StateVerifiedClientStatus"
	FilecoinClientQueryAsk            = "Filecoin.ClientQueryAsk"
	FilecoinWalletBalance             = "Filecoin.WalletBalance"
	FilecoinWalletNew                 = "Filecoin.WalletNew"
	FilecoinStateMinerInfo            = "Filecoin.StateMinerInfo"
	FilecoinStateMinerSectorCount     = "Filecoin.StateMinerSectorCount"
	FilecoinStateMinerPower           = "Filecoin.StateMinerPower"
	FilecoinStateAccountKey           = "Filecoin.StateAccountKey"
	FilecoinChainHead                 = "Filecoin.ChainHead"
	//miner
	FilecoinSynC2Result           = "Filecoin.SynC2Result"
	FilecoinPledgeSector          = "Filecoin.PledgeSector"
	FilecoinSectorNumAssignerMeta = "Filecoin.SectorNumAssignerMeta"
	FilecoinDealsImportData       = "Filecoin.DealsImportData"
	FilecoinSealingSchedDiag      = "Filecoin.SealingSchedDiag"
	FilecoinSectorStorage         = "Filecoin.SectorStorage"
	FilecoinSectorStoreFile       = "Filecoin.SectorStoreFile"
	FilecoinSealingAbort          = "Filecoin.SealingAbort"
	FilecoinSectorRemove          = "Filecoin.SectorRemove"
	//worker
	FilecoinStorageAddLocal = "Filecoin.StorageAddLocal"
)
