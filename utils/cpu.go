package utils

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// CPU information.
type CPU struct {
	Vendor      string  `json:"vendor,omitempty"`
	Model       string  `json:"model,omitempty"`
	Speed       string  `json:"speed,omitempty"`       // CPU clock rate in GHz
	Cache       uint    `json:"cache,omitempty"`       // CPU cache size in KB
	Cpus        uint    `json:"cpus,omitempty"`        // number of physical CPUs
	Cores       uint    `json:"cores,omitempty"`       // number of physical CPU cores
	Threads     uint    `json:"threads,omitempty"`     // number of logical (HT) CPU cores
	UsedPercent float64 `json:"usedPercent,omitempty"` // 使用率
}

var (
	reTwoColumns = regexp.MustCompile("\t+: ")
	reExtraSpace = regexp.MustCompile(" +")
	reCacheSize  = regexp.MustCompile(`^(\d+) KB$`)
)

func (si *OpServer) GetCPUInfo() {
	si.CPU.Threads = uint(runtime.NumCPU())

	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return
	}
	defer f.Close()

	cpu := make(map[string]bool)
	core := make(map[string]bool)

	var cpuID string

	s := bufio.NewScanner(f)
	for s.Scan() {
		if sl := reTwoColumns.Split(s.Text(), 2); sl != nil {
			switch sl[0] {
			case "physical id":
				cpuID = sl[1]
				cpu[cpuID] = true
			case "core id":
				coreID := fmt.Sprintf("%s/%s", cpuID, sl[1])
				core[coreID] = true
			case "vendor_id":
				if si.CPU.Vendor == "" {
					si.CPU.Vendor = sl[1]
				}
			case "model name":
				if si.CPU.Model == "" {
					// CPU model, as reported by /proc/cpuinfo, can be a bit ugly. Clean up...
					model := reExtraSpace.ReplaceAllLiteralString(sl[1], " ")
					si.CPU.Model = strings.Replace(model, "- ", "-", 1)
				}
			case "cache size":
				if si.CPU.Cache == 0 {
					if m := reCacheSize.FindStringSubmatch(sl[1]); m != nil {
						if cache, err := strconv.ParseUint(m[1], 10, 64); err == nil {
							si.CPU.Cache = uint(cache)
						}
					}
				}
			}
		}
	}
	if s.Err() != nil {
		return
	}

	var outFloat float64
	var usedPercent float64

	out, err := exec.Command("bash", "-c", `lscpu | grep 'CPU MHz:' | awk '{printf ("%.2f",$3/1024)}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `lscpu | grep 'CPU MHz:' | awk '{printf (\"%.2f\",$3/1024)}'` Filed", err.Error())
	} else {
		outFloat, _ = strconv.ParseFloat(string(out), 64)
	}

	usedPercentOut, err := exec.Command("bash", "-c", `top -bn1 | fgrep "Cpu(s)" | awk '{printf ("%.2f",100-$8)}'`).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("cmd `top -bn1 | fgrep \"Cpu(s)\" | awk '{print 100-$8}'` Filed", err.Error())
		usedPercent = 1
	} else {
		usedPercent, _ = strconv.ParseFloat(string(usedPercentOut), 64)
		if usedPercent < 1 {
			usedPercent = 1
		}
	}

	si.CPU.Speed = strconv.FormatFloat(math.Floor(outFloat*10)/10, 'f', -1, 64)
	si.CPU.UsedPercent = usedPercent
	si.CPU.Cpus = uint(len(cpu))
	si.CPU.Cores = uint(len(core))
}
