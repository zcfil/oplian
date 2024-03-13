package system

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"log"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/system"
	"oplian/utils"
)

type UserService struct{}

func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.ZC_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) {
		return userInter, errors.New("user name registered")
	}

	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.NewV4()
	err = global.ZC_DB.Create(&u).Error
	return u, err
}

func (userService *UserService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {

	log.Println("Login:", u)
	log.Println("global.ZC_DB:", global.ZC_DB)
	if nil == global.ZC_DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.SysUser
	err = global.ZC_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("password error")
		}
		MenuServiceApp.UserAuthorityDefaultRouter(&user)
	}
	return &user, err
}

func (userService *UserService) ChangePassword(u *system.SysUser, newPassword string) (userInter *system.SysUser, err error) {
	var user system.SysUser
	if err = global.ZC_DB.Where("id = ?", u.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.New("the original password is incorrect")
	}
	user.Password = utils.BcryptHash(newPassword)
	err = global.ZC_DB.Save(&user).Error
	return &user, err

}

func (userService *UserService) GetUserInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&system.SysUser{})
	var userList []system.SysUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").Find(&userList).Error
	return userList, total, err
}

func (userService *UserService) GetUserInfoPullList() (list []system.SysUser, err error) {
	db := global.ZC_DB.Model(&system.SysUser{})
	err = db.Find(&list).Error
	return list, err
}

func (userService *UserService) SetUserAuthority(id uint, authorityId uint) (err error) {
	assignErr := global.ZC_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&system.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("this user does not have this role")
	}
	err = global.ZC_DB.Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityId).Error
	return err
}

func (userService *UserService) SetUserAuthorities(id uint, authorityIds []uint) (err error) {
	return global.ZC_DB.Transaction(func(tx *gorm.DB) error {
		TxErr := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		var useAuthority []system.SysUserAuthority
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, system.SysUserAuthority{
				SysUserId: id, SysAuthorityAuthorityId: v,
			})
		}
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}
		TxErr = tx.Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityIds[0]).Error
		if TxErr != nil {
			return TxErr
		}
		return nil
	})
}

func (userService *UserService) DeleteUser(id int) (err error) {
	var user system.SysUser
	err = global.ZC_DB.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	err = global.ZC_DB.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error
	return err
}

func (userService *UserService) SetUserInfo(req system.SysUser) error {
	return global.ZC_DB.Updates(&req).Error
}

func (userService *UserService) GetUserInfo(uuid uuid.UUID) (user system.SysUser, err error) {
	var reqUser system.SysUser
	err = global.ZC_DB.Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return reqUser, err
	}
	MenuServiceApp.UserAuthorityDefaultRouter(&reqUser)
	return reqUser, err
}

func (userService *UserService) FindUserById(id int) (user *system.SysUser, err error) {
	var u system.SysUser
	err = global.ZC_DB.Where("`id` = ?", id).First(&u).Error
	return &u, err
}

func (userService *UserService) FindUserByUuid(uuid string) (user *system.SysUser, err error) {
	var u system.SysUser
	if err = global.ZC_DB.Where("`uuid` = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("user does not exist")
	}
	return &u, nil
}

func (userService *UserService) ResetPassword(ID uint) (err error) {
	err = global.ZC_DB.Model(&system.SysUser{}).Where("id = ?", ID).Update("password", utils.BcryptHash("123456")).Error
	return err
}

func (userService *UserService) GetUserInfoTidy(uuid string) (user system.SysUser, err error) {
	var reqUser system.SysUser
	err = global.ZC_DB.First(&reqUser, "uuid = ?", uuid).Error
	return reqUser, err
}
