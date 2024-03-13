package utils

var (
	IdVerify                 = Rules{"ID": []string{NotEmpty()}}
	ApiVerify                = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	MenuVerify               = Rules{"Path": {NotEmpty()}, "ParentId": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify           = Rules{"Title": {NotEmpty()}}
	LoginVerify              = Rules{"CaptchaId": {NotEmpty()}, "Captcha": {NotEmpty()}, "Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify           = Rules{"Username": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}, "AuthorityId": {NotEmpty()}}
	PageInfoVerify           = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	CustomerVerify           = Rules{"CustomerName": {NotEmpty()}, "CustomerPhoneData": {NotEmpty()}}
	AutoCodeVerify           = Rules{"Abbreviation": {NotEmpty()}, "StructName": {NotEmpty()}, "PackageName": {NotEmpty()}, "Fields": {NotEmpty()}}
	AutoPackageVerify        = Rules{"PackageName": {NotEmpty()}}
	AuthorityVerify          = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}}
	AuthorityIdVerify        = Rules{"AuthorityId": {NotEmpty()}}
	OldAuthorityVerify       = Rules{"OldAuthorityId": {NotEmpty()}}
	ChangePasswordVerify     = Rules{"Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	SetUserAuthorityVerify   = Rules{"AuthorityId": {NotEmpty()}}
	GateWayInfoVerify        = Rules{"GateWayId": {NotEmpty()}, "Ip": {NotEmpty()}, "Port": {NotEmpty()}, "Token": {NotEmpty()}}
	MachineRoomVerify        = Rules{"RoomName": {NotEmpty()}}
	HostRecordVerify         = Rules{"HostName": {NotEmpty()}, "HostType": {NotEmpty()}, "RoomId": {NotEmpty()}}
	BindSysHostRecordsVerify = Rules{"RoomId": {NotEmpty()}, "HostUUIDs": {NotEmpty()}}
	RequesOpVerify           = Rules{"GateId": {NotEmpty()}, "OpId": {NotEmpty()}}
	RequesFileTypeReqVerify  = Rules{"GateId": {NotEmpty()}, "FileType": {NotEmpty()}}
	GatewayVerify            = Rules{"GateId": {NotEmpty()}}
	//SectorVerify             = Rules{"GateId": {NotEmpty()}, "OpId": {NotEmpty()}, "SectorType": {NotEmpty()}, "Actor": {NotEmpty()}, "Count": {NotEmpty()}}
	ActorVerify            = Rules{"Actor": {NotEmpty()}}
	AddrVerify             = Rules{"Addr": {NotEmpty()}}
	MinerIdVerify          = Rules{"GateWayId": {NotEmpty()}, "MinerId": {NotEmpty()}}
	DealPageVerify         = Rules{"Actor": {NotEmpty()}, "QueueId": {NotEmpty()}, "SectorType": {NotEmpty()}}
	RecoverDetailVerify    = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}, "Id": {NotEmpty()}}
	ActorSectorIdVerify    = Rules{"Actor": {NotEmpty()}, "SectorId": {NotEmpty()}}
	SectorTaskVerify       = Rules{"GateId": {NotEmpty()}, "OpId": {NotEmpty()}, "SectorType": {NotEmpty()}, "SectorTotal": {NotEmpty()}}
	IdStatusVerify         = Rules{"ID": {NotEmpty()}, "Status": {NotEmpty()}}
	IdCountVerify          = Rules{"ID": {NotEmpty()}, "Count": {NotEmpty()}}
	ActorIdStatusVerify    = Rules{"Actor": {NotEmpty()}, "ID": {NotEmpty()}, "Status": {NotEmpty()}}
	DealCarVerify          = Rules{"Actor": {NotEmpty()}, "ID": {NotEmpty()}, "CarDir": {NotEmpty()}}
	ClassifyVerify         = Rules{"Classify": {NotEmpty()}}
	PateVerify             = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	DistributionVerify     = Rules{"FilePath": {NotEmpty()}, "SendType": {NotEmpty()}, "GateWayId": {NotEmpty()}}
	ExecuteScriptVerify    = Rules{"GateWayId": {NotEmpty()}, "ScriptContent": {NotEmpty()}}
	AddFileVerify          = Rules{"GateWayId": {NotEmpty()}, "FileType": {NotEmpty()}, "AddType": {NotEmpty()}}
	ModifyFileVerify       = Rules{"Id": {NotEmpty()}, "GateWayId": {NotEmpty()}, "FileType": {NotEmpty()}, "FileUrl": {NotEmpty()}}
	ForcedToStopVerify     = Rules{"GateWayId": {NotEmpty()}}
	OpIdVerify             = Rules{"OpId": {NotEmpty()}}
	ScriptStopVerify       = Rules{"ID": {NotEmpty()}, "GateWayId": {NotEmpty()}}
	SectorTaskStatusVerify = Rules{"Id": {NotEmpty()}, "Status": {NotEmpty()}}
	WorkerCarStatusVerify  = Rules{"Id": {NotEmpty()}, "TaskStatus": {NotEmpty()}}
	WorkerNumVerify        = Rules{"Id": {NotEmpty()}, "WorkerTaskNum": {NotEmpty()}}
)
