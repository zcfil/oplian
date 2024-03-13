package utils

import (
	"fmt"
	"log"
	"oplian/define"
	"os/exec"
	"strings"
)

type HostDiskInfo struct {
	Mounted string `json:"mounted"`
	Size    string `json:"size"`
	Avail   string `json:"avail"`
}

type StorageOpDiskInfo struct {
	Filesystem string `json:"filesystem"`
	Size       string `json:"size"`
	Used       string `json:"used"`
	UsePercent string `json:"usePercent"`
	Mountpoint string `json:"mountpoint"`
}

// GetOpDiskInfo Gets the hard disk information related to the mnt mounted on the host
func GetOpDiskInfo() []HostDiskInfo {
	out, err := exec.Command("bash", "-c", `df -h | awk '{ print $6,$2,$4 }' | grep 'T'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `lsblk` Filed", err.Error())
		return []HostDiskInfo{}
	}

	strBs := strings.Split(string(out), "\n")
	var diskInfo []HostDiskInfo

	for _, val := range strBs {
		var info HostDiskInfo
		if strings.Contains(val, "T") {
			strBs := strings.Split(val, " ")
			if len(strBs) < 3 {
				continue
			}
			if !strings.Contains(strBs[1], "T") {
				continue
			}
			info = HostDiskInfo{
				Mounted: strBs[0],
				Size:    strBs[1],
				Avail:   strBs[2],
			}
		} else {
			continue
		}
		diskInfo = append(diskInfo, info)
	}
	return diskInfo
}

// GetStorageOpDiskInfo Get the disk information of the storage drive
func GetStorageOpDiskInfo() []StorageOpDiskInfo {
	out, err := exec.Command("bash", "-c", `df -h | grep '/dev/' | awk '{ print $1,$2,$3,$5,$6 }'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `df -h | grep '/dev/' | awk '{ print $1,$2,$3,$5,$6 }'` Filed", err.Error())
		return []StorageOpDiskInfo{}
	}
	strBs := strings.Split(string(out), "\n")

	var storageInfo []StorageOpDiskInfo

	for _, v := range strBs {
		if len(v) == 0 || !strings.Contains(v, "/dev") {
			continue
		}
		str := strings.Split(v, " ")
		if len(str) < 5 {
			continue
		}
		info := StorageOpDiskInfo{
			Filesystem: str[0],
			Size:       str[1],
			Used:       str[2],
			UsePercent: str[3][:len(str[3])-1],
			Mountpoint: str[4],
		}
		storageInfo = append(storageInfo, info)
	}
	return storageInfo
}

// GetStorageOpOsDiskInfo Obtain information about the disk corresponding to the OS system disk
func GetStorageOpOsDiskInfo() OpDiskInfo {
	var diskInfo OpDiskInfo

	out, err := exec.Command("bash", "-c", `lsblk -l | grep -w / | awk '{printf ("%s",$1)}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `lsblk -l | grep -w / | awk '{printf (\"%s\",$1)}'` Filed", err.Error())
		return diskInfo
	}

	osDiskName := string(out)
	diskInfo.IsOs = 1
	diskInfo.Name = osDiskName

	out, err = exec.Command("bash", "-c", `df -h | grep -w / | awk '{print $2,$3,$5}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `df -h | grep -w / | awk '{print $1,$2,$3,$5}'` Filed", err.Error())
		return diskInfo
	}

	outStr := string(out)
	if len(outStr) > 0 {
		outStrs := strings.Split(outStr[:len(outStr)-1], "\n")
		// 相关信息赋值
		strBs := strings.Split(outStrs[len(outStrs)-1], " ")
		if len(strBs) >= 3 {
			diskInfo.TotalSize = strBs[0]
			diskInfo.UsedSize = strBs[1]
			if len(strBs[2]) > 1 {
				diskInfo.UsedPercent = strBs[2][:len(strBs[2])-1]
			}
		}
	}

	// 获取系统盘硬件信息
	out, err = exec.Command("bash", "-c", `smartctl -i /dev/`+osDiskName+` | grep -E 'Device Model:|Model Number'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println(`cmd smartctl -i /dev/`+osDiskName+` | grep -E 'Device Model:|Model Number' Filed`, err.Error())
		return diskInfo
	}

	infoStr := string(out)
	if len(infoStr) > 0 {

		infos := strings.Split(infoStr[:len(infoStr)-1], " ")
		if len(infos) > 3 {
			diskInfo.Brand = infos[len(infos)-2]
			diskInfo.Model = infos[len(infos)-1]
		}
	}

	return diskInfo
}

// OpDiskMd0Info Obtain information about the disk corresponding to md0
func GetDiskMd0Info() OpDiskMd0Info {
	var md0Info OpDiskMd0Info

	out, err := exec.Command("bash", "-c", `mdadm --detail /dev/md0 | grep "Active Devices" | awk '{print $NF}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `mdadm --detail /dev/md0 | grep \"Active Devices\" | awk '{print $NF}'` Filed", err.Error())
		return md0Info
	}
	str := string(out)
	if strings.Contains(str, "No such file or directory") {
		log.Println("/dev/md0 does not exist")
	} else {
		if len(str) > 0 {
			md0Info.DiskNum = str[:len(str)-1]
		} else {
			md0Info.DiskNum = "0"
		}
	}

	out, err = exec.Command("bash", "-c", `df -h | grep md0 | awk '{printf ("%s %s %s",$2,$3,$5)}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `df -h | grep md0 | awk '{printf (\"%s %s %s\",$2,$3,$5)}'` Filed", err.Error())
		return md0Info
	}

	md0Str := string(out)
	if len(md0Str) == 0 {
		log.Println("cmd `df -h | grep md0 | awk '{printf (\"%s %s %s\",$2,$3,$5)}'` 获取结果为空")
		return md0Info
	}
	// 相关信息赋值
	strBs := strings.Split(md0Str[:len(md0Str)-1], " ")
	if len(strBs) >= 3 {
		md0Info.TotalSize = strBs[0]
		md0Info.UsedSize = strBs[1]
		md0Info.UsedPercent = strBs[2]
	}
	return md0Info
}

// GetDiskIOStatus Obtain disk I/O information
func GetDiskIOStatus() (string, error) {

	out, err := exec.Command("bash", "-c", `dmesg  | grep -i "I/O error" | wc -l`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `dmesg  | grep -i \"I/O error\" | wc -l` Filed", err.Error())
		return "", err
	}
	status := string(out)
	if len(status) > 0 {
		return status[:len(status)-1], err
	} else {
		return "0", nil
	}
}

// GetMd0DiskSize Gets the disk size of md0
func GetMd0DiskSize() (string, error) {

	out, err := exec.Command("bash", "-c", `df -h | grep '/md0' | awk '{print $4}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `df -h | grep '/md0' | awk '{print $4}'` Filed", err.Error())
		return "", err
	}
	status := string(out)
	if len(status) > 0 {
		return status[:len(status)-1], err
	} else {
		return "0T", nil
	}
}

// DiskReassemblyArray The disk is regrouped into the disk array
func DiskReassemblyArray() error {

	diskNum, diskStr, done := GetDiskArrayList()
	if !done {
		return nil
	}

	if out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathDiskReassemblyArray+" %d %s",
		diskNum, diskStr[:len(diskStr)-1])).Output(); err != nil {
		log.Println(string(out))
		return fmt.Errorf("执行磁盘重新组阵列的脚本脚本失败:%s", err.Error())
	}
	return nil
}

// GetDiskArrayList Gets the disk information of the array to be grouped
func GetDiskArrayList() (int, string, bool) {
	var s OpServer

	s.GetStorageInfo()

	out, err := exec.Command("bash", "-c", "lsblk -d -o name,rota").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `lsblk -d -o name,rota` Filed", err.Error())
		return 0, "", false
	}

	rotas := strings.Split(string(out), "\n")

	newRotas := rotas[1 : len(rotas)-1]
	diskStatus := make(map[string]string)
	for _, val := range newRotas {
		if len(val) == 0 {
			continue
		}
		val = strings.ReplaceAll(val, " ", "")
		diskStatus[val[:len(val)-1]] = val[len(val)-1 : len(val)]
	}

	var diskSums []DiskSum
	for _, val := range s.Storage {
		if diskStatus[val.Name] == "1" {
			continue
		}
		var sumInfo DiskSum
		if len(diskSums) == 0 {
			sumInfo.Name = val.Name
			sumInfo.Size = val.Size
			sumInfo.SizeSum += val.Size
			diskSums = append(diskSums, sumInfo)
			continue
		}
		isContain := false
		for i := 0; i < len(diskSums); i++ {
			if val.Size == diskSums[i].Size && !strings.Contains(diskSums[i].Name, val.Name) {
				isContain = true
				diskSums[i].Name = diskSums[i].Name + "," + val.Name
				diskSums[i].SizeSum += val.Size
				break
			}
		}
		if !isContain {
			sumInfo.Name = val.Name
			sumInfo.Size = val.Size
			sumInfo.SizeSum += val.Size
			diskSums = append(diskSums, sumInfo)
		}
	}

	// 处理初始化硬盘脚本所需信息
	var maxSum uint
	for _, val := range diskSums {
		if val.SizeSum > maxSum {
			maxSum = val.SizeSum
		}
	}

	var choseDisks []string
	var diskNum int
	for _, val := range diskSums {
		if val.SizeSum == maxSum {
			choseDisks = strings.Split(val.Name, ",")
			diskNum = len(choseDisks)
		}
	}

	diskStr := ""
	for _, val := range choseDisks {
		diskStr += "/dev/" + val + ","
	}
	return diskNum, diskStr, true
}

// ChmodDirectory Folder path file recursive authorization
func ChmodDirectory(dirPath string) error {

	_, err := exec.Command("bash", "-c", fmt.Sprintf("chmod -R 777 %s", dirPath)).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd "+fmt.Sprintf("chmod -R 777 %s", dirPath)+" Filed", err.Error())
		return err
	}
	return nil
}
