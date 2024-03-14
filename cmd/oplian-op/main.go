package main

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download
import (
	"github.com/urfave/cli/v2"
	"log"
	"oplian/build"
	"oplian/cmd/oplian-op/cmd"
	"os"
)

// @title                       Swagger Example API
// @version                     0.0.1
// @description                 This is a sample Server pets
// @securityDefinitions.apikey  ApiKeyAuthw
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {

	local := []*cli.Command{
		cmd.Init,
		cmd.Run,
	}
	app := &cli.App{
		Name:                 "oplian-op",
		Usage:                "oplian-op",
		Version:              build.UserVersion(),
		EnableBashCompletion: true,
		Commands:             local,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Println("FAIL TO START:", err)
	}
}
