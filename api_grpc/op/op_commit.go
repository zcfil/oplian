package op

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"oplian/api_grpc/gateway"
	"oplian/define"
	"oplian/global"
	"oplian/service/gateway/commit2"
	"oplian/service/lotus/deploy"
	"oplian/service/op/commit"
	"oplian/service/pb"
	"strings"
	"sync"
	"time"
)

var opC2ConnMap = make(map[string]*grpc.ClientConn)
var opC2ConnLock sync.RWMutex

// OpC2Connect Connect op-c2
func (p *OpServiceImpl) OpC2Connect(ctx context.Context, args *pb.RequestConnect) (*emptypb.Empty, error) {

	f, _ := peer.FromContext(ctx)
	IP, _, err := net.SplitHostPort(strings.TrimSpace(f.Addr.String()))
	if err != nil {
		return &empty.Empty{}, err
	}
	opC2ConnLock.RLock()
	opC2Conn := gateway.OpConnMap[IP+":"+args.Port]
	opC2ConnLock.RUnlock()
	if opC2Conn == nil {
		opC2Conn, err = grpc.Dial(IP+":"+args.Port, grpc.WithInsecure())
		opC2ConnLock.Lock()
		opC2ConnMap[IP+":"+args.Port] = opC2Conn
		opC2ConnLock.Unlock()
		if err != nil {
			return &empty.Empty{}, err
		}
	}

	client := pb.NewOpC2ServiceClient(opC2Conn)
	if _, err = client.Heartbeat(ctx, &pb.String{Value: "opc2"}); err != nil {
		opC2ConnLock.Lock()
		opC2Conn.Close()
		opC2ConnMap[IP+":"+args.Port] = nil
		opC2ConnLock.Unlock()
	}

	global.OpC2Clients.LockRW.Lock()
	s, _ := json.Marshal(args)
	log.Println("OpC2Connect", zap.String("info", string(s)))
	global.OpC2Clients.Info[args.OpId] = &global.OpC2Info{
		client,
		IP,
		args.Port,
		false,
	}
	global.OpC2Clients.LockRW.Unlock()

	host := IP + ":" + args.Port
	if ok, _ := commit2.GetWorkersClient(args.OpId); !ok {
		commit2.SetWorkersClient(&commit2.WorkerInfo{OpC2Id: args.OpId, Host: host, GpuUse: false})
	}

	return &empty.Empty{}, nil
}

// OpC2Heartbeat OPC2 heartbeat check call
func (p *OpServiceImpl) OpC2Heartbeat(ctx context.Context, args *pb.String) (*pb.String, error) {
	reply := &pb.String{Value: "opc2 OpC2Heartbeat check succeeded !" + args.GetValue()}

	_, dis := global.OpC2Clients.GetOpC2Client(args.GetValue())
	if dis {
		return nil, errors.New("Op：" + args.GetValue() + " not online！")
	}

	return reply, nil
}

// Commit2TaskResult Obtain C2 task results
func (p *OpServiceImpl) Commit2TaskResult(ctx context.Context, args *pb.SectorID) (*pb.String, error) {

	return global.OpToGatewayClient.Commit2TaskResult(ctx, args)
}

// C2FileSynStatus Obtain C2 task results
func (p *OpServiceImpl) C2FileSynStatus(ctx context.Context, args *pb.C2SectorID) (*pb.String, error) {

	return global.OpToGatewayClient.C2FileSynStatus(ctx, args)
}

// AddC2Task Add C2 task
func (p *OpServiceImpl) AddC2Task(ctx context.Context, args *pb.SectorID) (*pb.String, error) {

	return global.OpToGatewayClient.AddC2Task(ctx, args)
}

// DelC2Task Delete C2 task
func (p *OpServiceImpl) DelC2Task(ctx context.Context, args *pb.SectorID) (*pb.String, error) {

	_, err := global.OpToGatewayClient.DelC2Task(ctx, args)
	if err != nil {
		return &pb.String{}, err
	}

	return &pb.String{}, nil
}

// CompleteCommit2 Complete C2 task
func (p *OpServiceImpl) CompleteCommit2(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	ok, _ := commit2.GetWorkersClient(args.OpId)
	if ok {
		log.Println("CompleteCommit2 clear occupied OPC2:", args.OpId)
		commit2.SetWorkersClient(&commit2.WorkerInfo{OpC2Id: args.OpId, Host: args.Host, GpuUse: false})
	} else {
		log.Println("CompleteCommit2 clearing occupied OPC2 failed:", args.OpId)
		return res, fmt.Errorf("CompleteCommit2 clearing occupied OPC2 failed:%s", args.OpId)
	}

	_, err := global.OpToGatewayClient.CompleteCommit2(ctx, args)
	if err != nil {
		return res, err
	}

	return res, nil

}

// Commit2TaskRun Running C2 tasks
func (p *OpServiceImpl) Commit2TaskRun(ctx context.Context, args *pb.SealerParam) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	client, dis := global.OpC2Clients.GetOpC2Client(args.OpC2Id)
	if dis {
		log.Println("Commit2TaskRun OpC2Clients Connection failed:" + args.OpC2Id)
		return res, errors.New("Commit2TaskRun OpC2Clients Connection failed:" + args.OpC2Id)
	}

	c2Port := strings.Split(args.Host, ":")
	sealPort := commit.Sl.GetSealPort(c2Port[1])
	log.Println(fmt.Sprintf("Commit2TaskRun c2Port:%s,sealPort:%s", c2Port, sealPort))
	if sealPort != "" {

		ok, wc := commit2.GetWorkersClient(args.OpC2Id)
		if ok && wc.GpuUse {
			return res, nil
		}

		commit2.SetWorkersClient(&commit2.WorkerInfo{OpC2Id: args.OpC2Id, Host: args.Host, GpuUse: true})
		ok, wc = commit2.GetWorkersClient(args.OpC2Id)
		if ok {
			log.Println("Commit2TaskRun Occupy NewOPC2:", wc.OpC2Id, wc.GpuUse, args.Sector.Id.Number)
			args.OpMainDisk = define.MainDisk
			args.SealPort = sealPort

			_, err := client.Commit2TaskRun(ctx, args)
			if err != nil {
				log.Println("Commit2TaskRun err:", err)
				return nil, err
			}
		} else {
			return res, errors.New("Commit2TaskRun GetWorkersClient failed:" + args.OpC2Id)
		}

	} else {
		log.Println("Commit2TaskRun Connection oplianSectorSeal failed:" + args.OpC2Id)
		return res, errors.New("Commit2TaskRun Connection oplianSectorSeal failed:" + args.OpC2Id)
	}

	return res, nil
}

// GetOpC2Client Return to opC2 client
func (p *OpServiceImpl) GetOpC2Client(ctx context.Context, args *pb.OpC2Client) (*pb.OpC2Client, error) {

	res := make([]*pb.OpInfo, 0)
	wcList := commit2.GetWorkersClientList()
	for _, v := range wcList {
		if v.GpuUse {
			if time.Now().Sub(v.TimeOut).Minutes() > 5 {
				log.Println(fmt.Sprintf("Commit2TaskRun Overtime occupying OPC2:%s, reassigning tasks", v.OpC2Id))
				commit2.SetWorkersClient(&commit2.WorkerInfo{OpC2Id: v.OpC2Id, Host: v.Host, GpuUse: false})
				continue
			}
		}
		hostAr := strings.Split(v.Host, ":")
		res = append(res, &pb.OpInfo{OpId: v.OpC2Id, Ip: hostAr[0], Port: hostAr[1], GpuUse: v.GpuUse})
	}

	return &pb.OpC2Client{OpInfo: res}, nil

}

// C2FileSynLotus C2 task file synchronization to lotus
func (p *OpServiceImpl) C2FileSynLotus(ctx context.Context, args *pb.FileInfo) (*pb.String, error) {

	args.OpId = global.OpUUID.String()
	return global.OpToGatewayClient.C2FileSynLotus(ctx, args)
}

// RedoSectorsTask Redo sector task
func (p *OpServiceImpl) RedoSectorsTask(ctx context.Context, args *pb.SectorsTask) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}

	_, err := deploy.SectorsRecoverServiceApi.WorkerRedoTask(args)
	if err != nil {
		return res, err
	}

	return res, nil
}

// RunOpC2 Start OPC2
func (p *OpServiceImpl) RunOpC2(ctx context.Context, args *pb.String) (*pb.String, error) {

	res := &pb.String{}
	return res, commit.SealerService.RunOpC2Client()
}

// StopOpC2 Stop OPC2
func (p *OpServiceImpl) StopOpC2(ctx context.Context, args *pb.String) (*pb.String, error) {
	return &pb.String{}, commit.SealerService.StopOpC2Client()
}

// RedoC2Task Redo C2 task
func (p *OpServiceImpl) RedoC2Task(ctx context.Context, args *pb.String) (*pb.String, error) {
	return global.OpToGatewayClient.RedoC2Task(ctx, args)
}
