package utils

import (
	"errors"
	"fmt"
	"log"
	"oplian/config"
	"oplian/define"
	"os/exec"
	"strconv"
	"strings"
)

// GetLogOvertimeStatus Checks whether the log times out
func GetLogOvertimeStatus() (bool, error) {

	out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathServerLogOvertime+" %s", define.MainDisk)).Output()
	if err != nil {
		return false, err
	}
	if string(out)[:len(string(out))-1] == config.PatrolTestTrue {
		return true, nil
	} else {
		return false, nil
	}
}

// GetWdpostBalance Get wallet balance
func GetWdpostBalance() (int, error) {

	out, err := exec.Command("bash", "-c", define.PathServerWdpostBalance).Output()
	if err != nil {
		return 0, err
	}
	if !strings.Contains(string(out), "post_balnce:") {
		return 0, errors.New("get wdpost balance failed")
	}
	strs := strings.Split(string(out)[:len(string(out))-1], ":")
	if len(strs) > 1 {
		balance, _ := strconv.Atoi(strs[1])
		return balance, nil
	} else {
		return 0, nil
	}
}

// GetLotusHigh Get lotus height
func GetLotusHigh() (bool, error) {
	out, err := exec.Command("bash", "-c", define.PathServerLotusHigh).Output()
	if err != nil {
		return false, err
	}
	if string(out)[:len(string(out))-1] == config.PatrolTestTrue {
		return true, nil
	} else {
		return false, nil
	}
}

type LotusPackageVersion struct {
	LotusVersion  string `json:"lotusVersion"`
	MinerVersion  string `json:"minerVersion"`
	BoostdVersion string `json:"boostdVersion"`
}

// GetLotusPackageVersion Gets the lotus package version
func GetLotusPackageVersion() (LotusPackageVersion, error) {
	var packageVersion LotusPackageVersion

	outLotus, err := exec.Command("bash", "-c", `/usr/local/sbin/lotus -v | awk '{print $3}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `/usr/local/sbin/lotus -v | awk '{print $3}'` Filed", err.Error())
		return packageVersion, err
	}
	if strings.Contains(string(outLotus), "/usr/local/sbin") {
		return packageVersion, errors.New("failed to obtain lotus version")
	}
	if len(string(outLotus)) > 1 {
		packageVersion.LotusVersion = string(outLotus)[:len(string(outLotus))-1]
	}

	minerLotus, err := exec.Command("bash", "-c", `/usr/local/sbin/lotus-miner -v | awk '{print $3}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `/usr/local/sbin/lotus-miner -v | awk '{print $3}'` Filed", err.Error())
		return packageVersion, err
	}
	if strings.Contains(string(minerLotus), "/usr/local/sbin") {
		return packageVersion, errors.New("failed to obtain miner version")
	}
	if len(string(minerLotus)) > 1 {
		packageVersion.MinerVersion = string(minerLotus)[:len(string(minerLotus))-1]
	}

	boostdLotus, err := exec.Command("bash", "-c", `/usr/local/sbin/boostd -v | awk '{print $3}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `/usr/local/sbin/boostd -v | awk '{print $3}'` Filed", err.Error())
		return packageVersion, err
	}
	if strings.Contains(string(boostdLotus), "/usr/local/sbin") {
		return packageVersion, errors.New("failed to obtain boostd version")
	}
	if len(string(boostdLotus)) > 1 {
		packageVersion.BoostdVersion = string(boostdLotus)[:len(string(boostdLotus))-1]
	}

	return packageVersion, nil
}

// GetLogInformationStatus Get log message status
func GetLogInformationStatus() (bool, error) {

	out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathServerLogInformation+" %s", define.MainDisk)).Output()
	if err != nil {
		return false, err
	}
	if string(out)[:len(string(out))-1] == config.PatrolTestTrue {
		return true, nil
	} else {
		return false, nil
	}
}

// GetBlockLogStatus Status of the block output log
func GetBlockLogStatus() (int, error) {

	script := fmt.Sprintf(`cat %s/ipfs/log*/*miner.log | grep took | grep "mined new block"  | wc -l`, define.MainDisk)
	out, err := exec.Command("bash", "-c", script).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `/usr/local/sbin/boostd -v | awk '{print $3}'` Filed", err.Error())
		return 0, err
	}
	blockNum, _ := strconv.Atoi(string(out)[:len(string(out))-1])
	return blockNum, nil
}

// GetDataCatalogStatus Gets the data directory status
func GetDataCatalogStatus() (bool, error) {
	// 获取数据目录状态
	out, err := exec.Command("bash", "-c", define.PathServerDataCatalog).Output()
	if err != nil {
		return false, err
	}
	if string(out)[:len(string(out))-1] == config.PatrolTestTrue {
		return true, nil
	} else {
		return false, nil
	}
}

// GetTimeSyncStatus Get the time synchronization status
func GetTimeSyncStatus() (bool, error) {

	out, err := exec.Command("bash", "-c", define.PathServerTimeSync).Output()
	if err != nil {
		return false, err
	}
	if string(out)[:len(string(out))-1] == config.PatrolTestTrue {
		return true, nil
	} else {
		return false, nil
	}
}

// GetHostDownStatus Get whether the system is down
func GetHostDownStatus() (bool, error) {

	out, err := exec.Command("bash", "-c", define.PathServerHostDown).Output()
	if err != nil {
		return false, err
	}
	if string(out)[:len(string(out))-1] == config.PatrolTestTrue {
		return true, nil
	} else {
		return false, nil
	}
}

// GetHostPingStatus Check whether the ping function of the host is normal
func GetHostPingStatus(ip string) bool {

	out, err := exec.Command("bash", "-c", `ping `+ip+` -c3 | grep '3 received'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `ping www.baidu.com -c3 | grep '3 received'` Filed", err.Error())
		return false
	}
	return len(string(out)) > 0
}

// GetHostBadSectorInfo Gets the corresponding miner id and the sector information for the error
func GetHostBadSectorInfo() (string, []string, error) {

	out, err := exec.Command("bash", "-c", define.PathIpfsProgram+`lotus-miner proving faults | awk '{print $NF }' | egrep -v '^sectors'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `lotus-miner proving faults | awk '{print $NF }' | egrep -v '^sectors'` Filed", err.Error())
		return "", []string{}, err
	}

	if !strings.Contains(string(out), "ERROR:") || len(string(out)) == 0 {
		log.Println("cmd `lotus-miner proving faults | awk '{print $NF }' | egrep -v '^sectors'` Filed")
		return "", []string{}, errors.New("failed to obtain miner information")
	}

	infos := strings.Split(string(out)[:len(string(out))-1], "\n")
	if len(infos) == 1 {
		return infos[0], []string{}, nil
	}
	return infos[0], infos[1:], nil
}

// GetHostBadSectorType Gets the classification of the corresponding error sector
func GetHostBadSectorType(sectorId string) (string, error) {

	out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathServerSectorType+" %s %s", sectorId, define.MainDisk)).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("failed to obtain the classification of incorrect sector", err.Error())
		return "", err
	}

	if strings.Contains(string(out), "ERROR:") || len(string(out)) == 0 {
		return config.SectorTypeElseType, nil
	}

	return string(out)[:len(string(out))-1], nil
}

// GetHostBadSectorSize Gets the size of the corresponding error sector
func GetHostBadSectorSize() (string, error) {

	out, err := exec.Command("bash", "-c", define.PathIpfsProgram+`lotus-miner info | grep Miner | grep GiB | awk '{print $3 }'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `lotus-miner info | grep Miner | grep GiB | awk '{print $3 }'` Filed", err.Error())
		return "", err
	}

	if strings.Contains(string(out), "ERROR:") {
		return "", nil
	}

	return string(out)[1 : len(string(out))-1], nil
}

// GetHostBadSectorAddress Gets the address of the corresponding error sector
func GetHostBadSectorAddress(minerId, sectorId string) (string, error) {

	out, err := exec.Command("bash", "-c", fmt.Sprintf(define.PathServerSectorAddress+" %s %s %s",
		minerId, sectorId, define.MainDisk)).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("failed to obtain address for incorrect sector", err.Error())
		return "", err
	}

	if strings.Contains(string(out), "ERROR:") || len(string(out)) == 0 {
		return "", nil
	}

	strs := strings.Split(string(out)[:len(string(out))-1], "\n")
	var address string
	for key, val := range strs {
		if key == 0 {
			address = val
		}
		address = address + "," + val
	}
	return address, nil
}
