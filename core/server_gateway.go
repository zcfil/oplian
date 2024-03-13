package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"oplian/core/internal"
	"oplian/define"
	"oplian/global"
	"oplian/utils"
)

func RunGatewayServer() {
	fmt.Println("gateway server begin")
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"username": "name1",
			"data":     "data1",
		})
	})

	r.GET("/download/file", DownloadOpFile)
	r.GET("/download/slot", DownloadSlot)
	r.GET("/download/sh", DownloadOpScript)

	r.Run(":" + define.GatewayDownLoadPort)
}

// DownloadOpFile Download the required files for OP
func DownloadOpFile(c *gin.Context) {
	fileName := c.Query("filename")
	fileList := []string{define.OPlianOpFile, define.OPlianOpC2File, define.IPAddressFile, define.ProgramLotus.String(), define.ProgramMiner.String(),
		define.ProgramWorkerTask.String(), define.ProgramWorkerStorage.String(), define.ProgramBoost.String(), define.ProgramBoostd.String()}
	if !utils.IsInStrList(fileName, fileList) {
		c.String(http.StatusBadRequest, "This file is not downloadable")
		return
	}

	var path string
	switch fileName {
	case define.OPlianOpFile, define.OPlianOpC2File:
		path = define.PathOplian
	case define.IPAddressFile:
		path = define.PathIpfsConfig
	default:
		path = define.PathIpfsProgram
	}

	fileContent, err := ioutil.ReadFile(path + fileName)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Writer.Write(fileContent)
}

// DownloadSlot Download the required plugins for OP
func DownloadSlot(c *gin.Context) {
	fileName := c.Query("filename")
	filePath := c.Query("filePath")

	fileContent, err := ioutil.ReadFile(filePath + fileName)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Writer.Write(fileContent)
}

// DownloadOpScript Download op script
func DownloadOpScript(c *gin.Context) {
	var shContent string
	shContent = "#!/bin/bash\n"

	shContent = shContent + "if ps -ef |grep -v grep |grep -w oplian-gateway; then\n    echo \"Due to the existence of the GateWay process, no file pull operation will be performed\"\nelse\n"

	intranetIP := global.LocalIP
	if intranetIP == "" {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to query the IP information of the computer room's internal network!"))
		return
	}

	curlUrl := "    curl http://" + intranetIP + ":" + define.GatewayDownLoadPort + "/download/file?filename="

	curlOP := curlUrl + define.OPlianOpFile + " > " + define.OPlianOpFile
	curlOPC2 := curlUrl + define.OPlianOpC2File + " > " + define.OPlianOpC2File

	curlConfig := curlUrl + define.IPAddressFile + " > " + define.IPAddressFile

	curlLotus := curlUrl + define.ProgramLotus.String() + " > " + define.ProgramLotus.String()
	curlMiner := curlUrl + define.ProgramMiner.String() + " > " + define.ProgramMiner.String()
	curlWorker := curlUrl + define.ProgramWorkerTask.String() + " > " + define.ProgramWorkerTask.String()
	curlStorage := curlUrl + define.ProgramWorkerStorage.String() + " > " + define.ProgramWorkerStorage.String()

	curlBoost := curlUrl + define.ProgramBoost.String() + " > " + define.ProgramBoost.String()
	curlBoostd := curlUrl + define.ProgramBoostd.String() + " > " + define.ProgramBoostd.String()

	envGatewayIP := "grep -r '" + internal.ConfigGatewayIP + "=\"" + intranetIP + "\"' ~/.bashrc || echo  'export " + internal.ConfigGatewayIP +
		"=\"" + intranetIP + "\"' >>  ~/.bashrc"
	envGatewayPort := "grep -r '" + internal.ConfigGatewayPort + "=\"" + internal.ConfigGatewayPortInfo + "\"' ~/.bashrc || echo  'export " + internal.ConfigGatewayPort +
		"=\"" + internal.ConfigGatewayPortInfo + "\"' >>  ~/.bashrc"

	shContent = shContent + "    mkdir " + define.PathOplian + "\n" + "    mkdir " + define.PathOplian + "log/\n" + "    cd " + define.PathOplian + "\n" + curlOP + "\n" + curlOPC2 + "\n" +
		"    mkdir " + define.PathIpfsConfig + "\n" + "    cd " + define.PathIpfsConfig + "\n" + curlConfig + "\n" +
		"    mkdir " + define.PathIpfsProgram + "\n" + "    cd " + define.PathIpfsProgram + "\n" + curlLotus + "\n" + curlMiner + "\n" +
		curlWorker + "\n" + curlStorage + "\n" + curlBoost + "\n" + curlBoostd + "\n" +
		"    chmod 777 " + define.PathOplian + define.OPlianOpFile + " " + define.PathOplian + define.OPlianOpC2File + "\n" +
		"    chmod 777 -R " + define.PathIpfsConfig + " " + define.PathIpfsProgram + "\n" +
		"fi\n"

	shContent = shContent + envGatewayIP + "\n" + envGatewayPort + "\n" +
		"cd " + define.PathOplian + "\n"

	shContent = shContent + "source ~/.bashrc" + "\n" + "bash" + "\n\n\n"

	//shContent = shContent + "    cd " + define.PathOplian + "\n" + "./" + define.OPlianOpFile + " init" + "\n"

	fileContent := []byte(shContent)

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "op_initialization.sh"))
	c.Writer.Write(fileContent)
}
