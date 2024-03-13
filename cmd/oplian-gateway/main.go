package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"oplian/build"
	"oplian/cmd/oplian-gateway/cmd"
	_ "oplian/source/system"
	"os"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title                       Swagger Example API
// @version                     0.0.1
// @description                 This is a sample Server pets
// @securityDefinitions.apikey  ApiKeyAuthw
// @in                          headere
// @name                        x-token
// @BasePath                    /da
func main() {

	local := []*cli.Command{
		cmd.Run,
	}
	app := &cli.App{
		Name:                 "oplian-gateway",
		Usage:                "oplian-gateway",
		Version:              build.UserVersion(),
		EnableBashCompletion: true,
		Commands:             local,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Println("fail to start:", err)
	}
}
