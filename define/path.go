package define

import (
	"log"
	"os"
	"path/filepath"
)

const (
	PathOplian                             = "/root/oplian/"
	PathIpfsProgram                        = "/root/oplian/bin/"
	PathIpfsScript                         = "/root/oplian/script/"
	PathIpfsConfig                         = "/root/oplian/config/"
	PathIpfsLog                            = "/root/oplian/log/"
	PathDiskInitialization                 = PathIpfsScript + "disk_initialization.sh"
	PathDiskPowerInitialization            = PathIpfsScript + "disk_power_initialization.sh"
	PathDiskInitializationArray            = PathIpfsScript + "disk_initialization_array.sh"
	PathDiskInitializationNotArray         = PathIpfsScript + "disk_initialization_not_array.sh"
	PathDiskReassemblyArray                = PathIpfsScript + "disk_array.sh"
	PathIpfsScriptHostSystemInitialization = PathIpfsScript + "host_system_initialization.sh"
	PathDiskReadRate                       = PathIpfsScript + "disk_rw_rate.sh"
	PathStorageDiskReadRate                = PathIpfsScript + "storage_disk_rate_test.sh"
	PathServerPortMonitor                  = PathIpfsScript + "server_port_monitor.sh"
	PathServerPortRequest                  = PathIpfsScript + "server_port_request.sh"
	PathServerKillIperf3                   = PathIpfsScript + "kill_net_test.sh"
	PathServerKillBench                    = PathIpfsScript + "kill_gpu_test.sh"
	PathServerKillAllTest                  = PathIpfsScript + "kill_all_test.sh"
	PathScriptExecute                      = PathIpfsScript + "kill_script_execute.sh"
	PathServerLogOvertime                  = PathIpfsScript + "log_overtime.sh"
	PathServerWdpostBalance                = PathIpfsScript + "get_wdpost_balance.sh"
	PathServerLotusHigh                    = PathIpfsScript + "lotus_high.sh"
	PathServerLogInformation               = PathIpfsScript + "log_information.sh"
	PathServerDataCatalog                  = PathIpfsScript + "data_catalog.sh"
	PathServerTimeSync                     = PathIpfsScript + "time_sync.sh"
	PathServerHostDown                     = PathIpfsScript + "host_down.sh"
	PathServerSectorType                   = PathIpfsScript + "sector_type.sh"
	PathServerSectorAddress                = PathIpfsScript + "sector_address.sh"
	PathIpfsScriptRunBenchTest             = PathIpfsScript + "run_bench_test.sh"
	PathIpfsScriptRunLotus                 = PathIpfsScript + "lotusrun.sh"
	PathIpfsScriptImportLotus              = PathIpfsScript + "import_lotus.sh"
	PathIpfsScriptRunWorker                = PathIpfsScript + "workerrun.sh"
	PathIpfsScriptRunStorage               = PathIpfsScript + "storagerun.sh"
	PathIpfsScriptRunMiner                 = PathIpfsScript + "minerrun.sh"
	PathIpfsScriptInitMiner                = PathIpfsScript + "init_miner.sh"
	PathIpfsScriptInitMinerWorker          = PathIpfsScript + "init_miner_worker.sh"
	PathIpfsScriptRunBoost                 = PathIpfsScript + "boostrun.sh"
	PathIpfsScriptRunOpC2                  = PathIpfsScript + "run_op_c2.sh"
	PathIpfsScriptRunSectorC2              = PathIpfsScript + "run_sector_c2.sh"
	PathIpfsScriptRunSectorTask            = PathIpfsScript + "run_sector_task.sh"
	PathIpfsScriptRunWorkerP2              = PathIpfsScript + "workerrunp2.sh"
	PathIpfsScriptDeadlinesProven          = PathIpfsScript + "deadlines_proven.sh"

	FileRootDir         = "/root"
	FileUploadDir       = "/root/upload/files"
	OpCsPathC1          = "/op/phasecs"
	OpCsPathC2          = "/op/proof"
	GateWayCsPathC1     = "/gateway/phasecs"
	GateWayCsPathC2     = "/gateway/proof"
	PathUidConfig       = "./config/op_c2_uuid"
	PathScriptMountDisk = PathIpfsScript + "mount_disk.sh"
	PathDiskNFSSync     = PathIpfsScript + "disk_nfs_sync.sh"
)

const (
	//DownLoadAddressOut = "http://10.0.1.251:82/"
	DownLoadAddressOut    = "http://14.215.165.39:9535/"
	IPAddressFile         = "qqwry.dat"
	OPlianFile            = "oplian"
	OPlianOpFile          = "oplian-op"
	OPlianOpC2File        = "oplian-op-c2"
	OplianSectorsSealFile = "oplian-sectors-seal"
)

var MainDisk string

var (
	PathMntMd0             = ""
	PathIpfsData           = ""
	PathIpfsLogs           = ""
	PathIpfsDataCar        = ""
	PathIpfsDataWorkerCar  = ""
	PathIpfsLotus          = ""
	PathIpfsPARAMETER      = ""
	PathIpfsLotusDatastore = ""
	PathIpfsLotusKeystore  = ""
	PathIpfsLotusToken     = ""
	PathIpfsMiner          = ""
	PathIpfsMinerToken     = ""
	PathIpfsMinerKeystore  = ""
	PathIpfsWorker         = ""
	PathIpfsWorkerToken    = ""
	PathIpfsStorage        = ""
	PathIpfsStorageToken   = ""
	PathIpfsBoost          = ""
	PathIpfsBoostToken     = ""
	PathIpfsBoostConfig    = ""
	PathSpeedUpFile        = ""
	PathProveParameters    = ""
	PathProveParents       = ""
	FileGateWayDir         = ""
)

func PathInit() {
	if MainDisk == "" {
		log.Fatal(MainDisk, " absent！")
	}
	_, err := os.Lstat(MainDisk)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(MainDisk, " absent！")
	}
	PathMntMd0 = filepath.Join(MainDisk, "")
	PathIpfsData = filepath.Join(MainDisk, "/ipfs/data")
	PathIpfsLogs = filepath.Join(MainDisk, "/ipfs/logs")
	PathIpfsDataCar = filepath.Join(MainDisk, "/ipfs/data/car")
	PathIpfsDataWorkerCar = filepath.Join(MainDisk, "/ipfs/data/worker-car")
	PathIpfsLotus = filepath.Join(MainDisk, "/ipfs/data/lotus")
	PathIpfsPARAMETER = filepath.Join(MainDisk, "/filecoin-proof-parameters")
	PathIpfsLotusDatastore = filepath.Join(MainDisk, "/ipfs/data/lotus/datastore")
	PathIpfsLotusKeystore = filepath.Join(MainDisk, "/ipfs/data/lotus/keystore")
	PathIpfsLotusToken = filepath.Join(MainDisk, "/ipfs/data/lotus/token")
	PathIpfsMiner = filepath.Join(MainDisk, "/ipfs/data/lotusminer")
	PathIpfsMinerToken = filepath.Join(MainDisk, "/ipfs/data/lotusminer/token")
	PathIpfsMinerKeystore = filepath.Join(MainDisk, "/ipfs/data/lotusminer/keystore")
	PathIpfsWorker = filepath.Join(MainDisk, "/ipfs/data/lotusworker")
	PathIpfsWorkerToken = filepath.Join(MainDisk, "/ipfs/data/lotusworker/token")
	PathIpfsStorage = filepath.Join(MainDisk, "/ipfs/data/lotusstorage")
	PathIpfsStorageToken = filepath.Join(MainDisk, "/ipfs/data/lotusstorage/token")
	PathIpfsBoost = filepath.Join(MainDisk, "/ipfs/data/boost")
	PathIpfsBoostToken = filepath.Join(MainDisk, "/ipfs/data/boost/token")
	PathIpfsBoostConfig = filepath.Join(MainDisk, "/ipfs/data/boost/config.toml")
	PathSpeedUpFile = filepath.Join(MainDisk, "/lotusworker")
	PathProveParameters = filepath.Join(MainDisk, "/filecoin-proof-parameters")
	PathProveParents = filepath.Join(MainDisk, "/filecoin-parents")
	FileGateWayDir = filepath.Join(MainDisk, "/ipfs/data/gateway")
}
