package op

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/service/pb"
	"oplian/utils"
	"strconv"
	"strings"
	"time"
)

var OpWarnServiceApi = new(OpWarnService)

type OpWarnService struct {
}

// CheckOpWarn 检测OP告警信息
func (w *OpWarnService) CheckOpWarn() error {

	warnNameMap := make(map[string]string)
	warnNameMap[define.WindowedPostWallet] = "windowedPost钱包"
	warnNameMap[define.DiskIo] = "磁盘IO"
	warnNameMap[define.Date] = "日期问题"
	warnNameMap[define.SoftwarePackage] = "程序包"
	warnNameMap[define.SoftwarePath] = "程序包目录"
	warnNameMap[define.SpeedUpFile] = "加速文件"
	warnNameMap[define.ProofParam] = "证明参数"
	warnNameMap[define.LotusHeight] = "lotus高度"

	scriptMap := make(map[string]string)
	scriptMap[define.WindowedPostWallet] = define.WindowedPostBalance
	scriptMap[define.DiskIo] = define.DiskIoRes
	scriptMap[define.Date] = define.DateRes
	scriptMap[define.SoftwarePackage] = ""
	scriptMap[define.SoftwarePath] = ""
	scriptMap[define.SpeedUpFile] = ""
	scriptMap[define.ProofParam] = ""
	scriptMap[define.LotusHeight] = define.LotusHeightReq

	ticker := time.NewTicker(time.Hour * time.Duration(utils.THREE))
	for {

		res, err := global.OpToGatewayClient.HostType(context.TODO(), &pb.String{Value: global.OpUUID.String()})
		if err != nil {
			log.Println(fmt.Sprintf("HostType err:%s", err.Error()))
			time.Sleep(time.Second * time.Duration(utils.Five))
			continue
		}

		select {
		case <-ticker.C:

			checkWarn := false
			warnKey := ""
			warnInfo := ""
			for k, v := range scriptMap {

				value := ""
				if !utils.IsNull(v) {
					b, _ := utils.ExecScript(v)
					value = string(b)
				}
				switch k {
				case define.WindowedPostWallet:

					balance, _ := strconv.Atoi(value)
					if balance < 20 && res.Value == define.NodeMachine.String() {
						warnKey = define.WindowedPostWallet
						warnInfo = "windowPost钱包余额不足,请充值!"
						checkWarn = true
					}
					break
				case define.DiskIo:
					if value != "" && res.Value == define.NodeMachine.String() {
						warnKey = define.DiskIo
						warnInfo = "磁盘IO异常:" + value
						checkWarn = true
					}
					break
				case define.Date:
					if time.Now().Sub(utils.StrToTime(value)).Seconds() > utils.ONE {
						warnKey = define.Date
						warnInfo = "时间不同步:" + time.Now().Sub(utils.StrToTime(v)).String()
						checkWarn = true
					}
					break
				case define.SoftwarePackage:

					if res.Value == define.NodeMachine.String() {

						fileInfoList, err := ioutil.ReadDir(define.PathIpfsProgram)
						if err != nil {
							warnKey = define.SoftwarePackage
							warnInfo = fmt.Sprintf("程序包目录不存在,目录: %s", define.PathIpfsProgram)
							checkWarn = true
						} else {

							packageMap := make(map[string]int)
							packageMap[define.ProgramLotus.String()] = utils.ZERO
							packageMap[define.ProgramMiner.String()] = utils.ZERO
							packageMap[define.ProgramWorkerTask.String()] = utils.ZERO
							packageMap[define.ProgramBoost.String()] = utils.ZERO
							for i := range fileInfoList {
								if !fileInfoList[i].IsDir() {
									if _, ok := packageMap[fileInfoList[i].Name()]; ok {
										packageMap[fileInfoList[i].Name()] = utils.ONE
									}
								}
							}

							packageStr := ""
							for k, v := range packageMap {
								if v == utils.ZERO {
									packageStr += k + ","
								}
							}

							if packageStr != "" && !checkWarn {
								packageStr = utils.SubStr(packageStr, utils.ZERO, len(packageStr)-utils.ONE)
								warnKey = define.SoftwarePackage
								warnInfo = "程序包不存在,程序包:" + packageStr
								checkWarn = true
							}
						}
					}

					break
				case define.SoftwarePath:

					if res.Value == define.NodeMachine.String() {

						fileInfoList, err := ioutil.ReadDir(define.PathIpfsData)
						if err != nil {
							warnKey = define.SoftwarePath
							warnInfo = fmt.Sprintf("程序包根目录不存在,目录:%s", define.PathIpfsData)
							checkWarn = true
						} else {

							packageMap := make(map[string]int)
							packageMap[define.PathIpfsLotus] = utils.ZERO
							packageMap[define.PathIpfsMiner] = utils.ZERO
							for i := range fileInfoList {
								if fileInfoList[i].IsDir() {
									if _, ok := packageMap[fileInfoList[i].Name()]; ok {
										packageMap[fileInfoList[i].Name()] = utils.ONE
									}
								}
							}

							packageStr := ""
							for k, v := range packageMap {
								if v == utils.ZERO {
									packageStr += k + ","
								}
							}

							if packageStr != "" {
								packageStr = utils.SubStr(packageStr, utils.ZERO, len(packageStr)-utils.ONE)
								warnKey = define.SoftwarePath
								warnInfo = "程序包目录不存在,目录:" + packageStr
								checkWarn = true
							}
						}
					}
					break
				case define.SpeedUpFile:

					if res.Value == define.NodeMachine.String() {
						fileInfoList, err := ioutil.ReadDir(define.PathSpeedUpFile)
						if err != nil {
							warnKey = define.SpeedUpFile
							warnInfo = "加速文件根目录不存在,目录:" + define.PathSpeedUpFile
							checkWarn = true
						} else {

							fileMap := make(map[string]int)
							fileMap[define.SpeedUpFile32G.String()] = utils.ZERO
							fileMap[define.SpeedUpFile64G.String()] = utils.ZERO
							for i := range fileInfoList {
								if !fileInfoList[i].IsDir() {
									if _, ok := fileMap[fileInfoList[i].Name()]; ok {
										fileMap[fileInfoList[i].Name()] = utils.ONE
									}
								}
							}

							fileStr := ""
							for k, v := range fileMap {
								if v == utils.ZERO {
									fileStr += k + ","
								}
							}

							if fileStr != "" && !checkWarn {
								fileStr = utils.SubStr(fileStr, utils.ZERO, len(fileStr)-1)
								warnKey = define.SpeedUpFile
								warnInfo = "加速文件不存在,文件:" + fileStr
								checkWarn = true
							}

							if utils.IsNull(fileStr) && !checkWarn {

								for i := range fileInfoList {
									if !fileInfoList[i].IsDir() {
										if _, ok := fileMap[fileInfoList[i].Name()]; ok {
											str := fmt.Sprintf("%s %s | grep %s", define.DuSh, define.PathSpeedUpFile, fileInfoList[i].Name())
											b, _ := utils.ExecScript(str)
											fileSize := utils.SubStr(string(b), 0, strings.Index(string(b), "G")+1)
											if fileSize != "32G" && fileSize != "64G" {
												fileStr += fmt.Sprintf("%s:%s", fileInfoList[i].Name(), fileSize) + ","
											}
										}
									}
								}

								if fileStr != "" {
									fileStr = utils.SubStr(fileStr, utils.ZERO, len(fileStr)-1)
									warnKey = define.SpeedUpFile
									warnInfo = "加速文件大小异常,文件:" + fileStr
									checkWarn = true
								}
							}
						}
					}

					break
				case define.ProofParam:

					if res.Value == define.NodeMachine.String() {

						dirMap := make(map[string]struct{})
						dirMap[define.PathProveParameters] = struct{}{}
						dirMap[define.PathProveParents] = struct{}{}
						dirErrMsg := ""
						for k, _ := range dirMap {
							_, err := ioutil.ReadDir(k)
							if err != nil {
								dirErrMsg += k + ","
							}
						}

						if dirErrMsg != "" && !checkWarn {
							dirErrMsg = utils.SubStr(dirErrMsg, utils.ZERO, len(dirErrMsg)-utils.ONE)
							warnKey = define.ProofParam
							warnInfo = "证明参数目录不存在,目录:" + dirErrMsg
							checkWarn = true
						}

						if utils.IsNull(dirErrMsg) && !checkWarn {

							for k, _ := range dirMap {

								b, _ := utils.ExecScript(fmt.Sprintf("%s %s", define.DuSh, k))
								fileSize := utils.SubStr(string(b), 0, strings.Index(string(b), "G")+1)
								if !strings.Contains("168G,270G,322G", fileSize) {
									dirErrMsg += k + ":" + fileSize + ","
								}
							}

							if dirErrMsg != "" {
								dirErrMsg = utils.SubStr(dirErrMsg, utils.ZERO, len(dirErrMsg)-utils.ONE)
								warnKey = define.ProofParam
								warnInfo = "证明参数目录文件大小异常,目录:" + dirErrMsg
								checkWarn = true
							}
						}
					}

					break
				case define.LotusHeight:

					h, _ := strconv.Atoi(value)
					if h > utils.ONE && res.Value == define.NodeMachine.String() {
						warnKey = define.LotusHeight
						warnInfo = fmt.Sprintf("lotus高度异常,现在有高度:%d", h)
						checkWarn = true
					}

					break
				}

				if checkWarn {

					data := &pb.WarnInfo{
						WarnName:   warnNameMap[warnKey],
						WarnType:   int32(define.BusinessWarn.Int()),
						ComputerId: global.OpUUID.String(),
						WarnInfo:   warnInfo,
					}
					_, err := global.OpToGatewayClient.AddWarn(context.TODO(), data)
					if err != nil {
						return err
					}
				}
			}
		}
	}
}
