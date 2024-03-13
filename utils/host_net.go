package utils

import (
	"context"
	"fmt"
	"oplian/config"
	"oplian/define"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// StartHostPortMonitor Start the monitor script for the auxiliary test host - the server side
func StartHostPortMonitor() error {

	_, err := exec.Command("bash", "-c", define.PathServerPortMonitor).Output()
	if err != nil {
		return err
	}
	return nil
}

// StartHostPortRequest Start the bandwidth test script - client for the test host
func StartHostPortRequest(ip string, port, testTime int) string {
	ctx := context.Background()
	done := make(chan string, 1)

	go func(ctx context.Context, ip string, port, testTime int) {

		out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathServerPortRequest+" %s %d %d", ip, port, testTime)).Output()
		if err != nil {

			done <- ""
			return
		}

		if len(string(out)) < 8 {

			//return "", errors.New("script execution failed")
			done <- ""
			return
		}

		result := string(out)[:len(string(out))-1]
		result = strings.ReplaceAll(result, " ", "")
		result = result[:len(result)-8]

		sizeNum, _ := strconv.ParseFloat(result[:len(result)-1], 64)
		sizeUnit := result[len(result)-1:]

		done <- fmt.Sprintf("%.2f", sizeNum) + sizeUnit
		return
	}(ctx, ip, port, testTime)

	select {
	case speed := <-done:
		fmt.Println(ip + " " + strconv.Itoa(port) + " call successfully!!")
		return speed
	case <-time.After(config.HostNetTimeout):
		fmt.Println(ip + " " + strconv.Itoa(port) + "timeout!!")
		return ""
	}
}

func CloseHostPortMonitor() error {

	_, err := exec.Command("bash", "-c", define.PathServerKillIperf3).Output()
	if err != nil {
		return err
	}
	return nil
}

func KillBenchAndKillScript() error {

	_, err := exec.Command("bash", "-c", define.PathServerKillAllTest).Output()
	if err != nil {

		return err
	}

	return nil
}

func DealMountDir(mountDir string, nodeIP string) string {
	diskPath := strings.Split(mountDir, ",")
	if len(diskPath) == 0 {
		return ""
	}
	var dealPath string
	for key, val := range diskPath {
		var path string
		if strings.Contains(val, "/mnt/") {
			path = strings.Replace(val, "/mnt/", "/mnt/"+nodeIP+"/", 1)
		} else {
			path = "/mnt/" + nodeIP + val
		}
		if key == 0 {
			dealPath = path
		} else {
			dealPath = dealPath + "," + path
		}
	}
	return dealPath
}

func RestoreMountDir(mountDir string, nodeIP string) string {
	diskPath := strings.Split(mountDir, ",")
	if len(diskPath) == 0 {
		return ""
	}
	dealPath := strings.Replace(mountDir, "/"+nodeIP+"/", "/", -1)
	return dealPath
}
