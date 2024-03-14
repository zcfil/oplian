package opinit

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"oplian/define"
	"oplian/service/lotus/oplocal"
	"oplian/utils"
	"os"
	"path/filepath"
	"time"
)

func RangePathSectors() {
	buf, err := os.ReadFile(filepath.Join(define.PathIpfsWorker, define.StorageConfig))
	if err != nil {
		log.Println(filepath.Join(define.PathIpfsWorker, define.StorageConfig, err.Error()))
		return
	}
	type LocalPath struct {
		Path string
	}
	type StorageConfig struct {
		StoragePaths []LocalPath
	}

	var storages StorageConfig
	if err = json.Unmarshal(buf, &storages); err != nil {
		log.Println(filepath.Join(define.PathIpfsWorker, define.StorageConfig, err.Error()))
		return
	}
	var size int64
	size = define.Ss32GiB
	var minerSector = make(map[string]string)
	//actor := utils.MinerActorID(miner)

	var mp = make(map[string]struct{})
	var count int
	for _, path := range storages.StoragePaths {
		if _, ok := mp[path.Path]; ok {
			continue
		}
		mp[path.Path] = struct{}{}
		var meta oplocal.LocalStorageMeta
		cfgBuf, err := os.ReadFile(filepath.Join(path.Path, define.SectorStoreConfig))
		if err != nil {
			continue
		}
		if err = json.Unmarshal(cfgBuf, &meta); err != nil {
			continue
		}
		if !meta.CanSeal {
			continue
		}
		unSealedDir := filepath.Join(path.Path, define.FTUnsealed.String())
		if utils.ExistFileOrDir(unSealedDir) {
			filepath.Walk(unSealedDir, func(p string, info fs.FileInfo, err error) error {
				if info == nil {
					return err
				}

				if utils.CheckSectorNum(info.Name()) {
					if size < info.Size() {
						size = info.Size()
					}
					_, number := utils.StringToSectorID(info.Name())
					//if actor == a || actor == 0 {
					minerSector[info.Name()] = unSealedDir
					oplocal.IsDiskCount[number] = struct{}{}
					oplocal.SealSectorPath.Push(info.Name(), path.Path)
					//}
					if utils.ExistFileOrDir(filepath.Join(path.Path, define.FTSealed.String(), info.Name())) {
						count++
						oplocal.PathSealCount.Push(path.Path, info.Name(), false)
					}
				}
				return nil
			})
		}
	}
	for _, path := range storages.StoragePaths {
		var meta oplocal.LocalStorageMeta
		cfgBuf, err := os.ReadFile(filepath.Join(path.Path, define.SectorStoreConfig))
		if err != nil {
			continue
		}
		if err = json.Unmarshal(cfgBuf, &meta); err != nil {
			continue
		}
		if !meta.CanSeal {
			continue
		}
		oplocal.DiskCount += utils.DiskSpaceSufficientCount(path.Path, uint64(size), 1)
		if oplocal.DiskCount > 1024 {
			log.Fatal("DiskSpaceSufficientCount error:", oplocal.DiskCount)
		}
	}
	runCount := oplocal.Tasking.GetRunCount(define.SealPreCommit1.String())
	if runCount == 0 {
		time.Sleep(time.Second * 20)
		runCount = oplocal.Tasking.GetRunCount(define.SealPreCommit1.String())
	}

	if oplocal.DiskCount > uint64(runCount) {
		oplocal.DiskCount -= uint64(runCount)
	} else {
		oplocal.DiskCount = 0
	}

	if count > 0 {
		count--
		oplocal.DiskCount += uint64(count)
	}
	log.Println(fmt.Sprintf("RangeSectorsPath :%+v\n%v", oplocal.PathSealCount, oplocal.SealSectorPath), oplocal.DiskCount)
	return
}
