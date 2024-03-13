package utils

import (
	"github.com/yinheli/qqwry"
	"log"
	"net"
	"os"
)

func GetIpAddress(ip string) *qqwry.QQwry {
	info := &qqwry.QQwry{}
	address := net.ParseIP(ip)
	if ip == "" || address == nil {
		log.Println("get ip os ip is empty")
	} else {
		dir, err := os.Getwd()
		if err != nil {
			log.Println("get ip os dir err", err.Error())
		}
		info = qqwry.NewQQwry(dir + "/config/qqwry.dat")
		info.Find(ip)
	}
	return info
}
