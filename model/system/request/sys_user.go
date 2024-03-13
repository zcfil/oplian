package request

import (
	"oplian/model/system"
)

// Register User register structure
type Register struct {
	Username     string `json:"userName" example:"用户名"`
	Password     string `json:"passWord" example:"密码"`
	NickName     string `json:"nickName" example:"昵称"`
	HeaderImg    string `json:"headerImg" example:"头像链接"`
	AuthorityId  uint   `json:"authorityId" swaggertype:"string" example:"int 角色id"`
	Enable       int    `json:"enable" swaggertype:"string" example:"int 是否启用"`
	AuthorityIds []uint `json:"authorityIds" swaggertype:"string" example:"[]uint 角色id"`
	Phone        string `json:"phone" example:"电话号码"`
	Email        string `json:"email" example:"电子邮箱"`
}

// User login structure
type Login struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

// Modify password structure
type ChangePasswordReq struct {
	ID          uint   `json:"-"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

// Modify  user's auth structure
type SetUserAuth struct {
	AuthorityId uint `json:"authorityId"`
}

// Modify  user's auth structure
type SetUserAuthorities struct {
	ID           uint
	AuthorityIds []uint `json:"authorityIds"`
}

type ChangeUserInfo struct {
	ID           uint                  `gorm:"primarykey"`
	NickName     string                `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`
	Phone        string                `json:"phone"  gorm:"comment:用户手机号"`
	AuthorityIds []uint                `json:"authorityIds" gorm:"-"`
	Email        string                `json:"email"  gorm:"comment:用户邮箱"`
	HeaderImg    string                `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"`
	SideMode     string                `json:"sideMode"  gorm:"comment:用户侧边主题"`
	Enable       int                   `json:"enable" gorm:"comment:冻结用户"`
	Authorities  []system.SysAuthority `json:"-" gorm:"many2many:sys_user_authority;"`
}
