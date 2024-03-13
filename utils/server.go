package utils

import (
	"math"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type OpServer struct {
	Os         Os              `json:"os"`
	CPU        CPU             `json:"cpu"`
	GPU        GPU             `json:"gpu"`
	Ram        Ram             `json:"ram"`
	Disk       Disk            `json:"disk"`
	UUID       string          `json:"uuid"`
	SystemInfo LinuxSystemInfo `json:"systemInfo"`
	Storage    []StorageDevice `json:"storage"`
}

type Server struct {
	Os   Os   `json:"os"`
	Cpu  Cpu  `json:"cpu"`
	Ram  Ram  `json:"ram"`
	Disk Disk `json:"disk"`
}

type Os struct {
	GOOS         string `json:"goos"`
	NumCPU       int    `json:"numCpu"`
	Compiler     string `json:"compiler"`
	GoVersion    string `json:"goVersion"`
	NumGoroutine int    `json:"numGoroutine"`
}

type Cpu struct {
	Cpus  []float64 `json:"cpus"`
	Cores int       `json:"cores"`
}

type Ram struct {
	UsedMB      int `json:"usedMb"`
	TotalMB     int `json:"totalMb"`
	UsedPercent int `json:"usedPercent"`
}

type Disk struct {
	UsedMB      int `json:"usedMb"`
	UsedGB      int `json:"usedGb"`
	TotalMB     int `json:"totalMb"`
	TotalGB     int `json:"totalGb"`
	UsedPercent int `json:"usedPercent"`
}

type DiskSum struct {
	Name    string `json:"name"`
	Size    uint   `json:"sizeSum"`
	SizeSum uint   `json:"num"`
}

type OpCPUInfo struct {
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	CpuNum      string `json:"cpus"`    // number of physical CPUs
	Speed       string `json:"speed"`   // CPU clock rate in GHz
	Threads     string `json:"threads"` // number of logical (HT) CPU cores
	UsedPercent string `json:"usedPercent"`
}

type OpGPUInfo struct {
	Mark        string `json:"mark"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	TotalMB     string `json:"totalMb"`
	UseRate     string `json:"useRate"`
	IsQualified string `json:"isQualified,omitempty"`
}

type OpRamInfo struct {
	TotalMB     string `json:"totalMb"`
	Type        string `json:"Type"`
	RamNum      string `json:"ramNum"`
	UsedMB      string `json:"usedMb"`
	UsedPercent string `json:"usedPercent"`
}

type OpDiskInfo struct {
	IsOs        int    `json:"isOs"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	TotalSize   string `json:"totalSize"`
	UsedSize    string `json:"usedSize"`
	UsedPercent string `json:"usedPercent"`
}

type OpDiskMd0Info struct {
	DiskNum     string `json:"diskNum"`
	TotalSize   string `json:"totalSize"`
	UsedSize    string `json:"usedSize"`
	UsedPercent string `json:"usedPercent"`
}

type OpStorageInfo struct {
	OpCPUInfo     OpCPUInfo     `json:"opCPUInfo"`
	OpGPUInfo     []OpGPUInfo   `json:"opGPUInfo"`
	OpRamInfo     OpRamInfo     `json:"opRamInfo"`
	OpDiskInfo    []OpDiskInfo  `json:"opDiskInfo"`
	OpDiskMd0Info OpDiskMd0Info `json:"opDiskMd0Info"`
}

type HostDisk struct {
	DiskSize      string        `json:"diskSize"`
	OpDiskMd0Info OpDiskMd0Info `json:"opDiskMd0Info"`
}

//@function: InitCPU
//@description: OS information
//@return: o Os, err error

func InitOS() (o Os) {
	o.GOOS = runtime.GOOS
	o.NumCPU = runtime.NumCPU()
	o.Compiler = runtime.Compiler
	o.GoVersion = runtime.Version()
	o.NumGoroutine = runtime.NumGoroutine()
	return o
}

//@function: InitCPU
//@description: CPU information
//@return: c CPU, err error

func InitCPU() (c Cpu, err error) {
	if cores, err := cpu.Counts(false); err != nil {
		return c, err
	} else {
		c.Cores = cores
	}
	if cpus, err := cpu.Percent(time.Duration(200)*time.Millisecond, true); err != nil {
		return c, err
	} else {
		c.Cpus = cpus
	}
	return c, nil
}

//@function: InitRAM
//@description: RAM information
//@return: r Ram, err error

func InitRAM() (r Ram, err error) {
	if u, err := mem.VirtualMemory(); err != nil {
		return r, err
	} else {
		r.UsedMB = int(u.Used) / MB
		r.TotalMB = int(u.Total) / MB
		r.UsedPercent = int(math.Ceil(u.UsedPercent))
	}
	return r, nil
}

//@function: InitDisk
//@description: Hard disk information
//@return: d Disk, err error

func InitDisk() (d Disk, err error) {
	if u, err := disk.Usage("/"); err != nil {
		return d, err
	} else {
		d.UsedMB = int(u.Used) / MB
		d.UsedGB = int(u.Used) / GB
		d.TotalMB = int(u.Total) / MB
		d.TotalGB = int(u.Total) / GB
		d.UsedPercent = int(math.Ceil(u.UsedPercent))
	}
	return d, nil
}
