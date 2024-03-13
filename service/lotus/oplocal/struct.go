package oplocal

import (
	"oplian/define"
	"oplian/service/pb"
	"sync"
	"time"
)

type Preselection struct {
	Pre1      int
	Pre2      int
	Pre       int
	LastTime  time.Time
	LastTime2 time.Time
	Lock      sync.Mutex
}

func (acc *Preselection) Ok(taskType define.TaskType) bool {
	switch taskType {
	case define.TTAddPiece, define.TTPreCommit1:
		//pre := 1
		//if taskType == define.TTAddPiece {
		//	pre = APCount
		//}
		if acc.LastTime.Add(time.Minute).After(time.Now()) && acc.Pre1 >= 1 {
			return false
		}
		if acc.LastTime.Add(time.Minute).Before(time.Now()) {
			acc.Pre1 = 0
		}
		acc.Pre1++
	case define.TTPreCommit2:
		if acc.LastTime.Add(time.Minute).After(time.Now()) && acc.Pre2 >= 1 {
			return false
		}
		if acc.LastTime.Add(time.Minute).Before(time.Now()) {
			acc.Pre2 = 0
		}
		acc.Pre2++
	}
	acc.LastTime = time.Now()
	return true
}

func (acc *Preselection) OkNew() (bool, bool) {
	var p1, p2 = true, true
	if acc.LastTime.Add(time.Minute * 2).After(time.Now()) {
		p1 = false
	} else {
		acc.LastTime = time.Now()
	}
	if acc.LastTime2.Add(time.Second * 10).After(time.Now()) {
		p2 = false
	} else {
		acc.LastTime2 = time.Now()
	}
	return p1, p2
}

type LocalWorkerInfo struct {
	info   map[string]*pb.OpenWindow
	lockRW sync.RWMutex
}

func (lw *LocalWorkerInfo) GetOpenWindow(wid string) *pb.OpenWindow {
	lw.lockRW.RLock()
	defer lw.lockRW.RUnlock()
	return OpWorkers.info[wid]
}

func (lw *LocalWorkerInfo) Push(args *pb.OpenWindow) {
	lw.lockRW.Lock()
	defer lw.lockRW.Unlock()
	OpWorkers.info[args.WorkerId] = args
	return
}

type LocalStorageMeta struct {
	ID string

	// A high weight means data is more likely to be stored in this path
	Weight uint64 // 0 = readonly

	// Intermediate data for the sealing process will be stored here
	CanSeal bool

	// Finalized sectors that will be proved over time will be stored here
	CanStore bool
	DirPath  string
}

type SealSectorDisk map[string]string

var sdLock sync.RWMutex

func (s SealSectorDisk) Get(sector string) string {
	sdLock.RLock()
	defer sdLock.RUnlock()
	return s[sector]
}
func (s SealSectorDisk) Push(sector string, path string) {
	sdLock.Lock()
	defer sdLock.Unlock()
	s[sector] = path
}
func (s SealSectorDisk) Remove(sector string) {
	sdLock.Lock()
	defer sdLock.Unlock()
	delete(s, sector)
}


type DiskSealCount map[string]map[string]bool

var dscLock sync.RWMutex

func (d DiskSealCount) Get(path string) int {
	dscLock.RLock()
	defer dscLock.RUnlock()
	return len(d[path])
}

func (d DiskSealCount) GetRun(path string, sector string) int {
	dscLock.RLock()
	defer dscLock.RUnlock()
	count := 0
	for sect, run := range d[path] {

		if sector == sect {
			continue
		}
		if run {
			count++
		}
	}
	return count
}

func (d DiskSealCount) GetTotal() int {
	dscLock.RLock()
	defer dscLock.RUnlock()
	var total int
	for _, disk := range d {
		total += len(disk)
	}
	return total
}

func (d DiskSealCount) Push(path string, sector string, running bool) {
	dscLock.Lock()
	defer dscLock.Unlock()
	if _, ok := d[path]; !ok {
		d[path] = make(map[string]bool)
	}
	d[path][sector] = running
}

func (d DiskSealCount) Sub(path string, sector string) {
	dscLock.Lock()
	defer dscLock.Unlock()
	delete(d[path], sector)
}
