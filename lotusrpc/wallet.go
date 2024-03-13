package lotusrpc

import (
	"encoding/json"
	"errors"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/utils"
	"time"
)

type OpLotusService struct{}

// WalletBalance
// @author: nathan
// @function: WalletBalance
// @description: Get wallet balance
// @param: httpToken, address string
// @return: float64, error
func (l *LotusRpc) WalletBalance(token, ip, address string) (balance float64, err error) {
	param := []string{address}
	data, err := utils.RequestDo(ip+":"+define.LotusPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinWalletBalance, param), time.Second*15)
	if err != nil {
		log.Println(ip+":"+define.LotusPort, token, param)
		global.ZC_LOG.Error(err.Error())
		return 0, err
	}
	//log.Println("钱包结果：", string(data))
	var res JsonRpcResult
	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	if res.Error != nil {
		return 0, err
	}
	return utils.NanoOrAttoToFIL(res.Result.(string), utils.AttoFIL)
}

// WalletNew
// @author: nathan
// @function: WalletNew
// @description: Get wallet balance
// @param: httpToken, walletType
// @return: address string, err error
func (l *LotusRpc) WalletNew(token, ip, walletType string) (address string, err error) {
	var param []string
	switch walletType {
	case define.WalletTypeBls:
		param = append(param, define.WalletTypeBls)
	case define.WalletTypeSecp256k1:
		param = append(param, define.WalletTypeSecp256k1)
	default:
		return "", errors.New("Unknown wallet type！")
	}
	res, err := utils.RequestDo(ip+":"+define.LotusPort, define.ApiRouter, token, NewJsonRpc(define.FilecoinWalletNew, param), time.Second*15)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return "", errors.New("Wallet creation failure：" + err.Error())
	}
	var result JsonRpcResult
	err = json.Unmarshal(res, &result)
	if err != nil {
		return "", err
	}
	return result.Result.(string), nil
}
