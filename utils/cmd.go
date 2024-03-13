package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func ExecuteScript(script string) (string, error) {
	cmd := exec.Command("bash", "-c", script)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

func ProgramRun(name string) bool {

	b, err := exec.Command("bash", "-c", fmt.Sprintf("ps -ef |grep -v grep |grep %s |wc -l", name)).CombinedOutput()
	if err != nil {
		log.Println("ProgramRun err:", err)
		return false
	}
	str := strings.Replace(string(b), "\n", "", -1)
	str = strings.Replace(str, " ", "", -1)
	runTotal, _ := strconv.Atoi(str)
	if runTotal >= 1 {
		return true
	} else {
		return false
	}

}
