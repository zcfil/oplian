package main

import (
	"flag"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	c2cmd "oplian/cmd/oplian-op-c2/cmd"
	"os"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

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
			c2cmd.RunOpC2,
		}
		app := &cli.App{
			Name:                 "oplian-op-c2",
			Usage:                "oplian-op-c2",
			EnableBashCompletion: true,
			Commands:             local,
		}
		err := app.Run(os.Args)
		if err != nil {
			log.Println("fail to start:", err)
		}
	}
}
