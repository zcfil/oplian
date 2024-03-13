package v1

import (
	"oplian/api/v1/example"
	"oplian/api/v1/gateway"
	"oplian/api/v1/lotus"
	"oplian/api/v1/slot"
	"oplian/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	GateWayApiGroup gateway.ApiGroup
	LotusApiGroup   lotus.ApiGroup
	SlotApiGroup    slot.ApiGroup
	ExampleApiGroup example.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
