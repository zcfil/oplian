package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-address"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"log"
	"net"
	"oplian/config"
	"oplian/define"
	"oplian/global"
	"oplian/model/lotus"
	modelSystem "oplian/model/system"
	models "oplian/model/system"
	"oplian/service"
	"oplian/service/lotus/deploy"
	"oplian/service/pb"
	"oplian/service/system"
	"oplian/utils"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GateWayServiceImpl struct{}

var OpConnMap = make(map[string]*grpc.ClientConn)
var opConnnLock sync.RWMutex

func (g *GateWayServiceImpl) OpOnline(ctx context.Context, args *pb.String) (*pb.Bool, error) {
	// Link Node Host
	_, dis := global.OpClinets.GetOpClient(args.Value)
	return &pb.Bool{Value: !dis}, nil
}

// Connect op
func (g *GateWayServiceImpl) OpConnect(ctx context.Context, args *pb.RequestConnect) (*pb.String, error) {
	p, _ := peer.FromContext(ctx)
	IP, _, err := net.SplitHostPort(strings.TrimSpace(p.Addr.String()))
	if err != nil {
		return nil, err
	}
	opConnnLock.RLock()
	opConn := OpConnMap[IP+":"+args.Port]
	opConnnLock.RUnlock()
	if opConn == nil {
		opConn, err = grpc.Dial(IP+":"+args.Port, grpc.WithInsecure())
		opConnnLock.Lock()
		OpConnMap[IP+":"+args.Port] = opConn
		opConnnLock.Unlock()
		if err != nil {
			return nil, err
		}
		// Save OP Connection
		client := pb.NewOpServiceClient(opConn)
		sClient := pb.NewSlotGateServiceClient(opConn)

		// 查询op信息
		hostInfo := system.HostRecordService{}
		info, err := hostInfo.GetSysHostRecordByIPAndGatewayId(IP, global.GateWayID.String())
		if err != nil || info.ID == 0 {

			global.OpClinets.SetOpClient(args.OpId, &global.OpInfo{
				Clinet:     client,
				SlotClient: sClient,
				Ip:         IP,
				Port:       args.Port,
				OpId:       args.OpId,
			})

		} else {

			global.OpClinets.SetOpClient(info.UUID, &global.OpInfo{
				Clinet:     client,
				SlotClient: sClient,
				Ip:         IP,
				Port:       args.Port,
				OpId:       info.UUID,
			})
		}
	}
	return &pb.String{Value: global.GateWayID.String()}, nil
}

// GateWayConnect Connect Gateway
func (g *GateWayServiceImpl) GateWayConnect(ctx context.Context, port string) (pb.GateServiceClient, error) {
	p, _ := peer.FromContext(ctx)
	IP, _, err := net.SplitHostPort(strings.TrimSpace(p.Addr.String()))
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(IP+":"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return pb.NewGateServiceClient(conn), nil
}

// Oplian end heartbeat check call
func (g *GateWayServiceImpl) OplianHeartbeat(ctx context.Context, args *pb.String) (*emptypb.Empty, error) {
	log.Println("Oplian heartbeat check succeeded !")
	return &emptypb.Empty{}, nil
}

// OP end heartbeat check call
func (g *GateWayServiceImpl) OpHeartbeat(ctx context.Context, args *pb.String) (*emptypb.Empty, error) {
	// Web backend request
	if args.Value == "" {
		return nil, nil
	}
	client, dis := global.OpClinets.GetOpClient(args.Value)
	if dis {
		return nil, errors.New("op not online：" + args.Value)
	}
	// Get the number of tasks
	conf, err := lotusService.GetWorkerConfig(args.Value)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			global.ZC_LOG.Error("Failed to obtain worker configuration：" + err.Error())
		}
	} else {
		// Set the number of OP worker tasks
		_, err = client.SetPreNumber(ctx, &pb.PreNumber{P1: int32(conf.PreCount1), P2: int32(conf.PreCount2)})
		if err != nil {
			global.ZC_LOG.Error("Failed to set the number of tasks：" + err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

// StrategyProcess Strategy alarm
func (g *GateWayServiceImpl) StrategyProcess(ctx context.Context, args *pb.StrategyInfo) (*pb.ResponseMsg, error) {
	log.Println("GateWay StrategyProcess succeeded gateway!")

	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"
	client, b := global.OpClinets.GetOpClient(args.OpId)
	if b || client == nil {
		res.Code = -1
		res.Msg = "StrategyProcess grpc client is nil"
		return res, errors.New("StrategyProcess grpc client is nil")
	}
	res, _ = client.StrategyProcess(ctx, &pb.StrategyInfo{OpId: args.OpId, RoomId: args.RoomId, StrategiesId: args.StrategiesId})

	return res, nil
}

// AddHostRecord Add host information
func (g *GateWayServiceImpl) AddHostRecord(ctx context.Context, args *pb.HostInfo) (*emptypb.Empty, error) {
	log.Println("AddHostRecord succeeded gateway! ")

	sysHost := modelSystem.SysHostRecord{
		ZC_MODEL:         global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		HostName:         args.HostName,
		IntranetIP:       args.IntranetIP,
		InternetIP:       args.InternetIP,
		UUID:             args.UUID,
		DeviceSN:         args.DeviceSN,
		HostManufacturer: args.HostManufacturer,
		HostModel:        args.HostModel,
		OperatingSystem:  args.OperatingSystem,
		CPUCoreNum:       int(args.CPUCoreNum),
		CPUModel:         args.CPUModel,
		MemorySize:       int(args.MemorySize),
		DiskNum:          int(args.DiskNum),
		DiskSize:         utils.Decimal(float64(args.DiskSize), 1),
		HostShelfLife:    0,
		HostType:         0,
		HostClassify:     0,
		ServerDNS:        args.ServerDNS,
		SubnetMask:       args.SubnetMask,
		Gateway:          args.Gateway,
		HostGroupId:      0,
		RoomId:           "",
		MonitorTime:      0,
		GatewayId:        args.GatewayId,
		GPUNum:           int(args.GPUNum),
		AssetNumber:      "",
		SystemVersion:    args.SystemVersion,
		SystemBits:       int(args.SystemBits),
		RoomName:         "",
		IsGroupArray:     args.IsGroupArray,
		NetOccupyTime:    0,
	}

	hostInfo := system.HostRecordService{}
	info, err := hostInfo.GetSysHostRecord(sysHost.UUID)
	if err != nil || info.ID == 0 {
		if args.HostClassify != 0 {
			sysHost.HostClassify = int(args.HostClassify)
		}
		global.ZC_LOG.Error("This node has not been created or recorded!", zap.Error(err))
		err = hostInfo.CreateSysHostRecord(sysHost)
		if err != nil {
			global.ZC_LOG.Error("Creation failed!", zap.Error(err))
			return nil, err
		}
		log.Println("Created successfully!")
	} else {
		err = hostInfo.UpdateSysHostRecordAuto(&sysHost)
		if err != nil {
			global.ZC_LOG.Error("Update failed!", zap.Error(err))
			return nil, err
		}
		if args.HostClassify != 0 {
			sysHost.HostClassify = int(args.HostClassify)
			err = hostInfo.UpdateSysHostRecordClassify(&sysHost)
			if err != nil {
				global.ZC_LOG.Error("Update failed!", zap.Error(err))
				return nil, err
			}
		}
		log.Println("Update success!")
	}

	roomInfo := system.MachineRoomRecordService{}
	sysRoom, err := roomInfo.GetRoomByGatewayId(args.GatewayId)
	if err != nil {
		global.ZC_LOG.Error("Machine room information query failed!", zap.Error(err))
		return nil, err
	}

	err = hostInfo.HostBindRoomByGatewayId(&sysRoom)
	if err != nil {
		global.ZC_LOG.Error("Host binding to the data center failed!", zap.Error(err))
	}

	return &emptypb.Empty{}, nil
}

func (g *GateWayServiceImpl) GetWorkerTaskList(ctx context.Context, args *pb.RequestWorkerId) (*pb.TaskList, error) {
	client, dis := global.OpClinets.GetOpClient(args.WorkerId)
	if dis {
		return nil, nil
	}
	return client.WorkerTaksRunList(ctx, &pb.String{Value: args.MinerId})
}

func (g *GateWayServiceImpl) GetWorkerList(ctx context.Context, args *pb.RequestMinerId) (*pb.WorkerList, error) {
	global.OpClinets.LockRW.RLock()
	defer global.OpClinets.LockRW.RUnlock()
	opCount := len(global.OpClinets.Info)
	//var onlines = make([]*pb.WorkerList, opCount)
	var on sync.WaitGroup
	on.Add(opCount)
	on.Wait()

	return nil, nil
}

func (g *GateWayServiceImpl) DownloadSnapshot(ctx context.Context, args *pb.Downtown) (*pb.DownloadInfo, error) {
	return service.ServiceGroupApp.GatewayServiceGroup.DowloadFile(args.Url, args.Path, "")
}

// ExecuteScript Execute script
func (p *GateWayServiceImpl) ExecuteScript(ctx context.Context, args *pb.ScriptInfo) (*pb.String, error) {

	client, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return &pb.String{}, errors.New("opClient Connection failed:" + args.OpId)
	}

	return client.ExecuteScript(ctx, args)
}

// FileDistribution File Distribution
func (p *GateWayServiceImpl) FileDistribution(ctx context.Context, args *pb.ScriptInfo) (*pb.String, error) {

	f := &pb.FileInfo{
		Path:     args.Path,
		FileData: args.FileData,
		FileName: args.FileName,
	}

	res := &pb.String{}

	client, _ := global.OpClinets.GetOpClient(args.OpId)
	if client == nil {
		return res, errors.New("opClient Connection failed:" + args.OpId)
	}
	_, err := client.FileDistribution(ctx, f)
	if err != nil {
		return res, err
	}

	return res, nil
}

// FileSynOpHost File Distribution
func (p *GateWayServiceImpl) FileSynOpHost(ctx context.Context, args *pb.FileSynOp) (*pb.String, error) {
	return &pb.String{}, system.JobPlatformServiceApp.ExecuteFileSynOp(args)
}

// RunStopService Start stop service
func (g *GateWayServiceImpl) RunStopService(ctx context.Context, args *pb.RunStop) (*emptypb.Empty, error) {
	oclient, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return nil, errors.New("Op：" + args.OpId + " not online！")
	}
	switch define.ServiceType(args.ServiceType) {
	case define.ServiceLotus:
		_, err := oclient.RunAndStopService(ctx, &pb.RunStopType{ServiceType: args.ServiceType, IsRun: args.IsRun})
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			return nil, err
		}
		status := define.RunStatusStop
		if args.IsRun {
			status = define.RunStatusRunning
		}
		return nil, lotusService.UpdateLotusStatus(uint(args.Id), status.Int())
	case define.ServiceMiner:
		if !args.IsRun || args.IsRun && args.LinkId == 0 {
			_, err := oclient.RunAndStopService(ctx, &pb.RunStopType{ServiceType: args.ServiceType, IsRun: args.IsRun})
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				return nil, err
			}
			return nil, lotusService.UpdateMinerStatusAndLink(uint(args.Id), define.RunStatusBool(args.IsRun).Int(), args.LinkId)
		}

		miner, err := lotusService.GetMinerRun(args.Id)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			return nil, err
		}

		var lotusInfo lotus.LotusInfo
		lotusInfo, err = lotusService.DeployService.GetLotus(args.LinkId)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			return nil, err
		}
		param := &pb.MinerRun{
			Ip:         miner.Ip,
			Actor:      miner.Actor,
			IsManage:   miner.IsManage,
			IsWdpost:   miner.IsWdpost,
			IsWnpost:   miner.IsWnpost,
			Partitions: miner.Partitions,
			LotusToken: lotusInfo.Token,
			LotusIp:    lotusInfo.Ip,
		}
		_, err = oclient.RunMiner(ctx, param)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			return nil, err
		}
		return &emptypb.Empty{}, lotusService.UpdateMinerStatusAndLink(uint(args.Id), define.RunStatusRunning.Int(), args.LinkId)
	case define.ServiceWorkerTask, define.ServiceWorkerStorage:
		if !args.IsRun || args.IsRun && args.LinkId == 0 {

			_, err := oclient.RunAndStopService(ctx, &pb.RunStopType{ServiceType: args.ServiceType, IsRun: args.IsRun})
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				return nil, err
			}
			return &emptypb.Empty{}, lotusService.DeployService.UpdateWorkerStatusAndLink(uint(args.Id), define.RunStatusBool(args.IsRun).Int(), 0, args.LinkId)
		}

		miner, err := lotusService.DeployService.GetMiner(args.LinkId)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			return nil, err
		}
		log.Println("Replace the miner:", args.LinkId, miner.Token, miner.Ip)

		param := define.TaskWorker.Int()
		if define.ServiceType(args.ServiceType) == define.ServiceWorkerStorage {
			param = define.StorageWorker.Int()
		}
		_, err = oclient.RunWorker(ctx, &pb.FilParam{Token: miner.Token, Ip: miner.Ip, Param: strconv.Itoa(param)})
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			return nil, err
		}
		return &emptypb.Empty{}, lotusService.DeployService.UpdateWorkerStatusAndLink(uint(args.Id), define.RunStatusRunning.Int(), 0, args.LinkId)
	}

	return &emptypb.Empty{}, nil
}

// FileOpSynGateWay OP synchronizes files to GateWay
func (p *GateWayServiceImpl) FileOpSynGateWay(ctx context.Context, args *pb.AddFileInfo) (*pb.ResponseMsg, error) {

	log.Println("gateway FileOpSynGateWay success")
	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"

	if args.AddType == define.AddOnline.Int64() {
		err := system.JobPlatformServiceApp.OnlineDownload(args)
		if err != nil {
			res.Code = -1
			res.Msg = err.Error()
			return res, nil
		}
	} else {

		id, err := system.JobPlatformServiceApp.FileNodeCopy(args)
		if err != nil {
			res.Code = -1
			res.Msg = err.Error()
			return res, nil
		}
		args.Id = strconv.Itoa(id)
		client, _ := global.OpClinets.GetOpClient(args.OpId)
		if client == nil {
			res.Code = -1
			res.Msg = "opClient Connection failed:" + args.OpId
			return res, nil
		}

		res, err = client.FileOpSynGateWay(ctx, args)
		if err != nil {
			res.Code = -1
			res.Msg = "FileOpSynGateWay FileOpSynGateWay err:" + args.OpId
			return res, nil
		}

	}

	return res, nil
}

// AddGateWayFile Add GateWay file
func (p *GateWayServiceImpl) AddGateWayFile(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"

	dirAr := strings.Split(define.MainDisk, "/")
	dir := fmt.Sprintf("/%s/", dirAr[1])
	if !strings.Contains(args.Path, dir) {
		args.Path = define.MainDisk + args.Path
	}

	_, err := system.JobPlatformServiceApp.CreateFile(args)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}

	return res, nil
}

// SysFilePoint OP point-to-point replication
func (p *GateWayServiceImpl) SysFilePoint(ctx context.Context, args *pb.SynFileInfo) (*pb.ResponseMsg, error) {

	log.Println("gateway SysFilePoint success")
	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"

	client, _ := global.OpClinets.GetOpClient(args.OpId)
	if client == nil {
		res.Code = -1
		res.Msg = "opClient Connection failed"
		return res, nil
	}

	res, _ = client.SysFileFrom(ctx, args)

	return res, nil
}

// DownLoadFiles File Download
func (p *GateWayServiceImpl) DownLoadFiles(ctx context.Context, args *pb.DownLoadInfo) (*pb.ResponseMsg, error) {

	log.Println("gateway DownLoadFiles success")
	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"

	err := system.JobPlatformServiceApp.ReadGateWayFile(args)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}

	return res, nil
}

// GatewayFileExist Does the file exist
func (g *GateWayServiceImpl) GatewayFileExist(ctx context.Context, args *pb.String) (*pb.Bool, error) {
	return &pb.Bool{Value: utils.FileExist(args.Value)}, nil
}

// AddHostMonitorRecord Add host monitoring information
func (g *GateWayServiceImpl) AddHostMonitorRecord(ctx context.Context, args *pb.HostMonitorInfo) (*emptypb.Empty, error) {
	//log.Println("AddHostMonitorRecord succeeded gateway! ")

	gpuInfo := &[]utils.GPUMonitorInfo{}
	err := json.Unmarshal([]byte(args.GPUUseInfo), gpuInfo)
	if err != nil {
		log.Println("Failed to parse host GPU information,hostUUID:", args.HostUUID)
		return &emptypb.Empty{}, err
	}

	if len(*gpuInfo) > 0 {
		for _, val := range *gpuInfo {
			sysHost := modelSystem.SysHostMonitorRecord{}
			if val.ID == "0" {
				sysHost = modelSystem.SysHostMonitorRecord{
					ZC_MODEL:       global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
					HostUUID:       args.HostUUID,
					CPUUseRate:     args.CPUUseRate,
					DiskUseRate:    args.DiskUseRate,
					MemoryUseRate:  args.MemoryUseRate,
					GPUUseRate:     val.UseRate,
					GPUID:          val.ID,
					CPUTemperature: args.CPUTemperature,
					DiskSize:       args.DiskSize,
					DiskUseSize:    args.DiskUseSize,
					MemorySize:     args.MemorySize,
					MemoryUseSize:  args.MemoryUseSize,
				}
			} else {
				sysHost = modelSystem.SysHostMonitorRecord{
					ZC_MODEL:   global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
					HostUUID:   args.HostUUID,
					GPUUseRate: val.UseRate,
					GPUID:      val.ID,
				}
			}
			hostMonitorInfo := system.HostMonitorRecordService{}
			err := hostMonitorInfo.CreateSysHostMonitorRecord(sysHost)
			if err != nil {
				global.ZC_LOG.Error("Creation failed!", zap.Error(err))
				continue
			}
		}
	} else {
		sysHost := modelSystem.SysHostMonitorRecord{
			ZC_MODEL:       global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
			HostUUID:       args.HostUUID,
			CPUUseRate:     args.CPUUseRate,
			DiskUseRate:    args.DiskUseRate,
			MemoryUseRate:  args.MemoryUseRate,
			GPUUseRate:     0,
			GPUID:          "0",
			CPUTemperature: args.CPUTemperature,
			DiskSize:       args.DiskSize,
			DiskUseSize:    args.DiskUseSize,
			MemorySize:     args.MemorySize,
			MemoryUseSize:  args.MemoryUseSize,
		}
		hostMonitorInfo := system.HostMonitorRecordService{}
		err := hostMonitorInfo.CreateSysHostMonitorRecord(sysHost)
		if err != nil {
			global.ZC_LOG.Error("Creation failed!", zap.Error(err))
			return &emptypb.Empty{}, err
		}
	}

	sysHost := modelSystem.SysHostRecord{
		UUID:        args.HostUUID,
		MonitorTime: int(time.Now().Unix()),
	}
	hostInfo := system.HostRecordService{}
	err = hostInfo.UpdateSysHostRecordMonitorTime(&sysHost)
	if err != nil {
		global.ZC_LOG.Error("Failed to update host monitoring information time!", zap.Error(err))
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UnZipSynFile Unzip the file and delete it
func (p *GateWayServiceImpl) UnZipSynFile(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	log.Println("op UnZipSynFile success")
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

// DelGateWayFile delete file
func (p *GateWayServiceImpl) DelGateWayFile(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}

	dirAr := strings.Split(define.MainDisk, "/")
	dir := fmt.Sprintf("/%s/", dirAr[1])
	if !strings.Contains(args.FileName, dir) {
		args.FileName = define.MainDisk + args.FileName
	}

	err := utils.DelFile(args.FileName)
	if err != nil {
		return res, err
	}

	return res, nil
}

// OpInformationTest Gateway calls Op to perform host information testing
func (p *GateWayServiceImpl) OpInformationTest(ctx context.Context, args *pb.HostTestInfo) (*pb.String, error) {
	log.Println("Start the host information test")
	hostUUIDs := strings.Split(args.HostUUIDs, ",")
	if args.IsAddPower {
		for _, val := range hostUUIDs {
			client1, _ := global.OpClinets.GetOpClient(val)
			if client1 == nil {
				log.Println("opClient Connection failed:" + args.HostUUID)
				return nil, errors.New("opClient Connection failed")
			}
			go client1.OpServerPortControl(ctx, &pb.String{Value: val})
		}
	} else {
		client1, _ := global.OpClinets.GetOpClient(hostUUIDs[0])
		if client1 == nil {
			log.Println("opClient Connection failed:" + args.HostUUID)
			return nil, errors.New("opClient Connection failed")
		}
		go client1.OpServerPortControl(ctx, &pb.String{Value: hostUUIDs[0]})
	}

	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.OpInformationTest(ctx, args)
}

// LotusHeight Lotus height
func (p *GateWayServiceImpl) LotusHeight(ctx context.Context, args *pb.RequestOp) (*pb.LotusHeightInfo, error) {

	log.Println("gateway LotusHeight success")
	client, b := global.OpClinets.GetOpClient(args.OpId)
	if b {
		log.Println("opClient Connection failed:" + args.OpId)
		return nil, errors.New("opClient Connection failed:" + args.OpId)
	}

	h, err := client.LotusHeight(ctx, args)
	if err != nil {
		log.Println("client.LotusHeight:" + err.Error())
		return nil, err
	}

	return &pb.LotusHeightInfo{Height: h.Height}, nil
}

// UpdateHostTestRecord OP modifies the added test information (to modify the corresponding test item results into the corresponding test information)
func (g *GateWayServiceImpl) UpdateHostTestRecord(ctx context.Context, args *pb.UpdateHostTestInfo) (*emptypb.Empty, error) {
	log.Println(fmt.Sprintf("Update Host Test Record! HostUUID: %s", args.HostUUID))

	sysHostTest := modelSystem.SysHostTestRecord{
		HostUUID:         args.HostUUID,
		TestBeginAt:      args.TestBeginAt,
		TestType:         args.TestType,
		IsAddPower:       args.IsAddPower,
		SelectHostUUIDs:  args.SelectHostUUIDs,
		SelectHostIPs:    args.SelectHostIPs,
		TestEndAt:        time.Now().Unix(),
		TestResult:       args.TestResult,
		CPUHardInfo:      args.CPUHardInfo,
		CPUHardScore:     args.CPUHardScore,
		GPUHardInfo:      args.GPUHardInfo,
		GPUHardScore:     args.GPUHardScore,
		MemoryHardInfo:   args.MemoryHardInfo,
		MemoryHardScore:  args.MemoryHardScore,
		DiskHardInfo:     args.DiskHardInfo,
		DiskHardScore:    args.DiskHardScore,
		NetTestInfo:      args.NetTestInfo,
		NetTestScore:     args.NetTestScore,
		GPUTestInfo:      args.GPUTestInfo,
		GPUTestScore:     args.GPUTestScore,
		DiskIO:           args.DiskIO,
		DiskAllRate:      args.DiskAllRate,
		DiskAllRateScore: args.DiskAllRateScore,
		DiskSSDRate:      args.DiskSSDRate,
		DiskSSDRateScore: args.DiskSSDRateScore,
	}
	if args.TestResult == config.HostUnderTest {
		sysHostTest.TestEndAt = 0
	}
	hostTestInfo := system.HostTestRecordService{}
	err := hostTestInfo.UpdateSysHostTestRecord(&sysHostTest)
	if err != nil {
		global.ZC_LOG.Error("OP failed to modify the added test information!", zap.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

// UpdateHostNetOccupyTime OP modifies the time when this host is used to test the networks of other hosts
func (g *GateWayServiceImpl) UpdateHostNetOccupyTime(ctx context.Context, args *pb.UpdateHostNet) (*emptypb.Empty, error) {
	hostUUIDs := strings.Split(args.HostUUIDs, ",")
	for _, val := range hostUUIDs {
		sysHostTest := modelSystem.SysHostRecord{
			UUID:          val,
			NetOccupyTime: time.Now().Unix(),
		}
		hostTestInfo := system.HostRecordService{}
		err := hostTestInfo.UpdateHostNetOccupyTime(&sysHostTest)
		if err != nil {
			global.ZC_LOG.Error("Modify the time when this host is used to test the networks of other hosts!", zap.Error(err))
			continue
		}
		client, _ := global.OpClinets.GetOpClient(val)
		if client == nil {
			log.Println("opClient Connection failed:" + val)
			return nil, errors.New("opClient Connection failed")
		}

		client.CloseOpServerPortControl(ctx, &pb.String{Value: "Turn off network port monitoring"})
	}
	return &emptypb.Empty{}, nil
}

// OpInformationPatrol Gateway retrieves Op to perform host information patrol
func (p *GateWayServiceImpl) OpInformationPatrol(ctx context.Context, args *pb.HostPatrolInfo) (*pb.String, error) {
	log.Println("Start the host information patrol")
	sysHostPatrol := modelSystem.SysHostPatrolRecord{
		ZC_MODEL:      global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		HostUUID:      args.HostUUID,
		PatrolType:    args.HostClassify,
		PatrolBeginAt: time.Now().Unix(),
		PatrolEndAt:   0,
		PatrolResult:  config.HostUnderTest,
		PatrolMode:    args.PatrolMode,
	}
	hostPatrolInfo := system.HostPatrolRecordService{}
	err := hostPatrolInfo.CreateSysHostPatrolRecord(sysHostPatrol)
	if err != nil {
		global.ZC_LOG.Error("Creation failed!", zap.Error(err))
		return &pb.String{Value: "Creation failed!"}, err
	}

	args.PatrolBeginAt = sysHostPatrol.PatrolBeginAt
	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.OpInformationPatrol(ctx, args)
}

// UpdateHostPatrolRecord OP modifies the added inspection information (to modify the corresponding test item results into the corresponding inspection information)
func (g *GateWayServiceImpl) UpdateHostPatrolRecord(ctx context.Context, args *pb.UpdateHostPatrolInfo) (*emptypb.Empty, error) {
	sysHostPatrol := modelSystem.SysHostPatrolRecord{
		HostUUID:                    args.HostUUID,
		PatrolBeginAt:               args.PatrolBeginAt,
		PatrolEndAt:                 time.Now().Unix(),
		PatrolResult:                args.PatrolResult,
		DiskIO:                      args.DiskIO,
		DiskIODuration:              args.DiskIODuration,
		HostIsDown:                  args.HostIsDown,
		HostIsDownDuration:          args.HostIsDownDuration,
		HostNetStatus:               args.HostNetStatus,
		HostNetDuration:             args.HostNetDuration,
		LogInfoStatus:               args.LogInfoStatus,
		LogInfoDuration:             args.LogInfoDuration,
		LogOvertimeStatus:           args.LogOvertimeStatus,
		LogOvertimeDuration:         args.LogOvertimeDuration,
		WalletBalanceStatus:         args.WalletBalanceStatus,
		WalletBalance:               float64(args.WalletBalance),
		WalletBalanceDuration:       args.WalletBalanceDuration,
		LotusSyncStatus:             args.LotusSyncStatus,
		LotusSyncDuration:           args.LotusSyncDuration,
		GPUDriveStatus:              args.GPUDriveStatus,
		GPUDriveDuration:            args.GPUDriveDuration,
		PackageVersionStatus:        args.PackageVersionStatus,
		PackageVersion:              args.PackageVersion,
		PackageVersionDuration:      args.PackageVersionDuration,
		DataCatalogStatus:           args.DataCatalogStatus,
		DataCatalogDuration:         args.DataCatalogDuration,
		EnvironmentVariableStatus:   args.EnvironmentVariableStatus,
		EnvironmentVariableDuration: args.EnvironmentVariableDuration,
		BlockLogStatus:              args.BlockLogStatus,
		BlockLogDuration:            args.BlockLogDuration,
		TimeSyncStatus:              args.TimeSyncStatus,
		TimeSyncDuration:            args.TimeSyncDuration,
		PingNetStatus:               args.PingNetStatus,
		PingNetDuration:             args.PingNetDuration,
	}
	patrolPatrolInfo := system.HostPatrolRecordService{}
	err := patrolPatrolInfo.UpdateSysHostPatrolRecord(&sysHostPatrol)
	if err != nil {
		global.ZC_LOG.Error("OP failed to modify the added inspection information!", zap.Error(err))
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// CloseOpInformationTest Gateway calls Op to shut down host information testing
func (p *GateWayServiceImpl) CloseOpInformationTest(ctx context.Context, args *pb.CloseHostTest) (*pb.String, error) {
	log.Println("Start close the host information test, args:", args)

	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}
	_, err := client.CloseOpInformationTest(ctx, args)
	if err != nil {
		log.Println("close op information test failed:" + args.HostUUID)
		return nil, errors.New("close op information test failed")
	}

	hostUUIDs := strings.Split(args.SelectHostUUIDs, ",")
	log.Println("hostUUIDs:" + args.SelectHostUUIDs)
	for _, val := range hostUUIDs {

		sysHostTestRecord := modelSystem.SysHostRecord{
			UUID:          val,
			NetOccupyTime: time.Now().Unix(),
		}

		hostRecord := system.HostRecordService{}
		err = hostRecord.UpdateHostNetOccupyTime(&sysHostTestRecord)
		if err != nil {
			global.ZC_LOG.Error("Failed to modify host information! "+val, zap.Error(err))
		}

		client1, _ := global.OpClinets.GetOpClient(val)
		if client1 == nil {
			log.Println("opClient1 Connection failed:" + val)
			continue
		}
		_, err := client1.CloseOpServerPortControl(ctx, &pb.String{})
		if err != nil {
			log.Println("close op server port control failed:" + val)
			continue
		}
	}

	sysHostTest := modelSystem.SysHostTestRecord{
		ZC_MODEL: global.ZC_MODEL{
			ID: uint(args.ID),
		},
		HostUUID:   args.HostUUID,
		TestResult: config.HostTestFailed,
	}
	hostTestInfo := system.HostTestRecordService{}
	err = hostTestInfo.UpdateSysHostTestRecordClose(&sysHostTest)
	if err != nil {
		global.ZC_LOG.Error("Failed to close host testing information!", zap.Error(err))
		return &pb.String{Value: "Failed to close host testing information!"}, err
	}

	return &pb.String{}, nil
}

// RestartAddHostTest Restart host information testing
func (p *GateWayServiceImpl) RestartAddHostTest(ctx context.Context, args *pb.RestartHostTest) (*pb.String, error) {
	log.Println("Restart the host information test")
	sysHostTest := modelSystem.SysHostTestRecord{
		ZC_MODEL: global.ZC_MODEL{
			ID: uint(args.ID),
		},
		HostUUID:    args.HostUUID,
		TestBeginAt: time.Now().Unix(),
		TestEndAt:   0,
		TestResult:  config.HostUnderTest,
	}
	hostTestInfo := system.HostTestRecordService{}
	err := hostTestInfo.RestartUpdateSysHostTestRecord(&sysHostTest)
	if err != nil {
		global.ZC_LOG.Error("Failed to modify host test table!", zap.Error(err))
		return &pb.String{Value: "Failed to modify host test table!"}, err
	}

	args.TestBeginAt = sysHostTest.TestBeginAt
	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.OpInformationTestRestart(ctx, args)
}

// ModifyFileStatus Update file status
func (p *GateWayServiceImpl) ModifyFileStatus(ctx context.Context, args *pb.FileManage) (*pb.ResponseMsg, error) {

	log.Println("gateway ModifyFileStatus success")
	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"

	err := system.JobPlatformServiceApp.ModifyFileStatus(int(args.Id), int(args.FileSize))
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}

	return res, nil
}

// GetOpHardwareInfo Gateway retrieves Op to obtain host information
func (p *GateWayServiceImpl) GetOpHardwareInfo(ctx context.Context, args *pb.OpHardwareInfo) (*pb.String, error) {
	log.Println("get op hardware info")
	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.GetOpHardwareInfo(ctx, args)
}

// GetOpMonitorInfo Gateway retrieves Op to obtain host monitoring information
func (p *GateWayServiceImpl) GetOpMonitorInfo(ctx context.Context, args *pb.OpHardwareInfo) (*pb.MonitorInfo, error) {
	log.Println("get op monitor info")
	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.GetOpMonitorInfo(ctx, args)
}

// GetOpScriptInfo Gateway retrieves Op to obtain the execution result of the host script
func (p *GateWayServiceImpl) GetOpScriptInfo(ctx context.Context, args *pb.OpScriptInfo) (*pb.OpScriptInfoResp, error) {
	log.Println("get op script run info")
	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.GetOpScriptInfo(ctx, args)
}

// DiskReMounting Gateway retrieves Op and executes re mounting
func (p *GateWayServiceImpl) DiskReMounting(ctx context.Context, args *pb.DiskReMountReq) (*pb.String, error) {
	log.Println("get op script run info")
	opClient, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if opClient == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	var storageInfo lotus.LotusStorageInfo
	var err error
	storageInfo, err = lotusService.GetStorageByOpID(args.MountOpId)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return nil, errors.New("get storageInfo by actor failed")
	} else {

		if len(storageInfo.Ip) > 0 {
			if storageInfo.ColonyType != define.StorageTypeNFS {
				return nil, err
			}
			mountInfo := &pb.MountDiskInfo{OpIP: storageInfo.Ip}
			if _, err = opClient.UninstallMountDisk(context.Background(), mountInfo); err != nil {
				global.ZC_LOG.Info("uninstall nodeMountDisk error：" + err.Error())
			}

			opDir := utils.RestoreMountDir(storageInfo.NFSDisk, storageInfo.Ip)
			mountInfo = &pb.MountDiskInfo{OpIP: storageInfo.Ip, OpDir: opDir}
			if _, err = opClient.NodeMountDisk(context.Background(), mountInfo); err != nil {
				global.ZC_LOG.Info("NodeMountDisk error：" + err.Error())
				return nil, err
			}

			if args.HostClassify == config.HostMinerType {
				storageInfo := &pb.StorageInfo{MountDir: opDir, NodeIp: args.NodeIP, StorageIp: storageInfo.Ip}
				if _, err = opClient.AddNodeStorage(context.Background(), storageInfo); err != nil {
					global.ZC_LOG.Info("NodeMountDisk error：" + err.Error())
					return nil, err
				}
			}
		}
	}
	return &pb.String{}, nil
}

// GetOpLogInfo Gateway retrieves log information corresponding to Op
func (p *GateWayServiceImpl) GetOpLogInfo(ctx context.Context, args *pb.OpLogInfoReq) (*pb.OpLogInfoResp, error) {
	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.GetOpLogInfo(ctx, args)
}

// GetOpLogLen Gateway retrieves the number of log rows corresponding to Op
func (p *GateWayServiceImpl) GetOpLogLen(ctx context.Context, args *pb.OpLogInfoReq) (*pb.OpLogLenResp, error) {
	log.Println("get op log len")
	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.GetOpLogLen(ctx, args)
}

// GetDiskLetter Gateway retrieves Op to obtain drive letter information
func (p *GateWayServiceImpl) GetDiskLetter(ctx context.Context, args *pb.DiskLetterReq) (*pb.OpScriptInfoResp, error) {
	log.Println("get op disk letter info")
	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.GetDiskLetter(ctx, args)
}

// GetOpMountInfo Gateway retrieves Op to obtain mounted disk information
func (p *GateWayServiceImpl) GetOpMountInfo(ctx context.Context, args *pb.DiskLetterReq) (*pb.OpMountDiskList, error) {
	log.Println("get op disk mount info")
	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.GetOpMountInfo(ctx, args)
}

// DelOpFile delete file
func (p *GateWayServiceImpl) DelOpFile(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	log.Println("gateway DelOpFile success")
	res := &pb.ResponseMsg{}

	client, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		log.Println("opClient Connection failed:" + args.OpId)
		return res, errors.New("opClient Connection failed:" + args.OpId)
	}
	res, _ = client.DelOpFile(ctx, args)

	return res, nil
}

// GetHostGroupArray Obtain host disk group array information
func (g *GateWayServiceImpl) GetHostGroupArray(ctx context.Context, args *pb.OpHostUUID) (*pb.HostGroupArray, error) {
	hostInfo := system.HostRecordService{}
	info, err := hostInfo.GetSysHostRecord(args.HostUUID)
	if err != nil {
		global.ZC_LOG.Error("This node has not been created or recorded!", zap.Error(err))
		return &pb.HostGroupArray{IsGroupArray: false}, err
	}

	return &pb.HostGroupArray{IsGroupArray: info.IsGroupArray}, nil
}

// GetHostInfoByIPAndGatewayId Get host information
func (g *GateWayServiceImpl) GetHostInfoByIPAndGatewayId(ctx context.Context, args *pb.RequestOp) (*pb.OpHostUUID, error) {
	// 获取对应主机信息
	hostInfo := system.HostRecordService{}
	info, err := hostInfo.GetSysHostRecordByIPAndGatewayId(args.GetIp(), global.GateWayID.String())
	if err != nil {
		global.ZC_LOG.Error("This node has not been created or recorded!", zap.Error(err))
		return &pb.OpHostUUID{HostUUID: ""}, err
	}

	return &pb.OpHostUUID{HostUUID: info.UUID}, nil
}

// RedoSectorsTask Redo sector task
func (g *GateWayServiceImpl) RedoSectorsTask(ctx context.Context, args *pb.String) (*pb.ResponseMsg, error) {

	log.Println("gateway RedoSectorsTask success")
	res := &pb.ResponseMsg{}

	return res, nil
}

// AddWarn Add alarm information (4 strategy alarms, 5 inspection alarms, 6 business alarms)
func (g *GateWayServiceImpl) AddWarn(ctx context.Context, args *pb.WarnInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	res.Code = 200
	res.Msg = "ok"
	data := &models.SysWarnManage{
		WarnName:   args.WarnName,
		WarnType:   int(args.WarnType),
		ComputerId: args.ComputerId,
		WarnInfo:   args.WarnInfo,
	}
	err := system.WarnManageServiceApp.AddWarnInfo(*data)
	if err != nil {
		return res, err
	}
	return res, nil
}

// AddBadSector Add error sectors
func (g *GateWayServiceImpl) AddBadSector(ctx context.Context, args *pb.BadSectorId) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}

	_, err := address.NewFromString(args.MinerId)
	if err != nil {
		log.Println(fmt.Sprintf("Abnormal node information :%s,%s", args.MinerId, err))
		return res, errors.New(fmt.Sprintf("Abnormal node information :%s,%s", args.MinerId, err))
	}

	data := &lotus.LotusSectorRecover{
		MinerId:       args.MinerId,
		SectorId:      int(args.SectorId),
		SectorType:    int(args.SectorType),
		SectorSize:    int(args.SectorSize),
		BelongingNode: args.BelongingNode,
		AbnormalTime:  time.Now(),
		SectorAddress: args.SectorAddress,
	}

	if args.AddType == 1 {

		count, err := deploy.SectorsRecoverServiceApi.GetBadSectorCount(*data)
		if err != nil {
			return res, err
		}

		if count == utils.TWO {
			err = deploy.SectorsRecoverServiceApi.AddBadSector(*data)
			if err != nil {
				return res, err
			}
		}
	} else if args.AddType == 2 {
		err = deploy.SectorsRecoverServiceApi.AddBadSector(*data)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

// OpFileToGateWay Synchronize host files to gateWay
func (g *GateWayServiceImpl) OpFileToGateWay(ctx context.Context, args *pb.AddFileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	client, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return nil, errors.New("opClient Connection failed:" + args.OpId)
	}

	_, err := client.OpFileToGateWay(ctx, args)
	if err != nil {
		return nil, errors.New("OpFileToGateWay err:" + err.Error())
	}

	return res, nil
}

// CheckOpPath Check the path of the op file
func (g *GateWayServiceImpl) CheckOpPath(ctx context.Context, args *pb.DirFileReq) (*pb.ResponseMsg, error) {

	log.Println("gateway CheckOpPath success")
	res := &pb.ResponseMsg{}

	client, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		log.Println("opClient Connection failed:" + args.OpId)
		return nil, errors.New("opClient Connection failed")
	}

	_, err := client.CheckOpPath(ctx, args)
	if err != nil {
		log.Println("OpFileToGateWay err:" + err.Error())
		return nil, errors.New("OpFileToGateWay err:" + err.Error())
	}

	return res, nil
}

// CarFilePath Get the car file path
func (g *GateWayServiceImpl) CarFilePath(ctx context.Context, args *pb.CarFile) (*pb.CarFile, error) {

	log.Println("gateway CarFilePath success")
	res := &pb.CarFile{}

	client, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		log.Println("opClient Connection failed:" + args.OpId)
		return nil, fmt.Errorf("opClient Connection failed:%s", args.OpId)
	}

	res, err := client.CarFilePath(ctx, args)
	if err != nil {
		return res, err
	}

	return res, nil
}

// HostType Get Host Type
func (g *GateWayServiceImpl) HostType(ctx context.Context, args *pb.String) (*pb.String, error) {
	res := &pb.String{}

	r, err := deploy.SectorsRecoverServiceApi.GetHostType(args.Value)
	if err != nil {
		return res, err
	}
	res.Value = r
	return res, nil
}

// CarFileList Get car file
func (g *GateWayServiceImpl) CarFileList(ctx context.Context, args *pb.SectorID) (*pb.CarArray, error) {

	log.Println("gateway CarFileList success")
	res := &pb.CarArray{}

	r, err := deploy.SectorsRecoverServiceApi.GetCarFileList(args)
	if err != nil {
		return res, err
	}

	for _, v := range r {
		res.CarInfo = append(res.CarInfo, &pb.CarInfo{FileName: fmt.Sprintf("%s.car", v.PieceCid), PieceCid: v.PieceCid, PieceSize: int64(v.PieceSize)})
	}

	return res, nil
}

// ModifyOnlineFile Update online files
func (g *GateWayServiceImpl) ModifyOnlineFile(ctx context.Context, args *pb.AddFileInfo) (*pb.ResponseMsg, error) {

	log.Println("gateway ModifyOnlineFile success")
	res := &pb.ResponseMsg{}

	err := system.JobPlatformServiceApp.ModifyOnlineFile(args)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetFileName Get file name
func (g *GateWayServiceImpl) GetFileName(ctx context.Context, args *pb.FileNameInfo) (*pb.String, error) {

	log.Println("gateway GetFileName success")
	res := &pb.String{}

	s, err := system.JobPlatformServiceApp.GetFileName(args)
	if err != nil {
		return res, err
	}

	return &pb.String{Value: s}, nil
}

// SetJobPlatformStop Set the switch for the homework platform program
func (g *GateWayServiceImpl) SetJobPlatformStop(ctx context.Context, args *pb.JobPlatform) (*pb.String, error) {

	log.Println("gateway SetJobPlatformStop success")
	system.JobPlatformServiceApp.IsStop = args.IsStop

	return &pb.String{}, nil
}

// ScriptStop Script termination
func (g *GateWayServiceImpl) ScriptStop(ctx context.Context, args *pb.ScriptInfo) (*pb.String, error) {

	log.Println("gateway ScriptStop success")

	client, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		global.ZC_LOG.Error("ScriptStop opClient Connection failed:" + args.OpId)
		return nil, errors.New("ScriptStop opClient Connection failed:" + args.OpId)
	}

	return client.ScriptStop(ctx, args)
}

// CheckOpIsOnline Check if the corresponding OP node is online
func (g *GateWayServiceImpl) CheckOpIsOnline(ctx context.Context, args *pb.OpHostUUID) (*pb.String, error) {
	log.Println("gateway CheckOpPath success")

	_, dis := global.OpClinets.GetOpClient(args.HostUUID)
	if dis {
		log.Println("CheckOpIsOnline opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("op not online：" + args.HostUUID)
	}

	return &pb.String{}, nil
}

// GetC2WorkerInfo Obtain C2 worker machine
func (g *GateWayServiceImpl) GetC2WorkerInfo(ctx context.Context, args *pb.String) (*pb.String, error) {

	total, err := deploy.WorkerClusterServiceApi.GetC2workerInfo(args.Value)
	if err != nil {
		return &pb.String{Value: "0"}, err
	}

	return &pb.String{Value: strconv.Itoa(total)}, nil
}

func (g *GateWayServiceImpl) GetGateWayFile(ctx context.Context, args *pb.String) (*pb.String, error) {

	if args.Value == "" {
		return &pb.String{}, errors.New("filePath is nil")
	}
	dirAr := strings.Split(define.MainDisk, "/")
	dir := fmt.Sprintf("/%s/", dirAr[1])
	if !strings.Contains(args.Value, dir) {
		args.Value = define.MainDisk + args.Value
	}
	_, err := os.Stat(args.Value)
	if err != nil {
		return &pb.String{}, nil
	}
	return &pb.String{Value: args.Value}, nil
}

// GetHostTypeAndStatus Obtain the host type and whether it is enabled or not
func (g *GateWayServiceImpl) GetHostTypeAndStatus(ctx context.Context, args *pb.String) (*pb.HostRestartInfo, error) {
	log.Println("GetHostTypeAndStatus begin")
	res := &pb.HostRestartInfo{}

	hostType, err := deploy.SectorsRecoverServiceApi.GetHostType(args.Value)
	if err != nil {
		return res, err
	}

	var restartInfo []*pb.RestartInfo
	switch hostType {
	case strconv.Itoa(config.HostMinerType):

		opLotusInfo, err := service.ServiceGroupApp.LotusServiceGroup.GetLotusByOpID(args.GetValue())
		if err != nil {
			return res, err
		}
		s0 := &pb.RestartInfo{
			OpType:   define.ProgramLotus.String(),
			OpStatus: uint64(opLotusInfo.RunStatus),
		}
		restartInfo = append(restartInfo, s0)

		opMinerInfo, err := service.ServiceGroupApp.LotusServiceGroup.GetMinerByOpId(args.GetValue())
		if err != nil {
			return res, err
		}
		s1 := &pb.RestartInfo{
			OpType:   define.ProgramMiner.String(),
			OpStatus: uint64(opMinerInfo.RunStatus),
		}
		restartInfo = append(restartInfo, s1)
	case strconv.Itoa(config.HostWorkerType):

		opWorkerInfo, err := service.ServiceGroupApp.LotusServiceGroup.GetWorkerByOPId(args.GetValue())
		if err != nil {
			return res, err
		}
		s1 := &pb.RestartInfo{
			OpType:   define.ProgramWorkerTask.String(),
			OpStatus: uint64(opWorkerInfo.RunStatus),
		}
		restartInfo = append(restartInfo, s1)
	}
	res.Info = restartInfo
	return res, nil
}

func (g *GateWayServiceImpl) CarFileParam(ctx context.Context, args *pb.String) (*pb.CarInfo, error) {

	res := &pb.CarInfo{}
	if args.Value == "" {
		return res, errors.New("pieceCid is null")
	}

	carFile, err := deploy.SectorsRecoverServiceApi.GetCarFileParam(args.Value)

	if err != nil {
		return res, err
	}

	res.FileName = carFile.CarFileName
	res.InPutDir = carFile.InputDir
	res.FileStr = carFile.FileStr
	res.PieceCid = carFile.PieceCid
	res.PieceSize = int64(carFile.PieceSize)
	res.CarSize = int64(carFile.CarSize)
	res.DataCid = carFile.DataCid

	return res, nil
}

// GetNodeMinerInfo Gateway retrieves Op to obtain host monitoring information
func (p *GateWayServiceImpl) GetNodeMinerInfo(ctx context.Context, args *pb.OpHardwareInfo) (*pb.NodeMinerInfoResp, error) {
	log.Println("get op monitor info")
	client, _ := global.OpClinets.GetOpClient(args.HostUUID)
	if client == nil {
		log.Println("opClient Connection failed:" + args.HostUUID)
		return nil, errors.New("opClient Connection failed")
	}

	return client.GetNodeMinerInfo(ctx, args)
}

// GetMinerToken 获取miner Token
func (p *GateWayServiceImpl) GetMinerToken(ctx context.Context, args *pb.String) (*pb.String, error) {

	token, err := lotusService.GetMinerToken(args.Value)
	if err != nil || token == "" {
		return &pb.String{}, fmt.Errorf("GetMinerToken token:%s,err:%v", token, err)
	}
	return &pb.String{Value: token}, nil
}
