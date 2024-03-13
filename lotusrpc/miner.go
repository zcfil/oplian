package lotusrpc

import (
	"encoding/json"
	"fmt"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/utils"
	"time"
)

// PledgeSector
// @author: nathan
// @function: PledgeSector
// @description: Pledge sector
// @param: httpToken string,  number uint64
// @return: Sector, error)
func (l *LotusRpc) PledgeSector(token, ip string) (sector Sector, err error) {
	var in []interface{}
	//in = append(in, number)

	data, err := utils.RequestDo(ip+":"+define.MinerPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinPledgeSector, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return sector, err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return sector, err
	}
	result, _ := json.Marshal(res.Result)
	_ = json.Unmarshal(result, &sector)
	return
}

// ActorAddress
// @author: nathan
// @function: ActorAddress
// @description: Get node number
// @param: httpToken string,  number uint64
// @return: string, error
func (l *LotusRpc) ActorAddress(token, ip string) (actor string, err error) {
	var in []interface{}
	//in = append(in, number)

	data, err := utils.RequestDo(ip+":"+define.MinerPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinPledgeSector, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return actor, err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return actor, err
	}
	return res.Result.(string), nil
}

// MaxSectorNumber
// @author: nathan
// @function: MaxSectorNumber
// @description: Maximum sector number
// @param: httpToken string
// @return: uint64, error
func (l *LotusRpc) MaxSectorNumber(token, ip string) (number uint64, err error) {
	var in []interface{}

	data, err := utils.RequestDo(ip+":"+define.MinerPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinSectorNumAssignerMeta, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return 0, err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	var next SectorNext
	result, _ := json.Marshal(res.Result)
	_ = json.Unmarshal(result, &next)
	return next.Next - 1, nil
}

// DealsImportData
// @author: nathan
// @function: DealsImportData
// @description: Offline import
// @param: httpToken,ip, cid, file string
// @return: uint64, error)
func (l *LotusRpc) DealsImportData(token, ip, cid, file string) (result string, err error) {
	var in []interface{}
	msg := make(map[string]string)
	msg["/"] = cid
	in = append(in, msg)
	in = append(in, file)

	data, err := utils.RequestDo(ip+":"+define.MinerPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinDealsImportData, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return "", err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	//var next SectorNext
	//result, _ := json.Marshal(res.Result)
	//_ = json.Unmarshal(result, &next)
	return "", nil
}

// SealingSchedDiag
// @author: nathan
// @function: SealingSchedDiag
// @description: Ready queue task information
// @param: httpToken,ip string
// @return: sched SchedDiagInfo, err error
func (l *LotusRpc) SealingSchedDiag(token, ip string) (sched SchedDiagInfo, err error) {
	var in []interface{}
	in = append(in, false)
	data, err := utils.RequestDo(ip+":"+define.MinerPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinSealingSchedDiag, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return SchedDiagInfo{}, err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return SchedDiagInfo{}, err
	}
	result, _ := json.Marshal(res.Result)
	_ = json.Unmarshal(result, &sched)
	return
}

// SectorStorage
// @author: nathan
// @function: SectorStorage
// @description: Sector storage path
// @param: httpToken, ip string
// @return: paths []string, err error
func (l *LotusRpc) SectorStorage(token, ip string, miner, number uint64) (paths []string, err error) {
	var in []interface{}
	mp := make(map[string]uint64)
	mp["miner"] = miner
	mp["number"] = number
	in = append(in, mp)
	data, err := utils.RequestDo(ip+":"+define.MinerPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinSectorStorage, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return nil, err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	if res.Result == nil {
		return nil, nil
	}
	for _, v := range res.Result.([]interface{}) {
		paths = append(paths, v.(string))
	}
	return paths, nil
}

// SealingAbort
// @author: nathan
// @function: SealingAbort
// @description: Terminates the sector task
// @param: (token, ip string, miner, number uint64, ID string)
// @return:err error
func (l *LotusRpc) SealingAbort(token, ip string, miner, number uint64, ID string) (err error) {
	var in []interface{}
	param := make(map[string]interface{})
	mp := make(map[string]uint64)
	mp["miner"] = miner
	mp["number"] = number
	param["sector"] = mp
	param["ID"] = ID
	in = append(in, param)
	data, err := utils.RequestDo(ip+":"+define.MinerPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinSealingAbort, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return err
	}
	log.Println("SealingAbort:", res.Result, res.Error)
	log.Println("data:", string(data))
	if res.Result == nil {
		return nil
	}
	if res.Error != nil {
		return fmt.Errorf("%s", res.Error)
	}
	return nil
}

// SectorRemove
// @author: nathan
// @function: SectorRemove
// @description: Remove sector
// @param: (token, ip string, number uint64)
// @return: err error
func (l *LotusRpc) SectorRemove(token, ip string, number uint64) (err error) {
	var in []interface{}
	in = append(in, number)
	data, err := utils.RequestDo(ip+":"+define.MinerPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinSectorRemove, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return err
	}
	log.Println("SectorRemove:", res.Result, res.Error)
	if res.Result == nil {
		return nil
	}
	if res.Error != nil {
		return fmt.Errorf("%s", res.Error)
	}
	return nil
}

// SectorStoreFile
// @author: nathan
// @function: SectorStoreFile
// @description: Whether the sector file exists
// @param: token, ip string, miner, number uint64)
// @return: exist int, err error
func (l *LotusRpc) SectorStoreFile(token, ip string, miner, number uint64) (exist int, err error) {
	var in []interface{}
	mp := make(map[string]uint64)
	mp["miner"] = miner
	mp["number"] = number
	in = append(in, mp)
	data, err := utils.RequestDo(ip+":"+define.MinerPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinSectorStoreFile, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return 0, err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	log.Println(res.Result)
	if res.Result == nil {
		return 0, nil
	}

	return res.Result.(int), nil
}

// SynC2Result
// @author: nathan
// @function: SynC2Result
// @description: Synchronize C2 results
// @param: token, ip string, filePath, filename string, fileData []byte
// @return: error
func (l *LotusRpc) SynC2Result(token, ip, filePath, filename string, fileData []byte) error {

	var in []interface{}
	in = append(in, filePath, filename, fileData)
	data, err := utils.RequestDo(ip+":"+define.MinerPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinSynC2Result, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return err
	}

	log.Println(fmt.Sprintf("SynC2Result Synchronization result:%+v", res))

	return nil
}
