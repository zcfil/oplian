package opinit

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"log"
	"oplian/auth"
	"oplian/service/gateway/commit2"
	"oplian/service/pb"
	"time"

	"oplian/global"
)

// OpC2HeartBeat op Heartbeat Check op-c2
func OpC2HeartBeat(ctx context.Context) {

	for {
		select {
		case <-time.After(time.Second * 10):
			var i = 0
			if len(global.OpC2Clients.Info) > 0 {
				global.OpC2Clients.LockRW.Lock()
				for k, v := range global.OpC2Clients.Info {
					_, err := v.Client.Heartbeat(ctx, &pb.String{Value: "op-c2"})
					if err != nil {

						if ok, _ := commit2.GetWorkersClient(k); ok {
							commit2.DelWorkersClient(k)
						}

						v.Disconnect = true
						continue
					}

					if ok, _ := commit2.GetWorkersClient(k); ok {
						v.Disconnect = false
						i++
					}
				}
				global.OpC2Clients.LockRW.Unlock()
				log.Println("on-line OP-C2：", i)
			}

		case <-ctx.Done():
			return
		}
	}
}

// OpConnectGateWay op to GateWay
func OpConnectGateWay(ctx context.Context) {
	conn, err := grpc.Dial(global.ROOM_CONFIG.Gateway.IP+":"+global.ROOM_CONFIG.Gateway.Port, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth.Authentication{Token: global.ROOM_CONFIG.Op.Token}))
	if err != nil {
		log.Println("Abnormal network connection: Dial: ", global.ROOM_CONFIG.Gateway.IP+":"+global.ROOM_CONFIG.Gateway.Port)
	}

	var gatewayId *pb.String
	global.OpToGatewayClient.GateServiceClient = pb.NewGateServiceClient(conn)
	global.OpToGatewayClient.SlotGateServiceClient = pb.NewSlotGateServiceClient(conn)
	if gatewayId, err = global.OpToGatewayClient.OpConnect(ctx, &pb.RequestConnect{OpId: global.OpUUID.String(), Port: global.ROOM_CONFIG.Op.Port}); err != nil {
		log.Println("Abnormal network connection: OpConnect: ", global.ROOM_CONFIG.Gateway.IP+":"+global.ROOM_CONFIG.Gateway.Port, err.Error())
		if conn != nil {
			conn.Close()
		}
		conn = nil
	} else {
		global.GateWayID, _ = uuid.FromString(gatewayId.Value)
	}


	go func() {
		for {
			select {
			case <-time.After(time.Second * 10):
				//init
				if conn == nil {
					conn, err = grpc.Dial(global.ROOM_CONFIG.Gateway.IP+":"+global.ROOM_CONFIG.Gateway.Port, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth.Authentication{Token: global.ROOM_CONFIG.Op.Token}))
					if err != nil {
						log.Println("Abnormal network connection: Dial: ", global.ROOM_CONFIG.Gateway.IP+":"+global.ROOM_CONFIG.Gateway.Port)
					}
					global.OpToGatewayClient.GateServiceClient = pb.NewGateServiceClient(conn)
					global.OpToGatewayClient.SlotGateServiceClient = pb.NewSlotGateServiceClient(conn)
				}
				//update
				if _, err = global.OpToGatewayClient.OpHeartbeat(ctx, &pb.String{Value: global.OpUUID.String()}); err != nil {
					if gatewayId, err = global.OpToGatewayClient.OpConnect(ctx, &pb.RequestConnect{OpId: global.OpUUID.String(), Port: global.ROOM_CONFIG.Op.Port}); err != nil {
						log.Println("Abnormal network connection: OpConnect: ", global.ROOM_CONFIG.Gateway.IP+":"+global.ROOM_CONFIG.Gateway.Port, err.Error())
						if conn != nil {
							conn.Close()
						}
						conn = nil
					} else {
						global.GateWayID, _ = uuid.FromString(gatewayId.Value)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
