package response

import (
	"github.com/satori/go.uuid"
	"oplian/model/system"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

type LoginResponse struct {
	User      system.SysUser `json:"user"`
	Token     string         `json:"token"`
	ExpiresAt int64          `json:"expiresAt"`
}

type UserPullListResponse struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"userName"`
	NickName string    `json:"nickName"`
}
