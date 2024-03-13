package utils

import (
	"log"
	"os/exec"
	"strings"
)

func GetRamTypeAndNum() (int, string) {
	out, err := exec.Command("bash", "-c", `dmidecode -t 17  | grep 'Type:' | awk '{print $2}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `dmidecode -t 17  | grep 'Type:' | awk '{print $2}'` Filed", err.Error())
		return 0, ""
	}

	strBs := strings.Split(string(out), "\n")
	var ramNum int
	var ramType string

	for _, val := range strBs {
		if len(val) == 0 {
			continue
		}
		ramNum++
		if strings.Contains(val, "Unknown") {
			continue
		}
		ramType = val
	}

	return ramNum, ramType
}
