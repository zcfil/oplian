// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: service/proto/slot/slot_header.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BatchUnsealed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host   []*HostParam `protobuf:"bytes,10,rep,name=host,proto3" json:"host,omitempty"`
	GateId string       `protobuf:"bytes,20,opt,name=gateId,proto3" json:"gateId,omitempty"`
}

func (x *BatchUnsealed) Reset() {
	*x = BatchUnsealed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_slot_slot_header_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchUnsealed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchUnsealed) ProtoMessage() {}

func (x *BatchUnsealed) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_slot_slot_header_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchUnsealed.ProtoReflect.Descriptor instead.
func (*BatchUnsealed) Descriptor() ([]byte, []int) {
	return file_service_proto_slot_slot_header_proto_rawDescGZIP(), []int{0}
}

func (x *BatchUnsealed) GetHost() []*HostParam {
	if x != nil {
		return x.Host
	}
	return nil
}

func (x *BatchUnsealed) GetGateId() string {
	if x != nil {
		return x.GateId
	}
	return ""
}

type CarWorkerTaskNoInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint64 `protobuf:"varint,10,opt,name=Id,proto3" json:"Id,omitempty"`
	TaskId     uint64 `protobuf:"varint,11,opt,name=TaskId,proto3" json:"TaskId,omitempty"`
	MinerId    string `protobuf:"bytes,1,opt,name=MinerId,proto3" json:"MinerId,omitempty"`
	WorkerIp   string `protobuf:"bytes,2,opt,name=WorkerIp,proto3" json:"WorkerIp,omitempty"`
	CarNo      string `protobuf:"bytes,3,opt,name=CarNo,proto3" json:"CarNo,omitempty"`
	InputDir   string `protobuf:"bytes,4,opt,name=InputDir,proto3" json:"InputDir,omitempty"`
	OutputDir  string `protobuf:"bytes,5,opt,name=OutputDir,proto3" json:"OutputDir,omitempty"`
	WalletAddr string `protobuf:"bytes,6,opt,name=WalletAddr,proto3" json:"WalletAddr,omitempty"`
	StartNo    uint64 `protobuf:"varint,7,opt,name=StartNo,proto3" json:"StartNo,omitempty"`
	EndNo      uint64 `protobuf:"varint,8,opt,name=EndNo,proto3" json:"EndNo,omitempty"`
	TaskStatus uint64 `protobuf:"varint,9,opt,name=TaskStatus,proto3" json:"TaskStatus,omitempty"`
}

func (x *CarWorkerTaskNoInfo) Reset() {
	*x = CarWorkerTaskNoInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_slot_slot_header_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CarWorkerTaskNoInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarWorkerTaskNoInfo) ProtoMessage() {}

func (x *CarWorkerTaskNoInfo) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_slot_slot_header_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarWorkerTaskNoInfo.ProtoReflect.Descriptor instead.
func (*CarWorkerTaskNoInfo) Descriptor() ([]byte, []int) {
	return file_service_proto_slot_slot_header_proto_rawDescGZIP(), []int{1}
}

func (x *CarWorkerTaskNoInfo) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CarWorkerTaskNoInfo) GetTaskId() uint64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *CarWorkerTaskNoInfo) GetMinerId() string {
	if x != nil {
		return x.MinerId
	}
	return ""
}

func (x *CarWorkerTaskNoInfo) GetWorkerIp() string {
	if x != nil {
		return x.WorkerIp
	}
	return ""
}

func (x *CarWorkerTaskNoInfo) GetCarNo() string {
	if x != nil {
		return x.CarNo
	}
	return ""
}

func (x *CarWorkerTaskNoInfo) GetInputDir() string {
	if x != nil {
		return x.InputDir
	}
	return ""
}

func (x *CarWorkerTaskNoInfo) GetOutputDir() string {
	if x != nil {
		return x.OutputDir
	}
	return ""
}

func (x *CarWorkerTaskNoInfo) GetWalletAddr() string {
	if x != nil {
		return x.WalletAddr
	}
	return ""
}

func (x *CarWorkerTaskNoInfo) GetStartNo() uint64 {
	if x != nil {
		return x.StartNo
	}
	return 0
}

func (x *CarWorkerTaskNoInfo) GetEndNo() uint64 {
	if x != nil {
		return x.EndNo
	}
	return 0
}

func (x *CarWorkerTaskNoInfo) GetTaskStatus() uint64 {
	if x != nil {
		return x.TaskStatus
	}
	return 0
}

type CarWorkerTaskDetailInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             uint64 `protobuf:"varint,10,opt,name=Id,proto3" json:"Id,omitempty"`
	TaskId         uint64 `protobuf:"varint,1,opt,name=TaskId,proto3" json:"TaskId,omitempty"`
	TaskName       string `protobuf:"bytes,2,opt,name=TaskName,proto3" json:"TaskName,omitempty"`
	DealExpireDate string `protobuf:"bytes,3,opt,name=DealExpireDate,proto3" json:"DealExpireDate,omitempty"`
	CarName        string `protobuf:"bytes,4,opt,name=CarName,proto3" json:"CarName,omitempty"`
	PieceCid       string `protobuf:"bytes,5,opt,name=PieceCid,proto3" json:"PieceCid,omitempty"`
	PieceSize      uint64 `protobuf:"varint,6,opt,name=PieceSize,proto3" json:"PieceSize,omitempty"`
	CarSize        uint64 `protobuf:"varint,7,opt,name=CarSize,proto3" json:"CarSize,omitempty"`
	DataCid        string `protobuf:"bytes,8,opt,name=DataCid,proto3" json:"DataCid,omitempty"`
	MinerId        string `protobuf:"bytes,9,opt,name=MinerId,proto3" json:"MinerId,omitempty"`
	WalletAddr     string `protobuf:"bytes,11,opt,name=WalletAddr,proto3" json:"WalletAddr,omitempty"`
	TaskStatus     uint64 `protobuf:"varint,12,opt,name=TaskStatus,proto3" json:"TaskStatus,omitempty"`
	WorkerIp       string `protobuf:"bytes,13,opt,name=WorkerIp,proto3" json:"WorkerIp,omitempty"`
	DealId         string `protobuf:"bytes,14,opt,name=DealId,proto3" json:"DealId,omitempty"`
	SectorId       string `protobuf:"bytes,16,opt,name=SectorId,proto3" json:"SectorId,omitempty"`
	CarOutputDir   string `protobuf:"bytes,17,opt,name=CarOutputDir,proto3" json:"CarOutputDir,omitempty"`
	OriginalOpId   string `protobuf:"bytes,18,opt,name=OriginalOpId,proto3" json:"OriginalOpId,omitempty"`
	OriginalDir    string `protobuf:"bytes,19,opt,name=OriginalDir,proto3" json:"OriginalDir,omitempty"`
	ValidityDays   uint64 `protobuf:"varint,20,opt,name=ValidityDays,proto3" json:"ValidityDays,omitempty"`
}

func (x *CarWorkerTaskDetailInfo) Reset() {
	*x = CarWorkerTaskDetailInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_slot_slot_header_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CarWorkerTaskDetailInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarWorkerTaskDetailInfo) ProtoMessage() {}

func (x *CarWorkerTaskDetailInfo) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_slot_slot_header_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarWorkerTaskDetailInfo.ProtoReflect.Descriptor instead.
func (*CarWorkerTaskDetailInfo) Descriptor() ([]byte, []int) {
	return file_service_proto_slot_slot_header_proto_rawDescGZIP(), []int{2}
}

func (x *CarWorkerTaskDetailInfo) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CarWorkerTaskDetailInfo) GetTaskId() uint64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *CarWorkerTaskDetailInfo) GetTaskName() string {
	if x != nil {
		return x.TaskName
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetDealExpireDate() string {
	if x != nil {
		return x.DealExpireDate
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetCarName() string {
	if x != nil {
		return x.CarName
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetPieceCid() string {
	if x != nil {
		return x.PieceCid
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetPieceSize() uint64 {
	if x != nil {
		return x.PieceSize
	}
	return 0
}

func (x *CarWorkerTaskDetailInfo) GetCarSize() uint64 {
	if x != nil {
		return x.CarSize
	}
	return 0
}

func (x *CarWorkerTaskDetailInfo) GetDataCid() string {
	if x != nil {
		return x.DataCid
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetMinerId() string {
	if x != nil {
		return x.MinerId
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetWalletAddr() string {
	if x != nil {
		return x.WalletAddr
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetTaskStatus() uint64 {
	if x != nil {
		return x.TaskStatus
	}
	return 0
}

func (x *CarWorkerTaskDetailInfo) GetWorkerIp() string {
	if x != nil {
		return x.WorkerIp
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetDealId() string {
	if x != nil {
		return x.DealId
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetSectorId() string {
	if x != nil {
		return x.SectorId
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetCarOutputDir() string {
	if x != nil {
		return x.CarOutputDir
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetOriginalOpId() string {
	if x != nil {
		return x.OriginalOpId
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetOriginalDir() string {
	if x != nil {
		return x.OriginalDir
	}
	return ""
}

func (x *CarWorkerTaskDetailInfo) GetValidityDays() uint64 {
	if x != nil {
		return x.ValidityDays
	}
	return 0
}

type RandInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumberIndex uint64 `protobuf:"varint,1,opt,name=NumberIndex,proto3" json:"NumberIndex,omitempty"`
	Number      uint64 `protobuf:"varint,2,opt,name=Number,proto3" json:"Number,omitempty"`
}

func (x *RandInfo) Reset() {
	*x = RandInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_slot_slot_header_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RandInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RandInfo) ProtoMessage() {}

func (x *RandInfo) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_slot_slot_header_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RandInfo.ProtoReflect.Descriptor instead.
func (*RandInfo) Descriptor() ([]byte, []int) {
	return file_service_proto_slot_slot_header_proto_rawDescGZIP(), []int{3}
}

func (x *RandInfo) GetNumberIndex() uint64 {
	if x != nil {
		return x.NumberIndex
	}
	return 0
}

func (x *RandInfo) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

type WorkerCarParam struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	MinerId  string `protobuf:"bytes,2,opt,name=MinerId,proto3" json:"MinerId,omitempty"`
	WorkerIp string `protobuf:"bytes,3,opt,name=WorkerIp,proto3" json:"WorkerIp,omitempty"`
	PieceCID string `protobuf:"bytes,4,opt,name=PieceCID,proto3" json:"PieceCID,omitempty"`
	SectorId string `protobuf:"bytes,5,opt,name=SectorId,proto3" json:"SectorId,omitempty"`
	TaskType uint64 `protobuf:"varint,6,opt,name=TaskType,proto3" json:"TaskType,omitempty"`
}

func (x *WorkerCarParam) Reset() {
	*x = WorkerCarParam{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_slot_slot_header_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerCarParam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerCarParam) ProtoMessage() {}

func (x *WorkerCarParam) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_slot_slot_header_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerCarParam.ProtoReflect.Descriptor instead.
func (*WorkerCarParam) Descriptor() ([]byte, []int) {
	return file_service_proto_slot_slot_header_proto_rawDescGZIP(), []int{4}
}

func (x *WorkerCarParam) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *WorkerCarParam) GetMinerId() string {
	if x != nil {
		return x.MinerId
	}
	return ""
}

func (x *WorkerCarParam) GetWorkerIp() string {
	if x != nil {
		return x.WorkerIp
	}
	return ""
}

func (x *WorkerCarParam) GetPieceCID() string {
	if x != nil {
		return x.PieceCID
	}
	return ""
}

func (x *WorkerCarParam) GetSectorId() string {
	if x != nil {
		return x.SectorId
	}
	return ""
}

func (x *WorkerCarParam) GetTaskType() uint64 {
	if x != nil {
		return x.TaskType
	}
	return 0
}

type RandList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RandInfo []*RandInfo `protobuf:"bytes,1,rep,name=randInfo,proto3" json:"randInfo,omitempty"`
}

func (x *RandList) Reset() {
	*x = RandList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_slot_slot_header_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RandList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RandList) ProtoMessage() {}

func (x *RandList) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_slot_slot_header_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RandList.ProtoReflect.Descriptor instead.
func (*RandList) Descriptor() ([]byte, []int) {
	return file_service_proto_slot_slot_header_proto_rawDescGZIP(), []int{5}
}

func (x *RandList) GetRandInfo() []*RandInfo {
	if x != nil {
		return x.RandInfo
	}
	return nil
}

type CarWorkerTaskDetailList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CarWorkerTaskDetailInfo []*CarWorkerTaskDetailInfo `protobuf:"bytes,1,rep,name=carWorkerTaskDetailInfo,proto3" json:"carWorkerTaskDetailInfo,omitempty"`
}

func (x *CarWorkerTaskDetailList) Reset() {
	*x = CarWorkerTaskDetailList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_slot_slot_header_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CarWorkerTaskDetailList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarWorkerTaskDetailList) ProtoMessage() {}

func (x *CarWorkerTaskDetailList) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_slot_slot_header_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarWorkerTaskDetailList.ProtoReflect.Descriptor instead.
func (*CarWorkerTaskDetailList) Descriptor() ([]byte, []int) {
	return file_service_proto_slot_slot_header_proto_rawDescGZIP(), []int{6}
}

func (x *CarWorkerTaskDetailList) GetCarWorkerTaskDetailInfo() []*CarWorkerTaskDetailInfo {
	if x != nil {
		return x.CarWorkerTaskDetailInfo
	}
	return nil
}

type CarFiles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RelationId  uint64 `protobuf:"varint,10,opt,name=RelationId,proto3" json:"RelationId,omitempty"`
	FileName    string `protobuf:"bytes,20,opt,name=FileName,proto3" json:"FileName,omitempty"`
	FileIndex   uint64 `protobuf:"varint,30,opt,name=FileIndex,proto3" json:"FileIndex,omitempty"`
	FileStr     string `protobuf:"bytes,40,opt,name=FileStr,proto3" json:"FileStr,omitempty"`
	CarFileName string `protobuf:"bytes,50,opt,name=CarFileName,proto3" json:"CarFileName,omitempty"`
	PieceCid    string `protobuf:"bytes,60,opt,name=PieceCid,proto3" json:"PieceCid,omitempty"`
	PieceSize   uint64 `protobuf:"varint,70,opt,name=PieceSize,proto3" json:"PieceSize,omitempty"`
	CarSize     uint64 `protobuf:"varint,80,opt,name=CarSize,proto3" json:"CarSize,omitempty"`
	DataCid     string `protobuf:"bytes,90,opt,name=DataCid,proto3" json:"DataCid,omitempty"`
	InputDir    string `protobuf:"bytes,100,opt,name=InputDir,proto3" json:"InputDir,omitempty"`
	MinerId     string `protobuf:"bytes,110,opt,name=MinerId,proto3" json:"MinerId,omitempty"`
}

func (x *CarFiles) Reset() {
	*x = CarFiles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_slot_slot_header_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CarFiles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarFiles) ProtoMessage() {}

func (x *CarFiles) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_slot_slot_header_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarFiles.ProtoReflect.Descriptor instead.
func (*CarFiles) Descriptor() ([]byte, []int) {
	return file_service_proto_slot_slot_header_proto_rawDescGZIP(), []int{7}
}

func (x *CarFiles) GetRelationId() uint64 {
	if x != nil {
		return x.RelationId
	}
	return 0
}

func (x *CarFiles) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *CarFiles) GetFileIndex() uint64 {
	if x != nil {
		return x.FileIndex
	}
	return 0
}

func (x *CarFiles) GetFileStr() string {
	if x != nil {
		return x.FileStr
	}
	return ""
}

func (x *CarFiles) GetCarFileName() string {
	if x != nil {
		return x.CarFileName
	}
	return ""
}

func (x *CarFiles) GetPieceCid() string {
	if x != nil {
		return x.PieceCid
	}
	return ""
}

func (x *CarFiles) GetPieceSize() uint64 {
	if x != nil {
		return x.PieceSize
	}
	return 0
}

func (x *CarFiles) GetCarSize() uint64 {
	if x != nil {
		return x.CarSize
	}
	return 0
}

func (x *CarFiles) GetDataCid() string {
	if x != nil {
		return x.DataCid
	}
	return ""
}

func (x *CarFiles) GetInputDir() string {
	if x != nil {
		return x.InputDir
	}
	return ""
}

func (x *CarFiles) GetMinerId() string {
	if x != nil {
		return x.MinerId
	}
	return ""
}

var File_service_proto_slot_slot_header_proto protoreflect.FileDescriptor

var file_service_proto_slot_slot_header_proto_rawDesc = []byte{
	0x0a, 0x24, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x73, 0x6c, 0x6f, 0x74, 0x2f, 0x73, 0x6c, 0x6f, 0x74, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x47, 0x0a, 0x0d, 0x62, 0x61, 0x74, 0x63, 0x68, 0x55, 0x6e, 0x73, 0x65, 0x61,
	0x6c, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x0a, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0a, 0x2e, 0x68, 0x6f, 0x73, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x52, 0x04, 0x68,
	0x6f, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x61, 0x74, 0x65, 0x49, 0x64, 0x18, 0x14, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x61, 0x74, 0x65, 0x49, 0x64, 0x22, 0xb3, 0x02, 0x0a, 0x13,
	0x43, 0x61, 0x72, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x54, 0x61, 0x73, 0x6b, 0x4e, 0x6f, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x4d,
	0x69, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x69,
	0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49,
	0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49,
	0x70, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x61, 0x72, 0x4e, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x43, 0x61, 0x72, 0x4e, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x6e, 0x70, 0x75, 0x74,
	0x44, 0x69, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x49, 0x6e, 0x70, 0x75, 0x74,
	0x44, 0x69, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x44, 0x69, 0x72,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x44, 0x69,
	0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x41, 0x64, 0x64,
	0x72, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4e, 0x6f, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x07, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4e, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x45,
	0x6e, 0x64, 0x4e, 0x6f, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x45, 0x6e, 0x64, 0x4e,
	0x6f, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0xc5, 0x04, 0x0a, 0x17, 0x43, 0x61, 0x72, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x54,
	0x61, 0x73, 0x6b, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x54,
	0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x26, 0x0a, 0x0e, 0x44, 0x65, 0x61, 0x6c, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x44,
	0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x44, 0x65, 0x61, 0x6c, 0x45,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x44, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x61, 0x72,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x61, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x69, 0x65, 0x63, 0x65, 0x43, 0x69, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x69, 0x65, 0x63, 0x65, 0x43, 0x69, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x50, 0x69, 0x65, 0x63, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x09, 0x50, 0x69, 0x65, 0x63, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x43, 0x61, 0x72, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07,
	0x43, 0x61, 0x72, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x44, 0x61, 0x74, 0x61, 0x43,
	0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x44, 0x61, 0x74, 0x61, 0x43, 0x69,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x57,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x54,
	0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0a, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x57,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x70, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x57,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x44, 0x65, 0x61, 0x6c, 0x49,
	0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x44, 0x65, 0x61, 0x6c, 0x49, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x53, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x18, 0x10, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x53, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x43,
	0x61, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x44, 0x69, 0x72, 0x18, 0x11, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x43, 0x61, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x44, 0x69, 0x72, 0x12,
	0x22, 0x0a, 0x0c, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x4f, 0x70, 0x49, 0x64, 0x18,
	0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x4f,
	0x70, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x44,
	0x69, 0x72, 0x18, 0x13, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x61, 0x6c, 0x44, 0x69, 0x72, 0x12, 0x22, 0x0a, 0x0c, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x69, 0x74,
	0x79, 0x44, 0x61, 0x79, 0x73, 0x18, 0x14, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x69, 0x74, 0x79, 0x44, 0x61, 0x79, 0x73, 0x22, 0x44, 0x0a, 0x08, 0x52, 0x61, 0x6e,
	0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x20, 0x0a, 0x0b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22,
	0xaa, 0x01, 0x0a, 0x0e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x43, 0x61, 0x72, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x69, 0x65, 0x63,
	0x65, 0x43, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x69, 0x65, 0x63,
	0x65, 0x43, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x53, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x22, 0x31, 0x0a, 0x08,
	0x52, 0x61, 0x6e, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x08, 0x72, 0x61, 0x6e, 0x64,
	0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x52, 0x61, 0x6e,
	0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x72, 0x61, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x22,
	0x6d, 0x0a, 0x17, 0x43, 0x61, 0x72, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x54, 0x61, 0x73, 0x6b,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x52, 0x0a, 0x17, 0x63, 0x61,
	0x72, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x54, 0x61, 0x73, 0x6b, 0x44, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x43, 0x61,
	0x72, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x54, 0x61, 0x73, 0x6b, 0x44, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x17, 0x63, 0x61, 0x72, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x54, 0x61, 0x73, 0x6b, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0xc4,
	0x02, 0x0a, 0x08, 0x43, 0x61, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x52,
	0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0a, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x46,
	0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46,
	0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x46, 0x69, 0x6c, 0x65,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72,
	0x18, 0x28, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x12,
	0x20, 0x0a, 0x0b, 0x43, 0x61, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x32,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x43, 0x61, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x69, 0x65, 0x63, 0x65, 0x43, 0x69, 0x64, 0x18, 0x3c, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x69, 0x65, 0x63, 0x65, 0x43, 0x69, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x50, 0x69, 0x65, 0x63, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x46, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x50, 0x69, 0x65, 0x63, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x43,
	0x61, 0x72, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x50, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x43, 0x61,
	0x72, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x44, 0x61, 0x74, 0x61, 0x43, 0x69, 0x64,
	0x18, 0x5a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x44, 0x61, 0x74, 0x61, 0x43, 0x69, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x44, 0x69, 0x72, 0x18, 0x64, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x44, 0x69, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x4d,
	0x69, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x6e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x69,
	0x6e, 0x65, 0x72, 0x49, 0x64, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_service_proto_slot_slot_header_proto_rawDescOnce sync.Once
	file_service_proto_slot_slot_header_proto_rawDescData = file_service_proto_slot_slot_header_proto_rawDesc
)

func file_service_proto_slot_slot_header_proto_rawDescGZIP() []byte {
	file_service_proto_slot_slot_header_proto_rawDescOnce.Do(func() {
		file_service_proto_slot_slot_header_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_proto_slot_slot_header_proto_rawDescData)
	})
	return file_service_proto_slot_slot_header_proto_rawDescData
}

var file_service_proto_slot_slot_header_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_service_proto_slot_slot_header_proto_goTypes = []interface{}{
	(*BatchUnsealed)(nil),           // 0: batchUnsealed
	(*CarWorkerTaskNoInfo)(nil),     // 1: CarWorkerTaskNoInfo
	(*CarWorkerTaskDetailInfo)(nil), // 2: CarWorkerTaskDetailInfo
	(*RandInfo)(nil),                // 3: RandInfo
	(*WorkerCarParam)(nil),          // 4: WorkerCarParam
	(*RandList)(nil),                // 5: RandList
	(*CarWorkerTaskDetailList)(nil), // 6: CarWorkerTaskDetailList
	(*CarFiles)(nil),                // 7: CarFiles
	(*HostParam)(nil),               // 8: hostParam
}
var file_service_proto_slot_slot_header_proto_depIdxs = []int32{
	8, // 0: batchUnsealed.host:type_name -> hostParam
	3, // 1: RandList.randInfo:type_name -> RandInfo
	2, // 2: CarWorkerTaskDetailList.carWorkerTaskDetailInfo:type_name -> CarWorkerTaskDetailInfo
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_service_proto_slot_slot_header_proto_init() }
func file_service_proto_slot_slot_header_proto_init() {
	if File_service_proto_slot_slot_header_proto != nil {
		return
	}
	file_service_proto_header_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_service_proto_slot_slot_header_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchUnsealed); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_proto_slot_slot_header_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CarWorkerTaskNoInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_proto_slot_slot_header_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CarWorkerTaskDetailInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_proto_slot_slot_header_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RandInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_proto_slot_slot_header_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkerCarParam); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_proto_slot_slot_header_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RandList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_proto_slot_slot_header_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CarWorkerTaskDetailList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_proto_slot_slot_header_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CarFiles); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_proto_slot_slot_header_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_service_proto_slot_slot_header_proto_goTypes,
		DependencyIndexes: file_service_proto_slot_slot_header_proto_depIdxs,
		MessageInfos:      file_service_proto_slot_slot_header_proto_msgTypes,
	}.Build()
	File_service_proto_slot_slot_header_proto = out.File
	file_service_proto_slot_slot_header_proto_rawDesc = nil
	file_service_proto_slot_slot_header_proto_goTypes = nil
	file_service_proto_slot_slot_header_proto_depIdxs = nil
}