package lotusrpc

import (
	"encoding/json"
	"fmt"
	"oplian/define"
	"oplian/utils"
	"strconv"
	"time"
)

func (l *LotusRpc) LotusHeight(token, ip string) (height uint64, err error) {
	bodyMap := make(map[string]interface{})
	resMap := make(map[string]interface{})

	url := ip + ":" + define.LotusPort
	b, err := utils.RequestDo(url, define.ApiRouter, token, NewJsonRpc(define.FilecoinChainHead, nil), time.Second*15)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(b, &bodyMap)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal([]byte(utils.Strval(bodyMap["result"])), &resMap)
	if err != nil {
		return 0, err
	}

	height, _ = strconv.ParseUint(utils.Strval(resMap["Height"]), 10, 64)

	return height, err
}

func (l *LotusRpc) ClientQueryAsk(token, ip, peerId, actor string) (ask AskInfo, err error) {
	var in []interface{}
	in = append(in, peerId)
	in = append(in, actor)
	url := ip + ":" + define.LotusPort
	data, err := utils.RequestDo(url, define.ApiRouter, token, NewJsonRpc(define.FilecoinClientQueryAsk, in), time.Second*15)
	if err != nil {
		return ask, err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return ask, err
	}
	if res.Error != nil {
		return ask, fmt.Errorf("unknown address protocol:%v", res.Error)
	}
	result, _ := json.Marshal(res.Result)
	_ = json.Unmarshal(result, &ask)
	return ask, nil
}

func (l *LotusRpc) StateVerifiedClientStatus(token, ip, addr string) (StoragePower string, err error) {
	var in []interface{}
	in = append(in, addr)
	in = append(in, []string{})
	url := ip + ":" + define.LotusPort
	data, err := utils.RequestDo(url, define.ApiRouter, token, NewJsonRpc(define.FilecoinStateVerifiedClientStatus, in), time.Second*15)
	if err != nil {
		return "", err
	}
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	if res.Result == nil {
		return "", fmt.Errorf("Result为空")
	}

	return utils.SizeStr(res.Result.(string)), nil
}
