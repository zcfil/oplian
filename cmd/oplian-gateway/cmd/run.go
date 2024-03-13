package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"oplian/api_grpc/gateway"
	"oplian/core"
	"oplian/define"
	"oplian/global"
	"oplian/initialize"
	"oplian/service"
	"oplian/service/lotus/deploy"
	"oplian/service/pb"
	"oplian/utils"
)

var Run = &cli.Command{
	Name:  "run",
	Usage: "operating system",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "listen-ip",
			Usage: "enter the internal IP",
			Value: global.LocalIP,
		},
	},
	Action: func(cctx *cli.Context) error {

		core.ViperRoom("config/config_room.yaml")
		global.ZC_VP = core.Viper("config/config.yaml")
		global.ZC_LOG = core.Zap()
		zap.ReplaceGlobals(global.ZC_LOG)

		global.ZC_DB = initialize.Gorm()
		initialize.Timer()
		initialize.DBList()
		if global.ZC_DB != nil {
			initialize.RegisterTables(global.ZC_DB)
			log.Println(initialize.InitMysqlData())
			db, _ := global.ZC_DB.DB()
			defer db.Close()
		}

		global.LocalIP = cctx.String("listen-ip")
		if global.LocalIP == "" {
			inIp := utils.GetIntranetIPList()
			switch len(inIp) {
			case 0:
				fmt.Print("No local intranet IP was obtained. Please manually enter the intranet IP to start using listen ip")
				return errors.New("no local IP obtained")
			case 1:
				global.LocalIP = inIp[0]
			default:
				fmt.Println("Obtaining multiple local intranet IPs ", inIp)
				fmt.Print("Please manually enter the internal network IP to start using listen ip")
				return errors.New("to many local IP obtained")
			}
		}

		initialize.CheckDiskProofParametersGateway()
		initialize.OpGenerateFile(define.PathIpfsScript, initialize.ScriptList)
		isNew, err := initialize.OpGenerateConfigFile(define.PathIpfsConfig, initialize.ConfigList)
		if err != nil {
			log.Fatal("Error generating corresponding config file")
		} else {
			if isNew {
				log.Fatal("The corresponding config file has been generated. " +
					"Please configure the corresponding config file information correctly and restart again")
			}
		}
		log.Println("---step5---")
		
		service.ServiceGroupApp.GatewayServiceGroup.DownloadAddressFile([]string{define.IPAddressFile}, define.PathIpfsConfig, define.DownLoadAddressOut)
		
		service.ServiceGroupApp.GatewayServiceGroup.DownloadAddressFile([]string{define.ProgramLotus.String(), define.ProgramMiner.String(),
			define.ProgramWorkerTask.String(), define.ProgramBoost.String(), define.ProgramBoostd.String()}, define.PathIpfsProgram, define.DownLoadAddressOut)
		
		service.ServiceGroupApp.GatewayServiceGroup.DownloadAddressFile([]string{define.OPlianOpFile, define.OPlianOpC2File}, define.PathOplian, define.DownLoadAddressOut)

		initialize.GatewayGenerateUUID()
		initialize.OpChmodDirectory()
		initialize.UpdateMachineRoomInfo(context.Background())

		go initialize.OpHeartBeat(context.Background())

		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(initialize.GateWayFilter))
		pb.RegisterGateServiceServer(grpcServer, new(gateway.GateWayServiceImpl))
		pb.RegisterSlotGateServiceServer(grpcServer, new(slot_gateway.SlotGateWayServiceImpl))
		lis, err := net.Listen("tcp", ":"+global.ROOM_CONFIG.Gateway.Port)
		if err != nil {
			log.Fatal(err)
		}

		go initialize.ConnectWeb(context.Background())
		go deploy.SectorsRecoverServiceApi.RedoSectorsTask()
		go deploy.WorkerClusterServiceApi.RunC2Task()
		go initialize.BeginHostPatrol()

		go initialize.CheckPatrolAndTest()
		go initialize.BeginCheckBadSector()
		go core.RunGatewayServer()
		go initialize.CheckLotusHeart()

		service.ServiceGroupApp.GatewayServiceGroup.SynFileToGateWay()
		log.Println("---gateway successfully started---")
		grpcServer.Serve(lis)
		return nil
	},
}
