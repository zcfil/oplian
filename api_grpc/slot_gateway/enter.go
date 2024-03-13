package slot_gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/model/common/response"
	responseModel "oplian/model/system/response"
	"oplian/service"
	"oplian/service/gateway"
	"oplian/service/pb"
	systemModel "oplian/service/system"
	"oplian/utils"
	"strings"
	"time"
)

var (
	deployService = service.ServiceGroupApp.LotusServiceGroup.DeployService
)

type ExaFileUploadAndDownload struct {
	ID         uint   `json:"fileId"`     // 文件id
	Name       string `json:"name"`       // 文件名
	Url        string `json:"url"`        // 文件地址
	Version    string `json:"version"`    // 版本号
	SuitSystem string `json:"suitSystem"` // 适用的系统
	FileMd5    string `json:"fileMd5"`    // 文件md5编码
}

func (slot *SlotGateWayServiceImpl) ReplacePlugFile(ctx context.Context, args *pb.ReplaceFileInfo) (*emptypb.Empty, error) {
	defer func() {
		if err := recover(); err != nil {
			println(err.(string))
		}
	}()

	log.Println("ReplaceFile begin")

	var fileList []ExaFileUploadAndDownload

	switch args.FileName {
	case define.ProgramLotus.String(), define.ProgramMiner.String(), define.ProgramWorkerTask.String(), define.ProgramWorkerStorage.String(), "test":
		//os.Remove(path.Join(define.PathIpfsProgram, args.FileName))
		// 获取对应的产品id对应的文件id
		data, err := utils.RequestDo(args.DownloadUrl, define.SlotFileListRouter+"?productId="+utils.Int64ToString(int64(args.ProductId)),
			"", nil, time.Second*15)
		if err != nil {
			log.Println("get slot file List failed: ", err.Error())
			return &emptypb.Empty{}, err
		}
		// 解析返回信息
		type Response struct {
			Code int                 `json:"code"`
			Data response.PageResult `json:"data"`
			Msg  string              `json:"msg"`
		}
		var dataRes Response
		if err = json.Unmarshal(data, &dataRes); err != nil {
			log.Println(err.Error())
			return &emptypb.Empty{}, err
		}
		// 解析返回的file信息
		fileData, _ := json.Marshal(dataRes.Data.List)
		if err = json.Unmarshal(fileData, &fileList); err != nil {
			log.Println(err.Error())
			return &emptypb.Empty{}, err
		}
		// 循环拉取文件信息
		for _, val := range fileList {
			// 将对应的版本的文件都拉取到gateway上面,方便后面op拉取
			_, err := new(gateway.DownloadService).DowloadFile("http://"+args.DownloadUrl+define.SlotDownloadFileRouter+"?fileId="+utils.Int64ToString(int64(val.ID)),
				define.PathIpfsProgram, args.FileName+"-"+strings.Replace(val.SuitSystem, " ", "-", -1))
			if err != nil {
				log.Println("DownloadFile err:", err.Error())
			}
		}
	default:
		return &emptypb.Empty{}, errors.New("non system required plugin files")
	}

	log.Println("args", args)
	log.Println("len(fileList)", len(fileList))
	log.Println("fileList", fileList)

	fileInfoList := make([]*pb.ReplaceInfo, len(fileList))
	for i := 0; i < len(fileList); i++ {
		fileInfoList[i] = &pb.ReplaceInfo{System: fileList[i].SuitSystem, FileMd5: fileList[i].FileMd5}
	}

	log.Println("fileInfoList", fileInfoList)

	// 获取gateway下面的机器列表
	hostService := systemModel.HostRecordService{}
	hostList, err := hostService.GetSysHostRecordListForReplace(global.GateWayID.String())
	if err != nil {
		global.ZC_LOG.Error("获取节点下的主机列表失败!", zap.Error(err))
	}

	if len(hostList) == 0 {
		return &emptypb.Empty{}, errors.New("non system required plugin files")
	}

	fmt.Println("pb.OpReplaceFileInfo", &pb.OpReplaceFileInfo{FileName: args.FileName, FileInfo: fileInfoList})

	client, dis := global.OpClinets.GetOpClient("b4eb6358-3b86-4f99-8e01-786463e3c6ea")
	if client == nil || dis {
		log.Println("opClient Connection failed:" + "b4eb6358-3b86-4f99-8e01-786463e3c6ea" + " opIp: " + "10.0.8.220")
		return &emptypb.Empty{}, errors.New("non system required plugin files")
	}
	// 判断是否能成功连接
	_, err = client.OpReplacePlugFile1(ctx, &pb.String{Value: "test op connect"})
	if err != nil {
		log.Println("OpClinets error, opIp: ", "10.0.8.220", err.Error())
		return &emptypb.Empty{}, errors.New("non system required plugin files")
	}
	// 转到op,进行文件的拉取
	_, err = client.OpReplacePlugFile(ctx, &pb.OpReplaceFileInfo{FileName: args.FileName, FileInfo: fileInfoList})
	if err != nil {
		log.Println("OpReplacePlugFile error, opIp: ", "10.0.8.220", err.Error())
		return &emptypb.Empty{}, errors.New("non system required plugin files")
	}
	return &emptypb.Empty{}, nil

	// 循环gateway下面的op,根据对应linux版本下载不同的文件,如果文件存在,MD5相同则不拉取,不同则拉取替换
	for _, v := range hostList {
		go func(opInfo responseModel.SysHostRecordPatrol) {
			defer func() {
				if err := recover(); err != nil {
					println(err.(string))
				}
			}()
			client, dis := global.OpClinets.GetOpClient(opInfo.UUID)
			if client == nil || dis {
				log.Println("opClient Connection failed:" + opInfo.UUID + " opIp: " + opInfo.IntranetIP)
				return
			}
			// 判断是否能成功连接
			_, err := client.OpReplacePlugFile1(ctx, &pb.String{Value: "test op connect"})
			if err != nil {
				log.Println("OpClinets error, opIp: ", opInfo.IntranetIP, err.Error())
				return
			}
			// 转到op,进行文件的拉取
			_, err = client.OpReplacePlugFile(ctx, &pb.OpReplaceFileInfo{FileName: args.FileName, FileInfo: fileInfoList})
			if err != nil {
				log.Println("OpReplacePlugFile error, opIp: ", opInfo.IntranetIP, err.Error())
				return
			}
		}(v)
	}

	return &emptypb.Empty{}, nil
}
