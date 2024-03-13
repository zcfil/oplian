package system

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"oplian/define"
	requestSys "oplian/model/common/request"
	"oplian/model/common/response"
	"oplian/model/system/request"
	"oplian/utils"
	"time"
)

type AppStoreApi struct{}

// GetSlotList
// @Tags      GetSlotList
// @Summary   Get plugin list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /slot/getSlotList [get]
func (app *AppStoreApi) GetSlotList(c *gin.Context) {
	var param requestSys.PageInfo
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(param, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	type productListReq struct {
		IsSelf   bool   `json:"isSelf" form:"isSelf"`
		ChainUrl string `json:"chainUrl" form:"chainUrl"`
		requestSys.PageInfo
	}

	data, err := utils.RequestDo(define.ChainsysIp+":"+define.ChainsysPort, define.SlotListRouter, "",
		productListReq{IsSelf: true, ChainUrl: define.ChainsysIp + ":" + define.ChainsysPort, PageInfo: param}, time.Second*15)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	type Response struct {
		Code int                 `json:"code"`
		Data response.PageResult `json:"data"`
		Msg  string              `json:"msg"`
	}

	var dataRes Response
	if err = json.Unmarshal(data, &dataRes); err != nil {
		log.Println(err.Error())
		return
	}

	response.OkWithDetailed(dataRes.Data, "Successfully obtained", c)

}

// GetSlotFileList
// @Tags      GetSlotFileList
// @Summary   Get the list of files corresponding to the plugin
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /slot/getSlotFileList [get]
func (app *AppStoreApi) GetSlotFileList(c *gin.Context) {
	productId := c.Query("productId")

	data, err := utils.RequestDo(define.ChainsysIp+":"+define.ChainsysPort, define.SlotFileListRouter+"?productId="+productId,
		"", nil, time.Second*15)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	type Response struct {
		Code int                 `json:"code"`
		Data response.PageResult `json:"data"`
		Msg  string              `json:"msg"`
	}

	var dataRes Response
	if err = json.Unmarshal(data, &dataRes); err != nil {
		log.Println(err.Error())
		return
	}
	response.OkWithDetailed(dataRes.Data, "Successfully obtained", c)
}

// ReplaceSlotFile
// @Tags      ReplaceSlotFile
// @Summary   replace file
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /slot/replaceSlotFile [get]
func (app *AppStoreApi) ReplaceSlotFile(c *gin.Context) {
	var param request.ProductFileReq
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//var errMsg string
	// Loop through the corresponding gateway and copy the corresponding file into the gateway
	//for _, gclient := range global.GateWayClinets.Gets() {

	//_, err = gclient.ReplacePlugFile(context.Background(), &pb.ReplaceFileInfo{
	//	FileName:    param.Name,
	//	ProductId:   uint64(param.ProductId),
	//	DownloadUrl: define.ChainsysIp + ":" + define.ChainsysPort,
	//})
	//if err != nil {
	//	errMsg += fmt.Sprintf("%s\n", err.Error())
	//	continue
	//}
	//}

	response.OkWithMessage("Modified successfully", c)

}
