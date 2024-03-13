package system

import (
	"oplian/global"
)

type JwtBlacklist struct {
	global.ZC_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
