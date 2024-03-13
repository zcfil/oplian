package global

import (
	"oplian/service/pb"
	"sync"
)

var GateWayClinets = GateWayClient{
	Info: make(map[string]*GateWayInfo),
}

type GateWayClient struct {
	Info   map[string]*GateWayInfo
	LockRW sync.RWMutex
}
type GateClient struct {
	pb.GateServiceClient
	pb.SlotGateServiceClient
}
type GateWayInfo struct {
	GateClient
	GateWayId  string
	Ip         string
	Port       string
	Token      string
	Disconnect bool
}

func (gate *GateWayClient) Gets() []*GateWayInfo {
	gate.LockRW.RLock()
	defer gate.LockRW.RUnlock()
	var res []*GateWayInfo
	for _, v := range gate.Info {
		res = append(res, v)
	}
	return res
}

func (gate *GateWayClient) GetGateWayClinet(gid string) *GateClient {
	gate.LockRW.RLock()
	defer gate.LockRW.RUnlock()
	info, ok := gate.Info[gid]
	if !ok {
		return nil
	}
	return &info.GateClient
}

func (gate *GateWayClient) IsDisconnect(OpId string) bool {
	gate.LockRW.RLock()
	defer gate.LockRW.RUnlock()
	info, ok := gate.Info[OpId]
	if !ok {
		return true
	}
	return info.Disconnect
}

var OpToGatewayClient struct {
	pb.GateServiceClient
	pb.SlotGateServiceClient
}
