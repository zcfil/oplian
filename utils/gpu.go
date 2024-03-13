package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"oplian/config"
	"oplian/define"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type GPU struct {
	Gpus []GPUInfo `json:"gpus"`
}

type GPUInfo struct {
	Mark    string `json:"mark"`
	Brand   string `json:"brand"`
	GpuInfo string `json:"gpuInfo"`
	TotalMB string `json:"totalMb"`
	UseRate string `json:"useRate"`
}

func (si *OpServer) GetGPUInfo() {
	var gpuInfos []GPUInfo

	si.GPU.Gpus = gpuInfos

	out, err := exec.Command("bash", "-c", "lspci | grep VGA").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `lspci | grep VGA` Filed", err.Error())
		return
	}

	isContainNvidia := false
	if strings.Contains(string(out), "NVIDIA") {
		isContainNvidia = true
	}

	if isContainNvidia {

		ch := make(chan []GPUInfo, 1)
		ctx, cancel := context.WithCancel(context.Background())

		go func(ctx context.Context, cancel context.CancelFunc) {
			fmt.Println("begin nvidia-smi!")

			out, err := exec.Command("bash", "-c", "timeout 6 nvidia-smi --query-gpu=index,name,memory.total,power.draw,power.limit --format=csv,noheader,nounits").CombinedOutput()
			if err != nil && !strings.Contains(err.Error(), "exit status") {
				log.Println("cmd `timeout 6 nvidia-smi --query-gpu=index,name,memory.total,power.draw,power.limit --format=csv,noheader,nounits` Filed", err.Error())
				return
			}

			if !strings.Contains(string(out), "GeForce") {
				log.Println("该系统命令`timeout 6 nvidia-smi --query-gpu=index,name,memory.total,power.draw,power.limit --format=csv,noheader,nounits`无法正常运行!")
				return
			}

			strBs := strings.Split(string(out), "\n")

			for _, v := range strBs {
				if len(v) == 0 {
					continue
				}
				str := strings.Split(v, ",")
				if len(str) > 4 {
					powerDraw, _ := strconv.ParseFloat(strings.TrimSpace(str[3]), 64)
					powerLimit, _ := strconv.ParseFloat(strings.TrimSpace(str[4]), 64)
					useRate := math.Trunc(powerDraw/powerLimit*1e2 + 0.5)
					info := GPUInfo{
						Mark:    str[0],
						Brand:   "NVIDIA",
						GpuInfo: strings.TrimSpace(str[1]),
						TotalMB: strings.TrimSpace(str[2]) + "M",
						UseRate: Float64ToString(useRate),
					}
					gpuInfos = append(gpuInfos, info)
				}
			}
			fmt.Println("deal nvidia-smi success!")
			ch <- gpuInfos
			close(ch)
			cancel()
		}(ctx, cancel)

		select {
		case <-ctx.Done():
			fmt.Println("get GPU info successfully!")
			si.GPU.Gpus = <-ch
			return
		case <-time.After(time.Second * 8):
			fmt.Println("get GPU info timeout!")
			return
		}
	}
}

func GetGPURunBenchTime() string {
	ctx := context.Background()
	done := make(chan string, 1)

	go func(ctx context.Context) {

		pathIpfs := define.PathIpfsScript

		out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathIpfsScriptRunBenchTest+" %s %s", pathIpfs, define.MainDisk)).Output()
		if err != nil && !strings.Contains(err.Error(), "exit runTime") {

			done <- ""
			return
		}
		runTime := string(out)
		if len(runTime) > 0 {
			done <- runTime[:len(runTime)-1]
			return
		} else {
			done <- ""
			return
		}
	}(ctx)
	select {
	case speed := <-done:
		fmt.Println("GPU run bench successfully!!")
		return speed
	case <-time.After(config.HostGPUTestTimeout):
		fmt.Println("GPU run bench timeout!!")
		return ""
	}
}

func CloseGPURunBench() error {

	_, err := exec.Command("bash", "-c", define.PathServerKillBench).Output()
	if err != nil && !strings.Contains(err.Error(), "exit runTime") {

		return err
	}
	return nil
}

// 测试GPU的驱动是否正常
func CheckGPUDrive() error {
	// 获取对应的显卡存储和占用
	out, err := exec.Command("bash", "-c", "timeout 6 nvidia-smi --query-gpu=index,name,memory.total,utilization.gpu --format=csv,noheader,nounits").CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `timeout 6 nvidia-smi --query-gpu=index,name,memory.total,utilization.gpu --format=csv,noheader,nounits` Filed", err.Error())
		return err
	}

	if !strings.Contains(string(out), "GeForce") {
		log.Println("该系统命令`timeout 6 nvidia-smi --query-gpu=index,name,memory.total,utilization.gpu --format=csv,noheader,nounits`无法正常运行!")
		return errors.New("system commands cannot function properly")
	}
	return nil
}
