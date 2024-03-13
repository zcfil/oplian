package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system/request"
)

type DBApi struct{}

// InitDB
// @Tags     InitDB
// @Summary  Initialize user database
// @Produce  application/json
// @Param    data  body
// @Success  200   {object}  response.Response{data=string}
// @Router   /init/initdb [post]
func (i *DBApi) InitDB(c *gin.Context) {
	if global.ZC_DB != nil {
		global.ZC_LOG.Error("Database configuration already exists!")
		response.FailWithMessage("Database configuration already exists", c)
		return
	}
	var dbInfo request.InitDB
	if err := c.ShouldBindJSON(&dbInfo); err != nil {
		global.ZC_LOG.Error("Parameter verification failed!", zap.Error(err))
		response.FailWithMessage("Parameter verification failed", c)
		return
	}
	if err := initDBService.InitDB(dbInfo); err != nil {
		global.ZC_LOG.Error("Automatic database creation failed!", zap.Error(err))
		response.FailWithMessage("Automatic database creation failed. Please check the backend logs and initialize after checking", c)
		return
	}
	response.OkWithMessage("Successfully created database automatically", c)
}

// CheckDB
// @Tags     CheckDB
// @Summary  Initialize user database
// @Produce  application/json
// @Success  200  {object}  response.Response{data=map[string]interface{},msg=string}
// @Router   /init/checkdb [post]
func (i *DBApi) CheckDB(c *gin.Context) {
	var (
		message  = "Go to initialize the database"
		needInit = true
	)

	if global.ZC_DB != nil {
		message = "Database does not require initialization"
		needInit = false
	}
	log.Println(message)
	response.OkWithDetailed(gin.H{"needInit": needInit}, message, c)
}
