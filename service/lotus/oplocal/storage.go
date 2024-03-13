package oplocal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/fs"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type StorageService struct{}

// FindSectorFile Find a sector file
func (deploy StorageService) FindSectorFile(workerPath, miner string, sectorNumber uint64, fileType define.SectorFileType) (string, LocalStorageMeta, error) {
	buf, err := os.ReadFile(filepath.Join(workerPath, define.StorageConfig))
	if err != nil {
		return "", LocalStorageMeta{}, err
	}
	type LocalPath struct {
		Path string
	}
	type StorageConfig struct {
		StoragePaths []LocalPath
	}

	var storages StorageConfig
	if err = json.Unmarshal(buf, &storages); err != nil {
		return "", LocalStorageMeta{}, err
	}
	for _, path := range storages.StoragePaths {
		sectorPath := filepath.Join(path.Path, fileType.String(), utils.SectorNumString(miner, sectorNumber))
		if utils.ExistFileOrDir(sectorPath) {
			var meta LocalStorageMeta
			cfgBuf, err := os.ReadFile(filepath.Join(path.Path, define.SectorStoreConfig))
			if err != nil {
				return "", LocalStorageMeta{}, err
			}
			if err = json.Unmarshal(cfgBuf, &meta); err != nil {
				return "", LocalStorageMeta{}, err
			}

			meta.DirPath = path.Path
			return sectorPath, meta, nil
		}
	}
	return "", LocalStorageMeta{}, fmt.Errorf("sector %s not find ", utils.SectorNumString(miner, sectorNumber))
}

func (deploy StorageService) FindSectorsPath(workerPath string, sectorMap map[string]uint64, fileType define.SectorFileType) (map[uint64]string, error) {

	sectorPathMap := make(map[uint64]string)
	buf, err := os.ReadFile(filepath.Join(workerPath, define.StorageConfig))
	if err != nil {
		return nil, err
	}
	type LocalPath struct {
		Path string
	}
	type StorageConfig struct {
		StoragePaths []LocalPath
	}

	var storages StorageConfig
	if err = json.Unmarshal(buf, &storages); err != nil {
		return nil, err
	}
	for _, storagePath := range storages.StoragePaths {

		_ = filepath.Walk(storagePath.Path, func(path string, info os.FileInfo, err error) error {

			if fileType == define.FTCache && info.IsDir() {

				if number, ok := sectorMap[info.Name()]; ok {
					sectorPathMap[number] = path
				}
			} else if fileType == define.FTUnsealed && !info.IsDir() {
				if number, ok := sectorMap[info.Name()]; ok {
					sectorPathMap[number] = path
				}
			}

			return nil
		})
	}

	return sectorPathMap, nil
}

// FindStorage Find enough hard drive storage
func (deploy StorageService) FindStorage(workerPath string, moveCount int, storeSeal string) (string, LocalStorageMeta, error) {
	buf, err := os.ReadFile(filepath.Join(workerPath, define.StorageConfig))
	if err != nil {
		return "", LocalStorageMeta{}, err
	}
	type LocalPath struct {
		Path string
	}
	type StorageConfig struct {
		StoragePaths []LocalPath
	}

	var storages StorageConfig
	if err = json.Unmarshal(buf, &storages); err != nil {
		return "", LocalStorageMeta{}, err
	}
	if len(storages.StoragePaths) < moveCount*2 {
		return "", LocalStorageMeta{}, fmt.Errorf("This storage is written in full! Number of hard disks: %d indicates the number of disks being transferred：%d ", len(storages.StoragePaths), moveCount)
	}
	for _, path := range storages.StoragePaths {
		if utils.ExistFileOrDir(path.Path) {
			var meta LocalStorageMeta
			cfgBuf, err := os.ReadFile(filepath.Join(path.Path, define.SectorStoreConfig))
			if err != nil {
				continue
			}
			if err = json.Unmarshal(cfgBuf, &meta); err != nil {
				continue
			}
			preCount1 := 1
			if storeSeal == "seal" {
				preCount1 = 4
				if !meta.CanSeal {
					continue
				}
			}
			if storeSeal == "store" && !meta.CanStore {
				continue
			}
			if utils.DiskSpaceSufficient(path.Path, define.Ss32GiB, preCount1) {
				return path.Path, meta, nil
			}
			log.Println(path.Path, "Insufficient capacity")
		} else {
			log.Println(path.Path, "not exist")
		}
	}
	return "", LocalStorageMeta{}, fmt.Errorf("This storage is full！%s %d", workerPath, len(storages.StoragePaths))
}

// StorageAttachInit
// @author: nathan
// @function: StorageAttachInit
// @description: Initialize storage
// @param: path string, canSeal, canStore bool
// @return: error
func (p StorageService) StorageAttachInit(path string, canSeal, canStore bool) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	_, err := os.Stat(filepath.Join(path, utils.MetaFile))
	if !os.IsNotExist(err) {
		return err
	}
	cfg := &utils.SectorstoreConfig{
		ID:         uuid.New().String(),
		Weight:     10,
		CanSeal:    canSeal,
		CanStore:   canStore,
		MaxStorage: 0,
		Groups:     nil,
		AllowTo:    nil,
	}
	return cfg.WriteFile(path)
}

// FindSealStorage Find enough hard drive storage
func (deploy StorageService) FindSealStorage(workerPath string, sector *pb.SectorRef) (*pb.SectorPath, error) {
	buf, err := os.ReadFile(filepath.Join(workerPath, define.StorageConfig))
	if err != nil {
		return nil, err
	}
	type LocalPath struct {
		Path string
	}
	type StorageConfig struct {
		StoragePaths []LocalPath
	}

	var storages StorageConfig
	if err = json.Unmarshal(buf, &storages); err != nil {
		return nil, err
	}
	unsealed := ""
	Max := -1
	var sectorPath pb.SectorPath
	minerSector := utils.SectorNumString(sector.Id.Miner, sector.Id.Number)
	for _, path := range storages.StoragePaths {
		if utils.ExistFileOrDir(path.Path) {
			var meta LocalStorageMeta
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

			count := PathSealCount.GetRun(path.Path, minerSector)

			filepath.Walk(filepath.Join(path.Path, "/cache"), func(spath string, info fs.FileInfo, err error) error {
				if info == nil {
					return err
				}

				if utils.CheckSectorNum(info.Name()) {
					if define.CacheFile.ProgressP1(spath, sector.ProofType) > 40 {
						count--
					}
				}
				return nil
			})

			//if exist unsealed
			if utils.ExistFileOrDir(filepath.Join(path.Path, "/unsealed", minerSector)) {
				sectorPath = pb.SectorPath{
					Id:          sector.Id,
					StoreId:     meta.ID,
					DiskPath:    path.Path,
					Unsealed:    filepath.Join(path.Path, "/unsealed", minerSector),
					Cache:       filepath.Join(path.Path, "/cache", minerSector),
					Sealed:      filepath.Join(path.Path, "/sealed", minerSector),
					Update:      filepath.Join(path.Path, "/update", minerSector),
					UpdateCache: filepath.Join(path.Path, "/updateCache", minerSector),
				}

				if !sector.PreAfter && 0 == utils.DiskSpaceSufficientCount(path.Path, sector.ProofType, count) {
					return nil, fmt.Errorf("This storage is zcfull！%s ", path.Path)
				}
				log.Println("select1：", path.Path, sector.Id)
				return &sectorPath, nil
			}
			if sector.PreAfter {
				continue
			}

			canCount := utils.DiskSpaceSufficientCount(path.Path, sector.ProofType, count)
			if int(canCount) > Max {
				Max = int(canCount)
				unsealed = filepath.Join(path.Path, "/unsealed", minerSector)
				sectorPath = pb.SectorPath{
					Id:          sector.Id,
					DiskPath:    path.Path,
					Unsealed:    unsealed,
					Cache:       filepath.Join(path.Path, "/cache", minerSector),
					Sealed:      filepath.Join(path.Path, "/sealed", minerSector),
					Update:      filepath.Join(path.Path, "/update", minerSector),
					UpdateCache: filepath.Join(path.Path, "/updateCache", minerSector),
				}
			}
		}

	}
	if unsealed == "" {
		return nil, fmt.Errorf("This storage is full！%s ", workerPath)
	}
	log.Println("select2：", unsealed)
	return &sectorPath, nil
}

// RangeSectors Traversing all sectors in the local area
func (deploy StorageService) RangeSectors(workerPath string, miner string) []uint64 {

	buf, err := os.ReadFile(filepath.Join(workerPath, define.StorageConfig))
	if err != nil {
		log.Println(filepath.Join(workerPath, define.StorageConfig, err.Error()))
		return nil
	}
	type LocalPath struct {
		Path string
	}
	type StorageConfig struct {
		StoragePaths []LocalPath
	}

	var storages StorageConfig
	if err = json.Unmarshal(buf, &storages); err != nil {
		log.Println(filepath.Join(workerPath, define.StorageConfig, err.Error()))
		return nil
	}
	var minerSector []uint64
	actor := utils.MinerActorID(miner)
	for _, path := range storages.StoragePaths {
		cacheDir := filepath.Join(path.Path, define.FTUnsealed.String())
		if utils.ExistFileOrDir(cacheDir) {
			filepath.Walk(cacheDir, func(path string, info fs.FileInfo, err error) error {
				if info == nil {
					return err
				}

				if utils.CheckSectorNum(info.Name()) {
					a, number := utils.StringToSectorID(info.Name())
					if actor == a {
						minerSector = append(minerSector, number)
					}
				}
				return nil
			})
		}

	}

	return minerSector
}

// ClearWorker Clear all local sectors
func (deploy StorageService) ClearWorker(workerPath string) {
	buf, err := os.ReadFile(filepath.Join(workerPath, define.StorageConfig))
	defer func() {
		os.RemoveAll(workerPath)
		os.RemoveAll(define.PathIpfsDataWorkerCar)
	}()
	if err != nil {
		log.Println(filepath.Join(workerPath, define.StorageConfig, err.Error()))
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
		log.Println(filepath.Join(workerPath, define.StorageConfig, err.Error()))
		return
	}
	for _, path := range storages.StoragePaths {
		os.RemoveAll(path.Path)
	}
}

// FindSectorCount Traversing the number of all sectors in the local area
func (deploy StorageService) FindSectorCount(dir string, stype define.SectorFileType, actor string) (int, error) {
	buf, err := os.ReadFile(filepath.Join(dir, define.StorageConfig))
	if err != nil {
		return 0, err
	}
	type LocalPath struct {
		Path string
	}
	type StorageConfig struct {
		StoragePaths []LocalPath
	}
	count := 0
	var storages StorageConfig
	if err = json.Unmarshal(buf, &storages); err != nil {
		return 0, err
	}
	for _, paths := range storages.StoragePaths {
		storagePath := path.Join(paths.Path, stype.String())
		filepath.Walk(storagePath, func(path1 string, info fs.FileInfo, err error) error {
			if info == nil {
				return err
			}
			actorId, Number := utils.StringToSectorID(info.Name())
			if actorId == 0 {
				return nil
			}
			if stype == define.FTUnsealed {
				resSector, err := global.OpToGatewayClient.GetSectorStatus(context.Background(), &pb.SectorID{Miner: fmt.Sprintf("f0%d", actorId), Number: Number})
				if err != nil {
					return err
				}
				if resSector.Status == define.Removed {
					os.Remove(path.Join(storagePath, fmt.Sprintf("s-t0%d-%d", actorId, Number)))
				}

			}
			if strconv.FormatUint(actorId, 10) == actor[1:] {
				count++
			}
			return nil
		})
	}
	return count, err
}

// FindAllSector Traversing all sectors in the local area
func (deploy StorageService) FindAllSector(dir string, sectorPath []string, miner string) (map[string]int, map[uint64]string, error) {

	countMap := make(map[string]int)
	sectorPathMap := make(map[uint64]string)
	var sectorLock sync.RWMutex
	buf, err := os.ReadFile(filepath.Join(dir, define.StorageConfig))
	if err != nil {
		return nil, nil, err
	}
	type LocalPath struct {
		Path string
	}
	type StorageConfig struct {
		StoragePaths []LocalPath
	}

	var storages StorageConfig
	if err = json.Unmarshal(buf, &storages); err != nil {
		return nil, nil, err
	}

	apCount, preCount := 0, 0
	var apCountLock, preCountLock sync.RWMutex
	actor := utils.MinerActorID(miner)
	var wg sync.WaitGroup
	wg.Add(len(storages.StoragePaths))
	for _, paths1 := range storages.StoragePaths {

		go func(paths LocalPath) {

			defer func() { wg.Done() }()

			//扇区数量统计
			//startTime := time.Now()
			var wg1 sync.WaitGroup
			wg1.Add(len(sectorPath))
			for _, sectorPathStr1 := range sectorPath {

				go func(sectorPathStr string) {

					defer func() { wg1.Done() }()

					storagePath := path.Join(paths.Path, sectorPathStr)
					filepath.Walk(storagePath, func(path1 string, info fs.FileInfo, err error) error {
						if info == nil {
							return err
						}
						actorId, Number := utils.StringToSectorID(info.Name())
						if actorId == 0 {
							return nil
						}

						switch sectorPathStr {
						case define.FTUnsealed.String():
							go func() {
								if UnseledClearTime.Add(time.Minute * 5).After(time.Now()) {
									return
								}
								UnseledClearTime = time.Now()
								resSector, err := global.OpToGatewayClient.GetSectorStatus(context.Background(), &pb.SectorID{Miner: fmt.Sprintf("f0%d", actorId), Number: Number})
								if err != nil {
									return
								}
								if resSector.Status == define.Removed {
									os.Remove(path.Join(storagePath, fmt.Sprintf("s-t0%d-%d", actorId, Number)))
								}
							}()

							if strconv.FormatUint(actorId, 10) == miner[1:] {
								apCountLock.Lock()
								apCount++
								apCountLock.Unlock()
							}

							break

						case define.FTCache.String():

							if strconv.FormatUint(actorId, 10) == miner[1:] {
								preCountLock.Lock()
								preCount++
								preCountLock.Unlock()
							}
							break
						}

						return nil
					})

				}(sectorPathStr1)

			}
			wg1.Wait()
			//log.Println("FindAllSector 扇区数量统计 耗时:", time.Now().Sub(startTime))

			//扇区路径
			//startTime1 := time.Now()
			var meta LocalStorageMeta
			cfgBuf, err := os.ReadFile(filepath.Join(paths.Path, define.SectorStoreConfig))
			if err != nil {
				return
			}
			if err = json.Unmarshal(cfgBuf, &meta); err != nil {
				return
			}
			if !meta.CanSeal {
				return
			}
			unsealedDir := filepath.Join(paths.Path, define.FTUnsealed.String())
			if utils.ExistFileOrDir(unsealedDir) {
				filepath.Walk(unsealedDir, func(path string, info fs.FileInfo, err error) error {
					if info == nil {
						return err
					}

					if utils.CheckSectorNum(info.Name()) {
						a, number := utils.StringToSectorID(info.Name())
						if actor == a {
							sectorLock.Lock()
							sectorPathMap[number] = unsealedDir
							sectorLock.Unlock()
						}
					}
					return nil
				})
			}
			//log.Println("FindAllSector 扇区路径 耗时:", time.Now().Sub(startTime1))

		}(paths1)
	}

	wg.Wait()

	countMap["apCount"] = apCount
	countMap["preCount"] = preCount

	return countMap, sectorPathMap, err
}
