package define

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

type ProgramName string

const (
	ProgramLotus         ProgramName = "lotus"
	ProgramMiner         ProgramName = "lotus-miner"
	ProgramWorkerTask    ProgramName = "lotus-worker"
	ProgramWorkerP2Task  ProgramName = "worker-p2"
	ProgramWorkerStorage ProgramName = "lotus-storage"
	ProgramSlotUnsealed  ProgramName = "slot-unsealed"
	ProgramBoost         ProgramName = "boost"
	ProgramBoostd        ProgramName = "boostd"
	SpeedUpFile32G       ProgramName = "s-32"
	SpeedUpFile64G       ProgramName = "s-64"
)

func (p ProgramName) String() string {
	return string(p)
}

type CacheFileName string

var CacheFile CacheFileName

func (cache CacheFileName) Layers() []string {
	var layers = make([]string, 11)
	for i := 0; i < 11; i++ {
		layers[i] = fmt.Sprintf("sc-02-data-layer-%d.dat", i+1)
	}
	return layers
}

func (cache CacheFileName) ProgressP1(spath string, SectorSize uint64) int {
	layers := cache.Layers()
	var progress = 0
	for _, layerName := range layers {
		lpath := path.Join(spath, layerName)
		file, err := os.Stat(lpath)
		if err != nil {
			//log.Println("treeName:", err)
			break
		}
		
		if uint64(file.Size()) >= SectorSize {
			progress += 100 / len(layers)
			continue
		}
		p, _ := strconv.Atoi(fmt.Sprintf("%.0f", float64(file.Size())*100/float64(SectorSize)))
		progress += p
		break
	}
	return progress
}

func (cache CacheFileName) String() string {
	return string(cache)
}

func (cache CacheFileName) LayerSize32G() uint64 {
	return Ss32GiB
}

func (cache CacheFileName) LayerSize64G() uint64 {
	return Ss64GiB
}

func (cache CacheFileName) TreeCSize() uint64 {
	return 4908534048
}

func (cache CacheFileName) TreeRLastSize() uint64 {
	return 9586976
}

func (cache CacheFileName) Trees32G() []string {
	var trees = make([]string, 16)
	for i := 0; i < 8; i++ {
		trees[i] = fmt.Sprintf("sc-02-data-tree-c-%d.dat", i)
	}
	j := 0
	for i := 8; i < 16; i++ {
		trees[i] = fmt.Sprintf("sc-02-data-tree-r-last-%d.dat", j)
		j++
	}
	return trees
}

func (cache CacheFileName) Trees64G() []string {
	var trees = make([]string, 32)
	for i := 0; i < 16; i++ {
		trees[i] = fmt.Sprintf("sc-02-data-tree-c-%d.dat", i)
	}
	j := 0
	for i := 16; i < 32; i++ {
		trees[i] = fmt.Sprintf("sc-02-data-tree-r-last-%d.dat", j)
		j++
	}
	return trees
}

const (
	//P1
	Layer1  CacheFileName = "sc-02-data-layer-1.dat"
	Layer2  CacheFileName = "sc-02-data-layer-2.dat"
	Layer3  CacheFileName = "sc-02-data-layer-3.dat"
	Layer4  CacheFileName = "sc-02-data-layer-4.dat"
	Layer5  CacheFileName = "sc-02-data-layer-5.dat"
	Layer6  CacheFileName = "sc-02-data-layer-6.dat"
	Layer7  CacheFileName = "sc-02-data-layer-7.dat"
	Layer8  CacheFileName = "sc-02-data-layer-8.dat"
	Layer9  CacheFileName = "sc-02-data-layer-9.dat"
	Layer10 CacheFileName = "sc-02-data-layer-10.dat"
	Layer11 CacheFileName = "sc-02-data-layer-11.dat"
	//P2 am
	TreeC0 CacheFileName = "sc-02-data-tree-c-0.dat"
	TreeC1 CacheFileName = "sc-02-data-tree-c-1.dat"
	TreeC2 CacheFileName = "sc-02-data-tree-c-2.dat"
	TreeC3 CacheFileName = "sc-02-data-tree-c-3.dat"
	TreeC4 CacheFileName = "sc-02-data-tree-c-4.dat"
	TreeC5 CacheFileName = "sc-02-data-tree-c-5.dat"
	TreeC6 CacheFileName = "sc-02-data-tree-c-6.dat"
	TreeC7 CacheFileName = "sc-02-data-tree-c-7.dat"
	//64G
	TreeC8  CacheFileName = "sc-02-data-tree-c-8.dat"
	TreeC9  CacheFileName = "sc-02-data-tree-c-9.dat"
	TreeC10 CacheFileName = "sc-02-data-tree-c-10.dat"
	TreeC11 CacheFileName = "sc-02-data-tree-c-11.dat"
	TreeC12 CacheFileName = "sc-02-data-tree-c-12.dat"
	TreeC13 CacheFileName = "sc-02-data-tree-c-13.dat"
	TreeC14 CacheFileName = "sc-02-data-tree-c-14.dat"
	TreeC15 CacheFileName = "sc-02-data-tree-c-15.dat"
	//P2 pm
	TreeRLast0 CacheFileName = "sc-02-data-tree-r-last-0.dat"
	TreeRLast1 CacheFileName = "sc-02-data-tree-r-last-1.dat"
	TreeRLast2 CacheFileName = "sc-02-data-tree-r-last-2.dat"
	TreeRLast3 CacheFileName = "sc-02-data-tree-r-last-3.dat"
	TreeRLast4 CacheFileName = "sc-02-data-tree-r-last-4.dat"
	TreeRLast5 CacheFileName = "sc-02-data-tree-r-last-5.dat"
	TreeRLast6 CacheFileName = "sc-02-data-tree-r-last-6.dat"
	TreeRLast7 CacheFileName = "sc-02-data-tree-r-last-7.dat"
	//64G
	TreeRLast8  CacheFileName = "sc-02-data-tree-r-last-8.dat"
	TreeRLast9  CacheFileName = "sc-02-data-tree-r-last-9.dat"
	TreeRLast10 CacheFileName = "sc-02-data-tree-r-last-10.dat"
	TreeRLast11 CacheFileName = "sc-02-data-tree-r-last-11.dat"
	TreeRLast12 CacheFileName = "sc-02-data-tree-r-last-12.dat"
	TreeRLast13 CacheFileName = "sc-02-data-tree-r-last-13.dat"
	TreeRLast14 CacheFileName = "sc-02-data-tree-r-last-14.dat"
	TreeRLast15 CacheFileName = "sc-02-data-tree-r-last-15.dat"
)
