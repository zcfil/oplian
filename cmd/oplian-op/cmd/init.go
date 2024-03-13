package cmd

import (
	"github.com/urfave/cli/v2"
	"oplian/initialize"
	"oplian/utils"
)

var Init = &cli.Command{
	Name:  "init",
	Usage: "initialization",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		initialize.HostSystemInitialization()
		utils.InstallNumactl()
		return nil
	},
}
