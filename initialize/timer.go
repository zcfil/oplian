package initialize

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/robfig/cron/v3"
	uuidGo "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"oplian/auth"
	"oplian/define"
	"oplian/model/gateway/request"
	modelSystem "oplian/model/system"
	"oplian/model/system/response"
	"oplian/service"
	"oplian/service/pb"
	"oplian/service/system"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"oplian/config"
	"oplian/global"
	"oplian/utils"
)

func Timer() {
	if global.ZC_CONFIG.Timer.Start {
		for i := range global.ZC_CONFIG.Timer.Detail {
			go func(detail config.Detail) {
				var option []cron.Option
				if global.ZC_CONFIG.Timer.WithSeconds {
					option = append(option, cron.WithSeconds())
				}
				_, err := global.ZC_Timer.AddTaskByFunc("ClearDB", global.ZC_CONFIG.Timer.Spec, func() {
					err := utils.ClearTable(global.ZC_DB, detail.TableName, detail.CompareField, detail.Interval)
					if err != nil {
						fmt.Println("timer error:", err)
					}
				}, option...)
				if err != nil {
					fmt.Println("add timer error:", err)
				}
			}(global.ZC_CONFIG.Timer.Detail[i])
		}
	}
}

// ConnectWeb gateway Connection web
func ConnectWeb(ctx context.Context) error {
	//gid := uuid.New()
	for {
		select {
		case <-time.After(time.Second * 10):
			client := &http.Client{
				Timeout: time.Second * 5,
			}
			var g = request.GateWayInfo{
				GateWayId: global.GateWayID.String(),
				IP:        global.ROOM_CONFIG.Gateway.IP,
				Port:      global.ROOM_CONFIG.Gateway.Port,
				Token:     global.ROOM_CONFIG.Gateway.Token,
			}
			buf, _ := json.Marshal(g)
			path := "http://" + global.ROOM_CONFIG.Web.Addr + "/conn/connectGateWay"
			req, err := http.NewRequest("POST", path, bytes.NewReader(buf))
			if err != nil {
				log.Println("ConnectWeb Error creating request:", err)
				return err
			}
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))
			req.Header.Add(define.TOKEN_NAME, global.ROOM_CONFIG.Web.Token)
			// 在这里开始进行请求
			_, err = client.Do(req)
			if err != nil {
				log.Println("ConnectWeb client.Do:", err)
				continue
			}
			//body, err := ioutil.ReadAll(resp.Body)
			//if err != nil {
			//	return err
			//}
			//fmt.Println(string(body))
			//resp.Body.Close()

		case <-ctx.Done():
			return nil
		}

	}
	return nil
}

// OpHeartBeat gateway Indicates the op for heartbeat check
func OpHeartBeat(ctx context.Context) {
	for {
		select {
		case <-time.After(time.Second * 10):
			var i = 0
			var numLock sync.RWMutex
			if len(global.OpClinets.GetOpClientList()) > 0 {
				var wg sync.WaitGroup
				wg.Add(len(global.OpClinets.GetOpClientList()))
				for k1, v1 := range global.OpClinets.GetOpClientList() {

					go func(k string, v *global.OpInfo) {
						c, cancel := context.WithTimeout(context.Background(), time.Second*5)
						defer func() {
							wg.Done()
							cancel()
						}()
						_, err := v.Clinet.Heartbeat(c, &pb.String{Value: "nathan"})
						if err != nil {
							v.Disconnect = true
						} else {
							v.Disconnect = false
							numLock.Lock()
							i++
							numLock.Unlock()
						}
						global.OpClinets.SetOpClient(k, v)

					}(k1, v1)
				}
				wg.Wait()
				log.Println("on-line OP1：", i)
			}

		case <-ctx.Done():
			return
		}
	}
}

// GateWayHeartbeat oplian Heartbeat check on the gateway
func GateWayHeartbeat(ctx context.Context) {
	for {
		select {
		case <-time.After(time.Second * 10):
			if len(global.GateWayClinets.Info) > 0 {
				global.GateWayClinets.LockRW.Lock()
				var i = 0
				for _, v := range global.GateWayClinets.Info {
					_, err := v.OplianHeartbeat(ctx, &pb.String{Value: "Oplian心跳检查Gateway"})
					if err != nil {
						v.Disconnect = true
						log.Println(err.Error())
						continue
					}
					i++
					v.Disconnect = false
				}
				global.GateWayClinets.LockRW.Unlock()
				log.Println("on-line GateWay：", i)
			}

		case <-ctx.Done():
			return
		}
	}
}

// OpC2ConnectOp The opC2 is connected to the Op
func OpC2ConnectOp(ctx context.Context) {
	//First connection
	conn, err := grpc.Dial(global.ROOM_CONFIG.Op.IP+":"+global.ROOM_CONFIG.Op.Port, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth.Authentication{Token: global.ROOM_CONFIG.Op.Token}))
	if err != nil {
		log.Println("Abnormal network connection: Dial: ", global.ROOM_CONFIG.Op.IP+":"+global.ROOM_CONFIG.Op.Port)
	}
	global.OpC2ToOp = pb.NewOpServiceClient(conn)
	if _, err = global.OpC2ToOp.OpC2Connect(ctx, &pb.RequestConnect{OpId: global.OpC2UUID.String(), Port: global.ROOM_CONFIG.OpC2.Port}); err != nil {
		log.Println("Abnormal network connection: OpConnect: ", global.ROOM_CONFIG.Op.IP+":"+global.ROOM_CONFIG.Op.Port, err.Error())
		if conn != nil {
			conn.Close()
		}
		conn = nil
	}
	//Heartbeat check
	go func() {
		for {
			select {
			case <-time.After(time.Second * 10):
				//Initial connection
				if conn == nil {
					conn, err = grpc.Dial(global.ROOM_CONFIG.Op.IP+":"+global.ROOM_CONFIG.Op.Port, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth.Authentication{Token: global.ROOM_CONFIG.Op.Token}))
					if err != nil {
						log.Println("Abnormal network connection: Dial: ", global.ROOM_CONFIG.Op.IP+":"+global.ROOM_CONFIG.Op.Port)
					}
					global.OpC2ToOp = pb.NewOpServiceClient(conn)
				}

				//Update link
				if _, err = global.OpC2ToOp.OpC2Heartbeat(ctx, &pb.String{Value: global.OpC2UUID.String()}); err != nil {
					if _, err = global.OpC2ToOp.OpC2Connect(ctx, &pb.RequestConnect{OpId: global.OpC2UUID.String(), Port: global.ROOM_CONFIG.OpC2.Port}); err != nil {
						log.Println("Abnormal network connection: OpConnect: ", global.ROOM_CONFIG.Op.IP+":"+global.ROOM_CONFIG.Op.Port, err.Error())
						if conn != nil {
							conn.Close()
						}
						conn = nil
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// GatewayGenerateUUID Gateway initialization generates a unique uuid value (instead of reading from config.yaml, read from a solidified file)
func GatewayGenerateUUID() {
	var gatewayUUID uuidGo.UUID

	_ = os.MkdirAll(define.PathIpfsConfig, os.ModePerm)

	//path := define.PathIpfsConfig + "gateway_uuid"
	path := "config/gateway_uuid"
	_, err := os.Stat(path)
	if err != nil {

		gatewayUUID = uuidGo.NewV4()

		out, createErr := os.Create(path)
		if createErr != nil {
			log.Println("function os.Create() Filed", createErr.Error())
			return
		}
		defer out.Close()

		_, writeErr := out.WriteString(gatewayUUID.String())
		if writeErr != nil {
			log.Println("function os.WriteString() Filed", writeErr.Error())
			return
		}
	} else {

		f, openError := ioutil.ReadFile(path) // 读取文件
		if openError != nil {
			log.Println("function ioutil.ReadFile() Filed", openError.Error())
			return
		}

		gatewayUUID, err = uuidGo.FromString(string(f))
		if err != nil {
			log.Println("function uuid.FromBytes() Filed", "err", err.Error())
			return
		}
	}
	global.GateWayID = gatewayUUID
}

// OpInitGlobalUUID The op initializes global.OpUUID
func OpInitGlobalUUID() {
	global.OpUUID = uuidGo.NewV4()
	path := define.PathIpfsConfig + "op_uuid"
	_, err := os.Stat(path)
	if err == nil {

		f, openError := ioutil.ReadFile(path) // 读取文件
		if openError != nil {
			log.Println("function ioutil.ReadFile() Filed", openError.Error())
			return
		}
		opUUIDFile, err := uuidGo.FromString(string(f))
		if err != nil {
			log.Println("function uuid.FromBytes() Filed", "err", err.Error())
			return
		}
		global.OpUUID = opUUIDFile
	}
	return
}

// OpGenerateUUID The op end is initialized to generate a unique uuid value
func OpGenerateUUID() {
	var opUUID uuidGo.UUID
	var isNew bool

	out, err := global.OpToGatewayClient.GetHostInfoByIPAndGatewayId(context.Background(), &pb.RequestOp{Ip: global.LocalIP})
	if err != nil {
		isNew = true
	} else {
		opUUID, _ = uuidGo.FromString(out.GetHostUUID())
	}

	_ = os.MkdirAll(define.PathIpfsConfig, os.ModePerm)

	path := define.PathIpfsConfig + "op_uuid"
	//path := "config/op_uuid"
	_, err = os.Stat(path)
	if err != nil {

		if isNew || opUUID.String() == "00000000-0000-0000-0000-000000000000" {
			opUUID = uuidGo.NewV4()
		}

		out, createErr := os.Create(path)
		if createErr != nil {
			log.Println("function os.Create() Filed", createErr.Error())
			return
		}
		defer out.Close() // 创建文件 defer 关闭

		_, writeErr := out.WriteString(opUUID.String())
		if writeErr != nil {
			log.Println("function os.WriteString() Filed", writeErr.Error())
			return
		}
	} else {

		if global.OpUUID != opUUID && opUUID.String() != "00000000-0000-0000-0000-000000000000" {
			err := ioutil.WriteFile(path, []byte(opUUID.String()), 0644)
			if err != nil {
				log.Println("ioutil WriteFile Filed", err.Error())
				return
			}
		} else {
			opUUID = global.OpUUID
		}
	}
	global.OpUUID = opUUID
}

// DiskInitialization Initialize the hard disk mount
func DiskInitialization(groupArray, c2Worker, isPower, isMiner bool) {

	if isPower || isMiner {
		if out, err := exec.Command("bash", "-c", define.PathDiskPowerInitialization).Output(); err != nil {
			log.Println("Execute the disk initialization script"+define.PathDiskPowerInitialization+"failure", string(out), err.Error())
		} else {
			log.Println(string(out))
		}
	} else {
		if out, err := exec.Command("bash", "-c", define.PathDiskInitialization).Output(); err != nil {
			log.Println("Execute the disk initialization script"+define.PathDiskInitialization+"failure", string(out), err.Error())
		} else {
			log.Println(string(out))
		}
	}
}

// CheckDiskProofParameters Check whether the disk has been pulled for proof parameters
func CheckDiskProofParameters(notProof bool, storage bool, paramPath string) error {
	if paramPath != "" {

		// 初始化对应路径
		if err := service.ServiceGroupApp.LotusServiceGroup.CheckProofsParameterPath(paramPath); err != nil {
			return err
		}
		define.MainDisk = path.Join(paramPath, "..")
		define.PathInit()

		return nil
	}

	storageInfo := utils.GetStorageOpDiskInfo()

	var proofsPosition string
	var maxSize string
	for _, val := range storageInfo {
		if err := service.ServiceGroupApp.LotusServiceGroup.CheckProofsParameter(val.Mountpoint); err == nil {

			proofsPosition = val.Mountpoint
		}
		if utils.CompareTwoSizes(val.Size, maxSize) {
			maxSize = val.Size
		}
	}
	if len(proofsPosition) != 0 {

		define.MainDisk = proofsPosition

		define.PathInit()

	} else {

		if !utils.CompareTwoSizes(maxSize, config.HostC2WorkerDiskSize) {
			log.Println("The disk capacity of the host is insufficient. The proof parameters cannot be downloaded")
			return nil
		} else {
			for _, val := range storageInfo {
				if val.Size == maxSize {
					define.MainDisk = val.Mountpoint
				}
			}

			define.PathInit()

			if !notProof && !storage {
				return errors.New("请加入--not-proof-parameters参数启动")
			} else {
				if !storage {

					err := service.ServiceGroupApp.LotusServiceGroup.DownlodParameters(context.Background(), global.LocalIP)
					if err != nil {
						log.Println("Description The parameter download failed")
						return err
					}
				}
			}

		}
	}
	return nil
}

// PolicyWarn Policy alarm
func PolicyWarn(ctx context.Context) {
	for {
		warn := system.WarnManageService{}
		res, err := warn.StrategyWarnConfig()
		if err != nil {
			return
		}

		err = system.WarnManageServiceApp.StrategyProcessing(ctx, res)
		if err != nil {
			time.Sleep(30 * time.Second)
		}
	}
}

// UpdateHostInfo The op updates the hardware information to the database
func UpdateHostInfo(ctx context.Context, groupArray, dcType, c2Worker, isStorage bool) {

	var s utils.OpServer
	var err error

	s.Os = utils.InitOS()
	s.GetCPUInfo()
	s.GetGPUInfo()
	s.GetLinuxSystemInfo()
	if s.Ram, err = utils.InitRAM(); err != nil {
		log.Println("func utils.InitRAM() Failed", err.Error())
		return
	}

	if s.Disk, err = utils.InitDisk(); err != nil {
		log.Println("func utils.InitDisk() Failed", err.Error())
		return
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
		GatewayId:        global.GateWayID.String(),
		SystemVersion:    s.SystemInfo.SystemVersion,
		SystemBits:       int64(s.SystemInfo.SystemBits),
		GPUNum:           int64(len(s.GPU.Gpus)),
		HostName:         s.SystemInfo.HostName,
		IsGroupArray:     groupArray,
	}

	if dcType {
		hostInfo.HostClassify = config.HostDCStorageType
	}

	if c2Worker {
		hostInfo.HostClassify = config.HostC2WorkerType
	}

	if isStorage {
		hostInfo.HostClassify = config.HostStorageType
	}

	_, err = global.OpToGatewayClient.AddHostRecord(ctx, hostInfo)
	if err != nil {
		log.Println("client.AddHostRecord", "Failed to record the host information. Procedure!")
		return
	}

	log.Println("client.AddHostRecord", "The host information is successfully recorded. Procedure!")
	return
}

// GetHostMonitorInfo The op updates the hardware monitoring information to the database
func GetHostMonitorInfo(ctx context.Context) {

	time.Sleep(config.GoFunctionWaitTime)
	c := cron.New()
	spec := "*/5 * * * *"
	c.AddFunc(spec, func() {
		var opServer utils.OpServer
		var err error

		res, err := global.OpToGatewayClient.HostType(context.TODO(), &pb.String{Value: global.OpUUID.String()})
		if err != nil {
			log.Println(fmt.Sprintf("HostType err:%s", err.Error()))
			time.Sleep(time.Second * time.Duration(utils.Five))
			return
		}
		hostClassify, _ := utils.StringToInt64(res.Value)

		var s utils.HostMonitor

		s.GetHostMonitorInfo(hostClassify)

		if opServer.Ram, err = utils.InitRAM(); err != nil {
			log.Println("func utils.InitRAM() Failed", err.Error())
			return
		}

		gpuUseInfo, _ := json.Marshal(s.GPUUseInfo)

		hostMonitorInfo := &pb.HostMonitorInfo{
			HostUUID:       global.OpUUID.String(),
			CPUUseRate:     s.CPUUseRate,
			DiskUseRate:    float32(s.MonitorDisk.UsePercent),
			MemoryUseRate:  float32(opServer.Ram.UsedPercent),
			GPUUseInfo:     string(gpuUseInfo),
			CPUTemperature: s.CPUTemperature,
			DiskSize:       s.MonitorDisk.Size,
			DiskUseSize:    s.MonitorDisk.Used,
			MemorySize:     int64(opServer.Ram.TotalMB),
			MemoryUseSize:  int64(opServer.Ram.UsedMB),
		}

		_, err = global.OpToGatewayClient.AddHostMonitorRecord(ctx, hostMonitorInfo)
		if err != nil {
			log.Println("client.AddHostMonitorRecord", "The host hardware monitoring information failed to be recorded!", err.Error())
			return
		}
		log.Println("client.AddHostMonitorRecord", "The host hardware monitoring information is successfully recorded. Procedure!")
	})
	go c.Start()
	defer c.Stop()
	select {}
}

// HostSystemInitialization The op host is initialized
func HostSystemInitialization() {

	if out, err := exec.Command("bash", "-c", define.PathIpfsScriptHostSystemInitialization).Output(); err != nil {
		log.Println("Description Failed to execute the system initialization script", err.Error())
	} else {
		log.Println(string(out))
	}
	log.Println("End of executing system initialization script")
}

// UpdateMachineRoomInfo gateway启动的时候添加对应信息到机房信息表里面
func UpdateMachineRoomInfo(ctx context.Context) {

	internetIp, err := utils.GetInternetIP()
	if err != nil {
		log.Println("get host internet ip failed", err.Error())
		return
	}

	ipAddress := utils.GetIpAddress(internetIp)
	strUUID := strings.Split(uuidGo.NewV4().String(), "-")

	// 处理机房信息表
	sysRoom := modelSystem.SysMachineRoomRecord{
		ZC_MODEL:        global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		RoomId:          strUUID[0],
		RoomName:        ipAddress.Country + "机房",
		RoomType:        0,
		CabinetsNum:     0,
		RoomTemp:        0,
		RoomLeader:      "",
		RoomLeaderPhone: "",
		RoomSupplier:    "",
		SupplierContact: "",
		SupplierPhone:   "",
		RoomAdmin:       "",
		RoomOwner:       "",
		PhysicalAddress: "",
		RoomArea:        0,
		GatewayId:       global.GateWayID.String(),
	}

	intranetIP := global.LocalIP
	if intranetIP == "" {
		return
	}

	sysRoom.IntranetIP = intranetIP
	// 写入数据库信息
	roomService := system.MachineRoomRecordService{}
	roomInfo, err := roomService.GetRoomByGatewayId(sysRoom.GatewayId)

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}

		err = roomService.CreateSysMachineRoomRecord(sysRoom)
		if err != nil {
			return
		}

		hostInfo := system.HostRecordService{}
		err = hostInfo.HostBindRoomByGatewayId(&sysRoom)
		if err != nil {
			global.ZC_LOG.Error("HostBindRoomByGatewayId err ", zap.Error(err))
		}

	} else {

		hostInfo := system.HostRecordService{}
		err = hostInfo.HostBindRoomByGatewayId(&roomInfo)
		if err != nil {
			global.ZC_LOG.Error("HostBindRoomByGatewayId err", zap.Error(err))
		}
	}
	return
}

// GetHostGroupArray Get the host array information
func GetHostGroupArray(ctx context.Context) bool {

	out, err := global.OpToGatewayClient.GetHostGroupArray(ctx, &pb.OpHostUUID{HostUUID: global.OpUUID.String()})
	if err != nil {
		return false
	}
	return out.IsGroupArray
}

// BeginHostPatrol When the gateway starts, the corresponding information is added to the equipment room information table
func BeginHostPatrol() {

	time.Sleep(config.GoFunctionWaitTime)
	go HostMinerPatrol()
	go HostWorkerPatrol()
	go HostStoragePatrol()
}

func HostMinerPatrol() {
	for {
		var patrolTimeInterval int64

		patrolConfig := system.PatrolConfigService{}
		minerConfig, err := patrolConfig.GetSysPatrolConfigInfoByType(config.HostMinerTest)
		if err != nil {
			global.ZC_LOG.Error("GetSysPatrolConfigInfoByType err:", zap.Error(err))
			patrolTimeInterval = config.PatrolMinerTimeConfig
		} else {
			if minerConfig.IntervalTime == 0 {
				patrolTimeInterval = config.PatrolMinerTimeConfig
			} else {
				patrolTimeInterval = minerConfig.IntervalTime
			}
		}

		hostService := system.HostRecordService{}
		hostList, err := hostService.GetSysHostRecordListForPatrol([]int64{config.HostMinerType}, global.GateWayID.String())
		if err != nil {
			global.ZC_LOG.Error("GetSysHostRecordListForPatrol err:", zap.Error(err))
		}

		if len(hostList) == 0 {
			continue
		}

		for _, val := range hostList {
			go func(val response.SysHostRecordPatrol) {
				sysHostPatrol := modelSystem.SysHostPatrolRecord{
					ZC_MODEL:      global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
					HostUUID:      val.UUID,
					PatrolType:    val.HostClassify,
					PatrolBeginAt: time.Now().Unix(),
					PatrolEndAt:   0,
					PatrolResult:  config.HostUnderTest,
					PatrolMode:    config.AutomaticTrigger,
				}
				hostPatrolInfo := system.HostPatrolRecordService{}
				err := hostPatrolInfo.CreateSysHostPatrolRecord(sysHostPatrol)
				if err != nil {
					global.ZC_LOG.Error("Creation failed!", zap.Error(err))
					return
				}
				args := &pb.HostPatrolInfo{
					HostClassify:  val.HostClassify,
					HostUUID:      val.UUID,
					PatrolMode:    config.AutomaticTrigger,
					PatrolBeginAt: sysHostPatrol.PatrolBeginAt,
					PatrolHostIP:  "",
				}
				client, dis := global.OpClinets.GetOpClient(args.HostUUID)
				if client == nil || dis {
					log.Println("opClient Connection failed:" + args.HostUUID)
					updatePatrolInfo := modelSystem.SysHostPatrolRecord{
						HostUUID:      sysHostPatrol.HostUUID,
						PatrolBeginAt: sysHostPatrol.PatrolBeginAt,
						PatrolEndAt:   time.Now().Unix(),
						PatrolResult:  config.HostTestFailed,
					}
					err = hostPatrolInfo.UpdateSysHostPatrolRecordConnectError(&updatePatrolInfo)
					if err != nil {
						global.ZC_LOG.Error("UpdateSysHostPatrolRecordConnectError err:", zap.Error(err))
						return
					}

					return
				}
				_, err = client.OpInformationPatrol(context.Background(), args)
				if err != nil {

					updatePatrolInfo := modelSystem.SysHostPatrolRecord{
						HostUUID:      sysHostPatrol.HostUUID,
						PatrolBeginAt: sysHostPatrol.PatrolBeginAt,
						PatrolEndAt:   time.Now().Unix(),
						PatrolResult:  config.HostTestFailed,
					}
					err = hostPatrolInfo.UpdateSysHostPatrolRecordConnectError(&updatePatrolInfo)
					if err != nil {
						global.ZC_LOG.Error("UpdateSysHostPatrolRecordConnectError err:", zap.Error(err))
						return
					}
				}
			}(val)
		}

		time.Sleep(time.Duration(patrolTimeInterval) * time.Second)
	}
}

func HostWorkerPatrol() {
	for {
		var patrolTimeInterval int64

		patrolConfig := system.PatrolConfigService{}
		workerConfig, err := patrolConfig.GetSysPatrolConfigInfoByType(config.HostWorkerTest)
		if err != nil {
			global.ZC_LOG.Error("GetSysPatrolConfigInfoByType err:", zap.Error(err))
			patrolTimeInterval = config.PatrolWorkerTimeConfig
		} else {
			if workerConfig.IntervalTime == 0 {
				patrolTimeInterval = config.PatrolMinerTimeConfig
			} else {
				patrolTimeInterval = workerConfig.IntervalTime
			}
		}

		hostService := system.HostRecordService{}
		hostList, err := hostService.GetSysHostRecordListForPatrol([]int64{config.HostWorkerType, config.HostC2WorkerType}, global.GateWayID.String())
		if err != nil {
			global.ZC_LOG.Error("GetSysHostRecordListForPatrol err:", zap.Error(err))
		}
		if len(hostList) == 0 {
			continue
		}

		for _, val := range hostList {
			go func(val response.SysHostRecordPatrol) {
				sysHostPatrol := modelSystem.SysHostPatrolRecord{
					ZC_MODEL:      global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
					HostUUID:      val.UUID,
					PatrolType:    config.HostWorkerType,
					PatrolBeginAt: time.Now().Unix(),
					PatrolEndAt:   0,
					PatrolResult:  config.HostUnderTest,
					PatrolMode:    config.AutomaticTrigger,
				}
				hostPatrolInfo := system.HostPatrolRecordService{}
				err := hostPatrolInfo.CreateSysHostPatrolRecord(sysHostPatrol)
				if err != nil {
					global.ZC_LOG.Error("Creation failed!", zap.Error(err))
					return
				}
				args := &pb.HostPatrolInfo{
					HostClassify:  config.HostWorkerType,
					HostUUID:      val.UUID,
					PatrolMode:    config.AutomaticTrigger,
					PatrolBeginAt: sysHostPatrol.PatrolBeginAt,
					PatrolHostIP:  "",
				}

				client, dis := global.OpClinets.GetOpClient(args.HostUUID)
				if client == nil || dis {

					log.Println("opClient Connection failed:" + args.HostUUID)
					updatePatrolInfo := modelSystem.SysHostPatrolRecord{
						HostUUID:      sysHostPatrol.HostUUID,
						PatrolBeginAt: sysHostPatrol.PatrolBeginAt,
						PatrolEndAt:   time.Now().Unix(),
						PatrolResult:  config.HostTestFailed,
					}
					err = hostPatrolInfo.UpdateSysHostPatrolRecordConnectError(&updatePatrolInfo)
					if err != nil {
						global.ZC_LOG.Error("UpdateSysHostPatrolRecordConnectError err:", zap.Error(err))
						return
					}
					return
				}

				_, err = client.OpInformationPatrol(context.Background(), args)
				if err != nil {

					updatePatrolInfo := modelSystem.SysHostPatrolRecord{
						HostUUID:      sysHostPatrol.HostUUID,
						PatrolBeginAt: sysHostPatrol.PatrolBeginAt,
						PatrolEndAt:   time.Now().Unix(),
						PatrolResult:  config.HostTestFailed,
					}
					err = hostPatrolInfo.UpdateSysHostPatrolRecordConnectError(&updatePatrolInfo)
					if err != nil {
						global.ZC_LOG.Error("UpdateSysHostPatrolRecordConnectError err:", zap.Error(err))
						return
					}
				}
			}(val)
		}

		time.Sleep(time.Duration(patrolTimeInterval) * time.Second)
	}
}

func HostStoragePatrol() {
	for {
		var patrolTimeInterval int64

		patrolConfig := system.PatrolConfigService{}
		storageConfig, err := patrolConfig.GetSysPatrolConfigInfoByType(config.HostStorageTest)
		if err != nil {
			global.ZC_LOG.Error("GetSysPatrolConfigInfoByType err:", zap.Error(err))
			patrolTimeInterval = config.PatrolStorageTimeConfig
		} else {
			if storageConfig.IntervalTime == 0 {
				patrolTimeInterval = config.PatrolMinerTimeConfig
			} else {
				patrolTimeInterval = storageConfig.IntervalTime
			}
		}

		hostService := system.HostRecordService{}
		hostList, err := hostService.GetSysHostRecordListForPatrol([]int64{config.HostStorageType}, global.GateWayID.String())
		if err != nil {
			global.ZC_LOG.Error("GetSysHostRecordListForPatrol err:", zap.Error(err))
		}
		if len(hostList) == 0 {
			continue
		}

		for _, val := range hostList {
			go func(val response.SysHostRecordPatrol) {
				sysHostPatrol := modelSystem.SysHostPatrolRecord{
					ZC_MODEL:      global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
					HostUUID:      val.UUID,
					PatrolType:    val.HostClassify,
					PatrolBeginAt: time.Now().Unix(),
					PatrolEndAt:   0,
					PatrolResult:  config.HostUnderTest,
					PatrolMode:    config.AutomaticTrigger,
				}
				hostPatrolInfo := system.HostPatrolRecordService{}
				err := hostPatrolInfo.CreateSysHostPatrolRecord(sysHostPatrol)
				if err != nil {
					global.ZC_LOG.Error("Creation failed!", zap.Error(err))
					return
				}

				otherInfo, err := hostService.GetSysOtherHostRecord(val.UUID, val.GatewayId)
				if err != nil {

					updatePatrolInfo := modelSystem.SysHostPatrolRecord{
						HostUUID:      sysHostPatrol.HostUUID,
						PatrolBeginAt: sysHostPatrol.PatrolBeginAt,
						PatrolEndAt:   time.Now().Unix(),
						PatrolResult:  config.HostTestFailed,
					}
					err = hostPatrolInfo.UpdateSysHostPatrolRecordConnectError(&updatePatrolInfo)
					if err != nil {
						return
					}
					if err == gorm.ErrRecordNotFound {
						return
					}

					return
				}
				args := &pb.HostPatrolInfo{
					HostClassify:  val.HostClassify,
					HostUUID:      val.UUID,
					PatrolMode:    config.AutomaticTrigger,
					PatrolBeginAt: sysHostPatrol.PatrolBeginAt,
					PatrolHostIP:  otherInfo.IntranetIP,
				}

				client, dis := global.OpClinets.GetOpClient(args.HostUUID)
				if client == nil || dis {

					log.Println("opClient Connection failed:" + args.HostUUID)
					updatePatrolInfo := modelSystem.SysHostPatrolRecord{
						HostUUID:      sysHostPatrol.HostUUID,
						PatrolBeginAt: sysHostPatrol.PatrolBeginAt,
						PatrolEndAt:   time.Now().Unix(),
						PatrolResult:  config.HostTestFailed,
					}
					err = hostPatrolInfo.UpdateSysHostPatrolRecordConnectError(&updatePatrolInfo)
					if err != nil {
						return
					}
					return
				}

				_, err = client.OpInformationPatrol(context.Background(), args)
				if err != nil {

					updatePatrolInfo := modelSystem.SysHostPatrolRecord{
						HostUUID:      sysHostPatrol.HostUUID,
						PatrolBeginAt: sysHostPatrol.PatrolBeginAt,
						PatrolEndAt:   time.Now().Unix(),
						PatrolResult:  config.HostTestFailed,
					}
					err = hostPatrolInfo.UpdateSysHostPatrolRecordConnectError(&updatePatrolInfo)
					if err != nil {
						global.ZC_LOG.Error("UpdateSysHostPatrolRecordConnectError err:", zap.Error(err))
						return
					}
				}
			}(val)
		}

		time.Sleep(time.Duration(patrolTimeInterval) * time.Second)
	}
}

// BeginCheckBadSector Error sector information is queried and recorded
func BeginCheckBadSector() {

	time.Sleep(config.GoFunctionWaitTime)
	for {

		hostService := system.HostRecordService{}
		hostList, err := hostService.GetSysHostRecordListForPatrol([]int64{config.HostMinerType}, global.GateWayID.String())
		if err != nil {
			global.ZC_LOG.Error("GetSysHostRecordListForPatrol err:", zap.Error(err))
		}
		if len(hostList) == 0 {
			continue
		}

		for _, val := range hostList {
			client, _ := global.OpClinets.GetOpClient(val.UUID)
			if client == nil {
				log.Println("opClient Connection failed:" + val.UUID)
				continue
			}

			args := &pb.HostCheckDiskInfo{
				GateWayId: global.GateWayID.String(),
				HostUUID:  val.UUID,
			}

			client.OpCheckBadSector(context.Background(), args)
		}

		time.Sleep(300 * time.Second)
	}
}

// DealC2WorkerDisk c2Worker
func DealC2WorkerDisk() {

	md0Size, err := utils.GetMd0DiskSize()
	if err != nil || len(md0Size) == 0 {
		log.Println(err.Error())
		return
	}

	if !utils.CompareTwoSizes(md0Size, config.HostC2WorkerDiskSize) {
		log.Println("Insufficient capacity")
		return
	} else {

		err = service.ServiceGroupApp.LotusServiceGroup.DownlodParameters(context.Background(), global.LocalIP)
		if err != nil {
			return
		}
	}
}

func OpChmodDirectory() {
	dirPath := fmt.Sprintf("%s %s", define.PathIpfsProgram, define.PathIpfsScript)
	err := utils.ChmodDirectory(dirPath)
	if err != nil {
		return
	}
}

func DownloadNecessaryFiles() {

	time.Sleep(config.GoFunctionWaitTime)

	conn, err := ftp.Dial(global.ROOM_CONFIG.Download.Addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Println("link error:", err.Error())
		return
	}

	err = conn.Login(global.ROOM_CONFIG.Download.User, global.ROOM_CONFIG.Download.Password)
	if err != nil {
		log.Println("login error: ", err.Error())
		return
	}
	defer conn.Quit()

	path := define.PathIpfsConfig + config.DownloadTarFileName
	_, err = os.Stat(path)

	if err == nil {

		r, err := conn.Retr("/" + config.DownloadTarFileHash)
		if err != nil {
			log.Println("conn.Retr error: ", err.Error())
			return
		}
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			log.Println("read hash file error: ", err.Error())
			return
		}
		r.Close()

		localFileHash := Gethash(path)
		if string(buf) == localFileHash {
			log.Println("It is already the latest ")
			return
		}

		err = os.Remove(path)

	}

	r, err := conn.Retr("/" + config.DownloadTarFileName)
	if err != nil {
		log.Println("Download "+config.DownloadTarFileName+" error: ", err.Error())
		return
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("read "+config.DownloadTarFileName+" error: ", err.Error())
		return
	}
	r.Close()
	err = ioutil.WriteFile(define.PathIpfsConfig+config.DownloadTarFileName, buf, 0644)
	if err != nil {
		log.Println("write "+config.DownloadTarFileName+" error: ", err.Error())
		return
	}

}

func Gethash(path string) (hash string) {
	file, _ := os.Open(path)
	h_ob := sha256.New()
	_, err := io.Copy(h_ob, file)
	if err == nil {
		hash := h_ob.Sum(nil)
		hashvalue := hex.EncodeToString(hash)
		return hashvalue
	} else {
		return ""
	}
}

// DownloadProofParameters
func DownloadProofParameters() {

	time.Sleep(config.GoFunctionWaitTime)

	intranetIP := global.LocalIP
	if len(intranetIP) > 0 {

		err := service.ServiceGroupApp.LotusServiceGroup.DownlodParameters(context.Background(), intranetIP)
		if err != nil {
			log.Println("DownlodParameters error:", err.Error())
			return
		}
		return
	}
	return
}

func CheckPatrolAndTest() {
	time.Sleep(config.GoFunctionWaitTime)
	go CheckPatrolList()
	go CheckTestList()
}

// CheckPatrolList Check the patrol list regularly
func CheckPatrolList() {

	ticker := time.NewTicker(config.HostNetTimeout)
	defer ticker.Stop()

	for range ticker.C {

		patrolService := system.HostPatrolRecordService{}
		patrolList, err := patrolService.GetTimeoutPatrolList()
		if err != nil {
			continue
		}

		for _, val := range patrolList {
			val.PatrolEndAt = time.Now().Unix()
			val.PatrolResult = config.HostTestFailed
			err = patrolService.UpdateSysHostPatrolRecordConnectError(&val)
			if err != nil {
				global.ZC_LOG.Error("Modification failed!", zap.Error(err))
				continue
			}
		}
	}
}

// CheckTestList Check the test list regularly
func CheckTestList() {

	ticker := time.NewTicker(config.HostTestTimeout)
	defer ticker.Stop()

	for range ticker.C {

		testService := system.HostTestRecordService{}
		testlList, err := testService.GetTimeoutTestList()
		if err != nil {
			continue
		}

		for _, val := range testlList {

			val.TestEndAt = time.Now().Unix()
			val.TestResult = config.HostTestFailed
			err = testService.UpdateSysHostTestResult(&val)
			if err != nil {
				continue
			}

			client, _ := global.OpClinets.GetOpClient(val.HostUUID)
			if client == nil {
				log.Println("CheckTestList opClient Connection failed:" + val.HostUUID)
				continue
			}

			client.KillBenchAndScript(context.Background(), &pb.String{})
		}
	}
}

// CheckDiskProofParametersGateway Gateway Checks whether the disk has been pulled down certificate parameters
func CheckDiskProofParametersGateway() {

	storageInfo := utils.GetStorageOpDiskInfo()
	var maxSize string
	for _, val := range storageInfo {
		if utils.CompareTwoSizes(val.Size, maxSize) {
			maxSize = val.Size
		}
	}

	for _, val := range storageInfo {
		if val.Size == maxSize {
			define.MainDisk = val.Mountpoint
		}
	}

	define.PathInit()
}

// DownloadFileFromGateway
func DownloadFileFromGateway(loadPath string, fileList []string, isWait bool) {
	if isWait {
		time.Sleep(config.WorkerMountNFSWaitTime)
	}
	loadPath = loadPath[:len(loadPath)-1]

	intranetIP := global.LocalIP
	for _, val := range fileList {
		if _, err := os.Stat(loadPath + "/" + val); os.IsNotExist(err) {
			opMap := make([]*pb.OpInfo, 0)
			fileMap := make([]*pb.FileInfo, 0)
			opMap = append(opMap, &pb.OpInfo{Ip: intranetIP, Port: global.ROOM_CONFIG.Op.Port, OpId: global.OpUUID.String()})
			fileMap = append(fileMap, &pb.FileInfo{FileName: val})

			info := &pb.DownLoadInfo{
				DownloadPath: loadPath,
				OpInfo:       opMap,
				FileInfo:     fileMap,
				GateWayPath:  loadPath,
			}

			res, _ := global.OpToGatewayClient.DownLoadFiles(context.TODO(), info)
			if res.Code != 200 {
				log.Println("Download file from gateway failed: ", res.Msg)
				return
			}

			t1 := &pb.FileInfo{
				FileName: loadPath + "/new_zip.zip",
				Path:     loadPath,
			}
			global.OpToGatewayClient.DelGateWayFile(context.TODO(), t1)
		}
	}
	OpChmodDirectory()
}

// OpRemountDisk op
func OpRemountDisk(ctx context.Context) {

	time.Sleep(config.WorkerMountNFSWaitTime)
	var err error
	res, err := global.OpToGatewayClient.HostType(context.TODO(), &pb.String{Value: global.OpUUID.String()})
	if err != nil {
		log.Println(fmt.Sprintf("HostType err:%s", err.Error()))
		return
	}

	hostClassify, _ := utils.StringToInt64(res.Value)
	if hostClassify == config.HostWorkerType || hostClassify == config.HostMinerType {
		_, err = global.OpToGatewayClient.WorkerMountNFS(context.TODO(), &pb.OpHostUUID{HostUUID: global.OpUUID.String(), HostType: res.Value})
		if err != nil {
			log.Println(fmt.Sprintf("WorkerMountNFS err:%s", err.Error()))
			return
		}
	}
}

// OpRestartLotus lotus related is performed when the op starts
func OpRestartLotus(ctx context.Context) {

	hostInfo, err := global.OpToGatewayClient.GetHostTypeAndStatus(context.TODO(), &pb.String{Value: global.OpUUID.String()})
	if err != nil {
		log.Println(fmt.Sprintf("OpRestartLotus err:%s", err.Error()))
		return
	}

	for _, val := range hostInfo.Info {
		if val.OpStatus != 1 {
			continue
		}
		switch val.OpType {
		case define.ProgramLotus.String():
			_, err = exec.Command("bash", "-c", "supervisorctl start lotus").CombinedOutput()
			if err != nil && !strings.Contains(err.Error(), "exit status") {
				log.Println("cmd `supervisorctl restart lotus` Filed", err.Error())
				continue
			}
			//log.Println("cmd `lsb_release -a | grep Description:`", string(out))
		case define.ProgramMiner.String():
			_, err = exec.Command("bash", "-c", "supervisorctl start lotus-miner").CombinedOutput()
			if err != nil && !strings.Contains(err.Error(), "exit status") {
				log.Println("cmd `supervisorctl restart lotus-miner` Filed", err.Error())
				continue
			}
			//log.Println("cmd `supervisorctl restart lotus-miner`", string(out))
			_, err = exec.Command("bash", "-c", "supervisorctl start boost").CombinedOutput()
			if err != nil && !strings.Contains(err.Error(), "exit status") {
				log.Println("cmd `supervisorctl restart boost` Filed", err.Error())
				continue
			}
			//log.Println("cmd `supervisorctl restart boost`", string(out1))
		case define.ProgramWorkerTask.String():
			_, err = exec.Command("bash", "-c", "supervisorctl start lotus-worker").CombinedOutput()
			if err != nil && !strings.Contains(err.Error(), "exit status") {
				log.Println("cmd `supervisorctl restart lotus-worker` Filed", err.Error())
				continue
			}
		}
	}

	log.Println("OpRestartLotus success")
}

func CheckLotusHeart() {

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {

		select {
		case <-ticker.C:

			if len(define.Lh.LotusMap) > 0 {

				define.Lh.LotusLock.RLock()
				lotusList := make(map[string]int)
				for k, v := range define.Lh.LotusMap {
					if time.Now().Sub(v).Minutes() > 1 {
						lotusList[k] = 0
					} else {
						lotusList[k] = 1
					}
				}
				define.Lh.LotusLock.RUnlock()
				if len(lotusList) > 0 {
					err := service.ServiceGroupApp.LotusServiceGroup.ModifyLotusStatus(lotusList)
					if err != nil {
						return
					}
				}
			}

			minerList := make(map[string]int)
			if len(define.Lh.MinerMap) > 0 {

				define.Lh.MinerLock.RLock()
				for k, v := range define.Lh.MinerMap {
					if time.Now().Sub(v).Minutes() > 1 {
						minerList[k] = 0
					} else {
						minerList[k] = 1
					}
				}
				define.Lh.MinerLock.RUnlock()
				if len(minerList) > 0 {
					err := service.ServiceGroupApp.LotusServiceGroup.ModifyMinerStatus(minerList)
					if err != nil {
						return
					}
				}
			}

			workerList := make(map[string]int)
			if len(define.Lh.WorkerMap) > 0 {

				define.Lh.WorkerLock.RLock()
				for k, v := range define.Lh.WorkerMap {
					if time.Now().Sub(v).Minutes() > 1 {
						workerList[k] = 0
					} else {
						workerList[k] = 1
					}
				}
				define.Lh.WorkerLock.RUnlock()
				if len(workerList) > 0 {
					err := service.ServiceGroupApp.LotusServiceGroup.ModifyWorkerStatus(workerList)
					if err != nil {
						return
					}
				}
			}
		}
	}
}

// DeleteMonitorInfo Periodically clear hardware monitoring information
func DeleteMonitorInfo() {

	spec := "0 1 * * ?"
	c := cron.New()
	c.AddFunc(spec, func() {
		hostMonitorInfo := system.HostMonitorRecordService{}
		err := hostMonitorInfo.DeleteByTime()
		if err != nil {
			log.Println("Failed to delete data from the previous day", zap.Error(err))
			return
		}
		log.Println("Cleaning completed")
	})
	c.Start()
	defer c.Stop()
	select {}
}

func RunP1P2() {

	token, err := global.OpToGatewayClient.GetMinerToken(context.TODO(), &pb.String{Value: global.OpUUID.String()})
	if err != nil {
		log.Println(fmt.Sprintf("GetMinerToken：%s %v", token, err))
		return
	}

	startNo, p2StartNo, p2StartNo1 := 0, 0, 0
	p2EndNo2, p2EndNo := 0, 0
	cupNo := runtime.NumCPU()
	endNo := cupNo - 1
	if cupNo <= 32 {
		endNo -= 4
		p2StartNo = endNo + 1
		p2EndNo = cupNo - 1
	} else {
		endNo -= 16
		p2StartNo = endNo + 1
		p2EndNo = p2StartNo + 7
		p2StartNo1 = p2EndNo + 1
		p2EndNo2 = p2EndNo + 8
	}

	log.Println(fmt.Sprintf("CPUs:%d,P1 bind%d,%d", cupNo, startNo, endNo))
	outs, err := exec.Command("bash", "-c", fmt.Sprintf("%s %s %s %s %s %d %d", define.PathIpfsScriptRunWorker, define.WorkerPort, token.Value, define.MainDisk, global.ROOM_CONFIG.Gateway.IP+":"+define.SlotUnsealedPort, startNo, endNo)).CombinedOutput()
	if err != nil {
		log.Println(fmt.Sprintf("init error：%s %v", string(outs), err))
		return
	}

	outs, err = exec.Command("bash", "-c", "supervisorctl update").CombinedOutput()
	if err != nil {
		log.Println(fmt.Errorf("%s,error: %s", string(outs), err.Error()))
		return
	}

	outs, _ = exec.Command("bash", "-c", fmt.Sprintf("supervisorctl start %s", define.ProgramWorkerTask)).CombinedOutput()

	portMap := make(map[string]string)
	if p2StartNo > 0 && p2EndNo > 0 {
		portMap[define.OpP2Port1] = fmt.Sprintf("%d,%d", p2StartNo, p2EndNo)
	}
	if p2StartNo1 > 0 && p2EndNo2 > 0 {
		portMap[define.OpP2Port2] = fmt.Sprintf("%d,%d", p2StartNo1, p2EndNo2)
	}

	for k, v := range portMap {

		vAr := strings.Split(v, ",")
		start, _ := strconv.Atoi(vAr[0])
		end, _ := strconv.Atoi(vAr[1])
		log.Println(fmt.Sprintf("CPUs:%d,P2 bind %d,%d", cupNo, start, end))
		outs, err = exec.Command("bash", "-c", fmt.Sprintf("%s %s %d %d %s", define.PathIpfsScriptRunWorkerP2, define.MainDisk, start, end, k)).CombinedOutput()
		if err != nil {
			log.Println(fmt.Errorf("worker-p2 init script error：%s %v", string(outs), err))
			return
		}

		outs, err = exec.Command("bash", "-c", "supervisorctl update").CombinedOutput()
		if err != nil {
			log.Println(fmt.Errorf("worker-p2 update：%s,error: %s", string(outs), err.Error()))
			return
		}

		outs, _ = exec.Command("bash", "-c", fmt.Sprintf("supervisorctl restart %s", fmt.Sprintf("p2-%s", v))).CombinedOutput()
	}

}
