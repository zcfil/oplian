package main

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download
import (
	"flag"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"oplian/cmd/oplian-op/cmd"
	"os"
)

var (
	GitTag    = "2000.01.01.release"
	BuildTime = "2000-01-01T00:00:00+0800"
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
		fmt.Println("Git Tag: " + GitTag)
		fmt.Println("Build Time: " + BuildTime)
	} else {
		local := []*cli.Command{
			cmd.Init,
			cmd.Run,
		}
		app := &cli.App{
			Name:                 "oplian-op",
			Usage:                "oplian-op",
			EnableBashCompletion: true,
			Commands:             local,
		}
		err := app.Run(os.Args)
		if err != nil {
			log.Println("启动失败:", err)
		}
	}

}
