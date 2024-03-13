package lotusrpc

import (
	"encoding/json"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/utils"
	"time"
)

// StateMinerInfo
// @author: nathan
// @function: StateMinerInfo
// @description: Get miner information
// @param: httpToken, address string
// @return: ActorControl, error)
func (l *LotusRpc) StateMinerInfo(token, ip, actor string) (control ActorControl, err error) {
	var in []interface{}
	in = append(in, actor)
	in = append(in, make([]map[string]string, 0))

	data, err := utils.RequestDo(ip+":"+define.LotusPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinStateMinerInfo, in), time.Second*15)
	if err != nil {
		log.Println(token, ip, actor, "JsonRpcResult：", string(data))
		global.ZC_LOG.Error(err.Error())
		return control, err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		log.Println(token, ip, actor, "JsonRpcResult：", string(data))
		return ActorControl{}, err
	}
	result, _ := json.Marshal(res.Result)
	_ = json.Unmarshal(result, &control)
	return
}

// StateAccountKey
// @author: nathan
// @function: StateAccountKey
// @description: f01234 Transfer wallet address
// @param: httpToken, address string
// @return: ActorControl, error)
func (l *LotusRpc) StateAccountKey(token, ip, actor string) (address string, err error) {
	var in []interface{}
	in = append(in, actor)
	in = append(in, []string{})

	data, err := utils.RequestDo(ip+":"+define.LotusPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinStateAccountKey, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return "", err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return "", err
	}

	if res.Result == nil {
		return actor, err
	}
	return res.Result.(string), err
}

func (l *LotusRpc) PreeID(token, ip, port string) (ID string, err error) {
	bodyMap := make(map[string]interface{})

	url := ip + ":" + port
	b, err := utils.RequestDo(url, define.ApiRouter, token, NewJsonRpc(define.FilecoinID, nil), time.Second*15)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(b, &bodyMap)
	if err != nil {
		return "", err
	}

	return utils.Strval(bodyMap["result"]), err
}

func (l *LotusRpc) StateMinerSectorCount(token, ip, actor string) (sectors MinerSectors, err error) {
	var in []interface{}
	in = append(in, actor)
	in = append(in, []string{})
	data, err := utils.RequestDo(ip+":"+define.LotusPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinStateMinerSectorCount, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return sectors, err
	}

	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return MinerSectors{}, err
	}
	result, _ := json.Marshal(res.Result)
	_ = json.Unmarshal(result, &sectors)
	return
}

func (l *LotusRpc) StateMinerPower(token, ip, actor string) (power Power, err error) {
	var in []interface{}
	in = append(in, actor)
	in = append(in, []string{})
	data, err := utils.RequestDo(ip+":"+define.LotusPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinStateMinerPower, in), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return power, err
	}

	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return power, err
	}
	result, _ := json.Marshal(res.Result)
	var pw PowerMap
	_ = json.Unmarshal(result, &pw)
	if p, ok := pw.MinerPower["QualityAdjPower"]; ok {
		power.MinerPower = p
	}
	if p, ok := pw.TotalPower["QualityAdjPower"]; ok {
		power.TotalPower = p
	}
	return
}
