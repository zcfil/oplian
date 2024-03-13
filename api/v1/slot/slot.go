package slot

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/slot/request"
	"oplian/service/pb"
)

type SlotApi struct{}

// InstallSlot
// @Tags      InstallSlot
// @Summary   安装插件
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SlotParam
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /slot/installSlot [post]
func (slot *SlotApi) InstallSlot(c *gin.Context) {
	var param request.SlotParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var errMsg string
	switch param.Name {
	case "unsealed":
		for _, gclient := range global.GateWayClinets.Gets() {
			//新增存储
			_, err = gclient.InstallUnsealed(context.Background(), &pb.String{})
			if err != nil {
				errMsg += fmt.Sprintf("%s\n", err.Error())
				continue
			}
		}
	}
	response.OkWithMessage(errMsg, c)
}
