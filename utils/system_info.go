package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"oplian/global"
	"os/exec"
	"strconv"
	"strings"
)

const (
	T = 1
	G = 1024 * T
	M = 1024 * G
	K = 1024 * M
)

type LinuxSystemInfo struct {
	OperatingSystem string  `json:"operatingSystem"`
	SystemVersion   string  `json:"systemVersion"`
	SystemBits      int     `json:"systemBits"`
	IntranetIP      string  `json:"intranetIp"`
	InternetIP      string  `json:"internetIp"`
	SubnetMask      string  `json:"subnetMask"`
	DiskNum         int     `json:"diskNum"`
	DiskSizeSum     float64 `json:"diskSizeSum"`
	ServerDNS       string  `json:"serverDNS"`
	Gateway         string  `json:"gateway"`
	HostName        string  `json:"hostName"`
	DeviceInfo
}

type DeviceInfo struct {
	DeviceSN         string `json:"deviceSN"`
	HostManufacturer string `json:"hostManufacturer"`
	HostModel        string `json:"hostModel"`
}

func (si *OpServer) GetLinuxSystemInfo() {

	si.SystemInfo.OperatingSystem = "Linux"
	si.SystemInfo.SystemVersion, si.SystemInfo.SystemBits = GetSystemVersion()

	si.SystemInfo.InternetIP, _ = GetInternetIP()

	si.SystemInfo.IntranetIP, si.SystemInfo.SubnetMask = GetIntranetIP()

	si.SystemInfo.DiskNum, si.SystemInfo.DiskSizeSum = getDiskInfo()

	deviceInfo := getDeviceInfo()
	si.SystemInfo.DeviceSN = deviceInfo.DeviceSN
	si.SystemInfo.HostManufacturer = deviceInfo.HostManufacturer
	si.SystemInfo.HostModel = deviceInfo.HostModel

	si.SystemInfo.ServerDNS = getDNSInfo()

	si.SystemInfo.Gateway = getGatewayInfo()

	si.SystemInfo.HostName = getLinuxName()
	return
}

// GetSystemVersion Get linux system version information
func GetSystemVersion() (string, int) {
	out, err := exec.Command("bash", "-c", "lsb_release -a | grep Description:").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `lsb_release -a | grep Description:` Filed", err.Error())
		return "", 0
	}
	if !strings.Contains(string(out), "Description") {
		log.Println("该系统命令`lsb_release -a | grep Description:`无法正常运行!")
		return "", 0
	}
	versionInfo := string(out)
	// 判断是否是多行数据
	if strings.Contains(string(out), "\n") {
		strVersions := strings.Split(string(out), "\n")
		for _, val := range strVersions {
			if strings.Contains(val, "Description") {
				versionInfo = val
			}
		}
	}

	strBs := strings.Split(versionInfo, ":")

	out, err = exec.Command("bash", "-c", "getconf LONG_BIT").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `getconf LONG_BIT` Filed", err.Error())
		return strings.Trim(strBs[1], " "), 0
	}

	systemBits, _ := strconv.Atoi(string(out)[:len(string(out))-1])

	return strings.Trim(strings.Trim(strBs[1], " "), "\t"), systemBits
}

// GetInternetIP Get the external network IP
func GetInternetIP() (string, error) {
	out, err := exec.Command("bash", "-c", "curl ifconfig.me | grep .").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `curl ifconfig.me | grep .` Filed", err.Error())
		return "", err
	}

	strBs := strings.Split(string(out), "\n")
	if len(strBs) > 1 {
		strB := strings.Replace(strBs[len(strBs)-2], " ", "", -1)
		if len(strB) > 20 {
			return "--", nil
		}
		return strB, nil
	} else {
		return "--", nil
	}
}

// GetIntranetIP Gets the Intranet IP and subnet mask
func GetIntranetIP() (string, string) {
	var mask string
	maskByte, err := exec.Command("bash", "-c", "ifconfig | grep broadcast| grep  -v  $(ifconfig |grep docker -A3|awk  '{print $2}'|egrep '([0-9]{1,3}\\.){3}[0-9]{1,3}' || echo w)|grep -E '(10\\.|192\\.168\\.|172\\.)' |awk   '{print  $4}'").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `iifconfig | grep broadcast| grep  -v  $(if...` Filed", err.Error())
		return global.LocalIP, "--"
	}
	if strings.Contains(string(maskByte), "\n") {
		masks := strings.Split(string(maskByte), "\n")
		if len(masks) >= 2 {
			mask = masks[len(masks)-2]
		} else {
			mask = masks[0]
		}
	}
	if strings.Contains(mask, "command not found") {
		return global.LocalIP, "--"
	}
	return global.LocalIP, mask
}

// GetIntranetIPList Get the Intranet IP
func GetIntranetIPList() []string {
	var inIp []string
	inIpByte, err := exec.Command("bash", "-c", "ifconfig | grep broadcast |awk '{print $2}'|grep  -v  $(ifconfig |grep docker -A3|awk  '{print $2}'|egrep '([0-9]{1,3}\\.){3}[0-9]{1,3}'|| echo w) |grep -E '(10\\.|192\\.168\\.|172\\.)'").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `iifconfig | grep broadcast| grep  -v  $(if...` Filed", err.Error())
		return []string{}
	}
	if strings.Contains(string(inIpByte), "1") {
		list := strings.Split(string(inIpByte), "\n")
		for _, val := range list {
			if len(val) > 0 {
				inIp = append(inIp, val)
			}
		}
	}
	return inIp
}

type DiskInfo struct {
	Size float64 `json:"size"`
	Unit string  `json:"unit"`
}

// getDiskInfo Gets the number and size of hard disks
func getDiskInfo() (int, float64) {
	out, err := exec.Command("bash", "-c", "lsblk").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `lsblk` Filed", err.Error())
		return 0, 0
	}

	strBs := strings.Split(string(out), "\n")
	var diskInfo []string
	for key, val := range strBs {
		if key == 0 || strings.Contains(val, "NAME") {
			continue
		}
		if !strings.Contains(val, "loop") && !strings.Contains(val, "└─") && !strings.Contains(val, "├─") {
			diskInfo = append(diskInfo, val)
		}
	}
	var newDiskInfoStruct []DiskInfo
	for _, val := range diskInfo {
		if val == "" {
			continue
		}
		var newDiskInfo []string
		diskInfo := strings.Split(val, " ")
		for i := 0; i < len(diskInfo); i++ {
			if diskInfo[i] == "" {
				continue
			}
			newDiskInfo = append(newDiskInfo, diskInfo[i])
		}

		if len(newDiskInfo) > 3 {
			strLen := len(newDiskInfo[3])
			if strLen > 1 {
				size, _ := strconv.ParseFloat(newDiskInfo[3][:strLen-1], 64)
				unit := newDiskInfo[3][strLen-1 : strLen]
				info := DiskInfo{
					Size: size,
					Unit: unit,
				}
				newDiskInfoStruct = append(newDiskInfoStruct, info)
			}
		}
	}

	var sizeSum float64
	for _, val := range newDiskInfoStruct {
		sizeSum = sizeSum + DealSizeUnit(val.Size, val.Unit)
	}

	return len(newDiskInfoStruct), Decimal(sizeSum/G, 1)
}

// DealSizeUnit  The final result is G
func DealSizeUnit(size float64, unit string) float64 {
	var result float64
	switch unit {
	case "P":
		result = size * M
	case "T":
		result = size * G
	case "G":
		result = size
	case "M":
		result = size / G
	case "K":
		result = size / M
	case "B":
		result = 0
	}
	return result
}

// Decimal float64 Preserves the number of digits after the decimal point
// value float64 floating-point number
// prec int Specifies the number of digits after the decimal point
func Decimal(value float64, prec int) float64 {
	value, _ = strconv.ParseFloat(strconv.FormatFloat(value, 'f', prec, 64), 64)
	return value
}

// getDeviceInfo Get information about the system hardware
func getDeviceInfo() *DeviceInfo {
	var deviceInfo DeviceInfo

	out, err := exec.Command("bash", "-c", "dmidecode | grep \"System Information\" -A9 | egrep \"Manufacturer|Product|Serial\"").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `dmidecode | grep \"System Information\" -A9 | egrep \"Manufacturer|Product|Serial\"` Filed", err.Error())
		return &deviceInfo
	}
	systemInfos := strings.Split(string(out), "\n")
	for _, val := range systemInfos {
		if val == "" {
			continue
		}
		strs := strings.Split(val, ":")
		if len(strs) > 1 {
			if strings.Contains(strs[0], "Manufacturer") {
				deviceInfo.HostManufacturer = strings.Trim(strs[1], " ")
			}
			if strings.Contains(strs[0], "Product Name") {
				deviceInfo.HostModel = strings.Trim(strs[1], " ")
			}
			if strings.Contains(strs[0], "Serial Number") {
				deviceInfo.DeviceSN = strings.Trim(strs[1], " ")
			}
		}
	}
	return &deviceInfo
}

// getDNSInfo Get the Linux system DNS
func getDNSInfo() string {

	out, err := exec.Command("bash", "-c", "sudo systemd-resolve --status | grep 'DNS Servers' -A2").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `sudo systemd-resolve --status | grep 'DNS Servers' -A2` Filed", err.Error())
		return ""
	}
	str := strings.ReplaceAll(string(out), " ", "")
	strBs := strings.Split(str, "\n")

	ipStr := ""

	for _, val := range strBs {
		if val == "" {
			continue
		}
		if strings.Contains(val, "DNS") {
			strDNSs := strings.Split(val, ":")
			if len(strDNSs) > 1 {
				ipStr = strDNSs[1]
			}
			continue
		}
		ipStr = ipStr + "," + val
	}
	if strings.Contains(ipStr, "commandnotfound") {
		return "--"
	}
	return ipStr
}

// getGatewayInfo Get gateway IP
func getGatewayInfo() string {
	// 获取设备信息
	out, err := exec.Command("bash", "-c", "ip route show | grep \"default via\"").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `ip route show | grep \"default via\"` Filed", err.Error())
		return "--"
	}

	strs := strings.Split(string(out), " ")
	if len(strs) > 2 {
		return strs[2]
	} else {
		return "--"
	}
}

// getLinuxName Get host name
func getLinuxName() string {
	// 获取系统名称
	out, err := exec.Command("bash", "-c", `hostnamectl | grep 'Static hostname' | awk '{print $3}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `hostnamectl | grep 'Static hostname' | awk '{print $3}'` Filed", err.Error())
		return "--"
	}

	str := string(out)
	if len(str) > 0 {
		return str[:len(str)-1]
	} else {
		return "--"
	}

}

// CompareTwoSizes Compare the sizes of the two disks in the format of a number + uppercase unit, for example, 3T 4G 5M
func CompareTwoSizes(sizeA, sizeB string) bool {

	if sizeA == "" {
		sizeA = "0T"
	}
	if sizeB == "" {
		sizeB = "0T"
	}

	sizeANum, _ := strconv.ParseFloat(sizeA[:len(sizeA)-1], 64)
	sizeAUnit := strings.ToUpper(sizeA[len(sizeA)-1:])

	sizeBNum, _ := strconv.ParseFloat(sizeB[:len(sizeB)-1], 64)
	sizeBUnit := strings.ToUpper(sizeB[len(sizeB)-1:])

	sizeASum := DealSizeUnit(sizeANum, sizeAUnit)
	sizeBSum := DealSizeUnit(sizeBNum, sizeBUnit)
	return sizeASum >= sizeBSum
}

// AddInputSizes Add the size of the stored data, for example, 3G + 500M. The unit of the result is G and two decimal places are reserved
func AddInputSizes(sizes []string) string {
	var sizeSum float64
	for _, val := range sizes {

		if val == "" {
			val = "0G"
		}
		sizeNum, _ := strconv.ParseFloat(val[:len(val)-1], 64)
		sizeUnit := val[len(val)-1:]
		sizeSum += DealSizeUnit(sizeNum, sizeUnit)
	}
	return fmt.Sprintf("%.2f", sizeSum) + "G"
}

// HostNetStatus Test the node machine network
func HostNetStatus() bool {
	out, err := exec.Command("bash", "-c", `ping www.baidu.com -c3 | grep '3 received'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `ping www.baidu.com -c3 | grep '3 received'` Filed", err.Error())
		return false
	}
	return len(string(out)) > 0
}

// AddInputTwoSizes Add two stored data sizes, such as 3T + 500G, and return the result in T units, reserving two decimal places
func AddInputTwoSizes(sizeA, sizeB string) string {

	if sizeA == "" {
		sizeA = "0T"
	}
	if sizeB == "" {
		sizeB = "0T"
	}
	var sizeSum float64

	sizeNumA, _ := strconv.ParseFloat(sizeA[:len(sizeA)-1], 64)
	sizeUnitA := sizeA[len(sizeA)-1:]
	sizeSum += DealSizeUnit(sizeNumA, sizeUnitA)

	sizeNumB, _ := strconv.ParseFloat(sizeB[:len(sizeB)-1], 64)
	sizeUnitB := sizeB[len(sizeB)-1:]
	sizeSum += DealSizeUnit(sizeNumB, sizeUnitB)

	return fmt.Sprintf("%.2f", sizeSum/1024) + "T"
}

// PercentageOfTwoSize Percentage of the size of both disks
func PercentageOfTwoSize(sizeA, sizeB string) float64 {
	if sizeA == "" || sizeB == "" {
		return 0
	}
	// 前置处理,防止传空报错
	if sizeA == "" {
		sizeA = "0T"
	}
	if sizeB == "" {
		sizeB = "0T"
	}

	sizeANum, _ := strconv.ParseFloat(sizeA[:len(sizeA)-1], 64)
	sizeAUnit := strings.ToUpper(sizeA[len(sizeA)-1:])

	sizeBNum, _ := strconv.ParseFloat(sizeB[:len(sizeB)-1], 64)
	sizeBUnit := strings.ToUpper(sizeB[len(sizeB)-1:])
	if sizeANum == 0 || sizeBNum == 0 {
		return 0
	}

	sizeASum := DealSizeUnit(sizeANum, sizeAUnit)
	sizeBSum := DealSizeUnit(sizeBNum, sizeBUnit)

	percentage := decimal.NewFromFloat(math.Ceil(sizeASum / sizeBSum * 1e2))
	sizePer, _ := percentage.Mul(decimal.NewFromFloat(1e-2)).Float64()
	return sizePer
}

// InstallNumactl Install numactl software
func InstallNumactl() {

	b, err := ExecuteScript("which numactl")
	if err != nil {
		log.Println("Install numactl err:", err)
	}

	if !strings.Contains(b, "numactl") {
		outs, _ := exec.Command("bash", "-c", "apt-get install numactl -y").CombinedOutput()
		log.Println("安装numactl：" + string(outs))
	}

	b, err = ExecuteScript("which unzip")
	if err != nil {
		log.Println("Install unzip err:", err)
	}

	if !strings.Contains(b, "unzip") {
		outs, _ := exec.Command("bash", "-c", "apt-get install unzip -y").CombinedOutput()
		log.Println("安装unzip：" + string(outs))
	}
}
