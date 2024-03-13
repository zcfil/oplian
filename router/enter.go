package router

import (
	"oplian/router/example"
	"oplian/router/gateway"
	"oplian/router/lotus"
	"oplian/router/slot"
	"oplian/router/system"
)

type RouterGroup struct {
	System  system.RouterGroup
	Gateway gateway.RouterGroup
	Lotus   lotus.RouterGroup
	Slot    slot.RouterGroup
	Example example.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
