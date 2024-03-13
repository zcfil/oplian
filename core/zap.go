package core

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"oplian/core/internal"
	"oplian/global"
	"oplian/utils"
	"os"
)

func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.ROOM_CONFIG.Zap.Director); !ok {
		fmt.Printf("create %v directory\n", global.ROOM_CONFIG.Zap.Director)
		_ = os.Mkdir(global.ROOM_CONFIG.Zap.Director, os.ModePerm)
	}

	cores := internal.Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if global.ROOM_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
