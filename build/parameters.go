package build

import (
	_ "embed"
	"encoding/json"
)

//go:embed proof-params/parameters.json
var params []byte

//go:embed proof-params/srs-inner-product.json
var srs []byte

func ParametersMap() map[string]struct{} {
	paramsMap := make(map[string]struct{})
	json.Unmarshal(params, &paramsMap)
	return paramsMap
}

func SrsJSON() []byte {
	return srs
}
