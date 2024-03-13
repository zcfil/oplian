package utils

import (
	"log"
	"math"
	"oplian/config"
	"os/exec"
	"strconv"
	"strings"
)

type HostMonitor struct {
	CPUUseRate     float32         `json:"cpuUseRate"`
	CPUTemperature string          `json:"cpuTemperature"`
	MonitorDisk    MonitorDiskInfo `json:"diskUseRate"`
	//MemoryUseRate  float32          `json:"memoryUseRate"`
	GPUUseInfo []GPUMonitorInfo `json:"gpuUseRate"`
}

type GPUMonitorInfo struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	UseRate float32 `json:"useRate"`
}

type MonitorDiskInfo struct {
	Size       string  `json:"size"`
	Used       string  `json:"used"`
	UsePercent float64 `json:"usePercent"`
}

func (h *HostMonitor) GetHostMonitorInfo(hostClassify int64) {
	h.CPUUseRate = GetCPUUseRate()
	h.MonitorDisk = getDiskUseRate(hostClassify)
	//h.MemoryUseRate = getMemoryUseRate()
	h.GPUUseInfo = getGPUUseRate()
	h.CPUTemperature = getCPUTemperature()
}

// GetCPUUseRate Get the CPU usage ratio
func GetCPUUseRate() float32 {
	out, err := exec.Command("bash", "-c", `top -bn1 | head -n 3 | grep %Cpu | awk '{print $2}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd get CPU use rate Filed", err.Error())
		return 0
	}
	useRate, _ := strconv.ParseFloat(string(out), 64)
	if useRate == 0 {
		useRate = 1
	}
	return float32(math.Ceil(useRate))
}

// getDiskUseRate Get the disk usage ratio
func getDiskUseRate(hostClassify int64) MonitorDiskInfo {
	var monitorDiskInfo MonitorDiskInfo

	if hostClassify == config.HostStorageType || hostClassify == config.HostDCStorageType || hostClassify == 0 {

		storageInfo := GetStorageOpDiskInfo()
		for _, v := range storageInfo {
			monitorDiskInfo.Size = AddInputTwoSizes(monitorDiskInfo.Size, v.Size)
			monitorDiskInfo.Used = AddInputTwoSizes(monitorDiskInfo.Used, v.Used)
		}
		monitorDiskInfo.UsePercent = PercentageOfTwoSize(monitorDiskInfo.Used, monitorDiskInfo.Size) * 100
		if monitorDiskInfo.UsePercent == 0 {
			monitorDiskInfo.UsePercent = 1
		}
		monitorDiskInfo.UsePercent = math.Ceil(monitorDiskInfo.UsePercent)
	} else {

		diskMd0Info := GetDiskMd0Info()

		monitorDiskInfo.Size = diskMd0Info.TotalSize
		monitorDiskInfo.Used = diskMd0Info.UsedSize
		usedPercent, _ := strconv.ParseFloat(diskMd0Info.UsedPercent, 64)
		monitorDiskInfo.UsePercent = usedPercent
		if monitorDiskInfo.UsePercent == 0 {
			monitorDiskInfo.UsePercent = 1
		}
		monitorDiskInfo.UsePercent = math.Ceil(monitorDiskInfo.UsePercent)
	}
	return monitorDiskInfo
}

// getMemoryUseRate Get the memory usage ratio
func getMemoryUseRate() float32 {
	out, err := exec.Command("bash", "-c", `free -m | awk 'NR==2{printf "%.2f", $3*100/$2 }'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd get memory use rate Filed", err.Error())
		return 0
	}
	useRate, _ := strconv.ParseFloat(string(out), 64)
	if useRate == 0 {
		useRate = 1
	}
	return float32(math.Ceil(useRate))
}

// getGPUUseRate Get the GPU usage ratio
func getGPUUseRate() []GPUMonitorInfo {
	out, err := exec.Command("bash", "-c", `timeout 6 nvidia-smi --query-gpu=index,name,power.draw,power.limit --format=csv,noheader,nounits`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd get GPU use rate Filed", err.Error())
		return []GPUMonitorInfo{}
	}

	if strings.Contains(string(out), "failed") {
		return []GPUMonitorInfo{}
	}

	infos := strings.Split(strings.TrimSpace(string(out)[:len(string(out))-1]), "\n")

	var gpuInfos []GPUMonitorInfo
	for _, val := range infos {
		if len(val) == 0 {
			continue
		}
		strs := strings.Split(val, ",")
		if len(strs) > 3 {
			powerDraw, _ := strconv.ParseFloat(strings.TrimSpace(strs[2]), 64)
			powerLimit, _ := strconv.ParseFloat(strings.TrimSpace(strs[3]), 64)
			useRate := math.Trunc(powerDraw/powerLimit*1e2 + 0.5)
			if useRate == 0 {
				useRate = 1
			}

			gpuInfo := GPUMonitorInfo{
				ID:      strs[0],
				Name:    strings.TrimSpace(strs[1]),
				UseRate: float32(useRate),
			}
			gpuInfos = append(gpuInfos, gpuInfo)
		}
	}

	return gpuInfos
}

// getCPUTemperature Obtaining the CPU Temperature
func getCPUTemperature() string {
	out, err := exec.Command("bash", "-c", `sensors | grep 'Tdie:' | awk '{printf "%s",$2}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd get CPU temperature Filed", err.Error())
		return "--"
	}
	if len(string(out)) == 0 {
		return "--"
	}
	if strings.Contains(string(out), "not found") {
		log.Println("Command 'sensors' not found, but can be installed with")
		return "--"
	}
	return string(out)[1:]
}
