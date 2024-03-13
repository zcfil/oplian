package system

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system"
	"oplian/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AutoCodeApi struct{}

var caser = cases.Title(language.English)

// PreviewTemp
// @Tags      AutoCode
// @Summary   Preview the created code
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /autoCode/preview [post]
func (autoApi *AutoCodeApi) PreviewTemp(c *gin.Context) {
	var a system.AutoCodeStruct
	_ = c.ShouldBindJSON(&a)
	if err := utils.Verify(a, utils.AutoCodeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	a.Pretreatment() // Process go keywords
	a.PackageT = caser.String(a.Package)
	autoCode, err := autoCodeService.PreviewTemp(a)
	if err != nil {
		global.ZC_LOG.Error("Preview failed!", zap.Error(err))
		response.FailWithMessage("Preview failed", c)
	} else {
		response.OkWithDetailed(gin.H{"autoCode": autoCode}, "Preview successful", c)
	}
}

// CreateTemp
// @Tags      AutoCode
// @Summary   Automatic Code Template
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {string}  string          "{"success":true,"data":{},"msg":"Created successfully"}"
// @Router    /autoCode/createTemp [post]
func (autoApi *AutoCodeApi) CreateTemp(c *gin.Context) {
	var a system.AutoCodeStruct
	_ = c.ShouldBindJSON(&a)
	if err := utils.Verify(a, utils.AutoCodeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	a.Pretreatment()
	var apiIds []uint
	if a.AutoCreateApiToSql {
		if ids, err := autoCodeService.AutoCreateApi(&a); err != nil {
			global.ZC_LOG.Error("Automation creation failed! Please clear junk data by yourself!", zap.Error(err))
			c.Writer.Header().Add("success", "false")
			c.Writer.Header().Add("msg", url.QueryEscape("Automation creation failed! Please clear junk data by yourself!"))
			return
		} else {
			apiIds = ids
		}
	}
	a.PackageT = caser.String(a.Package)
	err := autoCodeService.CreateTemp(a, apiIds...)
	if err != nil {
		if errors.Is(err, system.ErrAutoMove) {
			c.Writer.Header().Add("success", "true")
			c.Writer.Header().Add("msg", url.QueryEscape(err.Error()))
		} else {
			c.Writer.Header().Add("success", "false")
			c.Writer.Header().Add("msg", url.QueryEscape(err.Error()))
			_ = os.Remove("./ginvueadmin.zip")
		}
	} else {
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "ginvueadmin.zip")) // fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
		c.Writer.Header().Add("Content-Type", "application/json")
		c.Writer.Header().Add("success", "true")
		c.File("./ginvueadmin.zip")
		_ = os.Remove("./ginvueadmin.zip")
	}
}

// GetDB
// @Tags      AutoCode
// @Summary   Get all current databases
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /autoCode/getDatabase [get]
func (autoApi *AutoCodeApi) GetDB(c *gin.Context) {
	businessDB := c.Query("businessDB")
	dbs, err := autoCodeService.Database(businessDB).GetDB(businessDB)
	var dbList []map[string]interface{}
	for _, db := range global.ZC_CONFIG.DBList {
		var item = make(map[string]interface{})
		item["aliasName"] = db.AliasName
		item["dbName"] = db.Dbname
		item["disable"] = db.Disable
		item["dbtype"] = db.Type
		dbList = append(dbList, item)
	}
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
	} else {
		response.OkWithDetailed(gin.H{"dbs": dbs, "dbList": dbList}, "Successfully obtained", c)
	}
}

// GetTables
// @Tags      AutoCode
// @Summary   Get all tables in the current database
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /autoCode/getTables [get]
func (autoApi *AutoCodeApi) GetTables(c *gin.Context) {
	dbName := c.DefaultQuery("dbName", global.ZC_CONFIG.Mysql.Dbname)
	businessDB := c.Query("businessDB")
	tables, err := autoCodeService.Database(businessDB).GetTables(businessDB, dbName)
	if err != nil {
		global.ZC_LOG.Error("Query table failed!", zap.Error(err))
		response.FailWithMessage("Query table failed", c)
	} else {
		response.OkWithDetailed(gin.H{"tables": tables}, "Successfully obtained", c)
	}
}

// GetColumn
// @Tags      AutoCode
// @Summary   Get all fields in the current table
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /autoCode/getColumn [get]
func (autoApi *AutoCodeApi) GetColumn(c *gin.Context) {
	businessDB := c.Query("businessDB")
	dbName := c.DefaultQuery("dbName", global.ZC_CONFIG.Mysql.Dbname)
	tableName := c.Query("tableName")
	columns, err := autoCodeService.Database(businessDB).GetColumn(businessDB, tableName, dbName)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
	} else {
		response.OkWithDetailed(gin.H{"columns": columns}, "Successfully obtained", c)
	}
}

// CreatePackage
// @Tags      AutoCode
// @Summary   Create package
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /autoCode/createPackage [post]
func (autoApi *AutoCodeApi) CreatePackage(c *gin.Context) {
	var a system.SysAutoCode
	_ = c.ShouldBindJSON(&a)
	if err := utils.Verify(a, utils.AutoPackageVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := autoCodeService.CreateAutoCode(&a)
	if err != nil {

		global.ZC_LOG.Error("Created successfully!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
	} else {
		response.OkWithMessage("Created successfully", c)
	}
}

// GetPackage
// @Tags      AutoCode
// @Summary   Get package
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /autoCode/getPackage [post]
func (autoApi *AutoCodeApi) GetPackage(c *gin.Context) {
	pkgs, err := autoCodeService.GetPackage()
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
	} else {
		response.OkWithDetailed(gin.H{"pkgs": pkgs}, "Successfully obtained", c)
	}
}

// DelPackage
// @Tags      AutoCode
// @Summary   Delete package
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /autoCode/delPackage [post]
func (autoApi *AutoCodeApi) DelPackage(c *gin.Context) {
	var a system.SysAutoCode
	_ = c.ShouldBindJSON(&a)
	err := autoCodeService.DelPackage(a)
	if err != nil {
		global.ZC_LOG.Error("Delete failed!", zap.Error(err))
		response.FailWithMessage("Delete failed", c)
	} else {
		response.OkWithMessage("Delete successful", c)
	}
}

// AutoPlug
// @Tags      AutoCode
// @Summary   Create plugin template
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}
// @Router    /autoCode/createPlug [post]
func (autoApi *AutoCodeApi) AutoPlug(c *gin.Context) {
	var a system.AutoPlugReq
	err := c.ShouldBindJSON(&a)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	a.Snake = strings.ToLower(a.PlugName)
	a.NeedModel = a.HasRequest || a.HasResponse
	err = autoCodeService.CreatePlug(a)
	if err != nil {
		global.ZC_LOG.Error("Preview failed!", zap.Error(err))
		response.FailWithMessage("Preview failed", c)
		return
	}
	response.Ok(c)
}

// InstallPlugin
// @Tags      AutoCode
// @Summary   Install plug-in
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     plug  formData  file
// @Success   200   {object}  response.Response{data=[]interface{},msg=string}
// @Router    /autoCode/createPlug [post]
func (autoApi *AutoCodeApi) InstallPlugin(c *gin.Context) {
	header, err := c.FormFile("plug")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	web, server, err := autoCodeService.InstallPlugin(header)
	webStr := "Web plugin installation successful"
	serverStr := "Server plugin installation successful"
	if web == -1 {
		webStr = "The web plugin was not successfully installed. Please decompress and install it according to the documentation. If it is a pure backend plugin, please ignore this prompt"
	}
	if server == -1 {
		serverStr = "The server-side plugin was not successfully installed. Please decompress and install it according to the documentation. If it is a pure front-end plugin, please ignore this prompt"
	}
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData([]interface{}{
		gin.H{
			"code": web,
			"msg":  webStr,
		},
		gin.H{
			"code": server,
			"msg":  serverStr,
		}}, c)
}
