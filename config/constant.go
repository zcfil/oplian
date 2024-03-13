package config

import "time"

const (
	AutomaticTrigger = 1
	ManualTrigger    = 2

	HostUnderTest       = 1
	HostCompliance      = 2
	HostNotUpToStandard = 3
	HostTestFailed      = 4
	HostInClose         = 5

	HostMinerTest    = 1
	HostWorkerTest   = 2
	HostStorageTest  = 3
	HostC2WorkerTest = 4

	HostParamTestPassStr            = "1"
	HostParamTestDissatisfactionStr = "2"
	HostParamTestFailedStr          = "3"

	HostParamTestPass            = 1
	HostParamTestDissatisfaction = 2
	HostParamTestFailed          = 3

	HostMinerType     = 1
	HostWorkerType    = 2
	HostStorageType   = 3
	HostDCStorageType = 4
	HostC2WorkerType  = 5
	HostLotusType     = 6
	//	HostUnsealedType  = 10

	TimeFormat = "2006-01-02 15:04:05"

	NetTestHostNum         = 1
	NetTestHostNumAddPower = 4

	HostDiskIOPass = "0"

	HostPatrolStatusSuccess = true
	HostPatrolStatusFailed  = false

	HostC2WorkerDiskSize = "300G"

	HostNetTimeout         = 10 * time.Minute
	HostTestTimeout        = 10 * time.Minute
	HostGPUTestTimeout     = 1*time.Hour + 30*time.Minute
	GoFunctionWaitTime     = 10 * time.Second
	WorkerMountNFSWaitTime = 2 * time.Minute
)

const (
	NetworkSpeedAdd      = "15G"
	NetworkSpeed         = "8G"
	NetworkC2WorkerSpeed = "10M"

	StorageDiskOverallSpeed = "100M"
	StorageNFSDiskSpeed     = "50M"
	StorageDiskSize         = "10T"

	WorkerSSDDiskSpeed = "1G"
	WorkerSSDDiskSize  = "1T"

	NodeSSDDiskSpeed = "500M"
	NodeSSDDiskSize  = "3T"

	StorageRamSize        = 32
	WorkerRamSize         = 128
	NodeRamSize           = 250
	NodeRamSizeAdditional = 500
	C2WorkerRamSize       = 240

	GPUModelStandard         = 2080
	GPURunTime       float64 = 60

	CPUThreads = 4

	TestDefaultFile = "md0"

	TestHostNetTimeAdd = 300
	TestHostNetTime    = 10
)

var (
	HostNetTestPort = []int{5201, 5202, 5203, 5204}
)

const (
	PatrolTestTrue  = "true"
	PatrolTestFalse = "false"

	PatrolWdpostBalance float32 = 10

	PatrolLotusVersion  = "1.20.0+mainnet+zc-v1.4.0"
	PatrolMinerVersion  = "1.20.0+mainnet+zc-v1.5.1"
	PatrolBoostdVersion = "1.6.0+git.8609359.dirty"

	PatrolMinerTimeConfig   = 3 * 3600
	PatrolWorkerTimeConfig  = 3 * 3600
	PatrolStorageTimeConfig = 3 * 3600

	BadSectorTimeConfig = 12 * 3600

	SectorTypeCCType   = "CC"
	SectorTypeDCType   = "DC"
	SectorTypeElseType = "ELSE"

	SectorTypeCCTypeNum   = 1
	SectorTypeDCTypeNum   = 2
	SectorTypeElseTypeNum = 3
)

const (
	DownloadTarFileName = "oplian-require-file.tar"
	DownloadTarFileHash = "oplian-require-file.yaml"
)

const (
	CarDirName  = "carDir"
	SealDirName = "sealDir"
	CarDir      = "ll * /mnt/md0/ipfs/data/worker-car"
	SealDir     = "ll * /mnt/md0/workerseal"
)

var AllowScriptList = []string{"df -h", "lsblk", CarDirName, SealDirName}

const (
	WorkerLogType        = "worker"
	MinerProveLogType    = "prove"
	MinerWdpostLogType   = "wdpost"
	MinerWnpostLogType   = "wnpost"
	MinerRealTimeLogType = "realTime"
	LotusLogType         = "lotus"
	BoostLogType         = "boost"
	MinerLogType         = "miner"

	DefaultGetNum = 30

	LogTimeDuration = 30

	WorkerLogDir = "/ipfs/logs/worker.log"
	LotusLogDir  = "/ipfs/logs/lotus.log"
	MinerLogDir  = "/ipfs/logs/miner.log"
	BoostLogDir  = "/ipfs/logs/boost.log"

	LotusRestartCmd  = "supervisorctl restart lotus"
	MinerRestartCmd  = "supervisorctl restart lotus-miner"
	BoostRestartCmd  = "supervisorctl restart boost"
	WorkerRestartCmd = "supervisorctl restart lotus-worker"

	LotusHeightNum = 3

	FilMinerUrl = "https://api.filutils.com/api/v2/miner/"
)
