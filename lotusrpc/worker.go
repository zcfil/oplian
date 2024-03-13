package lotusrpc

import (
	"encoding/json"
	"oplian/define"
	"oplian/global"
	"oplian/utils"
	"time"
)

// StorageAddLocal
// @author: nathan
// @function: StorageAddLocal
// @description: Add local storage
// @param: token, ip, path string, workerType define.WorkerType
// @return: uint64, error)
func (l *LotusRpc) StorageAddLocal(token, ip, path string, workerType define.WorkerType) (err error) {
	var in []interface{}
	var port string
	switch workerType {
	case define.TaskWorker:
		port = define.WorkerPort
	case define.StorageWorker:
		port = define.StoragePort
	case define.MinerWorker:
		port = define.MinerPort
	}
	in = append(in, path)
	data, err := utils.RequestDo(ip+":"+port, define.ApiRouter, token, NewJsonRpc(define.FilecoinStorageAddLocal, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return err
	}
	if res.Error != nil {
		return
	}
	return nil
}
