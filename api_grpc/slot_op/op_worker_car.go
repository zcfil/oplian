package slot_op

import (
	"context"
	"oplian/global"
	"oplian/service/pb"
)

// GetRunWorkerTask 获取在跑任务
func (slot *SlotOpServiceImpl) GetRunWorkerTask(ctx context.Context, args *pb.String) (*pb.CarWorkerTaskNoInfo, error) {
	return global.OpToGatewayClient.GetRunWorkerTask(ctx, args)
}

// ModifyCarTaskNo 更新car任务编号
func (slot *SlotOpServiceImpl) ModifyCarTaskNo(ctx context.Context, args *pb.CarWorkerTaskNoInfo) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.ModifyCarTaskNo(ctx, args)
}

// AddCarWorkerTaskDetail 新增workerTask明细
func (slot *SlotOpServiceImpl) AddCarWorkerTaskDetail(ctx context.Context, args *pb.CarWorkerTaskDetailInfo) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.AddCarWorkerTaskDetail(ctx, args)
}

// GetRunCarTaskDetail 获取在跑workerCar任务明细
func (slot *SlotOpServiceImpl) GetRunCarTaskDetail(ctx context.Context, args *pb.WorkerCarParam) (*pb.CarWorkerTaskDetailInfo, error) {
	return global.OpToGatewayClient.GetRunCarTaskDetail(ctx, args)
}

// ModifyCarTaskDetailInfo 更改任务明细信息
func (slot *SlotOpServiceImpl) ModifyCarTaskDetailInfo(ctx context.Context, args *pb.CarWorkerTaskDetailInfo) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.ModifyCarTaskDetailInfo(ctx, args)
}

// GetRand 获取随机数
func (slot *SlotOpServiceImpl) GetRand(ctx context.Context, args *pb.String) (*pb.RandList, error) {
	return global.OpToGatewayClient.GetRand(ctx, &pb.String{})
}

// GetWaitCarTaskDetail 获取待创建任务明细
func (slot *SlotOpServiceImpl) GetWaitCarTaskDetail(ctx context.Context, args *pb.WorkerCarParam) (*pb.CarWorkerTaskDetailInfo, error) {
	return global.OpToGatewayClient.GetWaitCarTaskDetail(ctx, args)
}

// GetAllCarTaskDetail 获取所有待创建任务明细
func (slot *SlotOpServiceImpl) GetAllCarTaskDetail(ctx context.Context, args *pb.WorkerCarParam) (*pb.CarWorkerTaskDetailList, error) {
	return global.OpToGatewayClient.GetAllCarTaskDetail(ctx, args)
}

// GetBoostConfig 获取boost信息配置
func (slot *SlotOpServiceImpl) GetBoostConfig(ctx context.Context, args *pb.String) (*pb.String, error) {
	return global.OpToGatewayClient.GetBoostConfig(ctx, args)
}

// GetMainDisk 获取根目录
func (slot *SlotOpServiceImpl) GetMainDisk(ctx context.Context, args *pb.String) (*pb.String, error) {
	return global.OpToGatewayClient.GetMainDisk(ctx, args)
}

// DistributeWorkerTask workerCar任务分发
func (slot *SlotOpServiceImpl) DistributeWorkerTask(ctx context.Context, args *pb.String) (*pb.String, error) {
	return global.OpToGatewayClient.DistributeWorkerTask(ctx, args)
}

// AddCarRand 生成随机数
func (slot *SlotOpServiceImpl) AddCarRand(ctx context.Context, args *pb.String) (*pb.String, error) {
	return global.OpToGatewayClient.AddCarRand(ctx, args)
}

// CarFileExist 判断car文件是否存在
func (slot *SlotOpServiceImpl) CarFileExist(ctx context.Context, args *pb.SectorID) (*pb.String, error) {
	return global.OpToGatewayClient.CarFileExist(ctx, args)
}

// AddCarFile 增加Car文件
func (slot *SlotOpServiceImpl) AddCarFile(ctx context.Context, args *pb.CarFiles) (*pb.String, error) {
	return global.OpToGatewayClient.AddCarFile(ctx, args)
}
