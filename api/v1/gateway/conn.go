package gateway

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/gateway/request"
	"oplian/service/pb"
	"oplian/utils"
	"sync"
)

type ConnApi struct{}

var gateConnMap = make(map[string]*grpc.ClientConn)
var gateConnLock sync.RWMutex

func (conn *ConnApi) ConnectGateWay(c *gin.Context) {
	var g request.GateWayInfo
	err := c.ShouldBindJSON(&g)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(g, utils.GateWayInfoVerify)
	if err != nil {
		log.Println(err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	gateConnLock.RLock()
	gateConn := gateConnMap[g.IP+":"+g.Port]
	gateConnLock.RUnlock()
	if gateConn == nil {
		gateConn, err = grpc.Dial(g.IP+":"+g.Port, grpc.WithInsecure())
		gateConnLock.Lock()
		gateConnMap[g.IP+":"+g.Port] = gateConn
		gateConnLock.Unlock()
		if err != nil {
			log.Println(err.Error())
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	client := pb.NewGateServiceClient(gateConn)
	sClient := pb.NewSlotGateServiceClient(gateConn)
	if _, err = client.OplianHeartbeat(context.Background(), &pb.String{Value: "Connection testing"}); err != nil {
		log.Println(err.Error())
		response.FailWithMessage(err.Error(), c)
		gateConnLock.Lock()
		gateConn.Close()
		gateConnMap[g.IP+":"+g.Port] = nil
		gateConnLock.Unlock()
		return
	}
	global.GateWayClinets.LockRW.Lock()
	global.GateWayClinets.Info[g.GateWayId] = &global.GateWayInfo{global.GateClient{client, sClient}, g.GateWayId, g.IP, g.Port, g.Token, false}
	global.GateWayClinets.LockRW.Unlock()
	response.Ok(c)
}

func (conn *ConnApi) Ping(c *gin.Context) {
	response.OkWithMessage("ping success", c)
}
