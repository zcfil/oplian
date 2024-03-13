package response

import "oplian/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
