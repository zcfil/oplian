package global

import (
	"oplian/service/pb"
	"sync"
)

var OpC2Clients = OpC2Client{
	Info: make(map[string]*OpC2Info),
}

type OpC2Client struct {
	Info   map[string]*OpC2Info
	LockRW sync.RWMutex
}

type OpC2Info struct {
	Client     pb.OpC2ServiceClient
	Ip         string
	Port       string
	Disconnect bool
}

func (op *OpC2Client) GetOpC2Client(OpId string) (pb.OpC2ServiceClient, bool) {
	op.LockRW.RLock()
	defer op.LockRW.RUnlock()
	info, ok := op.Info[OpId]
	if !ok {
		return nil, true
	}
	return info.Client, info.Disconnect
}
func (op *OpC2Client) GetOpC2IP(OpId string) string {
	op.LockRW.RLock()
	defer op.LockRW.RUnlock()
	info, ok := op.Info[OpId]
	if !ok {
		return ""
	}
	return info.Ip
}
func (op *OpC2Client) C2IsDisconnect(OpId string) bool {
	op.LockRW.RLock()
	defer op.LockRW.RUnlock()
	info, ok := op.Info[OpId]
	if !ok {
		return true
	}
	return info.Disconnect
}
