package op

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"oplian/config"
	"oplian/define"
	"oplian/global"
	"oplian/lotusrpc"
	"oplian/service"
	"oplian/service/lotus/deploy"
	"oplian/service/pb"
	"oplian/service/system"
	"oplian/utils"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type OpServiceImpl struct{}

// Heartbeat Gateway heartbeat check call
func (p *OpServiceImpl) Heartbeat(ctx context.Context, args *pb.String) (*pb.String, error) {
	reply := &pb.String{Value: "op Heartbeat check succeeded !" + args.GetValue()}
	return reply, nil
}

// GetSystemInfo Obtain and record system information
func (p *OpServiceImpl) GetSystemInfo(ctx context.Context, args *pb.String) (*pb.String, error) {
	var s utils.OpServer
	var err error
	s.Os = utils.InitOS()
	s.GetCPUInfo()
	s.GetGPUInfo()
	s.GetLinuxSystemInfo()

	if s.Ram, err = utils.InitRAM(); err != nil {
		global.ZC_LOG.Error("func utils.InitRAM() Failed", zap.String("err", err.Error()))
		return nil, err
	}
	if s.Disk, err = utils.InitDisk(); err != nil {
		global.ZC_LOG.Error("func utils.InitDisk() Failed", zap.String("err", err.Error()))
		return nil, err
	}

	hostInfo := &pb.HostInfo{
		IntranetIP:       s.SystemInfo.IntranetIP,
		InternetIP:       s.SystemInfo.InternetIP,
		UUID:             global.OpUUID.String(),
		DeviceSN:         s.SystemInfo.DeviceSN,
		HostManufacturer: s.SystemInfo.HostManufacturer,
		HostModel:        s.SystemInfo.HostModel,
		OperatingSystem:  s.SystemInfo.OperatingSystem,
		CPUCoreNum:       int64(s.CPU.Threads),
		CPUModel:         s.CPU.Model,
		MemorySize:       int64(s.Ram.TotalMB / 1024),
		DiskNum:          int64(s.SystemInfo.DiskNum),
		DiskSize:         float32(s.SystemInfo.DiskSizeSum),
		ServerDNS:        s.SystemInfo.ServerDNS,
		SubnetMask:       s.SystemInfo.SubnetMask,
		Gateway:          s.SystemInfo.Gateway,
		SystemVersion:    s.SystemInfo.SystemVersion,
		SystemBits:       int64(s.SystemInfo.SystemBits),
	}

	_, err = global.OpToGatewayClient.AddHostRecord(ctx, hostInfo)
	if err != nil {
		return &pb.String{Value: "Host information recording failed!"}, err
	}

	reply := &pb.String{Value: "Host information recorded successfully!"}
	return reply, nil
}

// StrategyProcess Strategy processing
func (p *OpServiceImpl) StrategyProcess(ctx context.Context, args *pb.StrategyInfo) (*pb.ResponseMsg, error) {
	log.Println("op StrategyProcess succeeded!")

	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"
	_, err := system.WarnManageServiceApp.StrategyProcessType(args.OpId, args.RoomId, args.StrategiesId)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}
	return res, nil
}

// ExecuteScript Execute script
func (p *OpServiceImpl) ExecuteScript(ctx context.Context, args *pb.ScriptInfo) (*pb.String, error) {

	str, err := system.JobPlatformServiceApp.ExecuteResult(args)
	if err != nil {
		return &pb.String{}, err
	}

	return &pb.String{Value: str}, nil
}

// FileDistribution File Distribution
func (p *OpServiceImpl) FileDistribution(ctx context.Context, args *pb.FileInfo) (*pb.String, error) {

	res := &pb.String{}
	dirAr := strings.Split(define.MainDisk, "/")
	dir := fmt.Sprintf("/%s/", dirAr[1])
	if !strings.Contains(args.Path, dir) {
		args.Path = define.MainDisk + args.Path
	}
	_, err := system.JobPlatformServiceApp.CreateFile(args)
	if err != nil {
		return res, err
	}

	return res, nil
}

// RunAndStopService Starting service
func (g *OpServiceImpl) RunAndStopService(ctx context.Context, args *pb.RunStopType) (*pb.ResponseMsg, error) {
	return service.ServiceGroupApp.OpServiceGroup.StopService(args)
}

//// Stop service
//func (p *OpServiceImpl) RunService(ctx context.Context, args *pb.LotusRun) (*pb.ResponseMsg, error) {
//	return service.ServiceGroupApp.LotusServiceGroup.R(args)
//}

func (g *OpServiceImpl) UpdateWorker(ctx context.Context, args *pb.ConnectInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	return res, nil
}

// FileOpSynGateWay OP synchronizes files to GateWay
func (p *OpServiceImpl) FileOpSynGateWay(ctx context.Context, args *pb.AddFileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"

	go func() {

		err := system.JobPlatformServiceApp.FileOpSynGateWay(args)
		if err != nil {
			global.ZC_LOG.Error("FileOpSynGateWay", zap.String("err:", err.Error()))
			return
		}
	}()

	return res, nil
}

// SysFileFrom OP point-to-point replication
func (p *OpServiceImpl) SysFileFrom(ctx context.Context, args *pb.SynFileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"

	err := system.JobPlatformServiceApp.FilePointProcess(args)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}

	return res, nil
}

// DownLoadFiles File Download
func (p *OpServiceImpl) DownLoadFiles(ctx context.Context, args *pb.DownLoadInfo) (*pb.ResponseMsg, error) {

	log.Println("op DownLoadFiles success")
	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"

	return res, nil
}

// UnZipSynFile Resolve files and delete them
func (p *OpServiceImpl) UnZipSynFile(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"

	err := system.JobPlatformServiceApp.UnZipSynFile(args)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}

	return res, nil
}

// OpInformationTest Op performs host information testing
func (p *OpServiceImpl) OpInformationTest(ctx context.Context, args *pb.HostTestInfo) (*pb.String, error) {
	log.Println("Host information test args", args)

	utils.KillBenchAndKillScript()
	opStorageInfo, opServer, err := getHostInfo(&pb.OpHardwareInfo{
		HostClassify: args.TestType,
		HostUUID:     args.HostUUID,
	})
	if err != nil {
		return &pb.String{Value: "Failed to obtain host information!"}, err
	}

	var isGPUPass bool
	if args.TestType != config.HostStorageTest {
		if len(opStorageInfo.OpGPUInfo) > 0 {
			for i := 0; i < len(opStorageInfo.OpGPUInfo); i++ {
				infos := utils.GetNumFromStr(opStorageInfo.OpGPUInfo[i].Model)
				gpuModel, _ := strconv.Atoi(infos)
				if gpuModel >= config.GPUModelStandard {
					isGPUPass = true
					opStorageInfo.OpGPUInfo[i].IsQualified = config.HostParamTestPassStr
				} else {
					opStorageInfo.OpGPUInfo[i].IsQualified = config.HostParamTestDissatisfactionStr
				}
			}
		}
	}

	cpuInfo, _ := utils.StructToJsonStr(opStorageInfo.OpCPUInfo)
	memoryInfo, _ := utils.StructToJsonStr(opStorageInfo.OpRamInfo)

	hostMonitorInfo := &pb.UpdateHostTestInfo{
		TestBeginAt:     args.TestBeginAt,
		TestType:        args.TestType,
		IsAddPower:      args.IsAddPower,
		SelectHostUUIDs: args.HostUUIDs,
		SelectHostIPs:   args.HostIPs,
		HostUUID:        args.HostUUID,
		CPUHardInfo:     cpuInfo,
		MemoryHardInfo:  memoryInfo,
	}

	if args.TestType != config.HostStorageTest {
		gpuInfo, _ := utils.StructToJsonStr(opStorageInfo.OpGPUInfo)
		hostMonitorInfo.GPUHardInfo = gpuInfo

		if isGPUPass {
			hostMonitorInfo.GPUHardScore = config.HostParamTestPass
		} else {
			hostMonitorInfo.GPUHardScore = config.HostParamTestDissatisfaction
		}
	}

	if opServer.CPU.Threads > config.CPUThreads {
		hostMonitorInfo.CPUHardScore = config.HostParamTestPass
	} else {
		hostMonitorInfo.CPUHardScore = config.HostParamTestDissatisfaction
	}

	if args.TestType != config.HostMinerTest {
		diskIO, err := utils.GetDiskIOStatus()
		if err != nil {
			hostMonitorInfo.DiskIO = config.HostParamTestFailed
			hostMonitorInfo.TestResult = config.HostTestFailed
		} else {
			if diskIO == config.HostDiskIOPass {
				hostMonitorInfo.DiskIO = config.HostParamTestPass
			} else {
				hostMonitorInfo.DiskIO = config.HostParamTestDissatisfaction
			}
		}
	}

	if args.IsAddPower {
		ips := strings.Split(args.HostIPs, ",")
		sizeStrs := []string{}

		c := make(chan string, len(ips))

		for key, val := range ips {
			go func(val string) {
				netRate := utils.StartHostPortRequest(val, config.HostNetTestPort[key], config.TestHostNetTimeAdd)
				if netRate == "" {
					c <- ""
					return
				} else {
					c <- netRate
				}
			}(val)
		}

		for i := 0; i < len(ips); i++ {
			sizeStrs = append(sizeStrs, <-c)
		}

		sizeSum := utils.AddInputSizes(sizeStrs)
		hostMonitorInfo.NetTestInfo = sizeSum
		if utils.CompareTwoSizes(sizeSum, config.NetworkSpeedAdd) {
			hostMonitorInfo.NetTestScore = config.HostParamTestPass
		} else {
			hostMonitorInfo.NetTestScore = config.HostParamTestDissatisfaction
		}
	} else {
		netRate := utils.StartHostPortRequest(args.HostIPs, config.HostNetTestPort[0], config.TestHostNetTime)
		if netRate == "" {
			hostMonitorInfo.NetTestInfo = "--"
			hostMonitorInfo.NetTestScore = config.HostParamTestFailed
			hostMonitorInfo.TestResult = config.HostTestFailed
		} else {
			hostMonitorInfo.NetTestInfo = netRate
			if args.TestType == config.HostC2WorkerTest {
				if utils.CompareTwoSizes(netRate, config.NetworkC2WorkerSpeed) {
					hostMonitorInfo.NetTestScore = config.HostParamTestPass
				} else {
					hostMonitorInfo.NetTestScore = config.HostParamTestDissatisfaction
				}
			} else {
				if utils.CompareTwoSizes(netRate, config.NetworkSpeed) {
					hostMonitorInfo.NetTestScore = config.HostParamTestPass
				} else {
					hostMonitorInfo.NetTestScore = config.HostParamTestDissatisfaction
				}
			}
		}
	}

	_, err = global.OpToGatewayClient.UpdateHostNetOccupyTime(ctx, &pb.UpdateHostNet{HostUUIDs: args.HostUUIDs})
	if err != nil {
		log.Println("client.OpInformationTest", "Failed to modify the time when the host is being used to test other host networks!", err.Error())
	}

	diskSize := ""
	memorySize := opServer.Ram.TotalMB / 1024

	if args.TestType == config.HostMinerTest {
		dealHostMinerTestInfo(opStorageInfo, memorySize, hostMonitorInfo, args.IsAddPower)
	} else if args.TestType == config.HostWorkerTest {
		dealHostWorkerTestInfo(opStorageInfo, memorySize, hostMonitorInfo)
	} else if args.TestType == config.HostStorageTest {
		diskSize = dealHostStorageTestInfo(opStorageInfo, opServer, memorySize, hostMonitorInfo)
	} else {
		if memorySize > config.C2WorkerRamSize {
			hostMonitorInfo.MemoryHardScore = config.HostParamTestPass
		} else {
			hostMonitorInfo.MemoryHardScore = config.HostParamTestDissatisfaction
		}
	}

	hostDisk := utils.HostDisk{
		DiskSize:      diskSize,
		OpDiskMd0Info: opStorageInfo.OpDiskMd0Info,
	}

	if args.TestType != config.HostC2WorkerTest {
		diskInfo, _ := utils.StructToJsonStr(hostDisk)
		hostMonitorInfo.DiskHardInfo = diskInfo
	}

	if args.TestType != config.HostStorageTest {
		hostMonitorInfo.TestResult = config.HostUnderTest
		_, err = global.OpToGatewayClient.UpdateHostTestRecord(ctx, hostMonitorInfo)
		if err != nil {
			log.Println("client.OpInformationTest", "Failed to record host test information before recording GPU information", err.Error())
		}
		out := utils.GetGPURunBenchTime()
		if out == "" {
			hostMonitorInfo.GPUTestInfo = "--"
			hostMonitorInfo.GPUTestScore = config.HostParamTestFailed
			hostMonitorInfo.TestResult = config.HostTestFailed
		} else {
			gpuRunTime, _ := strconv.ParseFloat(out, 64)
			hostMonitorInfo.GPUTestInfo = fmt.Sprintf("%.0f", gpuRunTime) + "Min"
			if gpuRunTime <= config.GPURunTime {
				hostMonitorInfo.GPUTestScore = config.HostParamTestPass
			} else {
				hostMonitorInfo.GPUTestScore = config.HostParamTestDissatisfaction
			}
		}
	}

	getTestResult(hostMonitorInfo, args)

	// 插入数据信息
	_, err = global.OpToGatewayClient.UpdateHostTestRecord(ctx, hostMonitorInfo)
	if err != nil {
		log.Println("client.OpInformationTest", "Host test information recording failed!", err.Error())
		return &pb.String{Value: "Host test information recording failed!"}, err
	}
	log.Println("client.OpInformationTest", "Host test information recorded successfully!")
	return &pb.String{Value: "Host test information recorded successfully!"}, nil
}

// 确定最后通过状态
func getTestResult(hostMonitorInfo *pb.UpdateHostTestInfo, args *pb.HostTestInfo) {
	if hostMonitorInfo.TestResult != config.HostTestFailed {
		isTestSuccess := false
		cpuHardStatus := hostMonitorInfo.CPUHardScore == config.HostParamTestPass
		gpuHardStatus := hostMonitorInfo.GPUHardScore == config.HostParamTestPass
		memoryHardStatus := hostMonitorInfo.MemoryHardScore == config.HostParamTestPass
		diskHardStatus := hostMonitorInfo.DiskHardScore == config.HostParamTestPass
		netTestStatus := hostMonitorInfo.NetTestScore == config.HostParamTestPass
		gpuTestStatus := hostMonitorInfo.GPUTestScore == config.HostParamTestPass
		diskIOStatus := hostMonitorInfo.DiskIO == config.HostParamTestPass
		diskSSDRateStatus := hostMonitorInfo.DiskSSDRateScore == config.HostParamTestPass
		diskAllRateStatus := hostMonitorInfo.DiskAllRateScore == config.HostParamTestPass
		if args.TestType == config.HostMinerTest {
			if cpuHardStatus && gpuHardStatus && memoryHardStatus && diskHardStatus && netTestStatus && gpuTestStatus && diskSSDRateStatus {
				isTestSuccess = true
			}
		} else if args.TestType == config.HostWorkerTest {
			if cpuHardStatus && gpuHardStatus && memoryHardStatus && diskHardStatus && netTestStatus && gpuTestStatus && diskIOStatus && diskSSDRateStatus {
				isTestSuccess = true
			}
		} else if args.TestType == config.HostStorageTest {
			if cpuHardStatus && memoryHardStatus && diskHardStatus && netTestStatus && diskIOStatus && diskAllRateStatus {
				isTestSuccess = true
			}
		} else if args.TestType == config.HostC2WorkerTest {
			if cpuHardStatus && gpuHardStatus && memoryHardStatus && netTestStatus && diskIOStatus && gpuTestStatus {
				isTestSuccess = true
			}
		}
		if isTestSuccess {
			hostMonitorInfo.TestResult = config.HostCompliance
		} else {
			hostMonitorInfo.TestResult = config.HostNotUpToStandard
		}
	}
}

func dealHostMinerTestInfo(opStorageInfo utils.OpStorageInfo, memorySize int, hostMonitorInfo *pb.UpdateHostTestInfo, isAddPower bool) {
	if utils.CompareTwoSizes(opStorageInfo.OpDiskMd0Info.TotalSize, config.NodeSSDDiskSize) {
		hostMonitorInfo.DiskHardScore = config.HostParamTestPass
	} else {
		hostMonitorInfo.DiskHardScore = config.HostParamTestDissatisfaction
	}
	if isAddPower {
		if memorySize > config.NodeRamSizeAdditional {
			hostMonitorInfo.MemoryHardScore = config.HostParamTestPass
		} else {
			hostMonitorInfo.MemoryHardScore = config.HostParamTestDissatisfaction
		}
	} else {
		if memorySize > config.NodeRamSize {
			hostMonitorInfo.MemoryHardScore = config.HostParamTestPass
		} else {
			hostMonitorInfo.MemoryHardScore = config.HostParamTestDissatisfaction
		}
	}

	diskRate, err := getMd0DiskRate()
	if err == nil {
		hostMonitorInfo.DiskSSDRate = diskRate[1]
		if utils.CompareTwoSizes(diskRate[1], config.NodeSSDDiskSpeed) {
			hostMonitorInfo.DiskSSDRateScore = config.HostParamTestPass
		} else {
			hostMonitorInfo.DiskSSDRateScore = config.HostParamTestDissatisfaction
		}
	} else {
		hostMonitorInfo.DiskSSDRateScore = config.HostParamTestFailed
		hostMonitorInfo.TestResult = config.HostTestFailed
	}
}

func dealHostWorkerTestInfo(opStorageInfo utils.OpStorageInfo, memorySize int, hostMonitorInfo *pb.UpdateHostTestInfo) {
	if utils.CompareTwoSizes(opStorageInfo.OpDiskMd0Info.TotalSize, config.WorkerSSDDiskSize) {
		hostMonitorInfo.DiskHardScore = config.HostParamTestPass
	} else {
		hostMonitorInfo.DiskHardScore = config.HostParamTestDissatisfaction
	}

	if memorySize > config.WorkerRamSize {
		hostMonitorInfo.MemoryHardScore = config.HostParamTestPass
	} else {
		hostMonitorInfo.MemoryHardScore = config.HostParamTestDissatisfaction
	}

	diskRate, err := getMd0DiskRate()
	if err == nil {
		hostMonitorInfo.DiskSSDRate = diskRate[1]
		if utils.CompareTwoSizes(diskRate[1], config.WorkerSSDDiskSpeed) {
			hostMonitorInfo.DiskSSDRateScore = config.HostParamTestPass
		} else {
			hostMonitorInfo.DiskSSDRateScore = config.HostParamTestDissatisfaction
		}
	} else {
		hostMonitorInfo.DiskSSDRateScore = config.HostParamTestFailed
		hostMonitorInfo.TestResult = config.HostTestFailed
	}
}

func dealHostStorageTestInfo(opStorageInfo utils.OpStorageInfo, opServer utils.OpServer, memorySize int, hostMonitorInfo *pb.UpdateHostTestInfo) string {

	diskSize := ""
	sizeSum := 0
	for _, val := range opServer.Storage {
		sizeSum += int(val.Size)
	}
	diskSize = utils.IntToString(sizeSum) + "G"
	if utils.CompareTwoSizes(diskSize, config.StorageDiskSize) {
		hostMonitorInfo.DiskHardScore = config.HostParamTestPass
	} else {
		hostMonitorInfo.DiskHardScore = config.HostParamTestDissatisfaction
	}

	if memorySize > config.StorageRamSize {
		hostMonitorInfo.MemoryHardScore = config.HostParamTestPass
	} else {
		hostMonitorInfo.MemoryHardScore = config.HostParamTestDissatisfaction
	}

	diskRate, err := getAllDiskRate()
	if err == nil {
		hostMonitorInfo.DiskAllRate = diskRate
		if utils.CompareTwoSizes(diskRate, config.StorageDiskOverallSpeed) {
			hostMonitorInfo.DiskAllRateScore = config.HostParamTestPass
		} else {
			hostMonitorInfo.DiskAllRateScore = config.HostParamTestDissatisfaction
		}
	} else {
		hostMonitorInfo.DiskAllRateScore = config.HostParamTestFailed
		hostMonitorInfo.TestResult = config.HostTestFailed
	}
	return diskSize
}

// getAllDiskRate Get MD0 disk read information
func getAllDiskRate() (string, error) {
	out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathStorageDiskReadRate+" %s", config.TestDefaultFile)).Output()
	if err != nil {
		log.Println("Failed to execute disk read/write speed test script", err.Error())
		return "", err
	}

	strs := strings.Split(string(out), "\n")
	maxRate := ""
	for _, val := range strs {
		if len(val) == 0 || val == "\n" {
			continue
		}
		strVals := strings.Split(val, " ")
		dealVals := []string{}
		for _, v := range strVals {
			if len(v) == 0 {
				continue
			}
			dealVals = append(dealVals, v)
		}

		if len(dealVals) < 3 {
			return "", nil
		}

		val = dealVals[2] + dealVals[3][:len(dealVals[3])-3]
		if len(maxRate) == 0 {
			maxRate = val
		} else {
			if !utils.CompareTwoSizes(maxRate, val) {
				maxRate = val
			}
		}
	}
	return maxRate, nil
}

// getMd0DiskRate Get MD0 disk read information
func getMd0DiskRate() ([]string, error) {
	out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathDiskReadRate+" %s", config.TestDefaultFile)).Output()
	if err != nil {
		log.Println("Failed to execute disk read/write speed test script", err.Error())
		return []string{}, err
	}
	strs := strings.Split(string(out), "\n")
	dealStrs := []string{}
	for _, val := range strs {
		val = strings.ReplaceAll(val, " ", "")
		if len(val) < 3 || val == "\n" {
			continue
		}
		val = val[:len(val)-3]
		dealStrs = append(dealStrs, val)
	}
	return dealStrs, nil
}

// KillBenchAndScript Kill the kill, net, and bench script processes in the OP node
func (p *OpServiceImpl) KillBenchAndScript(ctx context.Context, args *pb.String) (*pb.String, error) {

	err := utils.KillBenchAndKillScript()
	if err != nil {
		log.Println("The kill, net, and bench script processes in the host kill op node have failed")
		return &pb.String{Value: "The kill, net, and bench script processes in the host kill op node have failed"}, err
	}

	return &pb.String{Value: "The kill, net, and bench script processes in the host kill op node were successful"}, err
}

// LotusHeight Lotus height
func (p *OpServiceImpl) LotusHeight(ctx context.Context, args *pb.RequestOp) (*pb.LotusHeightInfo, error) {

	h, err := lotusrpc.FullApi.LotusHeight(args.Token, args.Ip)
	if err != nil {
		return nil, err
	}

	return &pb.LotusHeightInfo{Height: int64(h)}, nil
}

// OpInformationPatrol Op performs host information patrol
func (p *OpServiceImpl) OpInformationPatrol(ctx context.Context, args *pb.HostPatrolInfo) (*pb.String, error) {
	log.Println("Host information patrol args", args)

	hostPatroInfo := &pb.UpdateHostPatrolInfo{
		PatrolBeginAt: args.PatrolBeginAt,
		HostUUID:      args.HostUUID,
	}

	ioBeginTime := time.Now()
	diskIO, err := utils.GetDiskIOStatus()
	if err != nil {
		hostPatroInfo.DiskIO = config.HostPatrolStatusFailed
	} else {
		if diskIO == config.HostDiskIOPass {
			hostPatroInfo.DiskIO = config.HostPatrolStatusSuccess
		} else {
			hostPatroInfo.DiskIO = config.HostPatrolStatusFailed
		}
	}
	hostPatroInfo.DiskIODuration = DealDurationStr(time.Now().Sub(ioBeginTime).String())

	timeSyncBeginTime := time.Now()
	timeSync, err := utils.GetTimeSyncStatus()
	if err != nil {
		hostPatroInfo.TimeSyncStatus = config.HostPatrolStatusFailed
	} else {
		if timeSync {
			hostPatroInfo.TimeSyncStatus = config.HostPatrolStatusSuccess
		} else {
			hostPatroInfo.TimeSyncStatus = config.HostPatrolStatusFailed
		}
	}
	hostPatroInfo.TimeSyncDuration = DealDurationStr(time.Now().Sub(timeSyncBeginTime).String())

	if args.HostClassify == config.HostMinerTest {
		netBeginTime := time.Now()
		if utils.HostNetStatus() {
			hostPatroInfo.HostNetStatus = config.HostPatrolStatusSuccess
		} else {
			hostPatroInfo.HostNetStatus = config.HostPatrolStatusFailed
		}
		hostPatroInfo.HostNetDuration = DealDurationStr(time.Now().Sub(netBeginTime).String())

		logOvertimeBeginTime := time.Now()
		isNormal, err := utils.GetLogOvertimeStatus()
		if err != nil {
			hostPatroInfo.LogOvertimeStatus = config.HostPatrolStatusFailed
		} else {
			if isNormal {
				hostPatroInfo.LogOvertimeStatus = config.HostPatrolStatusSuccess
			} else {
				hostPatroInfo.LogOvertimeStatus = config.HostPatrolStatusFailed
			}
		}
		hostPatroInfo.LogOvertimeDuration = DealDurationStr(time.Now().Sub(logOvertimeBeginTime).String())

		logOvertimeBlock := time.Now()
		_, err = utils.GetBlockLogStatus()
		if err != nil {
			hostPatroInfo.BlockLogStatus = config.HostPatrolStatusFailed
		} else {
			hostPatroInfo.BlockLogStatus = config.HostPatrolStatusSuccess
		}
		hostPatroInfo.BlockLogDuration = DealDurationStr(time.Now().Sub(logOvertimeBlock).String())

		logInformationBeginTime := time.Now()
		isNormal, err = utils.GetLogInformationStatus()
		if err != nil {
			hostPatroInfo.LogInfoStatus = config.HostPatrolStatusFailed
		} else {
			if isNormal {
				hostPatroInfo.LogInfoStatus = config.HostPatrolStatusSuccess
			} else {
				hostPatroInfo.LogInfoStatus = config.HostPatrolStatusFailed
			}
		}
		hostPatroInfo.LogInfoDuration = DealDurationStr(time.Now().Sub(logInformationBeginTime).String())

		wdpostBalanceBeginTime := time.Now()
		balance, err := utils.GetWdpostBalance()
		hostPatroInfo.WalletBalance = float32(balance)
		if err != nil {
			hostPatroInfo.WalletBalanceStatus = config.HostPatrolStatusFailed
		} else {
			if hostPatroInfo.WalletBalance >= config.PatrolWdpostBalance {
				hostPatroInfo.WalletBalanceStatus = config.HostPatrolStatusSuccess
			} else {
				hostPatroInfo.WalletBalanceStatus = config.HostPatrolStatusFailed
			}
		}
		hostPatroInfo.WalletBalanceDuration = DealDurationStr(time.Now().Sub(wdpostBalanceBeginTime).String())

		lotusBeginTime := time.Now()
		isNormal, err = utils.GetLotusHigh()
		if err != nil {
			hostPatroInfo.LotusSyncStatus = config.HostPatrolStatusFailed
		} else {
			if isNormal {
				hostPatroInfo.LotusSyncStatus = config.HostPatrolStatusSuccess
			} else {
				hostPatroInfo.LotusSyncStatus = config.HostPatrolStatusFailed
			}
		}
		hostPatroInfo.LotusSyncDuration = DealDurationStr(time.Now().Sub(lotusBeginTime).String())

		versionBeginTime := time.Now()
		packageVersion, _ := utils.GetLotusPackageVersion()
		hostPatroInfo.PackageVersion, _ = utils.StructToJsonStr(packageVersion)
		if strings.EqualFold(packageVersion.LotusVersion, config.PatrolLotusVersion) && strings.EqualFold(packageVersion.MinerVersion, config.PatrolMinerVersion) &&
			strings.EqualFold(packageVersion.BoostdVersion, config.PatrolBoostdVersion) {
			hostPatroInfo.PackageVersionStatus = config.HostPatrolStatusSuccess
		} else {
			hostPatroInfo.PackageVersionStatus = config.HostPatrolStatusFailed
		}
		hostPatroInfo.PackageVersionDuration = DealDurationStr(time.Now().Sub(versionBeginTime).String())

		dataCatalogBeginTime := time.Now()
		isNormal, err = utils.GetDataCatalogStatus()
		if err != nil {
			hostPatroInfo.DataCatalogStatus = config.HostPatrolStatusFailed
		} else {
			if isNormal {
				hostPatroInfo.DataCatalogStatus = config.HostPatrolStatusSuccess
			} else {
				hostPatroInfo.DataCatalogStatus = config.HostPatrolStatusFailed
			}
		}
		hostPatroInfo.DataCatalogDuration = DealDurationStr(time.Now().Sub(dataCatalogBeginTime).String())
	}

	if args.HostClassify != config.HostStorageTest {
		gpuBeginTime := time.Now()
		err = utils.CheckGPUDrive()
		if err != nil {
			hostPatroInfo.GPUDriveStatus = config.HostPatrolStatusFailed
		} else {
			hostPatroInfo.GPUDriveStatus = config.HostPatrolStatusSuccess
		}
		hostPatroInfo.GPUDriveDuration = DealDurationStr(time.Now().Sub(gpuBeginTime).String())

		hostDownBeginTime := time.Now()
		downStatus, err := utils.GetHostDownStatus()
		if err != nil {
			hostPatroInfo.HostIsDown = config.HostPatrolStatusFailed
		} else {
			if downStatus {
				hostPatroInfo.HostIsDown = config.HostPatrolStatusSuccess
			} else {
				hostPatroInfo.HostIsDown = config.HostPatrolStatusFailed
			}
		}
		hostPatroInfo.HostIsDownDuration = DealDurationStr(time.Now().Sub(hostDownBeginTime).String())
	}

	if args.HostClassify == config.HostStorageTest {
		pingNetBeginTime := time.Now()
		if utils.GetHostPingStatus(args.PatrolHostIP) {
			hostPatroInfo.PingNetStatus = config.HostPatrolStatusSuccess
		} else {
			hostPatroInfo.PingNetStatus = config.HostPatrolStatusFailed
		}
		hostPatroInfo.PingNetDuration = DealDurationStr(time.Now().Sub(pingNetBeginTime).String())
	}

	if args.HostClassify == config.HostMinerTest {
		if hostPatroInfo.DiskIO && hostPatroInfo.HostNetStatus && hostPatroInfo.LogInfoStatus && hostPatroInfo.LogOvertimeStatus && hostPatroInfo.BlockLogStatus &&
			hostPatroInfo.LotusSyncStatus && hostPatroInfo.GPUDriveStatus && hostPatroInfo.PackageVersionStatus && hostPatroInfo.DataCatalogStatus && hostPatroInfo.TimeSyncStatus &&
			hostPatroInfo.HostIsDown {
			hostPatroInfo.PatrolResult = config.HostCompliance
		} else {
			hostPatroInfo.PatrolResult = config.HostNotUpToStandard
		}
	}
	if args.HostClassify == config.HostWorkerTest || args.HostClassify == config.HostC2WorkerTest {
		if hostPatroInfo.DiskIO && hostPatroInfo.HostIsDown && hostPatroInfo.GPUDriveStatus && hostPatroInfo.TimeSyncStatus {
			hostPatroInfo.PatrolResult = config.HostCompliance
		} else {
			hostPatroInfo.PatrolResult = config.HostNotUpToStandard
		}
	}
	if args.HostClassify == config.HostStorageTest {
		if hostPatroInfo.DiskIO && hostPatroInfo.TimeSyncStatus && hostPatroInfo.PingNetStatus {
			hostPatroInfo.PatrolResult = config.HostCompliance
		} else {
			hostPatroInfo.PatrolResult = config.HostNotUpToStandard
		}
	}

	_, err = global.OpToGatewayClient.UpdateHostPatrolRecord(ctx, hostPatroInfo)
	if err != nil {
		log.Println("client.OpInformationPatrol", "Host inspection information recording failed!", err.Error())
		return &pb.String{Value: "Host inspection information recording failed!"}, err
	}
	log.Println("client.OpInformationPatrol", "Host inspection information recorded successfully!")
	return &pb.String{Value: "Host inspection information recorded successfully!"}, nil
}

// DealDurationStr Processing interval time to two decimal places
func DealDurationStr(str string) string {
	strNum := utils.GetNumFromStr(str)
	strLetter := utils.GetLetterFromStr(str)
	num, _ := strconv.ParseFloat(strNum, 64)
	return strconv.FormatFloat(num, 'f', 2, 64) + strLetter
}

// DelOpFile delete file
func (p *OpServiceImpl) DelOpFile(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"

	dirAr := strings.Split(define.MainDisk, "/")
	dir := fmt.Sprintf("/%s/", dirAr[1])
	if !strings.Contains(args.FileName, dir) {
		args.FileName = define.MainDisk + args.FileName
	}
	err := utils.DelFile(args.FileName)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}

	return res, nil
}

// DelGateWayFile Delete gateway file
func (p *OpServiceImpl) DelGateWayFile(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	_, err := global.OpToGatewayClient.DelGateWayFile(ctx, args)
	if err != nil {
		return res, err
	}

	return res, nil
}

// CreateOpFile save the file
func (p *OpServiceImpl) CreateOpFile(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	if !strings.Contains(args.Path, define.MainDisk) {
		args.Path = define.MainDisk + args.Path
	}
	res := &pb.ResponseMsg{}
	err := utils.CreateFile(utils.FileInfo{Path: args.Path, FileName: args.FileName, FileData: args.FileData})
	if err != nil {
		return res, err
	}

	return res, nil
}

// CloseOpInformationTest Op Shutdown Host Information Test
func (p *OpServiceImpl) CloseOpInformationTest(ctx context.Context, args *pb.CloseHostTest) (*pb.String, error) {
	log.Println("Close Host information test args", args)
	go utils.CloseHostPortMonitor()
	time.Sleep(20 * time.Second)
	if args.TestType != config.HostStorageTest {
		utils.CloseGPURunBench()
	}
	log.Println("client.OpInformationTest", "close host test information recorded successfully!")
	return &pb.String{Value: "close host test information recorded successfully!"}, nil
}

// OpInformationTestRestart Op performs host information testing again
func (p *OpServiceImpl) OpInformationTestRestart(ctx context.Context, args *pb.RestartHostTest) (*pb.String, error) {
	log.Println("Restart host information test args", args)
	if args.HostClassify == config.HostMinerTest {
		// Script for executing host information testing
		//if out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathIpfsScriptOpInformationTest)).Output(); err != nil {
		//	log.Println("Failed to execute host information test script", err.Error())
		//	return nil, err
		//} else {
		//	log.Println(string(out))
		//}
		log.Println("-------------Test output results-------------")
	} else if args.HostClassify == config.HostStorageTest {
		// Script for executing host information testing
		//if out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathIpfsScriptOpInformationTest)).Output(); err != nil {
		//	log.Println("Failed to execute host information test script", err.Error())
		//	return nil, err
		//} else {
		//	log.Println(string(out))
		//}
	}

	hostMonitorInfo := &pb.UpdateHostTestInfo{
		TestBeginAt:     args.TestBeginAt,
		HostUUID:        args.HostUUID,
		TestResult:      config.HostCompliance,
		CPUHardInfo:     "testCPUHardInfo",
		CPUHardScore:    0,
		GPUHardInfo:     "testGPUHardInfo",
		GPUHardScore:    0,
		MemoryHardInfo:  "testMemoryHardInfo",
		MemoryHardScore: 0,
		DiskHardInfo:    "testDiskHardInfo",
		DiskHardScore:   0,
		NetTestInfo:     "testNetTestInfo",
		NetTestScore:    0,
		GPUTestInfo:     "",
		GPUTestScore:    0,
	}

	_, err := global.OpToGatewayClient.UpdateHostTestRecord(ctx, hostMonitorInfo)
	if err != nil {
		log.Println("client.OpInformationTest", "Re execute Host test information recording failed!", err.Error())
		return &pb.String{Value: "Re execute host test information recording failed!"}, err
	}
	log.Println("client.OpInformationTest", "Re execute host test information recorded successfully!")
	return &pb.String{Value: "Re execute host test information recorded successfully!"}, nil
}

// GetOpHardwareInfo Obtain Op hardware information
func (p *OpServiceImpl) GetOpHardwareInfo(ctx context.Context, args *pb.OpHardwareInfo) (*pb.String, error) {
	opStorageInfo, _, err := getHostInfo(args)
	if err != nil {
		return nil, err
	}
	diskInfoStr, _ := utils.StructToJsonStr(opStorageInfo)
	return &pb.String{Value: diskInfoStr}, nil
}

func getHostInfo(args *pb.OpHardwareInfo) (utils.OpStorageInfo, utils.OpServer, error) {
	var s utils.OpServer
	var err error
	s.Os = utils.InitOS()
	s.GetCPUInfo()
	s.GetGPUInfo()
	if s.Ram, err = utils.InitRAM(); err != nil {
		log.Println("func utils.InitRAM() Failed", err.Error())
		return utils.OpStorageInfo{}, s, err
	}
	s.GetStorageInfo()

	var opStorageInfo utils.OpStorageInfo

	opStorageInfo.OpCPUInfo = utils.OpCPUInfo{
		Brand:       s.CPU.Vendor,
		Model:       s.CPU.Model,
		CpuNum:      strconv.Itoa(int(s.CPU.Cores)),
		Speed:       s.CPU.Speed,
		Threads:     strconv.Itoa(int(s.CPU.Threads)),
		UsedPercent: utils.Float64ToString(s.CPU.UsedPercent),
	}

	var opGPUInfo []utils.OpGPUInfo
	if len(s.GPU.Gpus) > 0 {
		for _, val := range s.GPU.Gpus {
			info := utils.OpGPUInfo{
				Mark:    val.Mark,
				Brand:   val.Brand,
				Model:   val.GpuInfo,
				TotalMB: val.TotalMB,
				UseRate: val.UseRate,
			}
			opGPUInfo = append(opGPUInfo, info)
		}
	}
	opStorageInfo.OpGPUInfo = opGPUInfo

	opRamInfo := utils.OpRamInfo{
		TotalMB:     strconv.Itoa(s.Ram.TotalMB) + "M",
		UsedMB:      strconv.Itoa(s.Ram.UsedMB) + "M",
		UsedPercent: strconv.Itoa(s.Ram.UsedPercent),
	}
	ramNum, ramType := utils.GetRamTypeAndNum()
	opRamInfo.Type = ramType
	opRamInfo.RamNum = strconv.Itoa(ramNum)
	opStorageInfo.OpRamInfo = opRamInfo

	if args.HostClassify == config.HostStorageType || args.HostClassify == config.HostDCStorageType || args.HostClassify == 0 {
		var opDiskInfoList []utils.OpDiskInfo
		storageInfo := utils.GetStorageOpDiskInfo()
		for _, val := range s.Storage {
			opDiskInfo := utils.OpDiskInfo{
				Name:  val.Name,
				Brand: val.Model,
				Model: val.Serial,
			}
			for _, v := range storageInfo {
				if strings.Contains(v.Filesystem, val.Name) {
					opDiskInfo.TotalSize = v.Size
					opDiskInfo.UsedSize = v.Used
					opDiskInfo.UsedPercent = v.UsePercent
					if v.Mountpoint == "/" {
						opDiskInfo.IsOs = 1
					}
				}
			}
			opDiskInfoList = append(opDiskInfoList, opDiskInfo)
		}
		opStorageInfo.OpDiskInfo = opDiskInfoList
	} else {
		osInfo := utils.GetStorageOpOsDiskInfo()
		diskMd0Info := utils.GetDiskMd0Info()

		opStorageInfo.OpDiskInfo = []utils.OpDiskInfo{osInfo}
		opStorageInfo.OpDiskMd0Info = diskMd0Info
	}
	return opStorageInfo, s, nil
}

// AddGateWayFile Add GateWay file
func (p *OpServiceImpl) AddGateWayFile(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	_, err := global.OpToGatewayClient.AddGateWayFile(ctx, args)
	if err != nil {
		return res, err
	}

	return res, nil
}

// OpServerPortControl Op performs host information testing, executes scripts to enable iperf3 port monitoring
func (p *OpServiceImpl) OpServerPortControl(ctx context.Context, args *pb.String) (*pb.String, error) {
	err := utils.StartHostPortMonitor()
	if err != nil {
		return &pb.String{Value: "Failed to start testing network broadband server"}, err
	}
	return nil, nil
}

// GetOpFilePath OP file path
func (p *OpServiceImpl) GetOpFilePath(ctx context.Context, args *pb.OpFilePath) (*pb.String, error) {

	res := &pb.String{}
	r, err := deploy.SectorsRecoverServiceApi.GetOpFilePath(args)
	if err != nil {
		return res, err
	}

	return &pb.String{Value: r}, nil
}

// CloseOpServerPortControl Close Op to execute host information test script
func (p *OpServiceImpl) CloseOpServerPortControl(ctx context.Context, args *pb.String) (*pb.String, error) {
	err := utils.CloseHostPortMonitor()
	if err != nil {
		return &pb.String{Value: "Failed to close testing network broadband testing script"}, err
	}
	log.Println("Successfully closed testing network broadband testing script")
	return nil, nil
}

// AddWarn Add alarm information (4 strategy alarms, 5 inspection alarms, 6 business alarms)
func (p *OpServiceImpl) AddWarn(ctx context.Context, args *pb.WarnInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	args.ComputerId = global.OpUUID.String()
	_, err := global.OpToGatewayClient.AddWarn(ctx, args)
	if err != nil {
		return res, err
	}

	return res, nil
}

// AddBadSector Add error sectors
func (p *OpServiceImpl) AddBadSector(ctx context.Context, args *pb.BadSectorId) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	_, err := global.OpToGatewayClient.AddBadSector(ctx, args)
	if err != nil {
		return res, err
	}

	return res, nil
}

// OpFileToGateWay Synchronize host files to gateWay
func (p *OpServiceImpl) OpFileToGateWay(ctx context.Context, args *pb.AddFileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	_, err := system.JobPlatformServiceApp.OpLocalFileSynGateWay(args)
	if err != nil {
		return res, err
	}

	return res, nil
}

// CheckOpPath Check the path of the op file
func (p *OpServiceImpl) CheckOpPath(ctx context.Context, args *pb.DirFileReq) (*pb.ResponseMsg, error) {
	log.Println("CheckOpPath begin")
	res := &pb.ResponseMsg{}
	if !utils.IsNull(args.FileName) {
		_, err := ioutil.ReadFile(path.Join(args.Path, args.FileName))
		if err != nil {
			return res, err
		}
	} else {
		_, err := ioutil.ReadDir(args.Path)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

// CarFilePath Get the car file path
func (p *OpServiceImpl) CarFilePath(ctx context.Context, args *pb.CarFile) (*pb.CarFile, error) {

	res := &pb.CarFile{}
	res, err := deploy.SectorsRecoverServiceApi.GetCarFilePath(args)
	if err != nil {
		return res, err
	}

	return res, nil
}

// OpCheckBadSector Op execution query error sector information
func (p *OpServiceImpl) OpCheckBadSector(ctx context.Context, args *pb.HostCheckDiskInfo) (*pb.String, error) {
	log.Println("Host information test args", args)
	minerId, sectorIds, err := utils.GetHostBadSectorInfo()
	if err != nil {
		return &pb.String{Value: "Failed to obtain host information!"}, err
	}
	if len(sectorIds) == 0 {
		return &pb.String{Value: "This node has no error sectors!"}, nil
	}

	sectorSize, err := utils.GetHostBadSectorSize()
	if err != nil {
		return &pb.String{Value: "Failed to obtain the corresponding sector size!"}, err
	}
	size, _ := strconv.Atoi(sectorSize)

	for _, val := range sectorIds {
		badSectorInfo := &pb.BadSectorId{
			MinerId:       minerId,
			SectorSize:    uint64(size),
			BelongingNode: args.HostUUID,
			SectorAddress: "",
		}

		sectorId, _ := strconv.Atoi(val)
		badSectorInfo.SectorId = uint64(sectorId)

		sectorType, err := utils.GetHostBadSectorType(val)
		if err != nil {
			log.Println("client.AddBadSector", fmt.Sprintf("Error in obtaining sector type! uuid:%s,sectorIs:%s", args.HostUUID, val), err.Error())
		} else {
			switch sectorType {
			case config.SectorTypeCCType:
				badSectorInfo.SectorType = config.SectorTypeCCTypeNum
			case config.SectorTypeDCType:
				badSectorInfo.SectorType = config.SectorTypeDCTypeNum
			default:
				badSectorInfo.SectorType = config.SectorTypeElseTypeNum
			}
		}

		address, err := utils.GetHostBadSectorAddress(minerId, val)
		if err != nil {
			log.Println("client.AddBadSector", fmt.Sprintf("Error in obtaining sector address! uuid:%s,sectorIs:%s", args.HostUUID, val), err.Error())
		} else {
			badSectorInfo.SectorAddress = address
		}

		if val != "" {
			badSectorInfo.AddType = 2
			_, err = global.OpToGatewayClient.AddBadSector(ctx, badSectorInfo)
			if err != nil {
				log.Println("client.AddBadSector", fmt.Sprintf("Error sector information recording failed! uuid:%s,sectorIs:%s", args.HostUUID, val), err.Error())
				continue
			}
		}
	}
	return &pb.String{Value: "Error sector processed successfully!"}, nil
}

// ScriptStop Script termination
func (p *OpServiceImpl) ScriptStop(ctx context.Context, args *pb.ScriptInfo) (*pb.String, error) {

	b, err := utils.ExecuteScript(fmt.Sprintf(define.PathScriptExecute+" %s", args.Script))
	return &pb.String{Value: b}, err
}

// NodeMountDisk Mount storage machine script
func (p *OpServiceImpl) NodeMountDisk(ctx context.Context, args *pb.MountDiskInfo) (*pb.String, error) {
	log.Println("op NodeMountDisk begin")
	b, err := utils.ExecuteScript(fmt.Sprintf("bash  -x "+define.PathScriptMountDisk+" %s %s", args.OpIP, args.OpDir))
	log.Println("NodeMountDisk result: ", b)
	if err != nil {
		return &pb.String{Value: b}, err
	}
	return &pb.String{Value: "Storage machine successfully mounted!"}, nil
}

// UninstallMountDisk Uninstalling storage machine scripts
func (p *OpServiceImpl) UninstallMountDisk(ctx context.Context, args *pb.MountDiskInfo) (*pb.String, error) {
	log.Println("op UninstallMountDisk begin")
	_, err := exec.Command("bash", "-c", "umount -lnf  /mnt/"+args.OpIP+"/disk*").CombinedOutput()
	if err != nil {
		log.Println("umount -lnf  /mnt/"+args.OpIP+"/disk*", err.Error())
		return &pb.String{}, err
	}
	return &pb.String{Value: "Storage machine uninstallation successful!"}, nil
}

// NodeAddShareDir Host adds shared file directory
func (p *OpServiceImpl) NodeAddShareDir(ctx context.Context, args *pb.MountDiskInfo) (*pb.String, error) {
	log.Println("op NodeAddShareDir begin")
	out, err := exec.Command("bash", "-c", "lsblk -f | grep - | grep /$ -v|grep  /boot -v  | awk '{print $NF}' | grep ^/ | sort -nr | uniq").CombinedOutput()
	if err != nil {
		log.Println("cmd `lsblk -f | grep - | grep /$ -v|grep  /boot -v  | awk '{print $NF}' | grep ^/ | sort -nr | uniq` Filed", err.Error())
		return &pb.String{}, err
	}

	var strs = []string{}
	for _, val := range strings.Split(string(out), "\n") {
		if val != "" {
			strs = append(strs, val)
		}
	}

	if len(strs) == 0 {
		return &pb.String{}, errors.New("no shared directory")
	}

	b, err := utils.ExecuteScript(define.PathDiskNFSSync)
	if err != nil {
		return &pb.String{Value: b}, err
	}

	log.Println("Storage machine initialization successful!")

	var dirStr string
	for key, str := range strs {
		if key == 0 {
			dirStr = str
		} else {
			dirStr = dirStr + "," + str
		}
	}

	return &pb.String{Value: dirStr}, nil
}

// OpReplacePlugFile OP obtains the Linux version of the machine and retrieves the corresponding file package based on the version
func (p *OpServiceImpl) OpReplacePlugFile(ctx context.Context, args *pb.OpReplaceFileInfo) (*pb.String, error) {
	log.Println("OpReplacePlugFile begin")
	log.Println("args", args)
	gatewayUrl := "http://" + global.ROOM_CONFIG.Gateway.IP + ":" + define.GatewayDownLoadPort + "/download/slot?filePath=" + define.PathIpfsProgram + "&filename="
	systemVersion, _ := utils.GetSystemVersion()
	systemInfoList := strings.Split(systemVersion, " ")
	if systemVersion == "" || len(systemInfoList) <= 1 {
		return &pb.String{}, errors.New("failed to obtain system version")
	}

	fileName := args.FileName + "-" + systemInfoList[0] + "-" + systemInfoList[1][:2]
	systemStr := systemInfoList[0] + " " + systemInfoList[1][:2]
	for _, val := range args.FileInfo {
		if val.System == systemStr {
			fileMd5 := utils.GetFileMD5(define.PathIpfsProgram + args.FileName)
			fmt.Println("fileMd5", fileMd5)
			if fileMd5 == val.FileMd5 {
				return &pb.String{Value: "Same file, no need to perform file replacement"}, nil
			} else {
				service.ServiceGroupApp.GatewayServiceGroup.DownloadAddressFile([]string{fileName}, define.PathIpfsProgram, gatewayUrl)
				os.Remove(define.PathIpfsProgram + args.FileName)
				os.Rename(define.PathIpfsProgram+fileName, define.PathIpfsProgram+args.FileName)
				return &pb.String{Value: "File replacement successful"}, nil
			}
		}
	}
	service.ServiceGroupApp.GatewayServiceGroup.DownloadAddressFile([]string{fileName}, define.PathIpfsProgram, gatewayUrl)
	err := os.Rename(define.PathIpfsProgram+fileName, define.PathIpfsProgram+args.FileName)
	if err != nil {
		return &pb.String{}, errors.New("failed to retrieve file")
	}
	return &pb.String{Value: "File replacement successful"}, nil
}

// OpReplacePlugFile1 OP obtains the Linux version of the machine and retrieves the corresponding file package based on the version
func (p *OpServiceImpl) OpReplacePlugFile1(ctx context.Context, args *pb.String) (*pb.String, error) {
	return &pb.String{Value: "File replacement successful"}, nil
}

// GetOpMonitorInfo Obtain Op hardware monitoring information
func (p *OpServiceImpl) GetOpMonitorInfo(ctx context.Context, args *pb.OpHardwareInfo) (*pb.MonitorInfo, error) {
	var opMonitorInfo pb.MonitorInfo
	opMonitorInfo.OpId = args.HostUUID

	out, err := exec.Command("bash", "-c", `df -h / | awk '{print $4,$5}' | awk 'NR>1'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `df -h / | awk '{print $4,$5}' | awk 'NR>1'` Filed", err.Error())
	} else {
		if len(string(out)) > 0 {
			strs := strings.Split(string(out)[:len(string(out))-1], " ")
			if len(strs) > 1 {
				opMonitorInfo.SysDiskLeave = strs[0]
				opMonitorInfo.SysDiskUseRate = strs[1]
			}
		}
	}

	if args.GetHostClassify() != config.HostDCStorageType && args.GetHostClassify() != config.HostStorageType {
		out, err = exec.Command("bash", "-c", `df -h `+define.MainDisk+` | awk '{print $4,$5}' | awk 'NR>1'`).CombinedOutput()
		if err != nil && !strings.Contains(err.Error(), "exit status") {
			log.Println("cmd `df -h "+define.MainDisk+" | awk '{print $4,$5}' | awk 'NR>1'` Filed", err.Error())
		} else {
			if len(string(out)) > 0 {
				strs := strings.Split(string(out)[:len(string(out))-1], " ")
				if len(strs) > 1 {
					opMonitorInfo.MainDiskLeave = strs[0]
					opMonitorInfo.MainDiskUseRate = strs[1]
				}
			}
		}

		out, err = exec.Command("bash", "-c", `mount  -l  |awk  -F ':' '{print $1}' | egrep  ^[0-9]+ |sort  -nr |uniq`).CombinedOutput()
		if err != nil && !strings.Contains(err.Error(), "exit status") {
			log.Println("cmd "+`mount  -l  |awk  -F ':' '{print $1}' |egrep  ^[0-9]+ |sort  -nr |uniq`+" Filed", err.Error())
			opMonitorInfo.Ips = []string{}
		} else {
			if len(string(out)) > 0 {
				opMonitorInfo.Ips = strings.Split(string(out)[:len(string(out))], "\n")
			} else {
				opMonitorInfo.Ips = []string{}
			}
		}

		ch := make(chan bool, 1)
		ctxCancel, cancel := context.WithCancel(context.Background())
		go func(ctxCancel context.Context, cancel context.CancelFunc) {
			out, err := exec.Command("bash", "-c", "timeout 6 nvidia-smi --query-gpu=index,name,memory.total,power.draw,power.limit --format=csv,noheader,nounits").CombinedOutput()
			if err != nil && !strings.Contains(err.Error(), "exit status") {
				log.Println("cmd `timeout 6 nvidia-smi --query-gpu=index,name,memory.total,power.draw,power.limit --format=csv,noheader,nounits` Filed", err.Error())
			}
			if !strings.Contains(string(out), "GeForce") {
				log.Println("cmd `timeout 6 nvidia-smi --query-gpu=index,name,memory.total,power.draw,power.limit --format=csv,noheader,nounits` unable to function properly!")
			} else {
				ch <- true
				close(ch)
			}
			cancel()
		}(ctxCancel, cancel)

		select {
		case <-ctxCancel.Done():
			fmt.Println("call successfully!!!")
			opMonitorInfo.GpuStatus = <-ch
		case <-time.After(time.Duration(time.Second * 8)):
			fmt.Println("timeout!!!")
			cancel()
		}
	}

	if args.GetHostClassify() != config.HostMinerType && args.GetHostClassify() != config.HostWorkerType {
		out, err := exec.Command("bash", "-c", "systemctl status nfs-server.service | grep Active'").CombinedOutput()
		if err != nil && !strings.Contains(err.Error(), "exit status") {
			log.Println("cmd `systemctl status nfs-server.service | grep Active", err.Error())
		}
		if strings.Contains(string(out), "active (exited)") {
			opMonitorInfo.NfsStatus = true
		} else {
			opMonitorInfo.NfsStatus = false
		}
	}

	opMonitorInfo.CpuUseRate = utils.GetCPUUseRate()

	diskIO, _ := utils.GetDiskIOStatus()
	if diskIO == config.HostDiskIOPass {
		opMonitorInfo.DiskStatus = true
	}

	switch args.GetHostClassify() {
	case config.HostMinerType:
		out, err = exec.Command("bash", "-c", `ss -lntip|grep `+define.LotusPort).CombinedOutput()
		if err != nil && !strings.Contains(err.Error(), "exit status") {
			log.Println("cmd "+`ss -lntip|grep `+define.LotusPort+" Filed", err.Error())
		} else {
			if len(string(out)) > 0 {
				opMonitorInfo.LotusStatus = true
			}
		}

		out, err = exec.Command("bash", "-c", `ss -lntip|grep `+define.MinerPort).CombinedOutput()
		if err != nil && !strings.Contains(err.Error(), "exit status") {
			log.Println("cmd "+`ss -lntip|grep `+define.MinerPort+" Filed", err.Error())
		} else {
			if len(string(out)) > 0 {
				opMonitorInfo.MinerStatus = true
			}
		}

		out, err = exec.Command("bash", "-c", `ss -lntip|grep `+define.BoostPort).CombinedOutput()
		if err != nil && !strings.Contains(err.Error(), "exit status") {
			log.Println("cmd "+`ss -lntip|grep `+define.BoostPort+" Filed", err.Error())
		} else {
			if len(string(out)) > 0 {
				opMonitorInfo.BoostStatus = true
			}
		}

		out, err = exec.Command("bash", "-c", define.PathIpfsProgram+`lotus sync status | grep 'Height diff:' | awk '{print $3}' | head -n 1`).CombinedOutput()
		if err != nil && !strings.Contains(err.Error(), "ERROR:") {
			log.Println("cmd "+define.PathIpfsProgram+`lotus sync status | grep 'Height diff:' | awk '{print $3}' | head -n 1`+" Filed", err.Error())
		} else {
			if len(string(out)) > 0 {
				height, _ := strconv.Atoi(string(out))
				if height <= config.LotusHeightNum {
					opMonitorInfo.LotusHeightStatus = true
				}
			}
		}
	case config.HostWorkerType:
		out, err = exec.Command("bash", "-c", `ss -lntip|grep `+define.WorkerPort).CombinedOutput()
		if err != nil && !strings.Contains(err.Error(), "exit status") {
			log.Println("cmd "+`ss -lntip|grep `+define.WorkerPort+" Filed", err.Error())
		} else {
			if len(string(out)) > 0 {
				opMonitorInfo.WorkerStatus = true
			}
		}
	}

	return &opMonitorInfo, nil
}

// GetOpScriptInfo Get Op script information
func (p *OpServiceImpl) GetOpScriptInfo(ctx context.Context, args *pb.OpScriptInfo) (*pb.OpScriptInfoResp, error) {
	var opMonitorInfo pb.OpScriptInfoResp
	out, err := exec.Command("bash", "-c", args.ScriptInfo).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd "+args.ScriptInfo+" Filed", err.Error())
		opMonitorInfo.ScriptResp = "Script execution failed!"
	} else {
		if len(string(out)) > 0 {
			opMonitorInfo.ScriptResp = string(out)
		} else {
			opMonitorInfo.ScriptResp = "The output result is empty!"
		}
	}
	return &opMonitorInfo, nil
}

// GetDiskLetter Obtain Op drive letter information
func (p *OpServiceImpl) GetDiskLetter(ctx context.Context, args *pb.DiskLetterReq) (*pb.OpScriptInfoResp, error) {
	var opScriptInfo pb.OpScriptInfoResp
	out, err := exec.Command("bash", "-c", `sudo smartctl -i `+args.DiskLetter+` | grep 'Serial Number' | awk '{print $3}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd "+`sudo smartctl -i `+args.DiskLetter+` | grep 'Serial Number' | awk '{print $3}'`+" Filed", err.Error())
		opScriptInfo.ScriptResp = "Script execution failed!"
	} else {
		if len(string(out)) > 0 {
			opScriptInfo.ScriptResp = string(out)[:len(string(out))-1]
		} else {
			opScriptInfo.ScriptResp = "The output result is empty!"
		}
	}
	return &opScriptInfo, nil
}

// GetOpMountInfo Get mounted disk information
func (p *OpServiceImpl) GetOpMountInfo(ctx context.Context, args *pb.DiskLetterReq) (*pb.OpMountDiskList, error) {
	var opScriptInfo pb.OpMountDiskList
	out, err := exec.Command("bash", "-c", `mount  -l  |awk  -F ':' '{print $1}' | egrep  ^[0-9]+ |sort  -nr |uniq`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd "+`mount  -l  |awk  -F ':' '{print $1}' |egrep  ^[0-9]+ |sort  -nr |uniq`+" Filed", err.Error())
		opScriptInfo.Ips = []string{}
	} else {
		if len(string(out)) > 0 {
			opScriptInfo.Ips = strings.Split(string(out)[:len(string(out))], "\n")
		} else {
			opScriptInfo.Ips = []string{}
		}
	}
	return &opScriptInfo, nil
}

// GetOpLogInfo Obtain corresponding log information for Op
func (p *OpServiceImpl) GetOpLogInfo(ctx context.Context, args *pb.OpLogInfoReq) (*pb.OpLogInfoResp, error) {
	var opMonitorInfo pb.OpLogInfoResp
	if args.LogType == config.WorkerLogType && args.HostClassify == config.HostWorkerType {
		opMonitorInfo.LogBeginNum, opMonitorInfo.LogResp = getHostLog(args, define.MainDisk+config.WorkerLogDir)
	} else if args.HostClassify == config.HostMinerType {
		switch args.LogType {
		case config.MinerProveLogType:
			opMonitorInfo.LogBeginNum, opMonitorInfo.LogResp = getHostLog(args, define.MainDisk+config.MinerLogDir)
		case config.MinerWdpostLogType:
			opMonitorInfo.LogBeginNum, opMonitorInfo.LogResp = getHostLog(args, define.MainDisk+config.MinerLogDir)
		case config.MinerRealTimeLogType:
			opMonitorInfo.LogBeginNum, opMonitorInfo.LogResp = getHostLog(args, define.MainDisk+config.MinerLogDir)
		case config.LotusLogType:
			opMonitorInfo.LogBeginNum, opMonitorInfo.LogResp = getHostLog(args, define.MainDisk+config.LotusLogDir)
		case config.BoostLogType:
			opMonitorInfo.LogBeginNum, opMonitorInfo.LogResp = getHostLog(args, define.MainDisk+config.BoostLogDir)
		}
	}
	return &opMonitorInfo, nil
}

// GetOpLogLen Obtain corresponding log information
func (p *OpServiceImpl) GetOpLogLen(ctx context.Context, args *pb.OpLogInfoReq) (*pb.OpLogLenResp, error) {
	var opLogLenResp pb.OpLogLenResp
	var logAddress string

	if args.LogType == config.WorkerLogType && args.HostClassify == config.HostWorkerType {
		logAddress = define.MainDisk + config.WorkerLogDir
	} else if args.HostClassify == config.HostMinerType {
		switch args.LogType {
		case config.MinerProveLogType:
			logAddress = define.MainDisk + config.MinerLogDir
		case config.MinerWdpostLogType:
			logAddress = define.MainDisk + config.MinerLogDir
		case config.MinerRealTimeLogType:
			logAddress = define.MainDisk + config.MinerLogDir
		case config.LotusLogType:
			logAddress = define.MainDisk + config.LotusLogDir
		case config.BoostLogType:
			logAddress = define.MainDisk + config.BoostLogDir
		}
	}

	if len(logAddress) == 0 {
		return &opLogLenResp, errors.New("log address acquisition failed")
	}

	logLen, err := exec.Command("bash", "-c", "wc -l "+logAddress+" | awk '{print $1}'").CombinedOutput()
	if err != nil {
		log.Println("cmd `wc -l "+logAddress+" | awk '{print $1}'` Filed", err.Error())
		return &opLogLenResp, err
	}
	if len(string(logLen)) > 2 {
		opLogLenResp.LogLenNum, _ = utils.StringToInt64(string(logLen)[:len(string(logLen))-1])
	}
	return &opLogLenResp, nil
}

// getHostLog Obtain corresponding log information
func getHostLog(logArgs *pb.OpLogInfoReq, logAddress string) (int64, string) {
	var logRowNum int64
	logBeginNum := logArgs.LogBeginNum
	getNum := logArgs.GetNum
	var logEndNum = logBeginNum

	logLen, err := exec.Command("bash", "-c", "wc -l "+logAddress+" | awk '{print $1}'").CombinedOutput()
	if err != nil {
		log.Println("cmd `wc -l "+logAddress+" | awk '{print $1}'` Filed", err.Error())
		return logEndNum, ""
	}
	if len(string(logLen)) < 0 {
		log.Println("get log info failed", err.Error())
		return logEndNum, ""
	}

	logRowNum, _ = utils.StringToInt64(string(logLen)[:len(string(logLen))-1])

	if logRowNum > logBeginNum {
		if logRowNum > logBeginNum+getNum {
			logEndNum = logBeginNum + getNum
		} else {
			logEndNum = logRowNum
		}
		logCmd := "cat " + logAddress + " | head -n " + utils.Int64ToString(logEndNum) + " | tail -n +" + utils.Int64ToString(logBeginNum+1)
		if logArgs.LogType == config.MinerProveLogType {
			logCmd = logCmd + " | grep -a post |grep -v "
		} else if logArgs.LogType == config.MinerWdpostLogType {
			logCmd = logCmd + " | -a  block "
		}

		out, err := exec.Command("bash", "-c", logCmd).CombinedOutput()
		if err != nil {
			log.Println("cmd `"+logCmd+"` Filed", err.Error())
			return logEndNum, ""
		}
		return logEndNum, string(out)
	}

	return logEndNum, ""
}

// GetNodeMinerInfo Obtain node miner information
func (p *OpServiceImpl) GetNodeMinerInfo(ctx context.Context, args *pb.OpHardwareInfo) (*pb.NodeMinerInfoResp, error) {
	var opMonitorInfo pb.NodeMinerInfoResp

	out, err := exec.Command("bash", "-c", `ss -lntip|grep `+define.MinerPort).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd "+`ss -lntip|grep `+define.LotusPort+" Filed", err.Error())
	} else {
		if len(string(out)) > 0 {
			opMonitorInfo.MinerProcessStatus = true
		}
	}

	out, err = exec.Command("bash", "-c", `cat `+define.MainDisk+`/ipfs/logs/miner.log|grep -a block |awk -F '[T|.]' '{print $2}' | uniq|awk  -F [:] '{print $1":"$2}'|uniq |wc -l`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd cat "+define.MainDisk+"/ipfs/logs/miner.log Filed", err.Error())
	} else {
		if len(string(out)) > 0 {
			num, _ := strconv.ParseInt(string(out)[:len(string(out))-1], 10, 64)
			opMonitorInfo.DailyExplosiveNum = num
		}
	}

	outStr, err := utils.ExecuteScript(fmt.Sprintf("bash -x "+define.PathIpfsScriptDeadlinesProven+" %s", define.MainDisk))
	if err != nil {
		log.Println(fmt.Sprintf("bash -x "+define.PathIpfsScriptDeadlinesProven+" %s Filed", define.MainDisk), err.Error())
		return &opMonitorInfo, err
	} else {
		if len(string(out)) > 0 {
			if outStr == "true" {
				opMonitorInfo.MessageOut = true
				opMonitorInfo.MessageOut = true
			} else {
				opMonitorInfo.MessageOut = false
				opMonitorInfo.MessageOut = false
			}
		}
	}

	return &opMonitorInfo, nil
}
