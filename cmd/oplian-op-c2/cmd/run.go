package cmd

import (
	"context"
	"errors"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"log"
	"net"
	"oplian/api_grpc/op"
	"oplian/core"
	"oplian/global"
	"oplian/initialize"
	"oplian/service/op/commit"
	"oplian/service/pb"
)

var RunOpC2 = &cli.Command{
	Name:  "run",
	Usage: "operating system",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "port",
			Usage: "run Port number",
			Value: "1",
		},
	},
	Action: func(ct *cli.Context) error {

		log.Println("opc2 port:", ct.Args().Get(0))
		if ct.Args().Get(0) == "" {
			return errors.New("port is nil")
		}

		commit.SealerService.InitOpC2Uid()
		RunOpC2Server(ct.Args().Get(0))
		return nil
	},
}

func RunOpC2Server(port string) {

	grpcServer := grpc.NewServer()
	core.SetConfigRoom()
	global.ZC_LOG = core.Zap()

	global.ROOM_CONFIG.OpC2.Port = port
	log.Println("global.ROOM_CONFIG.OpC2.Port", global.ROOM_CONFIG.OpC2.Port)
	pb.RegisterOpC2ServiceServer(grpcServer, new(op.OpC2ServiceImpl))
	lis, err := net.Listen("tcp", ":"+global.ROOM_CONFIG.OpC2.Port)
	if err != nil {
		log.Fatal(err)
	}

	initialize.OpC2ConnectOp(context.Background())
	grpcServer.Serve(lis)

}
