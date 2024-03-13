package system

import (
	"strconv"

	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/common/response"
	"oplian/model/system"
	systemReq "oplian/model/system/request"
	systemRes "oplian/model/system/response"
	"oplian/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Login
// @Tags     Base
// @Summary  log in
// @Produce   application/json
// @Param    data  body
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}
// @Router   /base/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var l systemReq.Login
	err := c.ShouldBindJSON(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//err = utils.Verify(l, utils.LoginVerify)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	if store.Verify(l.CaptchaId, l.Captcha, true) {
		u := &system.SysUser{Username: l.Username, Password: l.Password}
		user, err := userService.Login(u)
		if err != nil {
			global.ZC_LOG.Error("Login failed! The username does not exist or the password is incorrect!", zap.Error(err))
			response.FailWithMessage("The username does not exist or the password is incorrect", c)
			return
		}
		if user.Enable != 1 {
			global.ZC_LOG.Error("Login failed! User is prohibited from logging in!")
			response.FailWithMessage("User is prohibited from logging in", c)
			return
		}
		b.TokenNext(c, *user)
		return
	}
	response.FailWithMessage("Verification code error", c)
}

// TokenNext Sign jwt after logging in
func (b *BaseApi) TokenNext(c *gin.Context, user system.SysUser) {
	j := &utils.JWT{SigningKey: []byte(global.ZC_CONFIG.JWT.SigningKey)} // Unique signature
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain token!", zap.Error(err))
		response.FailWithMessage("Failed to obtain token", c)
		return
	}
	if !global.ZC_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "Login succeeded", c)
		return
	}
}

// Register
// @Tags     SysUser
// @Summary  User registration account
// @Produce   application/json
// @Param    data  body
// @Success  200   {object}  response.Response{data=systemRes.SysUserResponse,msg=string}
// @Router   /user/admin_register [post]
func (b *BaseApi) Register(c *gin.Context) {
	var r systemReq.Register
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var authorities []system.SysAuthority
	for _, v := range r.AuthorityIds {
		authorities = append(authorities, system.SysAuthority{
			AuthorityId: v,
		})
	}
	user := &system.SysUser{Username: r.Username, NickName: r.NickName, Password: r.Password, HeaderImg: r.HeaderImg, AuthorityId: r.AuthorityId, Authorities: authorities, Enable: r.Enable, Phone: r.Phone, Email: r.Email}
	userReturn, err := userService.Register(*user)
	if err != nil {
		global.ZC_LOG.Error("Register has failed!", zap.Error(err))
		response.FailWithDetailed(systemRes.SysUserResponse{User: userReturn}, "Register has failed", c)
		return
	}
	response.OkWithDetailed(systemRes.SysUserResponse{User: userReturn}, "Register was successful", c)
}

// ChangePassword
// @Tags      SysUser
// @Summary   User changes password
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /user/changePassword [post]
func (b *BaseApi) ChangePassword(c *gin.Context) {
	var req systemReq.ChangePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.ChangePasswordVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	u := &system.SysUser{ZC_MODEL: global.ZC_MODEL{ID: uid}, Password: req.Password}
	_, err = userService.ChangePassword(u, req.NewPassword)
	if err != nil {
		global.ZC_LOG.Error("Modification failed!", zap.Error(err))
		response.FailWithMessage("Modification failed, original password does not match current account", c)
		return
	}
	response.OkWithMessage("Modified successfully", c)
}

// GetUserList
// @Tags      SysUser
// @Summary   Paging to obtain user list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /user/getUserList [post]
func (b *BaseApi) GetUserList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := userService.GetUserInfoList(pageInfo)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// SetUserAuthority
// @Tags      SysUser
// @Summary   Change user permissions
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /user/setUserAuthority [post]
func (b *BaseApi) SetUserAuthority(c *gin.Context) {
	var sua systemReq.SetUserAuth
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if UserVerifyErr := utils.Verify(sua, utils.SetUserAuthorityVerify); UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	userID := utils.GetUserID(c)
	err = userService.SetUserAuthority(userID, sua.AuthorityId)
	if err != nil {
		global.ZC_LOG.Error("Modification failed!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims := utils.GetUserInfo(c)
	j := &utils.JWT{SigningKey: []byte(global.ZC_CONFIG.JWT.SigningKey)} // 唯一签名
	claims.AuthorityId = sua.AuthorityId
	if token, err := j.CreateToken(*claims); err != nil {
		global.ZC_LOG.Error("Modification failed!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		c.Header("new-token", token)
		c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt, 10))
		response.OkWithMessage("Modified successfully", c)
	}
}

// SetUserAuthorities
// @Tags      SysUser
// @Summary   Set user permissions
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /user/setUserAuthorities [post]
func (b *BaseApi) SetUserAuthorities(c *gin.Context) {
	var sua systemReq.SetUserAuthorities
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.SetUserAuthorities(sua.ID, sua.AuthorityIds)
	if err != nil {
		global.ZC_LOG.Error("Modification failed!", zap.Error(err))
		response.FailWithMessage("Modification failed", c)
		return
	}
	response.OkWithMessage("Modified successfully", c)
}

// DeleteUser
// @Tags      SysUser
// @Summary   delete user
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /user/deleteUser [delete]
func (b *BaseApi) DeleteUser(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	jwtId := utils.GetUserID(c)
	if jwtId == uint(reqId.ID) {
		response.FailWithMessage("Delete failed", c)
		return
	}
	err = userService.DeleteUser(reqId.ID)
	if err != nil {
		global.ZC_LOG.Error("Delete failed!", zap.Error(err))
		response.FailWithMessage("Delete failed", c)
		return
	}
	response.OkWithMessage("Delete successful", c)
}

// SetUserInfo
// @Tags      SysUser
// @Summary   set user information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /user/setUserInfo [put]
func (b *BaseApi) SetUserInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(user, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if len(user.AuthorityIds) != 0 {
		err = userService.SetUserAuthorities(user.ID, user.AuthorityIds)
		if err != nil {
			global.ZC_LOG.Error("Setting failed!", zap.Error(err))
			response.FailWithMessage("Setting failed", c)
			return
		}
	}
	err = userService.SetUserInfo(system.SysUser{
		ZC_MODEL: global.ZC_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		SideMode:  user.SideMode,
		Enable:    user.Enable,
	})
	if err != nil {
		global.ZC_LOG.Error("Setting failed!", zap.Error(err))
		response.FailWithMessage("Setting failed", c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}

// SetSelfInfo
// @Tags      SysUser
// @Summary   set user information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /user/SetSelfInfo [put]
func (b *BaseApi) SetSelfInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user.ID = utils.GetUserID(c)
	err = userService.SetUserInfo(system.SysUser{
		ZC_MODEL: global.ZC_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		SideMode:  user.SideMode,
		Enable:    user.Enable,
	})
	if err != nil {
		global.ZC_LOG.Error("Setting failed!", zap.Error(err))
		response.FailWithMessage("Setting failed", c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}

// GetUserInfo
// @Tags      SysUser
// @Summary   Obtain user information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /user/getUserInfo [get]
func (b *BaseApi) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	ReqUser, err := userService.GetUserInfo(uuid)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "Successfully obtained", c)
}

// ResetPassword
// @Tags      SysUser
// @Summary   Reset user password
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /user/resetPassword [post]
func (b *BaseApi) ResetPassword(c *gin.Context) {
	var user system.SysUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.ResetPassword(user.ID)
	if err != nil {
		global.ZC_LOG.Error("Reset failed!", zap.Error(err))
		response.FailWithMessage("Reset failed"+err.Error(), c)
		return
	}
	response.OkWithMessage("Reset successful", c)
}

// GetUserPullList
// @Tags      SysUser
// @Summary   Get user list dropdown
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /user/getUserList [post]
func (b *BaseApi) GetUserPullList(c *gin.Context) {
	list, err := userService.GetUserInfoPullList()
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	resp := make([]systemRes.UserPullListResponse, len(list))
	for key, val := range list {
		resp[key].UUID = val.UUID
		resp[key].Username = val.Username
		resp[key].NickName = val.NickName
	}

	response.OkWithDetailed(resp, "Successfully obtained", c)
}
