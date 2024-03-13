package initialize

import (
	"context"
	"google.golang.org/grpc"
)

// interceptor
func GateWayFilter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//md, _ := metadata.FromIncomingContext(ctx)
	//log.Println("Hello I'm interceptor", md["token"])
	//return nil, err
	return handler(ctx, req)
}
