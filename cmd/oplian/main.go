package main

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download
import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"oplian/build"
	"oplian/core"
	"oplian/global"
	"oplian/initialize"
	_ "oplian/source/system"
)

// @title                       Swagger Example API
// @version                     0.0.1
// @description                 This is a sample Server pets
// @securityDefinitions.apikey  ApiKeyAuthw
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {

	version := flag.Bool("v", false, "version")
	flag.Parse()

	if *version {
		fmt.Printf("oplian version %s\n", build.UserVersion())
	} else {
		go initialize.GateWayHeartbeat(context.Background())
		core.ViperRoom("config/config_room.yaml")
		global.ZC_VP = core.Viper("config/config.yaml")
		// initialize.OtherInit()
		global.ZC_LOG = core.Zap()
		zap.ReplaceGlobals(global.ZC_LOG)
		global.ZC_DB = initialize.Gorm()
		//go initialize.PolicyWarn(context.TODO())
		if global.ZC_DB != nil {
			initialize.RegisterTables(global.ZC_DB)
			fmt.Println("---InitMysqlData---")
			log.Println(initialize.InitMysqlData())
			db, _ := global.ZC_DB.DB()
			defer db.Close()
		}

		go initialize.DeleteMonitorInfo()
		core.RunWebServer()
	}

}
