package define

import "errors"

var (
	ErrorNodeNotFound    = errors.New("the node does not exist！")
	ErrorGatewayNotFound = errors.New("the gateway does not exist！")
)
