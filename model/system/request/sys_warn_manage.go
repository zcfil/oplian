package request

import (
	"oplian/model/common/request"
	"oplian/model/system"
)

type WarnReq struct {
	Id            int64  `json:"id"`
	WarnKeyWord   string `json:"warn_key_word"`
	WarnType      int    `json:"warn_type"`
	WarnStatus    int    `json:"warn_status"`
	Ip            string `json:"ip"`
	AssetsKeyWord string `json:"assets_key_word"`
	StrategiesId  string `json:"strategies_id"`
	ComputerType  int    `json:"computer_type"`
	RoomId        string `json:"room_id"`
	BeginTime     string `json:"begin_time"`
	EndTime       string `json:"end_time"`
	request.PageInfo
}

type StrategyReq struct {
	Id              int64  `json:"id"`
	StrategyKeyWord string `json:"strategy_key_word"`
	StrategyType    int64  `json:"strategy_type"`
	StrategyStatus  int64  `json:"strategy_status"`
	Ip              string `json:"ip"`
	request.PageInfo
}

type WarnStrategiesReq struct {
	Ws system.SysWarnStrategies `json:"warn_strategies"`
	Sr []system.SysOpRelations  `json:"op_relations"`
}
