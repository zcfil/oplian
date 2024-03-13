package system

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/response"
	systemRes "oplian/model/system/response"
)

var store = base64Captcha.DefaultMemStore

type BaseApi struct{}

// Captcha
// @Tags      Base
// @Summary   Generate verification code
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysCaptchaResponse,msg=string}
// @Router    /base/captcha [post]
func (b *BaseApi) Captcha(c *gin.Context) {

	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	// cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		global.ZC_LOG.Error("Verification code acquisition failed!", zap.Error(err))
		response.FailWithMessage("Verification code acquisition failed", c)
		return
	}
	response.OkWithDetailed(systemRes.SysCaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: 6,
	}, "Verification code obtained successfully", c)
}
