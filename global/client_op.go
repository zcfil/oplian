package global

import (
	"oplian/service/pb"
	"sync"
)

var OpClinets = OpClient{
	Info: make(map[string]*OpInfo),
}

type OpClient struct {
	Info   map[string]*OpInfo
	LockRW sync.RWMutex
}

type OpInfo struct {
	Clinet     pb.OpServiceClient
	SlotClient pb.SlotOpServiceClient
	Ip         string
	Port       string
	Disconnect bool
	OpId       string
}

// 获取OP客户端
func (op *OpClient) GetOpClient(OpId string) (pb.OpServiceClient, bool) {
	op.LockRW.RLock()
	defer op.LockRW.RUnlock()
	info, ok := op.Info[OpId]
	if !ok {
		return nil, true
	}
	return info.Clinet, info.Disconnect
}

func (op *OpClient) SetOpClient(OpId string, info *OpInfo) {
	op.LockRW.Lock()
	defer op.LockRW.Unlock()
	op.Info[OpId] = info
}

func (op *OpClient) GetOpClientList() (f map[string]*OpInfo) {
	op.LockRW.RLock()
	defer op.LockRW.RUnlock()
	info := make(map[string]*OpInfo)
	if len(op.Info) > 0 {
		for k, v := range op.Info {
			if k != "" {
				info[k] = v
			}
		}
	}
	return info
}

// 获取在线OP客户端列表
func (op *OpClient) OnLineList() []pb.OpServiceClient {
	op.LockRW.RLock()
	defer op.LockRW.RUnlock()
	var ops []pb.OpServiceClient
	for _, info := range op.Info {
		if info.Disconnect {
			continue
		}
		ops = append(ops, info.Clinet)
	}
	return ops
}
func (op *OpClient) GetOpIP(OpId string) string {
	op.LockRW.RLock()
	defer op.LockRW.RUnlock()
	info, ok := op.Info[OpId]
	if !ok {
		return ""
	}
	return info.Ip
}
func (op *OpClient) IsDisconnect(OpId string) bool {
	op.LockRW.RLock()
	defer op.LockRW.RUnlock()
	info, ok := op.Info[OpId]
	if !ok {
		return true
	}
	return info.Disconnect
}

func (op *OpClient) SetDisconnect(OpId string, disconnect bool) {

	OpClinets.LockRW.Lock()
	defer OpClinets.LockRW.Unlock()
	if info, ok := OpClinets.Info[OpId]; ok {
		info.Disconnect = disconnect
		OpClinets.Info[OpId] = info
	}
}

var OpC2ToOp pb.OpServiceClient
