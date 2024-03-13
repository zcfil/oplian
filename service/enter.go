package service

import (
	"oplian/service/example"
	"oplian/service/gateway"
	"oplian/service/lotus"
	"oplian/service/op"
	"oplian/service/slot"
	"oplian/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	LotusServiceGroup   lotus.ServiceGroup
	GatewayServiceGroup gateway.ServiceGroup
	OpServiceGroup      op.ServiceGroup
	SlotServiceGroup    slot.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
