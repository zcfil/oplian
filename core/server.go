package core

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"oplian/global"
	"oplian/initialize"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunWebServer() {

	Router := initialize.WebRouters()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.ZC_CONFIG.System.Addr)
	s := initServer(address, Router)
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	log.Println("server run success on ", zap.String("address", address))

	global.ZC_LOG.Error(s.ListenAndServe().Error())
}
