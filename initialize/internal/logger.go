package internal

import (
	"fmt"
	"log"

	"gorm.io/gorm/logger"
	"oplian/global"
)

type writer struct {
	logger.Writer
}

// NewWriter writer

func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf

func (w *writer) Printf(message string, data ...interface{}) {
	var logZap bool
	switch global.ZC_CONFIG.System.DbType {
	case "mysql":
		logZap = global.ZC_CONFIG.Mysql.LogZap
	}
	if logZap {
		log.Println(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}
