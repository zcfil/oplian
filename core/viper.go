package core

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"oplian/config"
	"oplian/core/internal"
	"oplian/define"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	uuid "github.com/satori/go.uuid"
	"oplian/global"
	_ "oplian/packfile"
)

func Viper(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config == "" { // 判断命令行参数是否为空
			if configEnv := os.Getenv(internal.ConfigEnv); configEnv == "" { // 判断 internal.ConfigEnv 常量存储的环境变量是否为空
				switch gin.Mode() {
				case gin.DebugMode:
					config = internal.ConfigDefaultFile
					fmt.Printf("You are using the %s environment name in gin mode, and the path to config is %s\n", gin.EnvGinMode, internal.ConfigDefaultFile)
				case gin.ReleaseMode:
					config = internal.ConfigReleaseFile
					fmt.Printf("You are using the %s environment name in gin mode, and the path to config is %s\n", gin.EnvGinMode, internal.ConfigReleaseFile)
				case gin.TestMode:
					config = internal.ConfigTestFile
					fmt.Printf("You are using the %s environment name in gin mode, and the path to config is %s\n", gin.EnvGinMode, internal.ConfigTestFile)
				}
			} else {
				config = configEnv
				fmt.Printf("You are using the %s environment variable, and the path to config is %s\n", internal.ConfigEnv, config)
			}
		} else {
			fmt.Printf("You are using the value passed by the - c parameter on the command line, and the path to config is %s\n", config)
		}
	} else {
		config = path[0]
		fmt.Printf("You are using the value passed by Func Viper(), and the path to config is %s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.ZC_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.ZC_CONFIG); err != nil {
		fmt.Println(err)
	}

	global.ZC_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	return v
}

func ViperRoom(config string) {
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.ROOM_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.ROOM_CONFIG); err != nil {
		fmt.Println(err)
	}
	// jwt key
	global.ROOM_CONFIG.JWT.SigningKey = uuid.NewV4().String()
}

// SetConfigRoom Assign config with the help of environment variables_ Room information
func SetConfigRoom() {
	if os.Getenv(internal.ConfigGatewayIP) != "" {
		global.ROOM_CONFIG.Gateway.IP = os.Getenv(internal.ConfigGatewayIP)
	} else {
		global.ROOM_CONFIG.Gateway.IP = internal.ConfigGatewayIPInfo
	}
	if os.Getenv(internal.ConfigGatewayPort) != "" {
		global.ROOM_CONFIG.Gateway.Port = os.Getenv(internal.ConfigGatewayPort)
	} else {
		global.ROOM_CONFIG.Gateway.Port = internal.ConfigGatewayPortInfo
	}
	if os.Getenv(internal.ConfigGatewayToken) != "" {
		global.ROOM_CONFIG.Gateway.Token = os.Getenv(internal.ConfigGatewayToken)
	} else {
		global.ROOM_CONFIG.Gateway.Token = internal.ConfigGatewayTokenInfo
	}

	if os.Getenv(internal.ConfigOpIP) != "" {
		global.ROOM_CONFIG.Op.IP = os.Getenv(internal.ConfigOpIP)
	} else {
		intranetIP := global.LocalIP
		if len(intranetIP) != 0 {
			global.ROOM_CONFIG.Op.IP = intranetIP
		} else {
			global.ROOM_CONFIG.Op.IP = internal.ConfigOpIPInfo
		}
	}
	if os.Getenv(internal.ConfigOpPort) != "" {
		global.ROOM_CONFIG.Op.Port = os.Getenv(internal.ConfigOpPort)
	} else {
		global.ROOM_CONFIG.Op.Port = define.OpPort
	}
	if os.Getenv(internal.ConfigOpToken) != "" {
		global.ROOM_CONFIG.Op.Token = os.Getenv(internal.ConfigOpToken)
	} else {
		global.ROOM_CONFIG.Op.Token = internal.ConfigOpTokenInfo
	}

	if os.Getenv(internal.ConfigOpC2Port) != "" {
		global.ROOM_CONFIG.OpC2.Port = os.Getenv(internal.ConfigOpC2Port)
	} else {
		global.ROOM_CONFIG.OpC2.Port = internal.ConfigOpC2PortInfo
	}
	if os.Getenv(internal.ConfigOpC2Token) != "" {
		global.ROOM_CONFIG.OpC2.Token = os.Getenv(internal.ConfigOpC2Token)
	} else {
		global.ROOM_CONFIG.OpC2.Token = internal.ConfigOpC2TokenInfo
	}

	global.ROOM_CONFIG.Zap = config.Zap{
		Level:         "info",
		Prefix:        "[oplian]",
		Format:        "console",
		Director:      define.PathOplian + "log",
		EncodeLevel:   "LowercaseColorLevelEncoder",
		StacktraceKey: "stacktrace",
		MaxAge:        0,
		ShowLine:      true,
		LogInConsole:  true,
	}
}
