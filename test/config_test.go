package test

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"oplian/config"
	"oplian/define"
	"oplian/global"
	"os"
	"testing"
)

func TestBoostConfig(t *testing.T) {
	listenIP := "10.0.1.196"
	listenPort := "55555"
	var boostConfig config.Boost
	if _, err := toml.DecodeFile(define.ConfigName, &boostConfig); err != nil {
		global.ZC_LOG.Error(fmt.Sprintf("读取配置文件失败：%s", err.Error()))
		return
	}
	boostConfig.Libp2p.ListenAddresses = []string{
		fmt.Sprintf("/ip4/%s/tcp/%s", listenIP, listenPort),
		fmt.Sprintf("/ip6/::/tcp/%s", listenPort),
	}
	boostConfig.Libp2p.AnnounceAddresses = []string{fmt.Sprintf("/ip4/%s/tcp/%s", listenIP, listenPort)}
	buf := new(bytes.Buffer)
	e := toml.NewEncoder(buf)
	if err := e.Encode(boostConfig); err != nil {
		global.ZC_LOG.Error(fmt.Sprintf("修改配置文件失败：%s", err.Error()))
		return
	}
	os.Remove(define.ConfigName)
	err := os.WriteFile(define.ConfigName, []byte(buf.String()), 0644)
	if err != nil {
		fmt.Printf("writing config file: %v", err)
	}
	fmt.Println()
}
