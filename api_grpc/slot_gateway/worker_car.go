package slot_gateway

import (
	"context"
	"errors"
	"fmt"
	"oplian/define"
	"oplian/model/lotus"
	"oplian/model/slot"
	"oplian/model/slot/response"
	"oplian/service/pb"
	slot2 "oplian/service/slot"
	"oplian/utils"
	"sync"
)

type SlotGateWayServiceImpl struct {
	CarUrlMap    map[string]string
	CarUrlLockRW sync.RWMutex
}

var SlotGateWayService = SlotGateWayServiceImpl{
	CarUrlMap: make(map[string]string),
}

// GetRunWorkerTask 获取在跑任务 OK
func (sl *SlotGateWayServiceImpl) GetRunWorkerTask(ctx context.Context, args *pb.String) (*pb.CarWorkerTaskNoInfo, error) {

	resInfo := &pb.CarWorkerTaskNoInfo{}
	resData, err := slot2.WorkerCarServiceApi.GetRunWorkerTask(args)
	if err != nil {
		return resInfo, err
	}
	if (resData != slot.WorkerCarTaskNo{}) {

		resInfo.Id = uint64(resData.ID)
		resInfo.TaskId = uint64(resData.TaskId)
		resInfo.MinerId = resData.MinerId
		resInfo.WorkerIp = resData.WorkerIp
		resInfo.CarNo = resData.CarNo
		resInfo.StartNo = uint64(resData.StartNo)
		resInfo.EndNo = uint64(resData.EndNo)
		resInfo.InputDir = resData.InputDir
		resInfo.OutputDir = resData.OutputDir
		resInfo.TaskStatus = uint64(resData.TaskStatus)
	}

	return resInfo, nil
}

// ModifyCarTaskNo 更新car任务编号 OK
func (slot *SlotGateWayServiceImpl) ModifyCarTaskNo(ctx context.Context, args *pb.CarWorkerTaskNoInfo) (*pb.ResponseMsg, error) {

	resInfo := &pb.ResponseMsg{}
	err := slot2.WorkerCarServiceApi.ModifyCarTaskNo(args)
	if err != nil {
		return resInfo, err
	}
	return resInfo, nil
}

// AddCarWorkerTaskDetail 新增workerTask明细 OK
func (slot *SlotGateWayServiceImpl) AddCarWorkerTaskDetail(ctx context.Context, args *pb.CarWorkerTaskDetailInfo) (*pb.ResponseMsg, error) {

	resInfo := &pb.ResponseMsg{}

	Height := utils.BlockHeight()
	Height += define.DealDdValidity
	paramAr := make([]lotus.LotusSectorPiece, 0)
	paramData := &lotus.LotusSectorPiece{
		Actor:          args.MinerId,
		WorkerIp:       args.WorkerIp,
		QueueId:        args.TaskId,
		ExpirationTime: utils.BlockHeightToTime(int64(Height)),
		PieceCid:       args.PieceCid,
		PieceSize:      args.PieceSize,
		CarSize:        int(args.CarSize),
		DataCid:        args.DataCid,
	}

	paramAr = append(paramAr, *paramData)
	err := slot2.WorkerCarServiceApi.AddWorkerCarTaskDetail(args.MinerId, paramAr)
	if err != nil {
		return resInfo, err
	}
	return resInfo, nil
}

// GetRunCarTaskDetail 获取在跑workerCar任务明细 OK
func (slot *SlotGateWayServiceImpl) GetRunCarTaskDetail(ctx context.Context, args *pb.WorkerCarParam) (*pb.CarWorkerTaskDetailInfo, error) {

	resInfo := &pb.CarWorkerTaskDetailInfo{}
	resData, err := slot2.WorkerCarServiceApi.GetRunCarTaskDetail(args.WorkerIp, args.MinerId, int(args.TaskType))
	if err != nil {
		return resInfo, err
	}

	if (resData != response.CarTaskDetailInfo{}) {
		resInfo.Id = uint64(resData.Id)
		resInfo.PieceCid = resData.PieceCid
		resInfo.PieceSize = resData.PieceSize
		resInfo.CarSize = uint64(resData.CarSize)
		resInfo.DataCid = resData.DataCid
		resInfo.WalletAddr = resData.QuotaWallet
		resInfo.TaskStatus = uint64(resData.JobStatus)
		resInfo.OriginalOpId = resData.OriginalOpId
		resInfo.OriginalDir = resData.OriginalDir
		resInfo.ValidityDays = uint64(resData.ValidityDays)
	}

	return resInfo, nil
}

// ModifyCarTaskDetailInfo 更改任务明细信息 OK
func (slot *SlotGateWayServiceImpl) ModifyCarTaskDetailInfo(ctx context.Context, args *pb.CarWorkerTaskDetailInfo) (*pb.ResponseMsg, error) {

	resInfo := &pb.ResponseMsg{}
	err := slot2.WorkerCarServiceApi.ModifyCarTaskDetailInfo(args)
	if err != nil {
		return resInfo, err
	}
	return resInfo, nil
}

// GetRand 获取随机数 OK
func (slot *SlotGateWayServiceImpl) GetRand(ctx context.Context, args *pb.String) (*pb.RandList, error) {

	resInfo := &pb.RandList{}
	resData, err := slot2.WorkerCarServiceApi.GetRand(&pb.String{})
	if err != nil {

	}
	for _, v := range resData {

		data := &pb.RandInfo{
			NumberIndex: uint64(v.NumIndex),
			Number:      uint64(v.Number),
		}
		resInfo.RandInfo = append(resInfo.RandInfo, data)
	}

	return resInfo, nil
}

// GetWaitCarTaskDetail 获取待创建任务明细 OK
func (sl *SlotGateWayServiceImpl) GetWaitCarTaskDetail(ctx context.Context, args *pb.WorkerCarParam) (*pb.CarWorkerTaskDetailInfo, error) {

	resInfo := &pb.CarWorkerTaskDetailInfo{}
	detailData, err := slot2.WorkerCarServiceApi.GetWaitCarTaskDetail(args.WorkerIp, args.MinerId)
	if err != nil {
		return resInfo, err
	}

	if (detailData != slot.WorkerCarTaskDetail{}) {
		resInfo.Id = uint64(detailData.ID)
		resInfo.MinerId = detailData.MinerId
		resInfo.WorkerIp = detailData.WorkerIp
		resInfo.DealId = detailData.DealId
		resInfo.CarName = detailData.CarName
	}

	return resInfo, nil
}

// GetAllCarTaskDetail 获取所有待创建任务明细 OK
func (sl *SlotGateWayServiceImpl) GetAllCarTaskDetail(ctx context.Context, args *pb.WorkerCarParam) (*pb.CarWorkerTaskDetailList, error) {

	resInfo := &pb.CarWorkerTaskDetailList{}
	resData, err := slot2.WorkerCarServiceApi.GetExecuteWorkerCarTask(args.MinerId, define.CarTaskTypeManual)
	if err != nil {
		return resInfo, err
	}
	if (resData != slot.WorkerCarTask{}) {

		detailData, err := slot2.WorkerCarServiceApi.GetAllCarTaskDetail(int(resData.ID), args.MinerId)
		if err != nil {
			return resInfo, err
		}

		if len(detailData) > 0 {

			for _, v := range detailData {

				data := &pb.CarWorkerTaskDetailInfo{
					WorkerIp: v.WorkerIp,
					//DealId:   v.DealId,
				}
				resInfo.CarWorkerTaskDetailInfo = append(resInfo.CarWorkerTaskDetailInfo, data)
			}
		}
	}

	return resInfo, nil
}

// CarFileExist 判断car文件是否存在 OK
func (slot *SlotGateWayServiceImpl) CarFileExist(ctx context.Context, args *pb.SectorID) (*pb.String, error) {
	SlotGateWayService.CarUrlLockRW.RLock()
	defer SlotGateWayService.CarUrlLockRW.RUnlock()

	resInfo := &pb.String{}
	if carUrl, ok := SlotGateWayService.CarUrlMap[args.PieceCid]; ok {
		resInfo.Value = carUrl
	} else {
		return resInfo, fmt.Errorf("%s,car文件路径找不到", args.PieceCid)
	}

	return resInfo, nil
}

// GetBoostConfig 获取boost信息配置
func (slot *SlotGateWayServiceImpl) GetBoostConfig(ctx context.Context, args *pb.String) (*pb.String, error) {
	return slot2.WorkerCarServiceApi.GetBoostConfig(args.Value)
}

// GetMainDisk 获取根目录
func (slot *SlotGateWayServiceImpl) GetMainDisk(ctx context.Context, args *pb.String) (*pb.String, error) {
	if define.MainDisk == "" {
		return &pb.String{}, errors.New("MainDisk 配置目录为空")
	}
	return &pb.String{Value: define.MainDisk}, nil
}

// DistributeWorkerTask workerCar任务分发
func (slot *SlotGateWayServiceImpl) DistributeWorkerTask(ctx context.Context, args *pb.String) (*pb.String, error) {

	go func() {

		//记录Car文件
		res, err := slot2.WorkerCarServiceApi.GetCarUrlList(args.Value)
		if err != nil {
			return
		}

		for _, v := range res {
			if v.CarPath != "" && v.PieceCid != "" {
				SlotGateWayService.CarUrlLockRW.Lock()
				SlotGateWayService.CarUrlMap[v.PieceCid] = v.CarPath
				SlotGateWayService.CarUrlLockRW.Unlock()
			}
		}

	}()

	return &pb.String{}, slot2.WorkerCarServiceApi.WorkerCarTask(args.Value)
}

// AddCarRand 生成随机数
func (slot *SlotGateWayServiceImpl) AddCarRand(ctx context.Context, args *pb.String) (*pb.String, error) {
	return &pb.String{}, slot2.WorkerCarServiceApi.AddCarRand()
}

// AddCarFile 增加Car文件
func (slot *SlotGateWayServiceImpl) AddCarFile(ctx context.Context, args *pb.CarFiles) (*pb.String, error) {
	return &pb.String{}, slot2.WorkerCarServiceApi.AddCarFiles(args)
}
