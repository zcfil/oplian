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
	"oplian/api_grpc/op"
	"oplian/core"
	"oplian/define"
	"oplian/global"
	"oplian/initialize"
	"oplian/initialize/opinit"
	"oplian/service/lotus/deploy"
	"oplian/service/op/commit"
	"oplian/service/pb"
	"oplian/utils"
	"os"
)

var Run = &cli.Command{
	Name:  "run",
	Usage: "operating system",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "dc-type",
			Usage: "Is it a DC original value host? True sets the machine as the original value host type, " +
				"false does not set it, and defaults to a non original value host",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "storage",
			Usage: "is it a storage device",
			Value: false,
		},
		&cli.StringFlag{
			Name:  "listen-ip",
			Usage: "enter the internal network IP",
			Value: global.LocalIP,
		},
		&cli.StringFlag{
			Name:  "paramters-path",
			Usage: "specify proof parameter path",
		},
		&cli.BoolFlag{
			Name: "not-proof-parameters",
			Usage: "If there are no proof parameters, the program will be downloaded from the file management platform. " +
				"However, due to issues with batch downloads, it is recommended to prepare filecoin proof parameters in/mnt/md0 in advance",
		},
		&cli.BoolFlag{
			Name:  "worker",
			Usage: "is it a computing power machine",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "miner",
			Usage: "is it a miner machine",
			Value: false,
		},
	},
	Action: func(cctx *cli.Context) error {
		go utils.InstallNumactl()

		dcType := cctx.Bool("dc-type")
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

		log.Println("LocalIP: ", global.LocalIP)

		os.MkdirAll(define.PathIpfsData, 0664)
		os.MkdirAll(define.PathIpfsLogs, 0664)
		initialize.OpGenerateFile(define.PathIpfsScript, initialize.ScriptList)
		core.SetConfigRoom()

		grpcServer := grpc.NewServer()
		global.ZC_LOG = core.Zap()
		zap.ReplaceGlobals(global.ZC_LOG)

		pb.RegisterOpServiceServer(grpcServer, new(op.OpServiceImpl))
		lis, err := net.Listen("tcp", ":"+global.ROOM_CONFIG.Op.Port)
		if err != nil {
			log.Fatal(err)
		}

		go opinit.OpC2HeartBeat(context.Background())

		initialize.OpInitGlobalUUID()
		opinit.OpConnectGateWay(context.Background())
		initialize.OpGenerateUUID()
		initialize.OpChmodDirectory()

		os.MkdirAll(define.PathIpfsDataCar, 0664)
		os.MkdirAll(define.PathIpfsLogs, 0664)

		go commit.SealerService.RunOpC2Client()
		go deploy.SectorsRecoverServiceApi.RunSectorSealTask()

		groupArray := initialize.GetHostGroupArray(context.Background())

		initialize.DiskInitialization(groupArray, false, cctx.Bool("worker"), cctx.Bool("miner"))

		err = initialize.CheckDiskProofParameters(cctx.Bool("not-proof-parameters"), cctx.Bool("storage"), cctx.String("paramters-path"))
		if err != nil {
			fmt.Println("Prove that the parameters do not exist. Please add the -- not proof parameters parameter to start, "+
				"but it is not recommended to do so. It is recommended to do so in advance", define.MainDisk,
				" to prepare filecoin proof parameters in the middle, or add -- parameters path to specify the path")
			return err
		}

		initialize.UpdateHostInfo(context.Background(), groupArray, dcType, false, cctx.Bool("storage"))

		go initialize.GetHostMonitorInfo(context.Background())
		go initialize.OpRemountDisk(context.Background())

		//go initialize.RunWorkerP2()

		go initialize.OpRestartLotus(context.Background())
		go opinit.RangePathSectors()

		grpcServer.Serve(lis)
		return nil
	},
}
