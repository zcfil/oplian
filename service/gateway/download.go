package gateway

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"oplian/define"
	"oplian/global"
	"oplian/model/common/request"
	model "oplian/model/gateway"
	"oplian/model/system"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type DownloadService struct{}

var Downloading = make(map[string]*pb.DownloadInfo)

//@author: nathan
//@function: DowloadFile
//@description: Download file
//@param: url,path string
//@return: DownloadInfo, error

func (w DownloadService) DowloadFile(url, path, fileName string) (*pb.DownloadInfo, error) {
	// Get the data
	resp, err := http.Get(url) //58G
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("DownloadFile file not found")
	}
	err = utils.AddDir(path)
	if err != nil {
		return nil, err
	}
	if len(fileName) == 0 {
		fileName = utils.SubStr(resp.Request.URL.Path, strings.LastIndex(resp.Request.URL.Path, `/`)+1, len(resp.Request.URL.Path))
	}
	counter := &pb.DownloadInfo{
		Total:    uint64(resp.ContentLength),
		FilePath: filepath.Join(path, fileName),
		Url:      url,
	}
	out, err := os.Create(counter.FilePath)
	if err != nil {
		return nil, err
	}

	record := model.DownloadRecord{ZC_MODEL: global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, Url: url,
		FilePath: counter.FilePath, FileSize: counter.Total}
	if err = global.ZC_DB.Save(&record).Error; err != nil {
		return nil, err
	}

	go func() {
		defer resp.Body.Close()
		//io.TeeReader(resp.Body, counter)

		defer out.Close()

		defer func() {
			record.FileStatus = 1
			record.LoadSize = record.FileSize
			if err != nil {
				record.LoadSize = 0
				record.ErrorMsg = err.Error()
				record.FileStatus = 2
			}
			global.ZC_DB.Save(&record)
		}()

		_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
		if err != nil {
			log.Println("下载失败：", err.Error())
			return
		}

		delete(Downloading, resp.Request.URL.Path)

		if fileName != "" {
			global.ZC_DB.Model(&system.SysFileManage{}).Where("file_status=1 and file_name=?", fileName).Update("file_status", define.FileFinish.Int())
		}

	}()
	Downloading[resp.Request.URL.Path] = counter
	return counter, nil
}

//@author: nathan
//@function: DowloadList
//@description: Paging for data
//@param: info request.PageInfo
//@return: list interface{}, total int64, err error

func (w DownloadService) DowloadList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&model.DownloadRecord{})
	var downList []model.DownloadRecord
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&downList).Error

	for _, v := range downList {
		if !utils.FileExist(v.FilePath) && v.FileStatus == 1 {
			v.FileStatus = 3
			if err = global.ZC_DB.Save(&v).Error; err != nil {
				log.Println(err.Error())
			}
		}
		if d, ok := Downloading[v.FilePath]; ok {
			v.LoadSize = d.Load
		}
	}
	return downList, total, err
}

// DownloadAddressFile Download the IP address information file
func (w DownloadService) DownloadAddressFile(fileList []string, fileAddress string, loadAddress string) {

	for _, val := range fileList {
		path := fileAddress + val
		info, err := os.Stat(path)
		if err != nil && info == nil {

			url := loadAddress + val
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("1", err.Error())
				return
			}
			defer resp.Body.Close()

			file, err := os.Create(path)
			if err != nil {
				fmt.Println("2", err.Error())
				return
			}
			defer file.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				fmt.Println("3", err.Error())
				return
			}
		}
	}
}

func (w DownloadService) SynFileToGateWay() {

	time.Sleep(time.Second * 5)
	type DownLoadFile struct {
		FileName string
		FilePath string
		Version  string
	}

	fileList := make([]DownLoadFile, 0)
	fileList = append(fileList, DownLoadFile{FileName: define.IPAddressFile, FilePath: define.PathIpfsConfig})
	version := make([]string, 0)
	version = append(version, "20", "22")
	for _, v := range version {
		fileList = append(fileList, DownLoadFile{FileName: define.ProgramLotus.String(), Version: v, FilePath: define.PathIpfsProgram})
		fileList = append(fileList, DownLoadFile{FileName: define.ProgramMiner.String(), Version: v, FilePath: define.PathIpfsProgram})
		fileList = append(fileList, DownLoadFile{FileName: define.ProgramWorkerTask.String(), Version: v, FilePath: define.PathIpfsProgram})
		fileList = append(fileList, DownLoadFile{FileName: define.ProgramWorkerP2Task.String(), Version: v, FilePath: define.PathIpfsProgram})
		fileList = append(fileList, DownLoadFile{FileName: define.OPlianFile, Version: v, FilePath: define.PathOplian})
		fileList = append(fileList, DownLoadFile{FileName: define.OPlianOpFile, Version: v, FilePath: define.PathOplian})
		fileList = append(fileList, DownLoadFile{FileName: define.OPlianOpC2File, Version: v, FilePath: define.PathOplian})
		fileList = append(fileList, DownLoadFile{FileName: define.OplianSectorsSealFile, Version: v, FilePath: define.PathOplian})
	}

	var wg sync.WaitGroup
	wg.Add(len(fileList))
	for _, v := range fileList {

		go func(file DownLoadFile) {

			defer func() {
				wg.Done()
			}()

			name := fmt.Sprintf("%s-%s", file.FileName, file.Version)
			if _, err := os.Stat(path.Join(define.FileGateWayDir, name)); err != nil {

				var roomRecord system.SysMachineRoomRecord
				if err := global.ZC_DB.Model(&system.SysMachineRoomRecord{}).Where("gateway_id", global.GateWayID.String()).Find(&roomRecord).Error; err != nil {
					log.Println("SysMachineRoomRecord query data err:", err)
					return
				}

				fileUrl := fmt.Sprintf("%s%s-%s", define.DownLoadAddressOut, file.FileName, file.Version)
				log.Println(fmt.Sprintf("SynFileToGateWay 同步文件 fileUrl:%s,目录:%s", fileUrl, define.FileGateWayDir))
				res, err := w.DowloadFile(fileUrl, define.FileGateWayDir, "")
				if err != nil {
					log.Println("DownloadFile err:", err.Error())
					return
				}

				r := &system.SysFileManage{
					RoomId:         roomRecord.RoomId,
					RoomName:       roomRecord.RoomName,
					GateWayId:      global.GateWayID.String(),
					FileType:       1,
					FileUrl:        fileUrl,
					FileStatus:     utils.ONE,
					ComputerSystem: fmt.Sprintf("Ubuntu-%s", file.Version),
					Version:        file.Version,
				}

				if res != nil {
					beginIndex := strings.LastIndex(res.FilePath, `/`)
					if beginIndex == -1 {
						beginIndex = strings.LastIndex(res.FilePath, `\`)
					}
					fileName := utils.SubStr(res.FilePath, beginIndex+1, len(res.FilePath))
					r.FileName = fileName
					r.FileSize = int(res.Total)
				} else {
					r.FileName = "未知文件"
					r.FileSize = utils.ZERO
					r.FileStatus = define.FileError.Int()
				}

				err = global.ZC_DB.Model(&system.SysFileManage{}).Create(r).Error
				if err != nil {
					log.Println("SysFileManage add data err:", err)
					return
				}

			}

			computerSystem := ""
			version := ""
			b, err := utils.ExecuteScript("hostnamectl | grep 'Operating System:'|awk -F ':' '{print $2}'")
			if err != nil {
				log.Println("ExecuteScript err:", err)
			}

			bOut := utils.Replace(b)
			if strings.Contains(bOut, "Ubuntu") {
				if strings.Contains(bOut, file.Version) {
					computerSystem = bOut
					version = file.Version
				}
			}

			if _, err := os.Stat(path.Join(file.FilePath, file.FileName)); err != nil {

				if version != "" && computerSystem != "" {

					go func(f DownLoadFile) {

						var sysFile system.SysFileManage
						for {

							if err = global.ZC_DB.Model(&system.SysFileManage{}).Where("file_name", f.FileName).Find(&sysFile).Error; err != nil {
								return
							}

							if (sysFile != system.SysFileManage{}) {
								if sysFile.FileStatus != 2 {
									time.Sleep(time.Second * 5)
									continue
								}
							}

							break
						}

						f.FileName = fmt.Sprintf("%s-%s", f.FileName, version)
						openPath := path.Join(define.FileGateWayDir, f.FileName)
						log.Println("SynFileToGateWay Open file:", openPath)
						openFile, err := os.Open(openPath)
						if err != nil {
							log.Println("Open file err:", err)
							return
						}
						defer func(openFile *os.File) {
							err := openFile.Close()
							if err != nil {
								log.Println("openFile Close err:", err)
							}
						}(openFile)

						createPath := path.Join(f.FilePath, file.FileName)
						log.Println("SynFileToGateWay Create file:", createPath)
						file, err := os.Create(path.Join(f.FilePath, file.FileName))
						if err != nil {
							log.Println("Create file err:", err)
							return
						}
						defer func(file *os.File) {
							err := file.Close()
							if err != nil {
								log.Println("file Close err:", err)
								return
							}
						}(file)

						_, err = io.Copy(file, openFile)
						if err != nil {
							log.Println("Copy file err:", err)
							return
						} else {
							os.Chmod(createPath, 777)
						}

					}(file)
				}
			}

		}(v)
	}
	wg.Wait()
	log.Println("The installation package is successfully downloaded")

}
