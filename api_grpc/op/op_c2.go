package op

import (
	"context"
	"log"
	"oplian/service/gateway/commit2"
	"oplian/service/pb"
)

type OpC2ServiceImpl struct{}

func (p *OpC2ServiceImpl) Heartbeat(ctx context.Context, args *pb.String) (*pb.String, error) {
	reply := &pb.String{Value: "Heartbeat check succeeded !" + args.GetValue()}
	log.Println(reply.String())
	return reply, nil
}

// Commit2TaskRun Running C2 tasks
func (p *OpC2ServiceImpl) Commit2TaskRun(ctx context.Context, args *pb.SealerParam) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	err := commit2.Commit2ServiceApi.RunCommit2(args)
	if err != nil {
		return res, err
	}

	return res, nil
}
