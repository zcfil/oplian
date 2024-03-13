package system

import (
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

type AuthorityMenuApi struct{}

// GetMenu
// @Tags      AuthorityMenu
// @Summary   Obtain user dynamic routing
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=systemRes.SysMenusResponse,msg=string}
// @Router    /menu/getMenu [post]
func (a *AuthorityMenuApi) GetMenu(c *gin.Context) {

	menus, err := menuService.GetMenuTree(utils.GetUserAuthorityId(c))
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
	}
	if menus == nil {
		menus = []system.SysMenu{}
	}
	response.OkWithDetailed(systemRes.SysMenusResponse{Menus: menus}, "Successfully obtained", c)
}

// GetBaseMenuTree
// @Tags      AuthorityMenu
// @Summary   Obtain user dynamic routing
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=systemRes.SysBaseMenusResponse,msg=string}
// @Router    /menu/getBaseMenuTree [post]
func (a *AuthorityMenuApi) GetBaseMenuTree(c *gin.Context) {
	menus, err := menuService.GetBaseMenuTree()
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(systemRes.SysBaseMenusResponse{Menus: menus}, "Successfully obtained", c)
}

// AddMenuAuthority
// @Tags      AuthorityMenu
// @Summary   Add menu and role association relationships
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /menu/addMenuAuthority [post]
func (a *AuthorityMenuApi) AddMenuAuthority(c *gin.Context) {
	var authorityMenu systemReq.AddMenuAuthorityInfo
	err := c.ShouldBindJSON(&authorityMenu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(authorityMenu, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := menuService.AddMenuAuthority(authorityMenu.Menus, authorityMenu.AuthorityId); err != nil {
		global.ZC_LOG.Error("Add failed!", zap.Error(err))
		response.FailWithMessage("Add failed", c)
	} else {
		response.OkWithMessage("Added successfully", c)
	}
}

// GetMenuAuthority
// @Tags      AuthorityMenu
// @Summary   Get the specified role menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /menu/getMenuAuthority [post]
func (a *AuthorityMenuApi) GetMenuAuthority(c *gin.Context) {
	var param request.GetAuthorityId
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.AuthorityIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	menus, err := menuService.GetMenuAuthority(&param)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithDetailed(systemRes.SysMenusResponse{Menus: menus}, "Acquisition failed", c)
		return
	}
	response.OkWithDetailed(gin.H{"menus": menus}, "Successfully obtained", c)
}

// AddBaseMenu
// @Tags      Menu
// @Summary   Add menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /menu/addBaseMenu [post]
func (a *AuthorityMenuApi) AddBaseMenu(c *gin.Context) {
	var menu system.SysBaseMenu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu, utils.MenuVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu.Meta, utils.MenuMetaVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = menuService.AddBaseMenu(menu)
	if err != nil {
		global.ZC_LOG.Error("Add failed!", zap.Error(err))
		response.FailWithMessage("Add failed", c)
		return
	}
	response.OkWithMessage("Added successfully", c)
}

// DeleteBaseMenu
// @Tags      Menu
// @Summary   Delete Menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /menu/deleteBaseMenu [post]
func (a *AuthorityMenuApi) DeleteBaseMenu(c *gin.Context) {
	var menu request.GetById
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = baseMenuService.DeleteBaseMenu(menu.ID)
	if err != nil {
		global.ZC_LOG.Error("Delete failed!", zap.Error(err))
		response.FailWithMessage("Delete failed", c)
		return
	}
	response.OkWithMessage("Delete successful", c)
}

// UpdateBaseMenu
// @Tags      Menu
// @Summary   update menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /menu/updateBaseMenu [post]
func (a *AuthorityMenuApi) UpdateBaseMenu(c *gin.Context) {
	var menu system.SysBaseMenu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu, utils.MenuVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu.Meta, utils.MenuMetaVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = baseMenuService.UpdateBaseMenu(menu)
	if err != nil {
		global.ZC_LOG.Error("Update failed!", zap.Error(err))
		response.FailWithMessage("Update failed", c)
		return
	}
	response.OkWithMessage("Update success", c)
}

// GetBaseMenuById
// @Tags      Menu
// @Summary   Retrieve menu based on ID
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=systemRes.SysBaseMenuResponse,msg=string}
// @Router    /menu/getBaseMenuById [post]
func (a *AuthorityMenuApi) GetBaseMenuById(c *gin.Context) {
	var idInfo request.GetById
	err := c.ShouldBindJSON(&idInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	menu, err := baseMenuService.GetBaseMenuById(idInfo.ID)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(systemRes.SysBaseMenuResponse{Menu: menu}, "Successfully obtained", c)
}

// GetMenuList
// @Tags      Menu
// @Summary   Paging to obtain the basic menu list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /menu/getMenuList [post]
func (a *AuthorityMenuApi) GetMenuList(c *gin.Context) {
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
	menuList, total, err := menuService.GetInfoList()
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     menuList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}
